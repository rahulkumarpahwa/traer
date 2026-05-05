package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const backendBaseURL = "http://localhost:8080"

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
	ID          string   `json:"id"`
	Kind        string   `json:"kind"`
	Option      string   `json:"option"`
	Title       string   `json:"title"`
	Message     string   `json:"message"`
	Logs        []string `json:"logs"`
	CreatedAt   string   `json:"createdAt"`
	Status      string   `json:"status,omitempty"`
	Progress    float64  `json:"progress,omitempty"`
	Output      string   `json:"output,omitempty"`
	FileName    string   `json:"fileName,omitempty"`
	DownloadURL string   `json:"downloadUrl,omitempty"`
	SourceURL   string   `json:"sourceUrl,omitempty"`
}

type backendJob struct {
	ID       string  `json:"id"`
	URL      string  `json:"url"`
	Type     string  `json:"type"`
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
	Output   string  `json:"output"`
	Error    string  `json:"error"`
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
			"mockMode":     false,
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
	if payload.Kind != "audio" && payload.Kind != "video" {
		return a.makeMockJob(payload)
	}

	reqBody, _ := json.Marshal(map[string]string{
		"url":  payload.URL,
		"type": payload.Kind,
	})

	var created backendJob
	if err := doBackendJSON(http.MethodPost, "/jobs/create", bytes.NewReader(reqBody), &created); err != nil {
		return a.makeErrorJob(payload, err)
	}

	now := time.Now()
	return JobResponse{
		ID:          created.ID,
		Kind:        payload.Kind,
		Option:      payload.Option,
		Title:       labelForKind(payload.Kind) + " queued",
		Message:     fmt.Sprintf("%s for %s is now queued.", labelForKind(payload.Kind), payload.URL),
		CreatedAt:   now.Format(time.RFC3339),
		Status:      created.Status,
		Progress:    created.Progress,
		Output:      created.Output,
		FileName:    fileNameFromOutput(created.Output),
		DownloadURL: buildDownloadURL(created.ID),
		SourceURL:   payload.URL,
		Logs: []string{
			fmt.Sprintf("[%s] backend accepted %s request", now.Format("15:04:05"), payload.Kind),
			fmt.Sprintf("[%s] source queued: %s", now.Add(1*time.Second).Format("15:04:05"), payload.URL),
			fmt.Sprintf("[%s] preset selected: %s", now.Add(2*time.Second).Format("15:04:05"), payload.Option),
		},
	}
}

func (a *App) GetActiveJobs() []JobResponse {
	var jobs []backendJob
	if err := doBackendJSON(http.MethodGet, "/jobs/active", nil, &jobs); err != nil {
		return []JobResponse{}
	}

	active := make([]JobResponse, 0, len(jobs))
	for _, job := range jobs {
		active = append(active, mapBackendJob(job))
	}

	return active
}

func (a *App) GetJobStatus(id string) JobResponse {
	escapedID := url.QueryEscape(strings.TrimSpace(id))
	var job backendJob
	if err := doBackendJSON(http.MethodGet, "/jobs/status?id="+escapedID, nil, &job); err != nil {
		return JobResponse{
			ID:      id,
			Status:  "failed",
			Message: err.Error(),
		}
	}

	return mapBackendJob(job)
}

func (a *App) GetDownloadURL(id string) string {
	return buildDownloadURL(id)
}

func (a *App) DownloadJob(id string) (string, error) {
	job := a.GetJobStatus(id)
	if strings.TrimSpace(job.ID) == "" || strings.TrimSpace(job.Status) == "failed" {
		return "", fmt.Errorf("unable to load job details")
	}
	if strings.TrimSpace(job.Output) == "" && strings.TrimSpace(job.FileName) == "" {
		return "", fmt.Errorf("file is not ready yet")
	}

	defaultName := job.FileName
	if defaultName == "" {
		defaultName = fileNameFromOutput(job.Output)
	}
	if defaultName == "" {
		defaultName = "traer-output"
	}

	targetPath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save Output File",
		DefaultFilename: defaultName,
	})
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(targetPath) == "" {
		return "", nil
	}

	resp, err := http.Get(buildDownloadURL(id))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		var payload map[string]interface{}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&payload); decodeErr == nil {
			if message, ok := payload["error"].(string); ok && message != "" {
				return "", fmt.Errorf(message)
			}
			if message, ok := payload["message"].(string); ok && message != "" {
				return "", fmt.Errorf(message)
			}
		}
		return "", fmt.Errorf("download request failed with status %d", resp.StatusCode)
	}

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, resp.Body); err != nil {
		return "", err
	}

	return targetPath, nil
}

