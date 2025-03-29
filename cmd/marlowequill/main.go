package main

import (
	"fmt"
	"log"
	"os"

	"github.com/reek004/MarloweQuill/internal/api"
	"github.com/reek004/MarloweQuill/internal/config"
	"github.com/reek004/MarloweQuill/internal/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create Gemini client
	geminiClient := api.NewGeminiClient(cfg)
	defer geminiClient.Close()

	// Create and start HTTP server
	srv := server.NewServer(geminiClient)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	fmt.Printf("Starting MarloweQuill web server on port %s...\n", port)
	if err := srv.Start(port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
