package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	AuthorID  primitive.ObjectID  `bson:"author_id" json:"author_id"`
	PostID    primitive.ObjectID  `bson:"post_id" json:"post_id"`
	ParentID  *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Content   string              `bson:"content" json:"content"`
	CreatedAt time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	IsDeleted bool                `bson:"is_deleted,omitempty" json:"is_deleted,omitempty"`
}
