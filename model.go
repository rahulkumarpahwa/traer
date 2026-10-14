package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenLogin screen = iota
	screenMain
)

type focusArea int

const (
	focusURLs focusArea = iota
	focusJobs
)

// ---- async result messages ----

type loginResultMsg struct{ err error }
type signupResultMsg struct{ err error }
type queuedMsg struct {
	jobs []Job
	err  error
}
type refreshedMsg struct {
	jobs []Job
	err  error
}
type retryResultMsg struct {
	id  string
	err error
}
type tickMsg time.Time

type model struct {
	client *Client
	screen screen

	// login / signup
	identifier textinput.Model
	username   textinput.Model
	email      textinput.Model
	password   textinput.Model
	signupMode bool
	loginFocus int // 0 identifier/username, 1 email (signup), 2 password

	// main
	urls      textarea.Model
	jobs      []Job
	jobIndex  int
	mainFocus focusArea

	spinner spinner.Model
	loading bool
	err     string
	notice  string
	width   int
	height  int
}

func newModel(client *Client) model {
	identifier := textinput.New()
	identifier.Placeholder = "username or email"
	identifier.Focus()
	identifier.CharLimit = 128
	identifier.Width = 40
	username := textinput.New()
	username.Placeholder = "username"
	username.CharLimit = 64
	username.Width = 40

	email := textinput.New()
	email.Placeholder = "you@example.com"
	email.CharLimit = 128
	email.Width = 40

	password := textinput.New()
	password.Placeholder = "password"
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'
	password.CharLimit = 128
	password.Width = 40

	urls := textarea.New()
	urls.Placeholder = "Paste one or more YouTube / YouTube Music links (one per line)"
	urls.CharLimit = 0
	urls.SetWidth(72)
	urls.SetHeight(6)

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = statusUploading

	return model{
		client:     client,
		screen:     screenLogin,
		identifier: identifier,
		username:   username,
		email:      email,
		password:   password,
		urls:       urls,
		spinner:    sp,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// loginFieldCount returns the number of active login inputs for the current
// mode. It deliberately returns a count (not pointers into a copy) so callers
// never mutate a throwaway value.
func (m model) loginFieldCount() int {
	if m.signupMode {
		return 3 // username, email, password
	}
	return 2 // identifier, password
}

// focusLoginField focuses the real fields on the model. It must use the pointer
// receiver and address m's own fields directly; going through a value-receiver
// helper would return pointers into a discarded copy and focus nothing.
func (m *model) focusLoginField() {
	m.identifier.Blur()
	m.username.Blur()
	m.email.Blur()
	m.password.Blur()

	if m.signupMode {
		switch m.loginFocus {
		case 0:
			m.username.Focus()
		case 1:
			m.email.Focus()
		default:
			m.password.Focus()
		}
		return
	}

	switch m.loginFocus {
	case 0:
		m.identifier.Focus()
	default:
		m.password.Focus()
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width > 12 {
			m.urls.SetWidth(m.width - 12)
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tickMsg:
		// Poll active jobs on the main screen.
		if m.screen == screenMain {
			return m, tea.Batch(m.refreshJobs(), tickCmd())
		}
		return m, nil

	case loginResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.screen = screenMain
		m.notice = "Signed in. Paste links and press ctrl+s to queue."
		m.urls.Focus()
		return m, tea.Batch(m.refreshJobs(), tickCmd())

	case signupResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.notice = "Account created. Sign in to continue."
		m.signupMode = false
		m.loginFocus = 0
		m.focusLoginField()
		return m, nil

	case queuedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.notice = "Queued " + itoa(len(msg.jobs)) + " music job(s) for Nextcloud upload."
		m.mergeJobs(msg.jobs)
		return m, m.refreshJobs()

	case refreshedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.mergeJobs(msg.jobs)
		return m, nil

	case retryResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.notice = "Retrying job " + msg.id
		return m, m.refreshJobs()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// ctrl+c always quits. The bare 'q' key is intentionally NOT a global quit
	// here because it is a valid character in usernames/passwords/URLs; quitting
	// on 'q' is handled contextually in the job-list view instead.
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	if m.screen == screenLogin {
		return m.handleLoginKey(msg)
	}
	return m.handleMainKey(msg)
}

func (m model) handleLoginKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down", "up":
		dir := 1
		if msg.String() == "up" {
			dir = -1
		}
		count := m.loginFieldCount()
		m.loginFocus = (m.loginFocus + dir + count) % count
		m.focusLoginField()
		return m, textinput.Blink
	case "ctrl+t":
		m.signupMode = !m.signupMode
		m.loginFocus = 0
		m.err = ""
		m.focusLoginField()
		return m, textinput.Blink
	case "enter":
		return m.submitLogin()
	}

	var cmd tea.Cmd
	switch m.loginFocus {
	case 0:
		if m.signupMode {
			m.username, cmd = m.username.Update(msg)
		} else {
			m.identifier, cmd = m.identifier.Update(msg)
		}
	case 1:
		if m.signupMode {
			m.email, cmd = m.email.Update(msg)
		} else {
			m.password, cmd = m.password.Update(msg)
		}
	case 2:
		m.password, cmd = m.password.Update(msg)
	}
	return m, cmd
}

func (m model) submitLogin() (tea.Model, tea.Cmd) {
	m.err = ""
	m.loading = true
	client := m.client

	if m.signupMode {
		username := strings.TrimSpace(m.username.Value())
		email := strings.TrimSpace(m.email.Value())
		password := m.password.Value()
		return m, func() tea.Msg {
			err := client.Signup(username, email, password)
			return signupResultMsg{err: err}
		}
	}

	identifier := strings.TrimSpace(m.identifier.Value())
	password := m.password.Value()
	return m, func() tea.Msg {
		err := client.Login(identifier, password)
		return loginResultMsg{err: err}
	}
}

func (m model) handleMainKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		if m.mainFocus == focusURLs {
			m.mainFocus = focusJobs
			m.urls.Blur()
		} else {
			m.mainFocus = focusURLs
			m.urls.Focus()
		}
		return m, textinput.Blink
	case "ctrl+s":
		return m.submitMusic()
	case "ctrl+r":
		m.loading = true
		return m, m.refreshJobs()
	}

	if m.mainFocus == focusURLs {
		var cmd tea.Cmd
		m.urls, cmd = m.urls.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "q":
		// Only quit from the job list view (not while editing URLs).
		return m, tea.Quit
	case "up", "k":
		if m.jobIndex > 0 {
			m.jobIndex--
		}
	case "down", "j":
		if m.jobIndex < len(m.jobs)-1 {
			m.jobIndex++
		}
	case "r":
		if m.jobIndex >= 0 && m.jobIndex < len(m.jobs) {
			id := m.jobs[m.jobIndex].ID
			m.loading = true
			client := m.client
			return m, func() tea.Msg {
				err := client.RetryJob(id)
				return retryResultMsg{id: id, err: err}
			}
		}
	}
	return m, nil
}

