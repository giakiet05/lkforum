package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post represents the core structure of a post in the forum.
type Post struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID      primitive.ObjectID `bson:"author_id" json:"author_id"`
	CommunityID   primitive.ObjectID `bson:"community_id" json:"community_id"`
	Type          PostType           `bson:"type" json:"type"`
	Title         string             `bson:"title" json:"title"`
	Content       *PostContent       `bson:"content,omitempty" json:"content,omitempty"`
	VotesCount    *VotesCount        `bson:"votes_count,omitempty" json:"votes_count"`
	CommentsCount int                `bson:"comments_count" json:"comments_count"`
	CreatedAt     time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt     *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	IsDeleted     bool               `bson:"is_deleted,omitempty" json:"is_deleted"`
	IsHidden      bool               `bson:"is_hidden,omitempty" json:"is_hidden,omitempty"`
	Tags          []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	IsDraft       bool               `bson:"is_draft,omitempty" json:"is_draft,omitempty"`
}

type PostType string

const (
	PostTypeText  PostType = "text"
	PostTypePoll  PostType = "poll"
	PostTypeVideo PostType = "video"
	PostTypeImage PostType = "image"
)

// PostContent holds the actual content of the post, varying by type.
type PostContent struct {
	Text   string   `bson:"text,omitempty" json:"text,omitempty"`
	Images []Image  `bson:"images,omitempty" json:"images,omitempty"` // Uses model.Image from common.go
	Videos []*Video `bson:"videos,omitempty" json:"videos,omitempty"` // Uses model.Video from common.go
	Poll   *Poll    `bson:"poll,omitempty" json:"poll,omitempty"`
}

// Poll represents a poll within a post.
type Poll struct {
	Question      string       `bson:"question" json:"question"`
	Options       []PollOption `bson:"options" json:"options"`
	TotalVotes    int          `bson:"total_votes" json:"total_votes"`
	ExpiresAt     *time.Time   `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	AllowMultiple bool         `bson:"allow_multiple" json:"allow_multiple"`
}

// PollOption represents a single option in a poll.
type PollOption struct {
	ID    string `bson:"id" json:"id"`
	Text  string `bson:"text" json:"text"`
	Votes int    `bson:"votes" json:"votes"`
}

// VotesCount stores the up and down vote counts.
type VotesCount struct {
	Up   int `bson:"up" json:"up"`
	Down int `bson:"down" json:"down"`
}
