# syntax=docker/dockerfile:1

# ---- build stage ----
FROM golang:1.25-alpine AS build
WORKDIR /src

# Bring in build tooling needed by modernc.org/sqlite (pure-Go, needs gcc for cgo-free build only if enabled; we keep CGO off).
RUN apk add --no-cache git ca-certificates

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the backend source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/traer-backend .

# ---- runtime stage ----
FROM alpine:3.20
WORKDIR /app

# ffmpeg for audio extraction/conversion and yt-dlp for media downloading.
RUN apk add --no-cache ca-certificates ffmpeg python3 py3-pip tzdata \
    && pip3 install --no-cache-dir --break-system-packages yt-dlp \
    && mkdir -p /app/downloads /app/database /app/local

COPY --from=build /out/traer-backend /app/traer-backend

# Runtime configuration is provided exclusively via environment variables so no
# secrets are baked into the image or committed to source control.
ENV HTTP_ADDRESS=0.0.0.0:8080 \
    DOWNLOAD_PATH=downloads/ \
    DB_PATH=database/database.db

# Persisted volumes: downloaded media awaiting upload, the SQLite database, and
# an optional local config file.
VOLUME ["/app/downloads", "/app/database", "/app/local"]

# The service is intended to run on an internal Docker network only. It is not
# exposed publicly here; reach it through your reverse proxy / Nextcloud network.
EXPOSE 8080

ENTRYPOINT ["/app/traer-backend"]
