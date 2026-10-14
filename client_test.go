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

// newMockBackend mimics the Traer backend's auth contract: POST /login sets an
// HttpOnly "token" cookie (Secure:false, Path:/) and returns 200; protected
// routes require that cookie and 401 otherwise.
func newMockBackend(t *testing.T, validUser, validPass string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
			return
		}
		ident := req.Username
		if ident == "" {
			ident = req.Email
		}
		if req.Password == "" || (req.Username == "" && req.Email == "") {
			http.Error(w, `{"error":"must provide either username+password or email+password"}`, http.StatusBadRequest)
			return
		}
		if ident != validUser || req.Password != validPass {
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
		_, _ = w.Write([]byte(`{"message":"login successful","user_id":"1"}`))
	})

	mux.HandleFunc("/jobs/active", func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("token"); err != nil {
			http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	})

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
			return
		}
		if req.Username == "" || req.Email == "" || req.Password == "" {
			http.Error(w, `{"error":"username, email, and password cannot be empty"}`, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"user created"}`))
	})

	return httptest.NewServer(mux)
}

func TestClientLoginRoundTrip(t *testing.T) {
	srv := newMockBackend(t, "alice", "secret")
	defer srv.Close()

	client := NewClient(srv.URL)

	if err := client.Login("alice", "secret"); err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	// The session cookie must be sent automatically on the next request.
	jobs, err := client.ActiveJobs()
	if err != nil {
		t.Fatalf("ActiveJobs after login failed (cookie not persisted?): %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected empty active jobs, got %d", len(jobs))
	}
}

func TestClientLoginBadCredentials(t *testing.T) {
	srv := newMockBackend(t, "alice", "secret")
	defer srv.Close()

	client := NewClient(srv.URL)
	if err := client.Login("alice", "wrong"); err == nil {
		t.Fatal("expected error for bad credentials, got nil")
	}
}

func TestClientActiveJobsRequiresAuth(t *testing.T) {
	srv := newMockBackend(t, "alice", "secret")
	defer srv.Close()

	client := NewClient(srv.URL)
	if _, err := client.ActiveJobs(); err == nil {
		t.Fatal("expected 401 error when not authenticated, got nil")
	}
}
