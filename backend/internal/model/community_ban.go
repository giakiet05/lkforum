package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommunityBan struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id"`
	CommunityID primitive.ObjectID `bson:"community_id"`
	Type        CommunityBanType   `bson:"type"`
	Reason      string             `bson:"reason"`
	BannedBy    primitive.ObjectID `bson:"banned_by"`
	BannedAt    time.Time          `bson:"banned_at"`
	ExpiresAt   time.Time          `bson:"expires_at"`
	IsDeleted   bool               `bson:"is_deleted"`
	DeletedAt   time.Time          `bson:"deleted_at"`
}

type CommunityBanType string

const (
	Banned CommunityBanType = "banned"
	Muted  CommunityBanType = "muted"
)
