package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPostHistory struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID   string             `bson:"user_id" json:"user_id"`
	PostID   string             `bson:"post_id" json:"post_id"`
	ViewedAt time.Time          `bson:"viewed_at" json:"viewed_at"`
}
