package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not found in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("LKForum")
	notificationCol := db.Collection("notifications")

	// Delete all notifications with type "comment" (old moderator invitations)
	filter := bson.M{"type": "comment", "message": bson.M{"$regex": "moderator"}}

	result, err := notificationCol.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal("Failed to delete notifications:", err)
	}

	log.Printf("✅ Deleted %d old moderator invitation notifications", result.DeletedCount)
	log.Println("💡 Now re-invite moderators to create new notifications with metadata")
}
