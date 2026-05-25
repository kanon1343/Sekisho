package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sekisho/internal/handler"
	"sekisho/internal/server"
	"testing"
	"time"
)

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	r := server.NewRouter()
	
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	
	r.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var resp handler.HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	
	if resp.Status != "ok" {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Status, "ok")
	}
}

func TestServerGracefulShutdown(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cfg := server.Config{Port: 0} // port 0 will pick a random available port

	errCh := make(chan error, 1)

	go func() {
		errCh <- server.Run(ctx, cfg)
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Trigger graceful shutdown by canceling the context
	cancel()

	// Wait for the server to stop
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("expected no error on graceful shutdown, got: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}
