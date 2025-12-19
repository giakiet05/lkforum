package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load config
	config.LoadConfig()

	// Connect to MongoDB
	client := config.NewMongoClient()
	db := client.Database(config.Cfg.DBName)
	usersCollection := db.Collection(config.UserColName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Admin credentials
	adminEmail := "admin@lkforum.com"
	adminUsername := "admin"
	adminPassword := "admin123" // Change this to a secure password

	// Check if admin already exists
	var existingUser model.User
	err := usersCollection.FindOne(ctx, bson.M{"email": adminEmail}).Decode(&existingUser)
	if err == nil {
		log.Printf("Admin user already exists: %s (%s)\n", existingUser.Username, existingUser.Email)
		return
	}
	if err != mongo.ErrNoDocuments {
		log.Fatalf("Error checking for existing admin: %v", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	// Create admin user
	adminUser := model.User{
		ID:         primitive.NewObjectID(),
		Username:   adminUsername,
		Email:      adminEmail,
		Password:   string(hashedPassword),
		Provider:   model.ProviderLocal,
		Role:       model.AdminRole,
		Reputation: 0,
		IsVerified: true,
		IsBanned:   false,
		CreatedAt:  time.Now(),
		RoleContent: model.RoleContent{
			AsAdmin: &model.AdminRoleContent{
				Permissions: []string{"all"},
			},
		},
	}

	// Insert admin user
	_, err = usersCollection.InsertOne(ctx, adminUser)
	if err != nil {
		log.Fatalf("Error creating admin user: %v", err)
	}

	fmt.Println("✅ Admin user created successfully!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Email:    %s\n", adminEmail)
	fmt.Printf("Username: %s\n", adminUsername)
	fmt.Printf("Password: %s\n", adminPassword)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("⚠️  Please change the password after first login!")
}
