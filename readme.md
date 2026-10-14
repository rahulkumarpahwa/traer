# TRAER — YouTube → Nextcloud Music Pipeline (TUI)

Traer is a terminal application built in Go with a [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI. It downloads the audio from one or more YouTube / YouTube Music links, converts it to MP3 (yt-dlp + FFmpeg), and uploads the result directly to your Nextcloud Music folder over the WebDAV API.

> Only download content you own or are authorized to use. Traer does not grant any rights to media; it is a tool you operate against sources you are permitted to access.

---

## Architecture

Traer has two parts:

1. **Backend** (`backend/`) — a Go HTTP service that owns the job queue, workers, yt-dlp/FFmpeg orchestration, SQLite auth (JWT), and the Nextcloud WebDAV uploader.
2. **TUI** (root module) — a Bubbletea client (`main.go`, `client.go`, `model.go`, `view.go`, `styles.go`). It is only an HTTP client to the backend; all processing happens server-side.

```
TUI (this module)  --HTTP+JWT cookie-->  backend  --yt-dlp/FFmpeg--> MP3
                                                              \-WebDAV-> Nextcloud
```

The TUI keeps the backend session in a cookie jar exactly like a browser would (the backend issues an HttpOnly `token` cookie on login).

---

## Features

- Sign in / create an account (username-or-email + password).
- Paste one or many YouTube / YouTube Music links (one per line).
- Queue each as a music job: extract MP3 and upload to Nextcloud.
- Live job list with status (`queued` / `running` / `uploading` / `done` / `failed`), a progress bar, and errors.
- Retry failed jobs; the backend resumes from the already-downloaded file instead of re-downloading.
- Duplicate protection: re-submitting a link that is already queued/running/done is a no-op.
- Temporary MP3 files are deleted only after a successful Nextcloud upload.

---

## Requirements

- Go (>= 1.25)
- On the backend host: `yt-dlp` and `ffmpeg` (auto-installed on Linux/macOS/Windows if missing), or run the backend via Docker where they are preinstalled.

---

## Running

### 1. Start the backend

Locally:
```bash
cd backend
# create local/local.yaml from the example and set JWT_SECRET + Nextcloud creds,
# or export the equivalent env vars (see backend/local.example.yaml / .env.example)
go run .
# listens on http://localhost:8080
```

Or with Docker (from the repo root):
```bash
cp .env.example .env   # then fill in JWT_SECRET and Nextcloud creds
docker compose up -d --build
# backend is published on 127.0.0.1:8080 (loopback only - not public)
```

### 2. Start the TUI

```bash
go run . --backend http://localhost:8080
# or set TRAER_BACKEND_URL instead of the flag
```

---

## Usage / Keys

**Sign-in screen**
- `enter` — submit
- `tab` / `up` / `down` — move between fields
- `ctrl+t` — toggle between sign in and create account
- `q` / `ctrl+c` — quit

**Main screen**
- Paste links into the box (one per line)
- `ctrl+s` — queue all pasted links to Nextcloud
- `tab` — switch focus between the link box and the job list
- When the job list is focused: `up`/`down` or `j`/`k` to select, `r` to retry a failed job
- `ctrl+r` — refresh now (also auto-refreshes every 2s)
- `q` / `ctrl+c` — quit

---

## Configuration

All secrets live in environment variables or a git-ignored config file and are never shown in the TUI or logged:

| Variable | Purpose |
| --- | --- |
| `JWT_SECRET` | Signs session cookies (backend) |
| `NEXTCLOUD_URL` | Nextcloud base URL (e.g. `http://app:80` inside Docker) |
| `NEXTCLOUD_USERNAME` | Nextcloud user |
| `NEXTCLOUD_APP_PASSWORD` | Nextcloud app password (never your account password) |
| `NEXTCLOUD_BASE_PATH` | Destination folder (default `Music/Downloads`) |
| `TRAER_BACKEND_URL` | Backend URL the TUI connects to (default `http://localhost:8080`) |

See `.env.example` and `backend/local.example.yaml`.

---

## Project structure

```
.
├── main.go        # TUI entry point (flags + Bubbletea program)
├── client.go      # backend HTTP client (login, queue, status, retry)
├── model.go       # Bubbletea model + update logic
├── view.go        # rendering
├── styles.go      # lipgloss styles + progress bar
├── backend/       # Go HTTP service (queue, workers, yt-dlp/FFmpeg, WebDAV)
├── docker-compose.yml
└── .env.example
```

---

## Security notes

- Nextcloud credentials stay on the backend only; the TUI never sees them.
- The backend API is published on loopback only in the provided compose file.
- Filenames are sanitized and URL-escaped before the WebDAV upload.
- Authentication is required before any pipeline action.
