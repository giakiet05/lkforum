package controller

import (
	"net/http"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type DebugController struct {
	db *mongo.Database
}

func NewDebugController(db *mongo.Database) *DebugController {
	return &DebugController{
		db: db,
	}
}

// CreateAdminUser creates an admin user for development
// ⚠️ THIS SHOULD BE REMOVED IN PRODUCTION
func (c *DebugController) CreateAdminUser(ctx *gin.Context) {
	dbCtx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Admin credentials
	adminEmail := "admin@lkforum.com"
	adminUsername := "admin"
	adminPassword := "admin123"

	collection := c.db.Collection(config.UserColName)

	// Check if admin already exists
	var existingUser model.User
	err := collection.FindOne(dbCtx, bson.M{"email": adminEmail}).Decode(&existingUser)
	if err == nil {
		dto.SendSuccess(ctx, http.StatusOK, "Admin already exists", gin.H{
			"email":    adminEmail,
			"username": adminUsername,
			"password": adminPassword,
			"message":  "✅ Use these credentials to login",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to hash password", "HASH_ERROR")
		return
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
	_, err = collection.InsertOne(dbCtx, adminUser)
	if err != nil {
		dto.SendError(ctx, http.StatusInternalServerError, "Failed to create admin", "CREATE_ERROR")
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "✅ Admin user created successfully!", gin.H{
		"email":    adminEmail,
		"username": adminUsername,
		"password": adminPassword,
		"message":  "Login at http://localhost:5174",
	})
}
