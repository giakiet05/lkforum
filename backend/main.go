package main

import (
	"github.com/giakiet05/lkforum/internal/bootstrap"
	"github.com/giakiet05/lkforum/internal/config"
	"log"
	"os"
)

func main() {
	config.LoadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	r, err := bootstrap.Init()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	// Start the server
	log.Printf("Server is running at http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
