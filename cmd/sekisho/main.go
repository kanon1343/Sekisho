package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Sekisho OAuth2 Proxy starting on https://localhost:4180...")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Sekisho Proxy is running\n")
	})

	server := &http.Server{
		Addr:    ":4180",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	err := server.ListenAndServeTLS("certs/localhost.pem", "certs/localhost-key.pem")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
