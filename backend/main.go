package main

import (
	"log"
	"os"

	"github.com/giakiet05/lkforum/internal/bootstrap"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/gin-gonic/gin"
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

	// CORS middleware
	allowOrigin := os.Getenv("FRONTEND_URL")
	if allowOrigin == "" {
		allowOrigin = "http://localhost:5173"
	}

	r.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Start the
	log.Printf("Server is running at http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
