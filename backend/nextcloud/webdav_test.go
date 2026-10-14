package nextcloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDavURL(t *testing.T) {
	c := NewClient("https://cloud.example.com", "alice", "secret")

	got, err := c.davURL("Music", "Downloads", "My Song #1.mp3")
	if err != nil {
		t.Fatalf("davURL returned error: %v", err)
	}

	want := "https://cloud.example.com/remote.php/dav/files/alice/Music/Downloads/My%20Song%20%231.mp3"
	if got != want {
		t.Fatalf("davURL mismatch:\n got: %s\nwant: %s", got, want)
	}

	// A space must be encoded exactly once (never %2520).
	if strings.Contains(got, "%25") {
		t.Fatalf("davURL double-encoded a character: %s", got)
	}
}

func TestDavURLWithSubpath(t *testing.T) {
	c := NewClient("https://example.com/nextcloud", "bob", "secret")
	got, err := c.davURL("Music", "track.mp3")
	if err != nil {
		t.Fatalf("davURL returned error: %v", err)
	}
	want := "https://example.com/nextcloud/remote.php/dav/files/bob/Music/track.mp3"
	if got != want {
		t.Fatalf("davURL subpath mismatch:\n got: %s\nwant: %s", got, want)
	}
}

func TestSanitizeName(t *testing.T) {
	cases := map[string]string{
		"../../etc/passwd":  "passwd",
		"a/b/c.mp3":         "c.mp3",
		"weird\x00name.mp3": "weirdname.mp3",
		"..":                "download.mp3",
		"normal song.mp3":   "normal song.mp3",
	}
	for in, want := range cases {
		if got := sanitizeName(in); got != want {
			t.Errorf("sanitizeName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestConfigured(t *testing.T) {
	if NewClient("", "u", "p").Configured() {
		t.Error("client with empty base URL should not be configured")
	}
	if NewClient("https://x", "", "p").Configured() {
		t.Error("client with empty username should not be configured")
	}
	if NewClient("https://x", "u", "").Configured() {
		t.Error("client with empty app password should not be configured")
	}
	if !NewClient("https://x", "u", "p").Configured() {
		t.Error("fully populated client should be configured")
	}
}

// TestUploadFileAgainstTestServer exercises the full PUT flow (folder creation +
// upload) against an in-process WebDAV-like server.
func TestUploadFileAgainstTestServer(t *testing.T) {
	var putPath, putAuth string
	var mkcolPaths []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		putAuth = user + ":" + pass
		switch r.Method {
		case "MKCOL":
			mkcolPaths = append(mkcolPaths, r.URL.Path)
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut:
			putPath = r.URL.Path
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "alice", "app-pass")
	if !c.Configured() {
		t.Fatal("test client should be configured")
	}

	// Write a small temp file to upload.
	tmp := filepath.Join(t.TempDir(), "My Song.mp3")
	if err := os.WriteFile(tmp, []byte("id3-bytes"), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if err := c.UploadFile(context.Background(), tmp, "Music/Downloads", "My Song.mp3"); err != nil {
		t.Fatalf("UploadFile failed: %v", err)
	}

	if putAuth != "alice:app-pass" {
		t.Errorf("unexpected basic auth on PUT: %q", putAuth)
	}
	// r.URL.Path is decoded by the server, so a correctly-encoded space round
	// trips back to a literal space here.
	wantPath := "/remote.php/dav/files/alice/Music/Downloads/My Song.mp3"
	if putPath != wantPath {
		t.Errorf("unexpected PUT path:\n got: %q\nwant: %q", putPath, wantPath)
	}
	if len(mkcolPaths) == 0 {
		t.Error("expected at least one MKCOL to create the destination folder")
	}
}
