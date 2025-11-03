package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	ChannelID      primitive.ObjectID  `bson:"channel_id" json:"channel_id"`
	SenderID       *primitive.ObjectID `bson:"sender_id,omitempty" json:"sender_id,omitempty"` // nil for system messages
	SenderUsername string              `bson:"sender_username,omitempty" json:"sender_username,omitempty"`
	Type           MessageType         `bson:"type" json:"type"`
	Content        string              `bson:"content" json:"content"`
	IsRead         bool                `bson:"is_read" json:"is_read"`
	IsSend         bool                `bson:"is_send" json:"is_send"`
	CreatedAt      time.Time           `bson:"created_at" json:"created_at"`
	IsDeleted      bool                `bson:"is_deleted" json:"is_deleted"`
	DeletedAt      *time.Time          `bson:"deleted_at" json:"deleted_at"`
}

type MessageType string

const (
	UserMessage   MessageType = "user"
	SystemMessage MessageType = "system"
)
