package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rahulkumarpahwa/traer/config"
	"github.com/rahulkumarpahwa/traer/job"
	"github.com/rahulkumarpahwa/traer/storage"
	"github.com/rahulkumarpahwa/traer/storage/open"

	jRoutes "github.com/rahulkumarpahwa/traer/routes/job"
	uRoutes "github.com/rahulkumarpahwa/traer/routes/user"
)

func main() {

	// setting up the configurations
	cfg, err := config.Load("local/local.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Resolve required executables (yt-dlp, ffmpeg)
	executablePaths, err := cfg.ResolveExecutables()
	if err != nil {
		log.Fatal(err)
	}

	// Opening the database:
	openDB := open.OpenDB{Config: cfg}
	DB, err := openDB.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Setting up the Storage
	userStorage := storage.UserStorage{DB: DB}

	router := http.NewServeMux()

	// Creating the job worker with resolved executable paths
	jobworker := job.JobWorker{
		Config:          cfg,
		ExecutablePaths: executablePaths,
	}

	// starting the worker
	jobworker.StartWorkers()

	// Creating the new Service Hanlder
	serviceHandler := jRoutes.ServiceHandler{JW: &jobworker}
	// Creating the new user Handler
	userHandler := uRoutes.UserHandler{UserStorage: &userStorage, Config: cfg}

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
	router.Handle("GET /jobs/download", userHandler.AuthMiddleware(http.HandlerFunc(serviceHandler.HandleDownload)))
	router.Handle("GET /instances/get", userHandler.AuthMiddleware(http.HandlerFunc(serviceHandler.HandleGetInstances)))

	// User routes - auth required
	router.Handle("PUT /users", userHandler.AuthMiddleware(http.HandlerFunc(userHandler.HandleUpdateUser)))
	router.Handle("DELETE /users", userHandler.AuthMiddleware(http.HandlerFunc(userHandler.HandleDeleteUser)))
	router.Handle("GET /users/id", userHandler.AuthMiddleware(http.HandlerFunc(userHandler.HandleGetUserByID)))
	router.Handle("GET /users/email", userHandler.AuthMiddleware(http.HandlerFunc(userHandler.HandleGetUserByEmail)))
	router.Handle("GET /users/username", userHandler.AuthMiddleware(http.HandlerFunc(userHandler.HandleGetUserByUsername)))

	server := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.HTTPServer.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTPServer.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTPServer.IdleTimeout) * time.Second,
	}

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
