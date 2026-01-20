package main

import (
	"context"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

// Migration script to convert Moderator.Avatar from string to *Image
func main() {
	log.Println("Starting moderator avatar migration...")

	// Load config
	cfg := config.LoadEnv()
	mongoClient := config.NewMongoClient(cfg)
	db := mongoClient.Database(cfg.DBName)
	communitiesCollection := db.Collection("communities")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find all communities with moderators
	cursor, err := communitiesCollection.Find(ctx, bson.M{
		"moderators": bson.M{"$exists": true, "$ne": []interface{}{}},
	})
	if err != nil {
		log.Fatalf("Failed to find communities: %v", err)
	}
	defer cursor.Close(ctx)

	updatedCount := 0
	errorCount := 0

	for cursor.Next(ctx) {
		var community struct {
			ID         interface{}   `bson:"_id"`
			Moderators []interface{} `bson:"moderators"`
		}

		if err := cursor.Decode(&community); err != nil {
			log.Printf("Failed to decode community: %v", err)
			errorCount++
			continue
		}

		// Check if migration needed (first moderator has string avatar)
		if len(community.Moderators) == 0 {
			continue
		}

		needsMigration := false
		moderatorMap, ok := community.Moderators[0].(bson.M)
		if ok {
			if avatarVal, exists := moderatorMap["avatar"]; exists {
				// Check if avatar is string (needs migration)
				if _, isString := avatarVal.(string); isString {
					needsMigration = true
				}
			}
		}

		if !needsMigration {
			log.Printf("Community %v already migrated, skipping", community.ID)
			continue
		}

		// Convert moderators
		var newModerators []model.Moderator
		for _, modInterface := range community.Moderators {
			modMap, ok := modInterface.(bson.M)
			if !ok {
				continue
			}

			var moderator model.Moderator

			// Convert basic fields
			if userID, ok := modMap["user_id"].(interface{}); ok {
				moderator.UserID = userID.(interface{ Hex() string }).(interface{ Bytes() [12]byte }).(interface{}).(interface{})
			}
			if username, ok := modMap["username"].(string); ok {
				moderator.Username = username
			}
			if isActive, ok := modMap["is_active"].(bool); ok {
				moderator.IsActive = isActive
			}
			if assignedAt, ok := modMap["assigned_at"].(time.Time); ok {
				moderator.AssignedAt = assignedAt
			}

			// Convert avatar from string to *Image
			if avatarStr, ok := modMap["avatar"].(string); ok && avatarStr != "" {
				moderator.Avatar = &model.Image{
					URL:        avatarStr,
					PublicID:   "",
					UploadedAt: time.Now(),
				}
			} else {
				moderator.Avatar = nil
			}

			newModerators = append(newModerators, moderator)
		}

		// Update community with converted moderators
		update := bson.M{
			"$set": bson.M{
				"moderators": newModerators,
			},
		}

		result, err := communitiesCollection.UpdateOne(
			ctx,
			bson.M{"_id": community.ID},
			update,
		)

		if err != nil {
			log.Printf("Failed to update community %v: %v", community.ID, err)
			errorCount++
			continue
		}

		if result.ModifiedCount > 0 {
			log.Printf("✅ Migrated community %v (%d moderators)", community.ID, len(newModerators))
			updatedCount++
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}

	log.Printf("\n=== Migration Summary ===")
	log.Printf("✅ Successfully migrated: %d communities", updatedCount)
	log.Printf("❌ Errors: %d", errorCount)
	log.Println("Migration completed!")
}
