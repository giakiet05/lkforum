package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Membership struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	CommunityID primitive.ObjectID `bson:"community_id" json:"community_id"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	User        *UserInfo          `bson:"user,omitempty" json:"user,omitempty"`
}

type UserInfo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Avatar   *Image             `bson:"avatar,omitempty" json:"avatar,omitempty"`
}
