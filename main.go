package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	backendURL := flag.String("backend", envOr("TRAER_BACKEND_URL", "http://localhost:8080"), "Traer backend base URL")
	flag.Parse()

	client := NewClient(*backendURL)

	program := tea.NewProgram(
		newModel(client),
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "traer: failed to run TUI:", err)
		os.Exit(1)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
