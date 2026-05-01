package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type App struct {
	ctx      context.Context
	settings SettingsPayload
	profile  Profile
}

type AppState struct {
	AppName      string          `json:"appName"`
	RecentLinks  []string        `json:"recentLinks"`
	Settings     SettingsPayload `json:"settings"`
	Profile      Profile         `json:"profile"`
	Suggestions  []string        `json:"suggestions"`
	Capabilities map[string]bool `json:"capabilities"`
}

type SettingsPayload struct {
	WhisperURL string `json:"whisperUrl"`
	OllamaURL  string `json:"ollamaUrl"`
	CloudURL   string `json:"cloudUrl"`
	OutputDir  string `json:"outputDir"`
	Theme      string `json:"theme"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Profile struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Plan       string `json:"plan"`
	Bio        string `json:"bio"`
	LastActive string `json:"lastActive"`
}

type StartJobPayload struct {
	URL    string `json:"url"`
	Kind   string `json:"kind"`
	Option string `json:"option"`
}

type JobResponse struct {
	ID        string   `json:"id"`
	Kind      string   `json:"kind"`
	Option    string   `json:"option"`
	Title     string   `json:"title"`
	Message   string   `json:"message"`
	Logs      []string `json:"logs"`
	CreatedAt string   `json:"createdAt"`
}

func NewApp() *App {
	return &App{
		settings: SettingsPayload{
			WhisperURL: "http://localhost:9000/transcribe",
			OllamaURL:  "http://localhost:11434",
			CloudURL:   "https://cloud.traer.app/upload",
			OutputDir:  "outputs",
			Theme:      "midnight",
		},
		profile: Profile{
			Name:       "Traer Operator",
			Email:      "operator@traer.local",
			Plan:       "Studio",
			Bio:        "Runs clean media workflows with practical defaults.",
			LastActive: time.Now().Format(time.RFC3339),
		},
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetInitialState() AppState {
	return AppState{
		AppName:     "traer",
		RecentLinks: []string{},
		Settings:    a.settings,
		Profile:     a.profile,
		Suggestions: []string{
			"Paste a YouTube, Loom, or podcast URL to start.",
			"Use Transcribe for clean markdown notes.",
			"Ship audio or video to cloud storage once processing completes.",
		},
		Capabilities: map[string]bool{
			"backendBound": a.ctx != nil,
			"mockMode":     true,
		},
	}
}

func (a *App) Login(payload LoginPayload) User {
	name := strings.TrimSpace(payload.Username)
	if name == "" {
		name = "Operator"
	}

	return User{
		Name:  name,
		Email: fmt.Sprintf("%s@traer.local", strings.ToLower(strings.ReplaceAll(name, " ", "."))),
		Role:  "Admin",
	}
}

func (a *App) SaveSettings(payload SettingsPayload) SettingsPayload {
	a.settings = payload
	return a.settings
}

func (a *App) GetProfile() Profile {
	a.profile.LastActive = time.Now().Format(time.RFC3339)
	return a.profile
}

func (a *App) StartJob(payload StartJobPayload) JobResponse {
	now := time.Now()
	actionLabel := map[string]string{
		"transcribe": "Transcript",
		"audio":      "Audio package",
		"video":      "Video export",
		"cloud":      "Cloud sync",
	}[payload.Kind]
	if actionLabel == "" {
		actionLabel = "Workflow"
	}

	return JobResponse{
		ID:        fmt.Sprintf("%s-%d", payload.Kind, now.UnixMilli()),
		Kind:      payload.Kind,
		Option:    payload.Option,
		Title:     fmt.Sprintf("%s queued", actionLabel),
		Message:   fmt.Sprintf("%s for %s is now processing with preset %s.", actionLabel, payload.URL, payload.Option),
		CreatedAt: now.Format(time.RFC3339),
		Logs: []string{
			fmt.Sprintf("[%s] traer accepted %s request", now.Format("15:04:05"), payload.Kind),
			fmt.Sprintf("[%s] source resolved: %s", now.Add(1*time.Second).Format("15:04:05"), payload.URL),
			fmt.Sprintf("[%s] preset selected: %s", now.Add(2*time.Second).Format("15:04:05"), payload.Option),
		},
	}
}
