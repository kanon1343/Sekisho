package main

import (
	"context"
	"fmt"
	"log"
	"sekisho/internal/server"
)

func main() {
	fmt.Println("Sekisho OAuth2 Proxy starting...")

	// Create a root context
	ctx := context.Background()

	// Hardcoded config for now, will be replaced by SK-104 config package
	cfg := server.Config{
		Port: 4180,
	}

	if err := server.Run(ctx, cfg); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
