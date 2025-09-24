package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikedPost struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID  primitive.ObjectID `bson:"user_id" json:"user_id"`
	PostID  primitive.ObjectID `bson:"post_id" json:"post_id"`
	LikedAt time.Time          `bson:"liked_at" json:"liked_at"`
}
