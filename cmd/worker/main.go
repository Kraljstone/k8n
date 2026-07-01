package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kraljstone/k8n/internal/worker"
)

const (
	managerAddr = "http://localhost:8080"
	workerPort  = ":8081"
	workerAddr  = "http://localhost:8081"
	workerID    = "worker-node-1"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "worker error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Fprintf(w, "[Worker] Booting up %s...\n", workerID)

	// Register with the manager on startup.
	client := worker.NewClient(managerAddr)
	go func() {
		if err := client.Register(workerID, workerAddr); err != nil {
			fmt.Fprintf(os.Stderr, "[Worker] Registration failed: %v\n", err)
			return
		}
		fmt.Fprintln(w, "[Worker] Successfully registered with Manager cluster!")
	}()

	server := worker.NewServer()
	httpServer := &http.Server{
		Addr:    workerPort,
		Handler: server,
	}

	go func() {
		fmt.Fprintf(w, "[Worker] Listening for incoming tasks on port %s\n", workerPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "[Worker] HTTP server error: %v\n", err)
		}
	}()

	<-ctx.Done()
	fmt.Fprintln(w, "\n[Worker] Shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	return httpServer.Shutdown(shutdownCtx)
}