func (a *App) makeMockJob(payload StartJobPayload) JobResponse {
	now := time.Now()
	actionLabel := labelForKind(payload.Kind)

	return JobResponse{
		ID:        fmt.Sprintf("%s-%d", payload.Kind, now.UnixMilli()),
		Kind:      payload.Kind,
		Option:    payload.Option,
		Title:     fmt.Sprintf("%s queued", actionLabel),
		Message:   fmt.Sprintf("%s for %s is now processing with preset %s.", actionLabel, payload.URL, payload.Option),
		CreatedAt: now.Format(time.RFC3339),
		SourceURL: payload.URL,
		Logs: []string{
			fmt.Sprintf("[%s] traer accepted %s request", now.Format("15:04:05"), payload.Kind),
			fmt.Sprintf("[%s] source resolved: %s", now.Add(1*time.Second).Format("15:04:05"), payload.URL),
			fmt.Sprintf("[%s] preset selected: %s", now.Add(2*time.Second).Format("15:04:05"), payload.Option),
		},
	}
}

func (a *App) makeErrorJob(payload StartJobPayload, err error) JobResponse {
	now := time.Now()
	return JobResponse{
		ID:        fmt.Sprintf("%s-%d", payload.Kind, now.UnixMilli()),
		Kind:      payload.Kind,
		Option:    payload.Option,
		Title:     fmt.Sprintf("%s failed", labelForKind(payload.Kind)),
		Message:   err.Error(),
		CreatedAt: now.Format(time.RFC3339),
		Status:    "failed",
		SourceURL: payload.URL,
	}
}

func labelForKind(kind string) string {
	switch kind {
	case "audio":
		return "Audio package"
	case "video":
		return "Video export"
	case "transcribe":
		return "Transcript"
	case "cloud":
		return "Cloud sync"
	default:
		return "Workflow"
	}
}

func doBackendJSON(method, path string, body *bytes.Reader, out interface{}) error {
	var reqBody *bytes.Reader
	if body == nil {
		reqBody = bytes.NewReader(nil)
	} else {
		reqBody = body
	}

	req, err := http.NewRequest(method, backendBaseURL+path, reqBody)
	if err != nil {
		return err
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		var payload map[string]interface{}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&payload); decodeErr == nil {
			if message, ok := payload["error"].(string); ok && message != "" {
				return fmt.Errorf(message)
			}
			if message, ok := payload["message"].(string); ok && message != "" {
				return fmt.Errorf(message)
			}
		}
		return fmt.Errorf("backend request failed with status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func mapBackendJob(job backendJob) JobResponse {
	fileName := fileNameFromOutput(job.Output)
	title := labelForKind(job.Type)
	if fileName != "" {
		title = fileName
	}

	message := fmt.Sprintf("%s task is %s.", labelForKind(job.Type), job.Status)
	if job.Error != "" {
		message = job.Error
	} else if fileName != "" {
		message = fmt.Sprintf("Output file ready: %s", fileName)
	}

	return JobResponse{
		ID:          job.ID,
		Kind:        job.Type,
		Option:      strings.ToUpper(job.Type),
		Title:       title,
		Message:     message,
		Status:      job.Status,
		Progress:    job.Progress,
		Output:      job.Output,
		FileName:    fileName,
		DownloadURL: buildDownloadURL(job.ID),
		SourceURL:   job.URL,
	}
}

func fileNameFromOutput(output string) string {
	if strings.TrimSpace(output) == "" {
		return ""
	}
	trimmed := strings.TrimSpace(output)
	parts := strings.FieldsFunc(trimmed, func(r rune) bool {
		return r == '/' || r == '\\'
	})
	if len(parts) == 0 {
		return trimmed
	}
	return parts[len(parts)-1]
}

func buildDownloadURL(id string) string {
	if strings.TrimSpace(id) == "" {
		return ""
	}
	return backendBaseURL + "/jobs/download?id=" + url.QueryEscape(strings.TrimSpace(id))
}
