package model

import (
	"math"
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
	HotScore      float64            `bson:"hot_score" json:"hot_score"` // Score for "best" sorting (votes + time decay)
	CreatedAt     time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt     *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	IsDeleted     bool               `bson:"is_deleted,omitempty" json:"is_deleted"`
	IsHidden      bool               `bson:"is_hidden,omitempty" json:"is_hidden,omitempty"`
	IsEdited      bool               `bson:"is_edited,omitempty" json:"is_edited,omitempty"`
	Tags          []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	IsDraft       bool               `bson:"is_draft,omitempty" json:"is_draft,omitempty"`

	IsBan     bool    `bson:"is_ban,omitempty" json:"is_ban,omitempty"`
	BanReason *string `bson:"ban_reason,omitempty" json:"ban_reason,omitempty"`

	// Moderation fields
	ModerationStatus ModerationStatus  `bson:"moderation_status,omitempty" json:"moderation_status,omitempty"`
	ModerationResult *ModerationResult `bson:"moderation_result,omitempty" json:"moderation_result,omitempty"`
	ModeratedAt      *time.Time        `bson:"moderated_at,omitempty" json:"moderated_at,omitempty"`
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

// ModerationStatus represents the moderation state of a post/comment
type ModerationStatus string

const (
	ModerationPending  ModerationStatus = "pending"  // Waiting for moderation
	ModerationApproved ModerationStatus = "approved" // Approved, visible to all
	ModerationRejected ModerationStatus = "rejected" // Rejected, only author can see
	ModerationSkipped  ModerationStatus = "skipped"  // Skipped (admin/mod posts)
)

// ModerationResult contains the result of AI moderation
type ModerationResult struct {
	IsViolation  bool     `bson:"is_violation" json:"is_violation"`
	Confidence   float64  `bson:"confidence" json:"confidence"`       // 0.0 - 1.0
	Categories   []string `bson:"categories" json:"categories"`       // ["hate_speech", "violence", ...]
	Reason       string   `bson:"reason" json:"reason"`               // Explanation in Vietnamese
	CheckedText  bool     `bson:"checked_text" json:"checked_text"`   // Was text content checked
	CheckedMedia bool     `bson:"checked_media" json:"checked_media"` // Were images/videos checked
}

// CalculateHotScore calculates the hot score for a post based on Reddit's algorithm
// Score = log10(max(|score|, 1)) + sign(score) * (created_at - epoch) / 45000
// This gives newer posts a chance to compete with older popular posts
func CalculateHotScore(upvotes, downvotes int, createdAt time.Time) float64 {
	score := upvotes - downvotes

	// Use absolute score, minimum 1 to avoid log(0)
	absScore := score
	if absScore < 0 {
		absScore = -absScore
	}
	if absScore < 1 {
		absScore = 1
	}

	// Sign of score: 1 for positive, -1 for negative, 0 for zero
	sign := 0.0
	if score > 0 {
		sign = 1.0
	} else if score < 0 {
		sign = -1.0
	}

	// Seconds since epoch (Reddit uses Dec 8, 2005, we use a more recent date)
	epoch := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	seconds := createdAt.Sub(epoch).Seconds()

	// Hot score formula
	// The 45000 divisor means a post needs ~12.5 hours to drop one "point" in ranking
	hotScore := math.Log10(float64(absScore)) + sign*seconds/45000.0

	return hotScore
}
