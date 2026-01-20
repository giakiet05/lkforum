package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

// --- Request DTOs ---

// CreatePostRequest defines the structure for creating a new text or poll post.
type CreatePostRequest struct {
	CommunityID string             `json:"community_id" binding:"required"`
	Title       string             `json:"title" binding:"required,min=3,max=300"`
	Type        model.PostType     `json:"type" binding:"required,oneof=text poll image video"`
	Text        string             `json:"text,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
	Poll        *CreatePollRequest `json:"poll,omitempty"`
}

// CreatePollRequest defines the structure for creating a poll within a post.
type CreatePollRequest struct {
	Question      string     `json:"question" binding:"required,min=2,max=500"`
	Options       []string   `json:"options" binding:"required,min=2,dive,min=1,max=200"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	AllowMultiple bool       `json:"allow_multiple,omitempty"`
}

// UpdatePostRequest defines the structure for updating a post's simple fields.
type UpdatePostRequest struct {
	Title *string   `json:"title,omitempty"`
	Text  *string   `json:"text,omitempty"`
	Tags  *[]string `json:"tags,omitempty"`
}

// UpdatePollRequest defines the structure for updating a poll's fields.
type UpdatePollRequest struct {
	Question      *string    `json:"question,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	AllowMultiple *bool      `json:"allow_multiple,omitempty"`
}

// UpdatePollOptionRequest defines the structure for updating a single poll option.
type UpdatePollOptionRequest struct {
	Text string `json:"text" binding:"required,min=1,max=200"`
}

// AddPollOptionsRequest defines the structure for adding new options to a poll.
type AddPollOptionsRequest struct {
	Options []string `json:"options" binding:"required,min=1,dive,min=1"`
}

// RemovePollOptionsRequest defines the structure for removing options from a poll.
type RemovePollOptionsRequest struct {
	OptionIDs []string `json:"option_ids" binding:"required,min=1"`
}

// PostVoteRequest defines the structure for voting on a post.
type PostVoteRequest struct {
	Value *bool `json:"value" binding:"required"`
}

// PollVoteRequest defines the structure for voting on a poll.
type PollVoteRequest struct {
	OptionID string `json:"option_id" binding:"required"`
}

// RemoveImagesRequest defines the structure for removing images from a post.
type RemoveImagesRequest struct {
	PublicIDs []string `json:"public_ids" binding:"required,min=1"`
}

// ReportPostRequest defines the structure for reporting a post.
type ReportPostRequest struct {
	Reason      string `json:"reason" binding:"required"`
	Description string `json:"description,omitempty"`
}

// GetPostsQuery defines the query parameters for fetching posts.
type GetPostsQuery struct {
	CommunityID string `form:"community_id"`
	AuthorID    string `form:"author_id"`
	Type        string `form:"type"`
	Sort        string `form:"sort"`
	TimeFrame   string `form:"time"`
	FeedType    string `form:"feed_type"` // home, popular, explore, all
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
}

type GetBanPostsQuery struct {
	CommunityID string `form:"community_id"`
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
}

// --- Response DTOs ---

// PostResponse is the detailed post object returned to the client.
type PostResponse struct {
	ID               string                 `json:"id"`
	Author           AuthorResponse         `json:"author"`
	Community        CommunityShortResponse `json:"community"`
	Title            string                 `json:"title"`
	Type             model.PostType         `json:"type"`
	Content          PostContentResponse    `json:"content"`
	VotesCount       *VotesCountResponse    `json:"votes_count,omitempty"`
	UserVote         string                 `json:"user_vote,omitempty"`
	CommentsCount    int                    `json:"comments_count"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        *time.Time             `json:"updated_at,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
	ModerationStatus *string                `json:"moderation_status,omitempty"`
	ModerationReason *string                `json:"moderation_reason,omitempty"`
	ModeratedAt      *time.Time             `json:"moderated_at,omitempty"`
}

// AuthorResponse contains short public information about a user.
type AuthorResponse struct {
	ID       string       `json:"id"`
	Username string       `json:"username"`
	Avatar   *model.Image `json:"avatar,omitempty"`
}

// CommunityShortResponse contains short public information about a community.
type CommunityShortResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PostContentResponse holds the actual content of the post.
type PostContentResponse struct {
	Text   string         `json:"text,omitempty"`
	Images []model.Image  `json:"images,omitempty"`
	Videos []*model.Video `json:"videos,omitempty"`
	Poll   *PollResponse  `json:"poll,omitempty"`
}

// PollResponse is the detailed poll object for responses.
type PollResponse struct {
	Question      string               `json:"question"`
	Options       []PollOptionResponse `json:"options"`
	TotalVotes    int                  `json:"total_votes"`
	UserVoteIDs   []string             `json:"user_vote_ids,omitempty"`
	ExpiresAt     *time.Time           `json:"expires_at,omitempty"`
	AllowMultiple bool                 `json:"allow_multiple"`
}

// PollOptionResponse represents a single poll option in a response.
type PollOptionResponse struct {
	ID         string  `json:"id"`
	Text       string  `json:"text"`
	Votes      int     `json:"votes"`
	Percentage float64 `json:"percentage"`
}

// VotesCountResponse represents the vote counts for a post or comment.
type VotesCountResponse struct {
	Up    int `json:"up"`
	Down  int `json:"down"`
	Score int `json:"score"`
}

// --- Mapping Functions ---

// FromPost creates a PostResponse from a Post model and enriched data.
func FromPost(post *model.Post, author *model.User, community *model.Community, userVote string, userPollVoteIDs []string) *PostResponse {
	if post == nil {
		return nil
	}

	if author == nil {
		author = &model.User{Username: "[deleted]"} // Graceful handling of deleted user
	}
	if community == nil {
		community = &model.Community{Name: "[deleted]"} // Graceful handling of deleted community
	}

	// Safely access the nested avatar to prevent nil pointer dereference
	var authorAvatar *model.Image
	if author.RoleContent.AsUser != nil {
		authorAvatar = author.RoleContent.AsUser.Avatar
	}

	resp := &PostResponse{
		ID:            post.ID.Hex(),
		Author:        AuthorResponse{ID: author.ID.Hex(), Username: author.Username, Avatar: authorAvatar},
		Community:     CommunityShortResponse{ID: community.ID.Hex(), Name: community.Name},
		Title:         post.Title,
		Type:          post.Type,
		CommentsCount: post.CommentsCount,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
		UserVote:      userVote,
		Tags:          post.Tags,
	}

	if post.VotesCount != nil {
		resp.VotesCount = &VotesCountResponse{Up: post.VotesCount.Up, Down: post.VotesCount.Down, Score: post.VotesCount.Up - post.VotesCount.Down}
	}

	if post.Content != nil {
		resp.Content = PostContentResponse{
			Text:   post.Content.Text,
			Images: post.Content.Images,
			Videos: post.Content.Videos,
		}
		if post.Content.Poll != nil {
			resp.Content.Poll = FromPoll(post.Content.Poll, userPollVoteIDs)
		}
	}

	return resp
}

// FromPoll creates a PollResponse from a Poll model.
func FromPoll(poll *model.Poll, userVoteIDs []string) *PollResponse {
	if poll == nil {
		return nil
	}

	options := make([]PollOptionResponse, len(poll.Options))
	for i, opt := range poll.Options {
		percentage := 0.0
		if poll.TotalVotes > 0 {
			percentage = (float64(opt.Votes) / float64(poll.TotalVotes)) * 100
		}
		options[i] = PollOptionResponse{
			ID:         opt.ID,
			Text:       opt.Text,
			Votes:      opt.Votes,
			Percentage: percentage,
		}
	}

	return &PollResponse{
		Question:      poll.Question,
		Options:       options,
		TotalVotes:    poll.TotalVotes,
		UserVoteIDs:   userVoteIDs,
		ExpiresAt:     poll.ExpiresAt,
		AllowMultiple: poll.AllowMultiple,
	}
}

// FromPosts creates a slice of PostResponse with optimized data fetching.
func FromPosts(posts []*model.Post, authors map[string]*model.User, communities map[string]*model.Community, userVotes map[string]string, userPollVotes map[string][]string) []*PostResponse {
	var validPosts []*PostResponse
	for _, post := range posts {
		author := authors[post.AuthorID.Hex()]
		community := communities[post.CommunityID.Hex()]

		// Skip posts from banned communities
		if community == nil {
			continue
		}

		userVote := userVotes[post.ID.Hex()]
		userPollVote := userPollVotes[post.ID.Hex()]

		validPosts = append(validPosts, FromPost(post, author, community, userVote, userPollVote))
	}
	return validPosts
}

// FromPostsWithModeration creates a slice of PostResponse with moderation info included.
func FromPostsWithModeration(
	posts []*model.Post,
	authors map[string]*model.User,
	communities map[string]*model.Community,
	userVotes map[string]string,
	userPollVotes map[string][]string,
) []*PostResponse {
	responses := make([]*PostResponse, len(posts))

	for i, post := range posts {
		author := authors[post.AuthorID.Hex()]
		community := communities[post.CommunityID.Hex()]
		userVote := userVotes[post.ID.Hex()]
		userPollVote := userPollVotes[post.ID.Hex()]

		resp := FromPost(post, author, community, userVote, userPollVote)

		// Always include status
		status := string(post.ModerationStatus)
		resp.ModerationStatus = &status

		// ModerationResult may be nil
		if post.ModerationResult != nil {
			resp.ModerationReason = &post.ModerationResult.Reason
		} else {
			resp.ModerationReason = nil
		}

		// ModeratedAt may be nil
		resp.ModeratedAt = post.ModeratedAt

		responses[i] = resp
	}

	return responses
}
