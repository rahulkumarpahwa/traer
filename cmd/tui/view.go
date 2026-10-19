package main

import (
	"strconv"
	"strings"
)

func (m model) View() string {
	if m.screen == screenLogin {
		return m.viewLogin()
	}
	return m.viewMain()
}

// backendIndicator renders the base URL with a reachability hint.
func (m model) backendIndicator() string {
	state := "unreachable"
	if m.backendOK {
		state = "ok"
	}
	return "backend: " + m.client.BaseURL + " (" + state + ")"
}

func (m model) viewLogin() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("traer"))
	b.WriteString("\n")
	if m.signupMode {
		b.WriteString(subtitleStyle.Render("Create an account (ctrl+t to switch to sign in)"))
	} else {
		b.WriteString(subtitleStyle.Render("Sign in to your media workspace (ctrl+t to create an account)"))
	}
	b.WriteString("\n\n")

	b.WriteString(m.renderField("Email", m.email.View()))
	b.WriteString("\n")
	b.WriteString(m.renderField("Password", m.password.View()))

	b.WriteString("\n\n")
	if m.loading {
		b.WriteString(m.spinner.View() + " Working…")
		b.WriteString("\n")
	}
	if m.err != "" {
		b.WriteString(errorStyle.Render("✗ " + m.err))
		b.WriteString("\n")
	}
	if m.notice != "" {
		b.WriteString(successStyle.Render("✓ " + m.notice))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("enter submit · tab/↑↓ move · ctrl+t toggle sign in/up · ctrl+c quit"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(m.backendIndicator()))

	return b.String()
}

func (m model) renderField(label, value string) string {
	return labelStyle.Render(label) + value
}

func (m model) viewMain() string {
	var b strings.Builder

	// Header with the selected process type for new jobs.
	b.WriteString(titleStyle.Render("traer") + subtitleStyle.Render("  ·  YouTube → media"))
	b.WriteString("  ")
	b.WriteString(helpStyle.Render(m.backendIndicator()))
	b.WriteString("\n\n")

	// Process type row.
	b.WriteString(labelStyle.Render("Process") +
		activeTagStyle.Render(m.currentProcessType()) + " " +
		helpStyle.Render("ctrl+p to cycle"))
	b.WriteString("\n\n")

	// URL input box.
	urlBox := boxStyle
	if m.mainFocus == focusURLs {
		urlBox = focusedBoxStyle
	}
	b.WriteString(urlBox.Render(m.urls.View()))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("ctrl+s queue jobs · ctrl+p process type · tab switch focus · ctrl+r refresh · ctrl+o logout"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Only download content you own or are authorized to use."))
	b.WriteString("\n\n")

	// Jobs section.
	b.WriteString(subtitleStyle.Render("Jobs"))
	b.WriteString("\n")

	if len(m.jobs) == 0 {
		b.WriteString(helpStyle.Render("  No jobs yet. Paste links above and press ctrl+s."))
		b.WriteString("\n")
	} else {
		for i, job := range m.jobs {
			b.WriteString(m.renderJob(i, job))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	if m.loading {
		b.WriteString(m.spinner.View() + " Working…")
		b.WriteString("\n")
	}
	if m.err != "" {
		b.WriteString(errorStyle.Render("✗ " + m.err))
		b.WriteString("\n")
	}
	if m.notice != "" {
		b.WriteString(successStyle.Render("✓ " + m.notice))
		b.WriteString("\n")
	}

	return b.String()
}

func (m model) renderJob(index int, job Job) string {
	cursor := "  "
	if m.mainFocus == focusJobs && index == m.jobIndex {
		cursor = "▸ "
	}

	name := strings.TrimSpace(job.URL)
	if name == "" {
		name = job.ID
	}
	name = truncate(name, 56)

	status := statusStyle(job.Status).Render(padRight(job.Status, 12))
	kind := helpStyle.Render(padRight(job.ProcessType, 14))
	line := cursor + status + " " + kind + " " + name

	// Secondary line: upload state, retries and creation time. The live API
	// does not expose per-job progress or error text, so none is shown.
	var detail string
	switch {
	case job.Status == statusFailed:
		detail = errorStyle.Render("failed after " + itoa(job.Retries) + " retries")
	case job.Status == statusUploading:
		detail = warnStyle.Render("uploading…")
	case job.Uploaded:
		detail = successStyle.Render("uploaded")
	case job.Status == statusCompleted:
		detail = successStyle.Render("completed")
	}

	meta := job.CreatedAt.Format("2006-01-02 15:04")
	if job.Retries > 0 && job.Status != statusFailed {
		meta += " · retries " + itoa(job.Retries)
	}
	row2 := helpStyle.Render(meta)
	if detail != "" {
		row2 = detail + "  " + row2
	}
	line += "\n" + cursor + "  " + row2

	return line
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	if max == 1 {
		return "…"
	}
	return s[:max-1] + "…"
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
