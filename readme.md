# 🚀 TRAER — Desktop Media Toolkit

TRAER is a **desktop application** built with **Go (backend)** and a modern **React + Tailwind UI** using Wails.

It allows you to:

* Transcribe YouTube videos into notes
* Convert videos into audio
* Download videos in different qualities
* Send media to your cloud server
* Manage API keys and settings securely
* View logs via an integrated terminal

---

# 🧱 Tech Stack

## 🖥️ Desktop Layer

* Wails

## ⚙️ Backend

* Go (you will implement)
* Optional: SQLite / remote DB for auth

## 🎨 Frontend

* React (latest)
* Vite
* Tailwind CSS (v3)
* Zustand
* React Router

---

# 📁 Project Structure

```bash
traer/
├── frontend/          # React UI (provided)
├── backend/           # Go backend (you build)
├── wails.json         # Wails config
├── main.go            # Entry point
└── build/             # Generated binaries
```

---

# 📦 Frontend Structure

```bash
frontend/
├── src/
│   ├── components/
│   ├── pages/
│   ├── store/
│   ├── hooks/
│   ├── utils/
│   ├── App.jsx
│   └── main.jsx
├── index.html
├── package.json
├── tailwind.config.js
├── vite.config.js
```

---

# ✨ Features

## 🎯 Core Features

### 📝 Transcription

* Input YouTube URL
* Output formats:

  * `.md`
  * `.mdx`
  * `.txt`

### 🔊 Audio Conversion

* Convert video → audio
* Supported formats:

  * MP3
  * WAV
  * AAC

### 🎥 Video Download

* Quality selection:

  * 144p → 1080p

### ☁️ Open Cloud Integration

* Send:

  * Audio
  * Video
  * Notes

---

## 🧩 UI Features

* Collapsible sidebar (recent links)
* Bottom terminal panel (logs)
* Dark/light mode toggle
* Profile & settings pages
* Clean desktop-first layout

---

## 🔐 Authentication

### Recommended Approach (IMPORTANT)

Do NOT store auth in frontend.

✔ Use backend:

* API-based login
* Secure storage (SQLite or server DB)

Frontend flow:

```json
Login → API → Validate → Store session
```

---

# ⚙️ Setup Instructions

## 1️⃣ Install Requirements

* Go (>= 1.20)
* Node.js (>= 18)
* npm or yarn
* Wails CLI

Install Wails:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

---

## 2️⃣ Create Wails App

```bash
wails init -n traer -t react
cd traer
```

---

## 3️⃣ Replace Frontend

Delete default frontend and copy:

```bash
your/frontend → traer/frontend
```

---

## 4️⃣ Install Frontend Dependencies

```bash
cd frontend
npm install
```

---

## 5️⃣ Run App (Development)

```bash
wails dev
```

✔ Hot reload enabled
✔ UI updates instantly

---

## 6️⃣ Build Desktop App

```bash
wails build
```

Output:

```bash
build/bin/traer.exe
```

---

# 🔌 Backend Integration (Go ↔ UI)

Wails auto-generates bindings.

Example Go function:

```go
func (a *App) ConvertAudio(url string) string {
    return "processing..."
}
```

Use in React:

```js
import { ConvertAudio } from "../wailsjs/go/main/App";

ConvertAudio(url).then(res => console.log(res));
```

---

# 🧠 Suggested Backend Modules

You should implement:

### 1. YouTube Processing

* Extract audio/video
* Use tools like:

  * yt-dlp

---

### 2. Transcription

* Whisper (local)
* Output formatter (MD/MDX)

---

### 3. Audio Conversion

* ffmpeg integration

---

### 4. Cloud Upload

* HTTP client → your server

---

### 5. Authentication

* Login endpoint
* Session/token handling

---

### 6. Terminal Logs

* Stream logs → UI

---

# 🎨 UI Customization

## Run frontend independently

```bash
cd frontend
npm run dev
```

## Tailwind styling

Edit:

```bash
src/index.css
tailwind.config.js
```

---

# 🔄 State Management (Zustand)

Example:

```js
const { sidebarOpen, toggleSidebar } = useStore();
```

Used for:

* UI state
* user session
* recent links

---

# 🧪 Future Enhancements

* Progress bars (conversion)
* Notifications (success/error)
* Drag & drop links
* Batch processing
* File preview
* Multi-user auth

---

# ⚠️ Security Notes

* Never store API keys in plain text frontend
* Use backend encryption
* Validate all inputs
* Avoid exposing shell directly

---

# 🧰 Troubleshooting

## App not starting

* Check Node + Go versions

## UI not loading

```bash
cd frontend
npm install
```

## Wails issues

```bash
wails doctor
```

---

# 📌 Development Workflow

1. Run:

```bash
wails dev
```

2. Edit UI in `/frontend`
3. Connect Go functions
4. Test live

---

# 📄 License

MIT (or your choice)

---


# 🔥 Final Notes

This project is:

* Desktop-first
* Modular
* Scalable

You now have:
✔ Full UI
✔ Clean architecture
✔ Ready for backend integration
