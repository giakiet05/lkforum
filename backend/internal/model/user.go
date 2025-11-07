package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthProvider defines the source of user authentication.
type AuthProvider string

const (
	ProviderLocal  AuthProvider = "local"  // Registered with email and password
	ProviderGoogle AuthProvider = "google" // Registered via Google OAuth
)

type User struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username                  string             `bson:"username" json:"username"`
	Email                     string             `bson:"email" json:"email"`
	Reputation                int                `bson:"reputation" json:"reputation"`
	Password                  string             `bson:"password,omitempty" json:"-"`
	Provider                  AuthProvider       `bson:"provider" json:"provider"`
	ProviderID                string             `bson:"provider_id,omitempty" json:"-"`
	Role                      Role               `bson:"role" json:"role"`
	RoleContent               RoleContent        `bson:"role_content,omitempty" json:"role_content,omitempty"`
	IsVerified                bool               `bson:"is_verified" json:"is_verified"`
	VerificationCode          string             `bson:"verification_code,omitempty" json:"-"`
	VerificationCodeExpiresAt *time.Time         `bson:"verification_code_expires_at,omitempty" json:"-"`
	CreatedAt                 time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	DeletedAt                 *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

// RoleContent holds the role-specific data for a user.
type RoleContent struct {
	AsUser  *UserRoleContent  `bson:"as_user,omitempty" json:"as_user,omitempty"`
	AsAdmin *AdminRoleContent `bson:"as_admin,omitempty" json:"as_admin,omitempty"`
}

// UserRoleContent contains data specific to a regular user's profile and status.
type UserRoleContent struct {
	// Profile Fields
	Avatar   Image  `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Cover    Image  `bson:"cover,omitempty" json:"cover,omitempty"`
	Bio      string `bson:"bio,omitempty" json:"bio,omitempty"`
	Location string `bson:"location,omitempty" json:"location,omitempty"`
	Website  string `bson:"website,omitempty" json:"website,omitempty"`

	// Status Fields
	BanStart *time.Time `bson:"ban_start,omitempty" json:"ban_start,omitempty"`
	BanEnd   *time.Time `bson:"ban_end,omitempty" json:"ban_end,omitempty"`
}

// AdminRoleContent contains data specific to an admin user.
type AdminRoleContent struct {
	Permissions []string `bson:"permissions,omitempty" json:"permissions,omitempty"`
}
