package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rahulkumarpahwa/traer/config"
	"github.com/rahulkumarpahwa/traer/database"
	"github.com/rahulkumarpahwa/traer/internal/handlers"
	"github.com/rahulkumarpahwa/traer/internal/nextcloud"
	"github.com/rahulkumarpahwa/traer/internal/process"
	"github.com/rahulkumarpahwa/traer/internal/storage"
)

func main() {

	// database start
	err := database.InitiateDatabase()
	if err != nil {
		// fmt.Println(err)
		log.Fatal(err)
	}

	// setting up the configurations
	cfg := config.Load()

	// Resolve required executables (yt-dlp, ffmpeg)
	executablePaths, err := cfg.ResolveExecutables()
	if err != nil {
		log.Fatal(err)
	}

	// Opening the database:
	openDB := storage.OpenDB{Config: cfg}
	DB, err := openDB.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Setting up the Storage
	userStorage := storage.UserStorage{DB: DB}

	// setting up the admin user
	err = userStorage.CreateAdminUser(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Setting up the Nextcloud WebDAV uploader (used by the music pipeline).
	// Credentials come exclusively from config/env and are never logged.
	var uploader *nextcloud.Client
	if cfg.NextcloudConfigured() {
		uploader = nextcloud.NewClient(cfg.NextcloudURL, cfg.NextcloudUsername, cfg.NextcloudAppPassword)
		log.Printf("nextCloud uploader configured for user %q (base path %q)", cfg.NextcloudUsername, cfg.NextcloudBasePath)
	} else {
		log.Println("[nextcloud] uploader NOT configured; music jobs will be rejected until NEXTCLOUD_* env vars are set")
	}

	router := http.NewServeMux()

	// Creating the job worker with resolved executable paths
	jobworker := process.JobWorker{
		Config:          cfg,
		ExecutablePaths: executablePaths,
		Uploader:        uploader,
	}

	// starting the worker
	jobworker.StartWorkers()

	// Creating the new Service Hanlder
	serviceHandler := handlers.ServiceHandler{JW: &jobworker}
	// Creating the new user Handler
	userHandler := handlers.UserHandler{UserStorage: &userStorage, Config: cfg}

	// Health check - no auth needed
	router.HandleFunc("GET /", healthHandler)

	// Auth routes - public, no middleware
	router.Handle("POST /login", http.HandlerFunc(userHandler.HandleLogin))
	router.Handle("POST /users", http.HandlerFunc(userHandler.HandleCreateUser)) // Registration is public
	router.Handle("POST /logout", http.HandlerFunc(userHandler.HandleLogout))

	// Job routes - auth required
	router.Handle("POST /jobs/create", userHandler.AuthMiddleware(http.HandlerFunc(serviceHandler.HandleCreateJobs)))
	router.Handle("GET /jobs/active", userHandler.AuthMiddleware(http.HandlerFunc(serviceHandler.HandleActiveJobs)))
	router.Handle("GET /jobs/status", userHandler.AuthMiddleware(http.HandlerFunc(serviceHandler.HandleJobStatus)))
	router.Handle("POST /jobs/retry", userHandler.AuthMiddleware(http.HandlerFunc(serviceHandler.HandleRetryJob)))
	router.Handle("GET /jobs/download", userHandler.AuthMiddleware(http.HandlerFunc(serviceHandler.HandleDownload)))
	router.Handle("GET /instances/get", userHandler.AuthMiddleware(http.HandlerFunc(serviceHandler.HandleGetInstances)))

	// User routes - auth required
	router.Handle("PUT /users", userHandler.AuthMiddleware(http.HandlerFunc(userHandler.HandleUpdateUser)))
	router.Handle("DELETE /users", userHandler.AuthMiddleware(http.HandlerFunc(userHandler.HandleDeleteUser)))
	router.Handle("GET /users/id", userHandler.AuthMiddleware(http.HandlerFunc(userHandler.HandleGetUserByID)))
	router.Handle("GET /users/email", userHandler.AuthMiddleware(http.HandlerFunc(userHandler.HandleGetUserByEmail)))

	server := http.Server{
		Addr:         ":" + cfg.HTTPAddress,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.HTTPReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTPWriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTPIdleTimeout) * time.Second,
	}

	slog.Info("server started at", "port", cfg.HTTPAddress)

	notifyChan := make(chan os.Signal, 1)

	signal.Notify(notifyChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-notifyChan

	log.Println("Server Shutdown in 3..2..1!")

	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelCtx()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server Shutdown Successfully!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Simple liveness check
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
