package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Community struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name,omitempty" json:"name,omitempty"`
	Description    *string            `bson:"description,omitempty" json:"description,omitempty"`
	Avatar         *string            `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Banner         *string            `bson:"banner,omitempty" json:"banner,omitempty"`
	Setting        CommunitySetting   `bson:"setting,omitempty" json:"setting,omitempty"`
	Moderators     []Moderator        `bson:"moderators,omitempty" json:"moderators,omitempty"`
	MemberCount    int64              `bson:"member_count,omitempty" json:"member_count,omitempty"`
	PostCount      int64              `bson:"post_count,omitempty" json:"post_count,omitempty"`
	CreateAt       time.Time          `bson:"create_at,omitempty" json:"create_at,omitempty"`
	CreateByID     primitive.ObjectID `bson:"create_by_id,omitempty" json:"create_by_id,omitempty"`
	CreateByName   string             `bson:"create_by_name,omitempty" json:"create_by_name,omitempty"`
	CreateByAvatar string             `bson:"create_by_avatar,omitempty" json:"create_by_avatar,omitempty"`
	IsDeleted      bool               `bson:"is_deleted,omitempty" json:"is_deleted,omitempty"`
	IsBanned       bool               `bson:"is_banned,omitempty" json:"is_banned,omitempty"`
}

type CommunitySetting struct {
	IsPrivate           bool `bson:"isPrivate" json:"isPrivate"` // visible only to members
	AllowPosts          bool `bson:"allowPosts" json:"allowPosts"`
	AllowComments       bool `bson:"allowComments" json:"allowComments"`
	AllowMedia          bool `bson:"allowMedia" json:"allowMedia"`
	PostRequireApproval bool `bson:"requireApproval" json:"requireApproval"`         // new posts need moderator approval
	JoinRequireApproval bool `bson:"joinRequireApproval" json:"joinRequireApproval"` // new member need moderator approval
	MaxPostLength       int  `bson:"maxPostLength,omitempty" json:"maxPostLength,omitempty"`
}

type Moderator struct {
	UserID     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Username   string             `bson:"username,omitempty" json:"username,omitempty"`
	AssignedAt time.Time          `bson:"assigned_at,omitempty" json:"assigned_at,omitempty"`
}
