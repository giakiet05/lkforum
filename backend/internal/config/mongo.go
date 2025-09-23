package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

// NewMongoClient creates and returns a new MongoDB client
func NewMongoClient() *mongo.Client {
	uri := os.Getenv("MONGO_URI") // ví dụ: mongodb://user:pass@localhost:27017/dbname

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	// test connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	log.Println("Connected to MongoDB successfully!")
	return client
}
