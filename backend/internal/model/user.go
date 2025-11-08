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
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username"`
	Email       string             `bson:"email" json:"email"`
	Reputation  int                `bson:"reputation" json:"reputation"`
	Password    string             `bson:"password,omitempty" json:"-"`
	Provider    AuthProvider       `bson:"provider" json:"provider"`
	ProviderID  string             `bson:"provider_id,omitempty" json:"-"`
	Role        Role               `bson:"role" json:"role"`
	RoleContent RoleContent        `bson:"role_content,omitempty" json:"role_content,omitempty"`
	Settings    *UserSettings      `bson:"settings,omitempty" json:"settings,omitempty"`
	IsVerified  bool               `bson:"is_verified" json:"is_verified"` // Always true for local users after registration
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
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
	// Visual
	Avatar *Image `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Cover  *Image `bson:"cover,omitempty" json:"cover,omitempty"`

	// Personal Info
	Bio         *string    `bson:"bio,omitempty" json:"bio,omitempty"`
	Gender      *Gender    `bson:"gender,omitempty" json:"gender,omitempty"`
	DateOfBirth *time.Time `bson:"date_of_birth,omitempty" json:"date_of_birth,omitempty"`
	Location    *VNProvince `bson:"location,omitempty" json:"location,omitempty"`
	Interests   []Interest `bson:"interests,omitempty" json:"interests,omitempty"`

	// Social Links
	SocialLinks *SocialLinks `bson:"social_links,omitempty" json:"social_links,omitempty"`

	// Activity Stats
	Stats *ActivityStats `bson:"stats,omitempty" json:"stats,omitempty"`

	// Status Fields
	BanStart *time.Time `bson:"ban_start,omitempty" json:"ban_start,omitempty"`
	BanEnd   *time.Time `bson:"ban_end,omitempty" json:"ban_end,omitempty"`
}

// SocialLinks contains user's social media links.
type SocialLinks struct {
	Website  *string `bson:"website,omitempty" json:"website,omitempty"`
	Facebook *string `bson:"facebook,omitempty" json:"facebook,omitempty"`
	YouTube  *string `bson:"youtube,omitempty" json:"youtube,omitempty"`
	GitHub   *string `bson:"github,omitempty" json:"github,omitempty"`
}

// ActivityStats tracks user's activity statistics.
type ActivityStats struct {
	PostCount     int       `bson:"post_count" json:"post_count"`
	CommentCount  int       `bson:"comment_count" json:"comment_count"`
	TotalUpvotes  int       `bson:"total_upvotes" json:"total_upvotes"`
	JoinedAt      time.Time `bson:"joined_at" json:"joined_at"`
	LastActiveAt  time.Time `bson:"last_active_at" json:"last_active_at"`
}

// AdminRoleContent contains data specific to an admin user.
type AdminRoleContent struct {
	Permissions []string `bson:"permissions,omitempty" json:"permissions,omitempty"`
}