func (m model) submitMusic() (tea.Model, tea.Cmd) {
	urls := splitUrls(m.urls.Value())
	if len(urls) == 0 {
		m.err = "Add at least one YouTube or YouTube Music link."
		return m, nil
	}
	m.err = ""
	m.loading = true
	client := m.client
	return m, func() tea.Msg {
		jobs, err := client.CreateMusicJobs(urls)
		return queuedMsg{jobs: jobs, err: err}
	}
}

// refreshJobs fetches active jobs and returns them as a message.
func (m model) refreshJobs() tea.Cmd {
	client := m.client
	return func() tea.Msg {
		jobs, err := client.ActiveJobs()
		return refreshedMsg{jobs: jobs, err: err}
	}
}

// mergeJobs folds refreshed/queued jobs into the tracked list without dropping
// jobs that are no longer "active" (so completed/failed ones remain visible).
func (m *model) mergeJobs(incoming []Job) {
	byID := make(map[string]Job, len(incoming))
	for _, j := range incoming {
		byID[j.ID] = j
	}

	seen := map[string]bool{}
	merged := make([]Job, 0, len(m.jobs)+len(incoming))
	for _, existing := range m.jobs {
		if j, ok := byID[existing.ID]; ok {
			merged = append(merged, j)
			seen[existing.ID] = true
		} else {
			merged = append(merged, existing)
		}
	}
	for _, j := range incoming {
		if !seen[j.ID] {
			merged = append(merged, j)
		}
	}

	if len(merged) > 50 {
		merged = merged[len(merged)-50:]
	}
	m.jobs = merged
	if m.jobIndex >= len(m.jobs) {
		m.jobIndex = len(m.jobs) - 1
	}
	if m.jobIndex < 0 {
		m.jobIndex = 0
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
