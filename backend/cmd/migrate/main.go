package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	log.Println("🚀 Starting moderator avatar migration...")

	// Load config and connect to MongoDB
	config.LoadConfig()
	mongoClient := config.NewMongoClient()
	db := mongoClient.Database(config.Cfg.DBName)
	collection := db.Collection("communities")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Find all communities
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("❌ Failed to find communities: %v", err)
	}
	defer cursor.Close(ctx)

	migratedCount := 0
	skippedCount := 0
	errorCount := 0

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("❌ Failed to decode community: %v", err)
			errorCount++
			continue
		}

		communityID := doc["_id"]
		moderators, ok := doc["moderators"].(primitive.A)
		if !ok || len(moderators) == 0 {
			skippedCount++
			continue
		}

		// Check if already migrated
		firstMod, ok := moderators[0].(bson.M)
		if !ok {
			continue
		}

		// If avatar is already an object, skip
		if avatar, exists := firstMod["avatar"]; exists {
			if _, isMap := avatar.(bson.M); isMap {
				skippedCount++
				continue
			}
		}

		// Convert all moderators
		var newModerators primitive.A
		for _, modInterface := range moderators {
			modMap, ok := modInterface.(bson.M)
			if !ok {
				continue
			}

			// Convert avatar from string to object
			if avatarVal, exists := modMap["avatar"]; exists {
				if avatarStr, ok := avatarVal.(string); ok {
					if avatarStr != "" {
						modMap["avatar"] = bson.M{
							"url":         avatarStr,
							"public_id":   "",
							"uploaded_at": time.Now(),
						}
					} else {
						modMap["avatar"] = nil
					}
				}
			}

			newModerators = append(newModerators, modMap)
		}

		// Update community
		result, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": communityID},
			bson.M{"$set": bson.M{"moderators": newModerators}},
		)

		if err != nil {
			log.Printf("❌ Failed to update community %v: %v", communityID, err)
			errorCount++
			continue
		}

		if result.ModifiedCount > 0 {
			log.Printf("✅ Migrated community %v (%d moderators)", communityID, len(newModerators))
			migratedCount++
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("❌ Cursor error: %v", err)
	}

	fmt.Println("\n=== Migration Summary ===")
	fmt.Printf("✅ Migrated: %d communities\n", migratedCount)
	fmt.Printf("⏭️  Skipped (already migrated or no moderators): %d\n", skippedCount)
	fmt.Printf("❌ Errors: %d\n", errorCount)
	fmt.Println("✅ Migration completed successfully!")
}
