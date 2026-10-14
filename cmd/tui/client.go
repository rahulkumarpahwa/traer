package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// Client is a thin HTTP client for the Traer backend. It mirrors the endpoints
// used by the previous Wails bridge and keeps the JWT session in a cookie jar,
// exactly like the backend issues it (an HttpOnly "token" cookie).
type Client struct {
	BaseURL string
	http    *http.Client
}

// Job is the subset of the backend job model the TUI needs to render status,
// progress, upload state, and errors.
type Job struct {
	ID          string  `json:"id"`
	URL         string  `json:"url"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	Progress    float64 `json:"progress"`
	Output      string  `json:"output"`
	Error       string  `json:"error"`
	Uploaded    bool    `json:"uploaded"`
	UploadError string  `json:"uploadError"`
}

// FileName returns the base name of the job output, if any.
func (j Job) FileName() string {
	if strings.TrimSpace(j.Output) == "" {
		return ""
	}
	trimmed := strings.TrimSpace(j.Output)
	parts := strings.FieldsFunc(trimmed, func(r rune) bool {
		return r == '/' || r == '\\'
	})
	if len(parts) == 0 {
		return trimmed
	}
	return parts[len(parts)-1]
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

// Login authenticates with a username/email + password and stores the session
// cookie returned by the backend.
func (c *Client) Login(identifier, password string) error {
	payload := map[string]string{"password": password}
	if strings.Contains(identifier, "@") {
		payload["email"] = identifier
	} else {
		payload["username"] = identifier
	}
	_, err := c.doJSON(http.MethodPost, "/login", payload, nil)
	return err
}

// Signup registers a new account. It does not sign you in automatically.
func (c *Client) Signup(username, email, password string) error {
	payload := map[string]string{"password": password}
	if username != "" {
		payload["username"] = username
	}
	if email != "" {
		payload["email"] = email
	}
	_, err := c.doJSON(http.MethodPost, "/users", payload, nil)
	return err
}

// CreateMusicJobs queues one or more YouTube/YouTube Music URLs for MP3
// extraction and upload to Nextcloud. It returns the created jobs.
func (c *Client) CreateMusicJobs(urls []string) ([]Job, error) {
	payload := map[string]interface{}{
		"urls": urls,
		"type": "music",
	}

	var batch struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Jobs   []Job  `json:"jobs"`
	}
	status, err := c.doJSON(http.MethodPost, "/jobs/create", payload, &batch)
	if err != nil {
		return nil, err
	}
	_ = status

	if len(batch.Jobs) > 0 {
		return batch.Jobs, nil
	}

	// Backend returns a single job object when only one URL was accepted.
	var single Job
	if _, err := c.doJSON(http.MethodPost, "/jobs/create", payload, &single); err == nil && single.ID != "" {
		return []Job{single}, nil
	}
	return nil, fmt.Errorf("backend did not return any jobs")
}

// ActiveJobs returns jobs that are queued, running, or uploading.
func (c *Client) ActiveJobs() ([]Job, error) {
	var jobs []Job
	_, err := c.doJSON(http.MethodGet, "/jobs/active", nil, &jobs)
	return jobs, err
}

// JobStatus returns the current state of a single job.
func (c *Client) JobStatus(id string) (Job, error) {
	var job Job
	_, err := c.doJSON(http.MethodGet, "/jobs/status?id="+url.QueryEscape(id), nil, &job)
	return job, err
}

// RetryJob re-queues a previously failed job.
func (c *Client) RetryJob(id string) error {
	_, err := c.doJSON(http.MethodPost, "/jobs/retry?id="+url.QueryEscape(id), nil, nil)
	return err
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
