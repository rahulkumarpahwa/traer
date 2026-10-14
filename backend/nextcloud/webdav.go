package nextcloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Client is a minimal Nextcloud WebDAV client used to upload finished media
// files into a user's Nextcloud account. Credentials are supplied via config
// (environment variables) and are never logged or exposed to the frontend.
type Client struct {
	BaseURL     string // e.g. https://cloud.example.com
	Username    string
	AppPassword string
	HTTPClient  *http.Client
}

// NewClient builds a WebDAV client. It does not validate the credentials; use
// Configured to check whether the required settings are present.
func NewClient(baseURL, username, appPassword string) *Client {
	return &Client{
		BaseURL:     strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		Username:    strings.TrimSpace(username),
		AppPassword: appPassword,
		HTTPClient:  &http.Client{Timeout: 30 * time.Minute},
	}
}

// Configured reports whether the client has the minimum required settings to
// attempt an upload.
func (c *Client) Configured() bool {
	return c != nil && c.BaseURL != "" && c.Username != "" && c.AppPassword != ""
}

// setAuth attaches HTTP Basic auth using the Nextcloud app password. The secret
// is only ever placed on the Authorization header and never logged.
func (c *Client) setAuth(req *http.Request) {
	req.SetBasicAuth(c.Username, c.AppPassword)
}

// davURL builds the WebDAV endpoint for the given remote path segments. Each
// segment is escaped exactly once so filenames containing spaces or unicode
// remain safe (and are not double-encoded).
func (c *Client) davURL(segments ...string) (string, error) {
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid nextcloud base url: %w", err)
	}
	if base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("nextcloud base url must include scheme and host")
	}

	parts := []string{"remote.php", "dav", "files", c.Username}
	for _, seg := range segments {
		seg = strings.Trim(strings.TrimSpace(seg), "/")
		if seg == "" {
			continue
		}
		parts = append(parts, seg)
	}

	var b strings.Builder
	b.WriteString(base.Scheme)
	b.WriteString("://")
	b.WriteString(base.Host)
	// Preserve any sub-path install prefix already present in the base URL
	// (url.Parse keeps it correctly escaped in EscapedPath).
	if prefix := strings.TrimSuffix(base.EscapedPath(), "/"); prefix != "" {
		b.WriteString(prefix)
	}
	for _, seg := range parts {
		b.WriteString("/")
		b.WriteString(url.PathEscape(seg))
	}
	return b.String(), nil
}

// mkcol creates a single collection (folder). A 405 (already exists) or any 2xx
// response is treated as success.
func (c *Client) mkcol(ctx context.Context, target string) error {
	req, err := http.NewRequestWithContext(ctx, "MKCOL", target, nil)
	if err != nil {
		return err
	}
	c.setAuth(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return nil
	case resp.StatusCode == http.StatusMethodNotAllowed: // 405 already exists
		return nil
	default:
		return fmt.Errorf("failed to create collection %s: status %d", target, resp.StatusCode)
	}
}

// ensureCollection creates the folder path recursively via MKCOL, ignoring
// folders that already exist.
func (c *Client) ensureCollection(ctx context.Context, dir string) error {
	var acc []string
	for _, seg := range splitPath(dir) {
		acc = append(acc, seg)
		target, err := c.davURL(acc...)
		if err != nil {
			return err
		}
		if err := c.mkcol(ctx, target); err != nil {
			return err
		}
	}
	return nil
}

// UploadFile uploads the local file to <remoteDir>/<fileName> under the user's
// Nextcloud files. The destination folder is created if necessary. The file is
// streamed directly from disk so large MP3s are not buffered in memory.
func (c *Client) UploadFile(ctx context.Context, localPath, remoteDir, fileName string) error {
	if !c.Configured() {
		return fmt.Errorf("nextcloud client is not configured")
	}

	fileName = sanitizeName(fileName)
	if fileName == "" {
		return fmt.Errorf("invalid destination filename")
	}

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open local file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat local file: %w", err)
	}

	dirSegments := splitPath(remoteDir)
	if err := c.ensureCollection(ctx, remoteDir); err != nil {
		return fmt.Errorf("ensure remote folder: %w", err)
	}

	target, err := c.davURL(append(dirSegments, fileName)...)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, target, file)
	if err != nil {
		return err
	}
	c.setAuth(req)
	req.ContentLength = info.Size()
	req.Header.Set("Content-Type", "audio/mpeg")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("nextcloud upload failed: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

// splitPath normalizes a remote path into clean segments, tolerating both
// forward and back slashes and dropping empty/"." components.
func splitPath(p string) []string {
	p = strings.ReplaceAll(p, "\\", "/")
	var out []string
	for _, seg := range strings.Split(p, "/") {
		seg = strings.TrimSpace(seg)
		if seg != "" && seg != "." {
			out = append(out, seg)
		}
	}
	return out
}

// sanitizeName strips any path separators and control characters from a
// filename so it cannot escape the destination directory or corrupt the request.
func sanitizeName(name string) string {
	name = path.Base(strings.ReplaceAll(strings.TrimSpace(name), "\\", "/"))
	name = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, name)
	if name == "." || name == ".." || name == "/" {
		return "download.mp3"
	}
	return name
}
