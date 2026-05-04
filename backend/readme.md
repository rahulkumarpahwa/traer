# Traer Backend

## Overview

Traer Backend is a Go-based service that provides asynchronous media processing using yt-dlp. It allows users to submit YouTube URLs and convert them into audio (MP3) or video (MP4) files. The system is built with a worker queue architecture to handle concurrent jobs efficiently while tracking progress in real time.

---

## Features

* Asynchronous job processing using worker pool
* Support for audio (MP3) and video (MP4) conversion
* Real-time progress tracking via polling
* File download endpoint for completed jobs
* Configurable worker concurrency
* Graceful server shutdown
* Thread-safe job state management

---

## Architecture

The system is composed of three main components:

### 1. HTTP Layer

Handles incoming API requests and responses.

### 2. Job Service

Responsible for:

* Queue management
* Worker pool execution
* Progress tracking
* Job lifecycle management

### 3. External Tools

* yt-dlp: media downloading
* ffmpeg: media conversion

---

## Project Structure

```
.
├── main.go          # Application entry point
├── routes/          # HTTP handlers
├── job/             # Worker, queue, and processing logic
├── types/           # Shared structs and enums
├── utils/           # Helper utilities (JSON response, errors, IDs)
└── downloads/       # Output files directory
```

---

## Worker Model

* Fixed-size worker pool (default: 2 workers)
* Jobs are pushed into a buffered channel queue
* Each worker:

  * Picks a job
  * Executes yt-dlp command
  * Parses progress in real time
  * Updates job state
  * Moves to next job after completion

Concurrency is controlled by limiting the number of workers.

---

## Job Lifecycle

1. Job Created → Status: queued
2. Worker picks job → Status: running
3. Progress updates (0–100%)
4. Job completes → Status: done
5. Output file available for download

Failure scenarios:

* Invalid input
* yt-dlp execution failure
* Queue overflow

---

## API Endpoints

### Health Check

```
GET /
```

Response:

```
200 OK
```

---

### Create Job

```
POST /jobs/create
```

Request:

```json
{
  "url": "https://youtube.com/...",
  "type": "audio" | "video"
}
```

Response:

```json
{
  "id": "job-id",
  "status": "queued",
  "progress": 0
}
```

---

### Get Active Jobs

```
GET /jobs/active
```

Returns all jobs with status:

* queued
* running

---

### Get Job Status

```
GET /jobs/status?id=<job_id>
```

Response:

```json
{
  "id": "job-id",
  "status": "running",
  "progress": 42.3,
  "output": "",
  "error": ""
}
```

---

### Download File

```
GET /jobs/download?id=<job_id>
```

Returns:

* File stream if job is completed
* Error if job is not ready or failed

---

### Get Worker Instances

```
GET /instances/get
```

Response:

```json
{
  "instances": 2
}
```

---

## Concurrency Model

* Job queue implemented using Go channels
* Worker pool implemented using goroutines
* Shared job state protected via mutex locks
* One job per worker at a time
* Maximum parallel jobs = number of workers

---

## Progress Tracking

* yt-dlp outputs progress to stdout/stderr
* Output is parsed using regex
* Progress is updated in shared job struct
* Clients poll `/jobs/status` to retrieve latest state

---

## Error Handling

* Centralized JSON error responses via utility package
* Graceful failure states recorded per job
* Queue overflow handled with fallback status

---

## Shutdown Handling

* Graceful shutdown using context and OS signals
* Ensures HTTP server stops accepting new requests
* Allows in-flight operations to complete within timeout

---

## Dependencies

* Go (1.20+ recommended)
* yt-dlp
* ffmpeg

Ensure both yt-dlp and ffmpeg are installed and available in system PATH.

---

## Running the Server

```
go run main.go
```

Server runs on:

```
http://localhost:8080
```

---

## Limitations

* No persistent storage (in-memory job store)
* No authentication or authorization
* No job cancellation support
* Worker scaling is static after startup
* Polling-based progress tracking (no WebSocket support)

---

## Future Improvements

* Persistent storage (SQLite / Redis)
* WebSocket-based real-time updates
* Job cancellation and retry mechanisms
* Dynamic worker scaling
* File cleanup and retention policies
* Multi-user support with authentication

---

## Summary

This backend implements a scalable and concurrent job processing system using Go. It demonstrates production-grade concepts such as worker pools, concurrency control, process orchestration, and API-driven job management while maintaining simplicity and extensibility.
