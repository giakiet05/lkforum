package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

const (
	UserColName            = "users"
	PostColName            = "posts"
	CommunityColName       = "communities"
	CommentColName         = "comments"
	ConversationColName    = "conversations"
	MessageColName         = "messages"
	VoteColName            = "votes"
	NotificationColName    = "notifications"
	ReportColName          = "reports"
	MembershipColName      = "memberships"
	LikedPostColName       = "liked_posts"
	SavedPostColName       = "saved_posts"
	UserPostHistoryColName = "user_post_history"
)

// NewMongoClient creates and returns a new MongoDB client
func NewMongoClient() *mongo.Client {
	uri := os.Getenv("MONGO_URI") // e.g. mongodb://user:pass@localhost:27017
	if uri == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	log.Println("Connected to MongoDB successfully!")
	Client = client

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is not set")
	}

	Database = client.Database(dbName)

	// Verify required collections exist
	if err := verifyCollections(ctx, Database); err != nil {
		log.Fatalf("Collection verification failed: %v", err)
	}

	log.Printf("Using database: %s\n", dbName)
	return client
}

func verifyCollections(ctx context.Context, db *mongo.Database) error {
	collections, err := db.ListCollectionNames(ctx, struct{}{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	required := []string{
		UserColName,
		PostColName,
		CommunityColName,
		CommentColName,
		ConversationColName,
		MessageColName,
		VoteColName,
		NotificationColName,
		ReportColName,
		MembershipColName,
		LikedPostColName,
		SavedPostColName,
		UserPostHistoryColName,
	}

	existing := make(map[string]bool, len(collections))
	for _, c := range collections {
		existing[c] = true
	}

	for _, name := range required {
		if !existing[name] {
			return fmt.Errorf("required collection %q does not exist in database", name)
		}
	}

	log.Println("All required collections verified")
	return nil
}
