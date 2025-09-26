package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Type      string               `bson:"type" json:"type"` // direct or group
	Members   []primitive.ObjectID `bson:"members" json:"members"`
	Name      string               `bson:"name,omitempty" json:"name,omitempty"`
	Avatar    string               `bson:"avatar,omitempty" json:"avatar,omitempty"`
	CreatedBy primitive.ObjectID   `bson:"created_by" json:"created_by"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}
