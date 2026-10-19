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

	activeTagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	// One style per job status returned by GET /jobs (types.JobStatus).
	styleNotStarted = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	styleQueued     = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	styleRunning    = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	styleUploading  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	styleUploaded   = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	styleCompleted  = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	styleFailed     = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
)

// statusStyle returns the colorized style for a job status pill.
func statusStyle(status string) lipgloss.Style {
	switch status {
	case statusRunning:
		return styleRunning
	case statusUploading:
		return styleUploading
	case statusUploaded:
		return styleUploaded
	case statusCompleted:
		return styleCompleted
	case statusFailed:
		return styleFailed
	case statusNotStarted:
		return styleNotStarted
	default:
		return styleQueued
	}
}
