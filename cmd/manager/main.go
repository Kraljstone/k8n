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

	"github.com/Kraljstone/k8n/internal/manager"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	server := manager.NewServer()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	go func() {
		fmt.Fprintf(w, "Manager listening on %s...\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "HTTP listen error: %v\n", err)
		}
	}()

	<-ctx.Done()
	fmt.Fprintln(w, "\nShutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	return httpServer.Shutdown(shutdownCtx)
}
