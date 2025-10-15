package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vote struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	TargetType VoteTargetType     `bson:"target_type,omitempty" json:"target_type,omitempty"`
	TargetID   primitive.ObjectID `bson:"target_id" json:"target_id"`
	Value      bool               `bson:"value" json:"value"`
	CreateAt   time.Time          `bson:"create_at,omitempty" json:"create_at,omitempty"`
}

type PollVote struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID    primitive.ObjectID `bson:"post_id" json:"post_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	OptionID  primitive.ObjectID `bson:"option_id" json:"option_id"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
type VoteType string

const (
	VoteTypeUp   VoteType = "up"
	VoteTypeDown VoteType = "down"
)

type VoteTargetType string

const (
	VoteTargetPost    VoteTargetType = "post"
	VoteTargetComment VoteTargetType = "comment"
)
