package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vote struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	TargetType VoteTargetType     `bson:"target_type,omitempty" json:"target_type,omitempty"`
	TargetID   string             `bson:"target_id" json:"target_id"`
	Value      bool               `bson:"value" json:"value"`
	CreateAt   time.Time          `bson:"create_at,omitempty" json:"create_at,omitempty"`
}

type VoteTargetType string

const (
	VoteTargetPost    VoteTargetType = "post"
	VoteTargetComment VoteTargetType = "comment"
)
