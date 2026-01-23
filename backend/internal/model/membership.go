package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MembershipStatus string

const (
	MembershipStatusPending  MembershipStatus = "pending"
	MembershipStatusApproved MembershipStatus = "approved"
	MembershipStatusRejected MembershipStatus = "rejected"
)

type Membership struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	CommunityID primitive.ObjectID `bson:"community_id" json:"community_id"`
	Status      MembershipStatus   `bson:"status,omitempty" json:"status,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	User        *UserInfo          `bson:"user,omitempty" json:"user,omitempty"`
}

type UserInfo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Avatar   *Image             `bson:"avatar,omitempty" json:"avatar,omitempty"`
}
