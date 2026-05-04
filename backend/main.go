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

	"github.com/rahulkumarpahwa/traer/job"
	"github.com/rahulkumarpahwa/traer/routes"
)

func main() {

	router := http.NewServeMux()

	// Creating the job worker
	jobworker := job.JobWorker{}

	// starting the worker
	jobworker.StartWorkers()

	// Creating the new Service Hanlder
	serviceHandler := routes.ServiceHandler{JW: &jobworker}

	router.HandleFunc("GET /", healthHandler)
	router.HandleFunc("POST /jobs/create", serviceHandler.HandleCreateJobs)
	router.HandleFunc("GET /jobs/active", serviceHandler.HandleActiveJobs)
	router.HandleFunc("GET /jobs/status", serviceHandler.HandleJobStatus)
	router.HandleFunc("GET /jobs/download", serviceHandler.HandleDownload)
	router.HandleFunc("GET /instances/get", serviceHandler.HandleGetInstances)

	server := http.Server{
		Addr:        ":8080",
		Handler:     router,
		ReadTimeout: time.Duration(time.Second * 20),
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

	log.Println("Server Shutdown in 5..4..3..2..1!")

	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelCtx()

	err := server.Shutdown(ctx)
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
