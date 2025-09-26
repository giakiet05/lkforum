package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	ConversationID primitive.ObjectID   `bson:"conversation_id" json:"conversation_id"`
	SenderID       *primitive.ObjectID  `bson:"sender_id,omitempty" json:"sender_id,omitempty"` // nil for system messages
	Type           MessageType          `bson:"type" json:"type"`
	Content        string               `bson:"content" json:"content"`
	CreatedAt      time.Time            `bson:"create_at" json:"create_at"`
	ReadBy         []primitive.ObjectID `bson:"read_by,omitempty" json:"read_by,omitempty"`
}

type MessageType string

const (
	MessageTypeUser   MessageType = "user"
	MessageTypeSystem MessageType = "system"
)
