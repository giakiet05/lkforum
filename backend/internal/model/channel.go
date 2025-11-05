package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Members   []ChannelMember    `bson:"members" json:"members"`
	Settings  []ChannelSetting   `bson:"settings" json:"settings"`
	Status    ChannelStatus      `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type ChannelStatus string

const (
	ChannelStatusActive ChannelStatus = "active"
	ChannelStatusBlock  ChannelStatus = "block"
)

type ChannelSetting struct {
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	Nickname        *string            `bson:"nickname,omitempty" json:"nickname,omitempty"`
	Notification    bool               `bson:"notification" json:"notification"`
	TypingIndicator bool               `bson:"typing_indicator" json:"typing_indicator"`
	IsDeleted       bool               `bson:"is_deleted" json:"is_deleted"`
}

type ChannelMember struct {
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username string             `bson:"username" json:"username"`
	Avatar   string             `bson:"avatar" json:"avatar"`
}
