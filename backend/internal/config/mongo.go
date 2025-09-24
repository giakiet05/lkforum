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
	Client                    *mongo.Client
	Database                  *mongo.Database
	UserCollection            *mongo.Collection
	PostCollection            *mongo.Collection
	CommunityCollection       *mongo.Collection
	CommentCollection         *mongo.Collection
	VoteCollection            *mongo.Collection
	NotificationCollection    *mongo.Collection
	ReportCollection          *mongo.Collection
	MembershipCollection      *mongo.Collection
	LikedPostCollection       *mongo.Collection
	SavedPostCollection       *mongo.Collection
	UserPostHistoryCollection *mongo.Collection
)

// NewMongoClient creates and returns a new MongoDB client
func NewMongoClient() *mongo.Client {
	uri := os.Getenv("MONGO_URI") // e.g. mongodb://user:pass@localhost:27017

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

	// Check if database actually exists
	dbs, err := client.ListDatabaseNames(ctx, struct{}{})
	if err != nil {
		log.Fatalf("Could not list databases: %v", err)
	}

	found := false
	for _, db := range dbs {
		if db == dbName {
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("Database %q does not exist on this MongoDB server", dbName)
	}

	Database = client.Database(dbName)
	log.Printf("Using database: %s\n", dbName)

	if err := initCollections(ctx); err != nil {
		log.Fatalf("Collection initialization failed: %v", err)
	}

	return client
}

func initCollections(ctx context.Context) error {
	collections, err := Database.ListCollectionNames(ctx, struct{}{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	required := map[string]**mongo.Collection{
		"users":             &UserCollection,
		"posts":             &PostCollection,
		"communities":       &CommunityCollection,
		"comments":          &CommentCollection,
		"votes":             &VoteCollection,
		"notifications":     &NotificationCollection,
		"reports":           &ReportCollection,
		"memberships":       &MembershipCollection,
		"liked_posts":       &LikedPostCollection,
		"saved_posts":       &SavedPostCollection,
		"user_post_history": &UserPostHistoryCollection,
	}

	existing := make(map[string]bool, len(collections))
	for _, c := range collections {
		existing[c] = true
	}

	// Assign collection handles if they exist, otherwise fail
	for name, ref := range required {
		if !existing[name] {
			return fmt.Errorf("required collection %q does not exist in database", name)
		}
		*ref = Database.Collection(name)
	}

	log.Println("All required collections verified and initialized")
	return nil
}
