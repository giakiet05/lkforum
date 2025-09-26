package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Community struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Avatar      string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Banner      string             `bson:"banner,omitempty" json:"banner,omitempty"`
	Setting     *CommunitySetting  `bson:"setting,omitempty" json:"setting,omitempty"`
	Moderators  *Moderator         `bson:"moderators,omitempty" json:"moderators,omitempty"`
	CreateAt    time.Time          `bson:"create_at,omitempty" json:"create_at,omitempty"`
}

type CommunitySetting struct {
	IsPrivate            bool `bson:"is_private,omitempty" json:"is_private,omitempty"`
	PostApprovalRequired bool `bson:"post_approval_required,omitempty" json:"post_approval_required,omitempty"`
	JoinApprovalRequired bool `bson:"join_approval_required,omitempty" json:"join_approval_required,omitempty"`
}

type Moderator struct {
	UserID     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	AssignedAt time.Time          `bson:"assigned_at,omitempty" json:"assigned_at,omitempty"`
}
