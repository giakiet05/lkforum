package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Draft represents an unfinished post.
type Draft struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID    primitive.ObjectID `bson:"author_id" json:"author_id"`
	CommunityID *string            `bson:"community_id,omitempty" json:"community_id,omitempty"`
	Type        *PostType          `bson:"type,omitempty" json:"type,omitempty"`
	Title       *string            `bson:"title,omitempty" json:"title,omitempty"`
	Content     *PostContent       `bson:"content,omitempty" json:"content,omitempty"`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
