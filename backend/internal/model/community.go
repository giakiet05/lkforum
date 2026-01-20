package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Community struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string             `bson:"name,omitempty" json:"name,omitempty"`
	Description    *string            `bson:"description,omitempty" json:"description,omitempty"`
	Avatar         *string            `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Banner         *string            `bson:"banner,omitempty" json:"banner,omitempty"`
	Setting        CommunitySetting   `bson:"setting,omitempty" json:"setting,omitempty"`
	Rules          []CommunityRule    `bson:"rules,omitempty" json:"rules,omitempty"`
	Moderators     []Moderator        `bson:"moderators,omitempty" json:"moderators,omitempty"`
	MemberCount    int64              `bson:"member_count,omitempty" json:"member_count,omitempty"`
	PostCount      int64              `bson:"post_count,omitempty" json:"post_count,omitempty"`
	CreateAt       time.Time          `bson:"create_at,omitempty" json:"create_at,omitempty"`
	CreateByID     primitive.ObjectID `bson:"create_by_id,omitempty" json:"create_by_id,omitempty"`
	CreateByName   string             `bson:"create_by_name,omitempty" json:"create_by_name,omitempty"`
	CreateByAvatar string             `bson:"create_by_avatar,omitempty" json:"create_by_avatar,omitempty"`
	Is18Plus       bool               `bson:"is_18_plus" json:"is_18_plus"`
	IsDeleted      bool               `bson:"is_deleted" json:"is_deleted"`
	IsBanned       bool               `bson:"is_banned" json:"is_banned"`
	BanReason      *string            `bson:"ban_reason,omitempty" json:"ban_reason,omitempty"`
}

type CommunitySetting struct {
	IsPrivate           bool `bson:"is_private" json:"is_private"`                       //post visible only to members
	PostRequireApproval bool `bson:"post_require_approval" json:"post_require_approval"` // new posts need moderator approval
	JoinRequireApproval bool `bson:"join_require_approval" json:"join_require_approval"` // new member need moderator approval
	MaxPostLength       int  `bson:"max_post_length" json:"max_post_length"`
}

type CommunityRule struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
}

type Moderator struct {
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username   string             `bson:"username" json:"username"`
	Avatar     *Image             `bson:"avatar" json:"avatar"`
	IsActive   bool               `bson:"is_active" json:"is_active"`
	AssignedAt time.Time          `bson:"assigned_at" json:"assigned_at"`
}
