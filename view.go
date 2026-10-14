package main

import (
	"strings"
)

func (m model) View() string {
	if m.screen == screenLogin {
		return m.viewLogin()
	}
	return m.viewMain()
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

	if m.signupMode {
		b.WriteString(m.renderField("Username", m.username.View()))
		b.WriteString("\n")
		b.WriteString(m.renderField("Email", m.email.View()))
		b.WriteString("\n")
		b.WriteString(m.renderField("Password", m.password.View()))
	} else {
		b.WriteString(m.renderField("Identifier", m.identifier.View()))
		b.WriteString("\n")
		b.WriteString(m.renderField("Password", m.password.View()))
	}

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
	b.WriteString(helpStyle.Render("enter submit · tab/↑↓ move · ctrl+t toggle sign in/up · q quit"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("backend: " + m.client.BaseURL))

	return b.String()
}

func (m model) renderField(label, value string) string {
	return labelStyle.Render(label) + value
}

func (m model) viewMain() string {
	var b strings.Builder

	// Header
	header := titleStyle.Render("traer") + subtitleStyle.Render("  ·  Music → Nextcloud")
	b.WriteString(header)
	b.WriteString("\n\n")

	// URL input box.
	urlBox := boxStyle
	if m.mainFocus == focusURLs {
		urlBox = focusedBoxStyle
	}
	b.WriteString(urlBox.Render(m.urls.View()))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("ctrl+s queue to Nextcloud · tab switch focus · ctrl+r refresh · q quit"))
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
		barWidth := 24
		if m.width > 0 && m.width-60 > 8 {
			barWidth = m.width - 60
			if barWidth > 48 {
				barWidth = 48
			}
		}
		for i, job := range m.jobs {
			b.WriteString(m.renderJob(i, job, barWidth))
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

func (m model) renderJob(index int, job Job, barWidth int) string {
	cursor := "  "
	if m.mainFocus == focusJobs && index == m.jobIndex {
		cursor = "▸ "
	}

	name := job.FileName()
	if name == "" {
		name = job.URL
	}
	if name == "" {
		name = job.ID
	}
	// Keep names from wrapping the layout.
	if len(name) > 44 {
		name = name[:41] + "…"
	}

	status := statusStyle(job.Status).Render(padRight(job.Status, 9))

	var line string
	if job.Status == "done" || job.Status == "failed" {
		line = cursor + padRight(name, 44) + " " + status
	} else {
		line = cursor + padRight(name, 44) + " " + status + " " +
			progressBar(job.Progress, barWidth) + " " + itoa(int(job.Progress)) + "%"
	}

	var detail string
	switch {
	case job.Uploaded:
		detail = successStyle.Render("uploaded to Nextcloud")
	case job.UploadError != "":
		detail = errorStyle.Render(job.UploadError)
	case job.Error != "":
		detail = errorStyle.Render(job.Error)
	case job.Status == "uploading":
		detail = warnStyle.Render("uploading to Nextcloud…")
	case job.FileName() != "":
		detail = jobSourceStyle.Render(job.FileName())
	}

	if detail != "" {
		line += "\n" + cursor + "  " + detail
	}
	if job.Status == "failed" {
		line += "\n" + cursor + "  " + helpStyle.Render("press r to retry")
	}
	return line
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
