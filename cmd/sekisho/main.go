package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4180"
	}

	certFile := os.Getenv("TLS_CERT_FILE")
	if certFile == "" {
		certFile = "certs/localhost.pem"
	}

	keyFile := os.Getenv("TLS_KEY_FILE")
	if keyFile == "" {
		keyFile = "certs/localhost-key.pem"
	}

	fmt.Printf("Sekisho OAuth2 Proxy starting on https://localhost:%s...\n", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Sekisho Proxy is running\n")
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
