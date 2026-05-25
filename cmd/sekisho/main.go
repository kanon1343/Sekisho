package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"sekisho/internal/server"
)

func main() {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "4180"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid PORT environment variable: %v", err)
	}

	certFile := os.Getenv("TLS_CERT_FILE")
	if certFile == "" {
		certFile = "certs/localhost.pem"
	}

	keyFile := os.Getenv("TLS_KEY_FILE")
	if keyFile == "" {
		keyFile = "certs/localhost-key.pem"
	}

	fmt.Printf("Sekisho OAuth2 Proxy starting on https://localhost:%d...\n", port)

	// Create a root context
	ctx := context.Background()

	cfg := server.Config{
		Port:     port,
		CertFile: certFile,
		KeyFile:  keyFile,
	}

	if err := server.Run(ctx, cfg); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
