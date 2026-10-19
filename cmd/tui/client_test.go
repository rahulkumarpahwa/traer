package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func timeNowPlusDay() time.Time {
	return time.Now().Add(24 * time.Hour)
}

// newMockBackend mirrors the LIVE Traer backend contract (cmd/api/main.go +
// internal/handlers): POST /login accepts only email+password and sets an
// HttpOnly "token" cookie (Secure:false, Path:/); protected routes require
// that cookie and 401 otherwise; POST /jobs/create expects {urls,
// process_type} and answers with a plain JSON string; GET /jobs returns a
// tdd.GetProcess array.
func newMockBackend(t *testing.T, validEmail, validPass string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
			return
		}
		if req.Password == "" || req.Email == "" {
			http.Error(w, `{"error":"must provide either username+password or email+password"}`, http.StatusBadRequest)
			return
		}
		if req.Email != validEmail || req.Password != validPass {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "test-jwt",
			Path:     "/",
			Expires:  timeNowPlusDay(),
			HttpOnly: true,
			Secure:   false,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"login successful","user_id":"00000000-0000-0000-0000-000000000001"}`))
	})

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
			return
		}
		if req.Email == "" || req.Password == "" {
			http.Error(w, `{"error":"must provide either username+password or email+password"}`, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"message":"User Signup successfully"}`))
	})

	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("token"); err != nil {
			http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-24 * time.Hour),
			HttpOnly: true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"logout successful"}`))
	})

	mux.HandleFunc("POST /jobs/create", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("token"); err != nil {
			http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
			return
		}
		var req struct {
			URLs        []string `json:"urls"`
			ProcessType string   `json:"process_type"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
			return
		}
		if len(req.URLs) == 0 {
			http.Error(w, `{"error":"url is required"}`, http.StatusBadRequest)
			return
		}
		switch req.ProcessType {
		case "audio", "video", "playlist", "upload", "transcription":
		default:
			http.Error(w, `{"error":"invalid process type"}`, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`"process has been created\n"`))
	})

	mux.HandleFunc("GET /jobs", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("token"); err != nil {
			http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{
			"id":"11111111-1111-1111-1111-111111111111",
			"user_id":"00000000-0000-0000-0000-000000000001",
			"url":"https://youtu.be/dQw4w9WgXcQ",
			"status":"queued",
			"process_type":"audio",
			"retries":0,
			"uploaded":false,
			"temp_path":"",
			"created_at":"2026-10-19T10:00:00Z",
			"updated_at":"2026-10-19T10:00:00Z"
		}]`))
	})

	return httptest.NewServer(mux)
}

func TestClientLoginRoundTrip(t *testing.T) {
	srv := newMockBackend(t, "alice@example.com", "secret")
	defer srv.Close()

	client := NewClient(srv.URL)

	userID, err := client.Login("alice@example.com", "secret")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if userID == "" {
		t.Fatal("expected a user ID from login, got empty string")
	}

	// The session cookie must be sent automatically on the next request.
	jobs, err := client.ListJobs()
	if err != nil {
		t.Fatalf("ListJobs after login failed (cookie not persisted?): %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job from the mock, got %d", len(jobs))
	}
	// The element must decode like tdd.GetProcess.
	if jobs[0].Status != statusQueued || jobs[0].ProcessType != "audio" || jobs[0].URL == "" {
		t.Errorf("unexpected job decode: %+v", jobs[0])
	}
}

func TestClientLoginBadCredentials(t *testing.T) {
	srv := newMockBackend(t, "alice@example.com", "secret")
	defer srv.Close()

	client := NewClient(srv.URL)
	if _, err := client.Login("alice@example.com", "wrong"); err == nil {
		t.Fatal("expected error for bad credentials, got nil")
	}
}

func TestClientListJobsRequiresAuth(t *testing.T) {
	srv := newMockBackend(t, "alice@example.com", "secret")
	defer srv.Close()

	client := NewClient(srv.URL)
	if _, err := client.ListJobs(); err == nil {
		t.Fatal("expected 401 error when not authenticated, got nil")
	}
}

func TestClientSignupAndHealth(t *testing.T) {
	srv := newMockBackend(t, "alice@example.com", "secret")
	defer srv.Close()

	client := NewClient(srv.URL)
	if err := client.Health(); err != nil {
		t.Fatalf("Health failed: %v", err)
	}
	if err := client.Signup("bob@example.com", "pw"); err != nil {
		t.Fatalf("Signup failed: %v", err)
	}
	// Signup does not create a session; jobs must still be protected.
	if _, err := client.ListJobs(); err == nil {
		t.Fatal("expected ListJobs to stay unauthorized after signup")
	}
}

func TestClientCreateJobsPayload(t *testing.T) {
	srv := newMockBackend(t, "alice@example.com", "secret")
	defer srv.Close()

	client := NewClient(srv.URL)
	if _, err := client.Login("alice@example.com", "secret"); err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	msg, err := client.CreateJobs([]string{"https://youtu.be/dQw4w9WgXcQ"}, "audio")
	if err != nil {
		t.Fatalf("CreateJobs failed: %v", err)
	}
	if msg == "" {
		t.Fatal("expected a message from POST /jobs/create")
	}

	// The mock rejects process types the backend does not accept.
	if _, err := client.CreateJobs([]string{"https://youtu.be/dQw4w9WgXcQ"}, "music"); err == nil {
		t.Fatal("expected error for unsupported process type, got nil")
	}
}

func TestClientLogout(t *testing.T) {
	srv := newMockBackend(t, "alice@example.com", "secret")
	defer srv.Close()

	client := NewClient(srv.URL)
	if _, err := client.Login("alice@example.com", "secret"); err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if err := client.Logout(); err != nil {
		t.Fatalf("Logout failed: %v", err)
	}
	// Cookie cleared: a second logout must 401 like the backend does.
	if err := client.Logout(); err == nil {
		t.Fatal("expected error when logging out without a cookie, got nil")
	}
}
