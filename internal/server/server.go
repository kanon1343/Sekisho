package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Config represents the server configuration.
type Config struct {
	Port int
}

// Run starts the HTTP server and handles graceful shutdown.
func Run(ctx context.Context, cfg Config) error {
	router := NewRouter()
	
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Starting server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// Block until we receive our signal or the server returns an error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		fmt.Printf("Received signal: %v. Initiating graceful shutdown...\n", sig)
		
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			// forceful shutdown
			if err := srv.Close(); err != nil {
				return fmt.Errorf("could not close server gracefully: %w", err)
			}
			return fmt.Errorf("could not shutdown server gracefully: %w", err)
		}
	case <-ctx.Done():
		fmt.Println("Context cancelled. Initiating graceful shutdown...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("could not shutdown server gracefully: %w", err)
		}
	}

	fmt.Println("Server stopped")
	return nil
}
