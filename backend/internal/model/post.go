package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID       primitive.ObjectID `bson:"author_id" json:"author_id"`
	AuthorUsername string             `bson:"author_username,omitempty" json:"author_username,omitempty"`
	AuthorAvatar   string             `bson:"author_avatar,omitempty" json:"author_avatar,omitempty"`
	CommunityID    primitive.ObjectID `bson:"community_id" json:"community_id"`
	CommunityName  string             `bson:"community_name,omitempty" json:"community_name,omitempty"`
	Type           PostType           `bson:"type" json:"type"`
	Title          string             `bson:"title" json:"title"`
	Content        *PostContent       `bson:"content,omitempty" json:"content,omitempty"`
	VotesCount     *VotesCount        `bson:"votes_count" json:"votes_count"`
	CommentsCount  int                `bson:"comment_count" json:"comment_count"`
	CreatedAt      time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	IsDeleted      bool               `bson:"is_deleted,omitempty" json:"is_deleted"`
}

type PostType string

const (
	PostTypeText  PostType = "text"
	PostTypePoll  PostType = "poll"
	PostTypeVideo PostType = "video"
	PostTypeImage PostType = "image"
)

type PostContent struct {
	Text   string  `bson:"text,omitempty" json:"text,omitempty"`
	Poll   *Poll   `bson:"poll,omitempty" json:"poll,omitempty"`
	Video  *Video  `bson:"video,omitempty" json:"video,omitempty"`
	Images []Image `bson:"images,omitempty" json:"images,omitempty"`
}

type Image struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL        string             `bson:"url" json:"url"`
	UploadedAt time.Time          `bson:"uploaded_at,omitempty" json:"uploaded_at,omitempty"`
}

type Poll struct {
	Question      string       `bson:"question,omitempty" json:"question,omitempty"`
	Options       []PollOption `bson:"options,omitempty" json:"options,omitempty"`
	TotalVotes    *int         `bson:"total_votes,omitempty" json:"total_votes,omitempty"`
	ExpiresAt     *time.Time   `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	AllowMultiple bool         `bson:"allow_multiple,omitempty" json:"allow_multiple,omitempty"`
}

type PollOption struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Text  string             `bson:"text" json:"text"`
	Votes int                `bson:"votes" json:"votes"`
}

type Video struct {
	Thumbnail string `bson:"thumbnail,omitempty" json:"thumbnail,omitempty"`
	URL       string `bson:"url,omitempty" json:"url,omitempty"`
}

type VotesCount struct {
	Up   int `bson:"up" json:"up"`
	Down int `bson:"down" json:"down"`
}
