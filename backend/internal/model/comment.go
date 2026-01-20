package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Author     CommentAuthor       `bson:"author" json:"author"`
	PostID     primitive.ObjectID  `bson:"post_id" json:"post_id"`
	ParentID   *primitive.ObjectID `bson:"parent_id" json:"parent_id"`
	Content    string              `bson:"content" json:"content"`
	VotesCount *VotesCount         `bson:"votes_count,omitempty" json:"votes_count,omitempty"`
	CreatedAt  time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
	DeletedAt  *time.Time          `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	IsDeleted  bool                `bson:"is_deleted" json:"is_deleted"`
}

type CommentAuthor struct {
	ID       primitive.ObjectID `bson:"id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Avatar   *Image             `bson:"avatar" json:"avatar"`
}
