package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"
)

// Client is a thin HTTP client for the live Traer backend routes (see
// postman_collection.json):
//
//	GET    /              health check (no auth)
//	POST   /login         {email,password} -> sets HttpOnly "token" cookie
//	POST   /users         {email,password} registration (no auth)
//	POST   /logout        clears the session cookie
//	POST   /jobs/create   {urls,process_type} (auth)
//	GET    /jobs          list all jobs (auth)
//
// The backend routes that are commented out in cmd/api/main.go
// (/jobs/status, /jobs/retry, /jobs/download, /instances/get) intentionally
// have no client methods here so the TUI never calls endpoints that 404.
type Client struct {
	BaseURL string
	http    *http.Client
}

// processTypes mirrors the process_type values accepted by POST /jobs/create
// (see internal/handlers/job.go and internal/types/job.go).
var processTypes = []string{"audio", "video", "playlist", "upload", "transcription"}

// Job status values returned by GET /jobs (internal/types.JobStatus).
const (
	statusNotStarted = "not_started"
	statusQueued     = "queued"
	statusRunning    = "running"
	statusUploading  = "uploading"
	statusUploaded   = "uploaded"
	statusCompleted  = "completed"
	statusFailed     = "failed"
)

// Job mirrors tdd.GetProcess, the element shape of the GET /jobs response.
type Job struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	URL         string    `json:"url"`
	Status      string    `json:"status"`
	ProcessType string    `json:"process_type"`
	Retries     int       `json:"retries"`
	Uploaded    bool      `json:"uploaded"`
	TempPath    string    `json:"temp_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// loginResponse is the JSON body of a successful POST /login.
type loginResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

// youTubeRegex mirrors internal/validators.IsYouTubeURL so obviously bad
// links are rejected before they hit the backend.
var youTubeRegex = regexp.MustCompile(`^(https?://)?(www\.|music\.)?(youtube\.com|youtu\.be)/.+`)

// IsYouTubeURL reports whether raw looks like a YouTube / YouTube Music link.
func IsYouTubeURL(raw string) bool {
	return youTubeRegex.MatchString(raw)
}

// NewClient builds a backend client with a cookie jar for session handling.
func NewClient(baseURL string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		BaseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		http: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
	}
}

// doJSON issues a request and decodes the JSON response into out (if non-nil).
// It returns the HTTP status code and any error.
func (c *Client) doJSON(method, path string, body interface{}, out interface{}) (int, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return 0, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		var payload map[string]interface{}
		message := fmt.Sprintf("backend request failed with status %d", resp.StatusCode)
		if decodeErr := json.NewDecoder(resp.Body).Decode(&payload); decodeErr == nil {
			if m, ok := payload["error"].(string); ok && m != "" {
				message = m
			} else if m, ok := payload["message"].(string); ok && m != "" {
				message = m
			}
		}
		return resp.StatusCode, fmt.Errorf("%s", message)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return resp.StatusCode, err
		}
	}
	return resp.StatusCode, nil
}

// Health pings GET / to verify the backend is reachable.
func (c *Client) Health() error {
	_, err := c.doJSON(http.MethodGet, "/", nil, nil)
	return err
}

// Login authenticates with email + password (the only credential shape the
// backend accepts) and stores the session cookie returned by the backend.
// It returns the user ID reported by the backend.
func (c *Client) Login(email, password string) (string, error) {
	payload := map[string]string{"email": email, "password": password}
	var out loginResponse
	if _, err := c.doJSON(http.MethodPost, "/login", payload, &out); err != nil {
		return "", err
	}
	return out.UserID, nil
}

// Signup registers a new account with email + password. It does not sign you
// in automatically; the backend only sets the session cookie on /login.
func (c *Client) Signup(email, password string) error {
	payload := map[string]string{"email": email, "password": password}
	_, err := c.doJSON(http.MethodPost, "/users", payload, nil)
	return err
}

// Logout clears the session cookie on the backend. The backend returns 401
// when no cookie is present, which callers may treat as "already signed out".
func (c *Client) Logout() error {
	_, err := c.doJSON(http.MethodPost, "/logout", nil, nil)
	return err
}

// CreateJobs queues one or more YouTube URLs for processing with the given
// process type (one of processTypes). The backend responds with a plain JSON
// string message rather than job objects, so callers should refresh the job
// list via ListJobs afterwards. It returns the backend's message.
func (c *Client) CreateJobs(urls []string, processType string) (string, error) {
	payload := map[string]interface{}{
		"urls":         urls,
		"process_type": processType,
	}
	var message string
	_, err := c.doJSON(http.MethodPost, "/jobs/create", payload, &message)
	return message, err
}

// ListJobs returns all jobs known to the backend (GET /jobs).
func (c *Client) ListJobs() ([]Job, error) {
	var jobs []Job
	_, err := c.doJSON(http.MethodGet, "/jobs", nil, &jobs)
	return jobs, err
}

// splitUrls parses a free-form blob of text (newlines, commas, spaces) into a
// de-duplicated list of trimmed URLs while preserving order.
func splitUrls(text string) []string {
	seen := map[string]bool{}
	var out []string
	for _, raw := range strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n' || r == '\r' || r == '\t' || r == ',' || r == ' '
	}) {
		u := strings.TrimSpace(raw)
		if u == "" || seen[u] {
			continue
		}
		seen[u] = true
		out = append(out, u)
	}
	return out
}
