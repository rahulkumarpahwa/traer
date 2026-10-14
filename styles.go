package main

import "github.com/charmbracelet/lipgloss"

var (
	appName = "traer"

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("78")).
			Bold(true)

	warnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Width(14)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("239")).
			Padding(0, 1)

	focusedBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")).
			Padding(0, 1)

	statusQueued    = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	statusRunning   = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	statusUploading = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	statusDone      = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	statusFailed    = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))

	jobSourceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

// statusStyle returns the colorized style for a job status pill.
func statusStyle(status string) lipgloss.Style {
	switch status {
	case "running":
		return statusRunning
	case "uploading":
		return statusUploading
	case "done":
		return statusDone
	case "failed":
		return statusFailed
	default:
		return statusQueued
	}
}

// progressBar renders a fixed-width textual progress bar for a 0–100 value.
func progressBar(progress float64, width int) string {
	if width < 4 {
		width = 4
	}
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	filled := int((progress / 100) * float64(width))
	if filled > width {
		filled = width
	}
	bar := lipgloss.NewStyle().Foreground(lipgloss.Color("78")).Render(repeat("█", filled))
	track := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(repeat("░", width-filled))
	return bar + track
}

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, 0, len(s)*n)
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}
