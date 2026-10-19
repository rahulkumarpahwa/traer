package main

import (
	"strconv"
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

// maxTrackedJobs caps how many jobs the TUI keeps in memory for display.
const maxTrackedJobs = 50

// ---- async result messages ----

type healthResultMsg struct{ err error }
type loginResultMsg struct {
	userID string
	err    error
}
type signupResultMsg struct{ err error }
type logoutResultMsg struct{ err error }
type queuedMsg struct {
	count int
	err   error
}
type refreshedMsg struct {
	jobs []Job
	err  error
}
type tickMsg time.Time

type model struct {
	client *Client
	screen screen

	// login / signup — the backend only accepts email+password credentials,
	// so both modes share the same two fields.
	email      textinput.Model
	password   textinput.Model
	loginFocus int // 0 email, 1 password
	signupMode bool

	// main screen
	urls         textarea.Model
	processIndex int
	jobs         []Job
	jobIndex     int
	mainFocus    focusArea

	// session
	userID    string
	backendOK bool

	spinner spinner.Model
	loading bool
	err     string
	notice  string
	width   int
	height  int
}

func newModel(client *Client) model {
	email := textinput.New()
	email.Placeholder = "you@example.com"
	email.CharLimit = 128
	email.Width = 40
	email.Focus()

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
	sp.Style = spinnerStyle

	return model{
		client:   client,
		screen:   screenLogin,
		email:    email,
		password: password,
		urls:     urls,
		spinner:  sp,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.checkBackend())
}

// focusLoginField focuses the real fields on the model. It must use the pointer
// receiver and address m's own fields directly; going through a value-receiver
// helper would return pointers into a discarded copy and focus nothing.
func (m *model) focusLoginField() {
	m.email.Blur()
	m.password.Blur()
	if m.loginFocus == 0 {
		m.email.Focus()
	} else {
		m.password.Focus()
	}
}

// currentProcessType returns the process_type selected for new jobs.
func (m model) currentProcessType() string {
	if m.processIndex < 0 || m.processIndex >= len(processTypes) {
		return processTypes[0]
	}
	return processTypes[m.processIndex]
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
		// Poll the job list on the main screen.
		if m.screen == screenMain {
			return m, tea.Batch(m.refreshJobs(), tickCmd())
		}
		return m, nil

	case healthResultMsg:
		m.backendOK = msg.err == nil
		return m, nil

	case loginResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.userID = msg.userID
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

	case logoutResultMsg:
		m.loading = false
		m.userID = ""
		m.jobs = nil
		m.jobIndex = 0
		m.screen = screenLogin
		m.signupMode = false
		m.loginFocus = 0
		m.email.SetValue("")
		m.password.SetValue("")
		m.focusLoginField()
		if msg.err != nil {
			// The backend 401s when the cookie is already gone; either way
			// there is no session left, so still land on the login screen.
			m.err = "Signed out (" + msg.err.Error() + ")"
		} else {
			m.err = ""
			m.notice = "Signed out."
		}
		return m, textinput.Blink

	case queuedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.notice = "Queued " + strconv.Itoa(msg.count) + " job(s) as " + m.currentProcessType() + "."
		m.urls.SetValue("")
		return m, m.refreshJobs()

	case refreshedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.setJobs(msg.jobs)
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// ctrl+c always quits. The bare 'q' key is intentionally NOT a global quit
	// here because it is a valid character in emails/passwords/URLs; quitting
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
		m.loginFocus = (m.loginFocus + dir + 2) % 2
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
	if m.loginFocus == 0 {
		m.email, cmd = m.email.Update(msg)
	} else {
		m.password, cmd = m.password.Update(msg)
	}
	return m, cmd
}

// submitLogin validates the email+password form and either registers
// (POST /users) or signs in (POST /login), depending on the active mode.
func (m model) submitLogin() (tea.Model, tea.Cmd) {
	m.err = ""
	email := strings.TrimSpace(m.email.Value())
	password := m.password.Value()

	if !strings.Contains(email, "@") {
		m.err = "Enter a valid email address."
		return m, nil
	}
	if password == "" {
		m.err = "Password is required."
		return m, nil
	}

	m.loading = true
	client := m.client
	if m.signupMode {
		return m, func() tea.Msg {
			return signupResultMsg{err: client.Signup(email, password)}
		}
	}
	return m, func() tea.Msg {
		userID, err := client.Login(email, password)
		return loginResultMsg{userID: userID, err: err}
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
		return m.submitJobs()
	case "ctrl+r":
		m.loading = true
		return m, tea.Batch(m.refreshJobs(), m.checkBackend())
	case "ctrl+p":
		m.processIndex = (m.processIndex + 1) % len(processTypes)
		m.notice = "Process type: " + m.currentProcessType()
		return m, nil
	case "ctrl+o":
		m.loading = true
		client := m.client
		return m, func() tea.Msg {
			return logoutResultMsg{err: client.Logout()}
		}
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
	}
	return m, nil
}

// submitJobs validates the pasted links and queues them through
// POST /jobs/create with the selected process_type.
func (m model) submitJobs() (tea.Model, tea.Cmd) {
	urls := splitUrls(m.urls.Value())
	if len(urls) == 0 {
		m.err = "Add at least one YouTube or YouTube Music link."
		return m, nil
	}
	for _, u := range urls {
		if !IsYouTubeURL(u) {
			m.err = "Not a YouTube link: " + u
			return m, nil
		}
	}
	m.err = ""
	m.loading = true
	client := m.client
	processType := m.currentProcessType()
	return m, func() tea.Msg {
		_, err := client.CreateJobs(urls, processType)
		return queuedMsg{count: len(urls), err: err}
	}
}

// refreshJobs fetches the full job list (GET /jobs) and returns it as a
// message.
func (m model) refreshJobs() tea.Cmd {
	client := m.client
	return func() tea.Msg {
		jobs, err := client.ListJobs()
		return refreshedMsg{jobs: jobs, err: err}
	}
}

// checkBackend pings GET / and reports reachability.
func (m model) checkBackend() tea.Cmd {
	client := m.client
	return func() tea.Msg {
		return healthResultMsg{err: client.Health()}
	}
}

// setJobs replaces the tracked job list with the server's view (the backend
// returns every job, not just active ones), keeping the list bounded and the
// selection index valid.
func (m *model) setJobs(jobs []Job) {
	if len(jobs) > maxTrackedJobs {
		jobs = jobs[len(jobs)-maxTrackedJobs:]
	}
	m.jobs = jobs
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
