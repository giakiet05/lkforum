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

	// Ban fields
	IsBanned  bool       `bson:"is_banned" json:"is_banned"`
	BanUntil  *time.Time `bson:"ban_until,omitempty" json:"ban_until,omitempty"` // null = permanent ban
	BanReason *string    `bson:"ban_reason,omitempty" json:"ban_reason,omitempty"`

	CreatedAt time.Time  `bson:"created_at,omitempty" json:"created_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // Soft delete by user
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
	Bio         *string     `bson:"bio,omitempty" json:"bio,omitempty"`
	Gender      *Gender     `bson:"gender,omitempty" json:"gender,omitempty"`
	DateOfBirth *time.Time  `bson:"date_of_birth,omitempty" json:"date_of_birth,omitempty"`
	Location    *VNProvince `bson:"location,omitempty" json:"location,omitempty"`
	Interests   []Interest  `bson:"interests,omitempty" json:"interests,omitempty"`

	// Social Links
	SocialLinks *SocialLinks `bson:"social_links,omitempty" json:"social_links,omitempty"`

	// Activity Stats
	Stats *ActivityStats `bson:"stats,omitempty" json:"stats,omitempty"`
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
	PostCount    int       `bson:"post_count" json:"post_count"`
	CommentCount int       `bson:"comment_count" json:"comment_count"`
	TotalUpvotes int       `bson:"total_upvotes" json:"total_upvotes"`
	JoinedAt     time.Time `bson:"joined_at" json:"joined_at"`
	LastActiveAt time.Time `bson:"last_active_at" json:"last_active_at"`
}

// AdminRoleContent contains data specific to an admin user.
type AdminRoleContent struct {
	Permissions []string `bson:"permissions,omitempty" json:"permissions,omitempty"`
}

// IsBannedNow checks if user is currently banned
func (u *User) IsBannedNow() bool {
	if !u.IsBanned {
		return false
	}

	// BanUntil = nil → permanent ban
	if u.BanUntil == nil {
		return true
	}

	// BanUntil > now → still banned
	return u.BanUntil.After(time.Now())
}

func CloneUser(u *User) *User {
	if u == nil {
		return nil
	}

	clone := *u // shallow copy

	// Deep copy DeletedAt
	if u.DeletedAt != nil {
		t := *u.DeletedAt
		clone.DeletedAt = &t
	}

	// Deep copy BanUntil
	if u.BanUntil != nil {
		t := *u.BanUntil
		clone.BanUntil = &t
	}

	// Deep copy BanReason
	if u.BanReason != nil {
		s := *u.BanReason
		clone.BanReason = &s
	}

	// Deep copy Settings
	if u.Settings != nil {
		s := *u.Settings
		clone.Settings = &s
	}

	// Deep copy RoleContent
	clone.RoleContent = CloneRoleContent(u.RoleContent)

	return &clone
}

func CloneRoleContent(rc RoleContent) RoleContent {
	newRC := RoleContent{}

	if rc.AsUser != nil {
		newRC.AsUser = CloneUserRoleContent(rc.AsUser)
	}

	if rc.AsAdmin != nil {
		newRC.AsAdmin = CloneAdminRoleContent(rc.AsAdmin)
	}

	return newRC
}

func CloneUserRoleContent(u *UserRoleContent) *UserRoleContent {
	if u == nil {
		return nil
	}

	clone := *u

	// Avatar
	if u.Avatar != nil {
		img := *u.Avatar
		clone.Avatar = &img
	}

	// Cover
	if u.Cover != nil {
		img := *u.Cover
		clone.Cover = &img
	}

	// Bio
	if u.Bio != nil {
		v := *u.Bio
		clone.Bio = &v
	}

	// Gender
	if u.Gender != nil {
		v := *u.Gender
		clone.Gender = &v
	}

	// DateOfBirth
	if u.DateOfBirth != nil {
		t := *u.DateOfBirth
		clone.DateOfBirth = &t
	}

	// Location
	if u.Location != nil {
		v := *u.Location
		clone.Location = &v
	}

	// Interests (slice copy)
	if len(u.Interests) > 0 {
		clone.Interests = append([]Interest(nil), u.Interests...)
	}

	// SocialLinks
	if u.SocialLinks != nil {
		sl := *u.SocialLinks

		// Deep copy all string pointers
		if u.SocialLinks.Website != nil {
			s := *u.SocialLinks.Website
			sl.Website = &s
		}
		if u.SocialLinks.Facebook != nil {
			s := *u.SocialLinks.Facebook
			sl.Facebook = &s
		}
		if u.SocialLinks.YouTube != nil {
			s := *u.SocialLinks.YouTube
			sl.YouTube = &s
		}
		if u.SocialLinks.GitHub != nil {
			s := *u.SocialLinks.GitHub
			sl.GitHub = &s
		}

		clone.SocialLinks = &sl
	}

	// Activity Stats
	if u.Stats != nil {
		stats := *u.Stats
		clone.Stats = &stats
	}

	return &clone
}

func CloneAdminRoleContent(a *AdminRoleContent) *AdminRoleContent {
	if a == nil {
		return nil
	}

	clone := *a

	// deep copy slice
	if len(a.Permissions) > 0 {
		clone.Permissions = append([]string(nil), a.Permissions...)
	}

	return &clone
}
