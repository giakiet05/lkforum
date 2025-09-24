package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Type      string             `bson:"type,omitempty" json:"type,omitempty"`
	Message   string             `bson:"message,omitempty" json:"message,omitempty"`
	Metadata  interface{}        `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
