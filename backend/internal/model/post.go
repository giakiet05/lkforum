package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID      primitive.ObjectID `bson:"author_id" json:"author_id"`
	CommunityID   primitive.ObjectID `bson:"community_id" json:"community_id"`
	Type          string             `bson:"type" json:"type"`
	Content       *PostContent       `bson:"content,omitempty" json:"content,omitempty"`
	VotesCount    *VotesCount        `bson:"votes_count" json:"votes_count"`
	CommentsCount int                `bson:"comments_count,omitempty" json:"comments_count,omitempty"`
	CreatedAt     time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt     *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	IsDeleted     bool               `bson:"is_deleted,omitempty" json:"is_deleted,omitempty"`
}

type PostContent struct {
	Text  string `bson:"text,omitempty" json:"text,omitempty"`
	Poll  *Poll  `bson:"poll,omitempty" json:"poll,omitempty"`
	Video *Video `bson:"video,omitempty" json:"video,omitempty"`
}

type Poll struct {
	Question string   `bson:"question,omitempty" json:"question,omitempty"`
	Options  []Option `bson:"options,omitempty" json:"options,omitempty"`
}

type Option struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content string             `bson:"text" json:"text"`
	Vote    int                `bson:"vote" json:"vote"`
}

type Video struct {
	Title     string `bson:"title,omitempty" json:"title,omitempty"`
	Thumbnail string `bson:"thumbnail,omitempty" json:"thumbnail,omitempty"`
	URL       string `bson:"url,omitempty" json:"url,omitempty"`
}

type VotesCount struct {
	Up   int `bson:"up" json:"up"`
	Down int `bson:"down" json:"down"`
}
