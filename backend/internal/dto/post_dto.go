package dto

import "time"

// Request DTO

type CreatePostRequest struct {
	CommunityID string `json:"community_id" validate:"required"`
	Title       string `json:"title" validate:"required,min=3,max=300"`
	Type        string `json:"type" validate:"required,oneof=text image video poll"`
	//
	Text   string               `json:"text,omitempty"`
	Images []ImageUploadRequest `json:"images,omitempty"`
	Video  *VideoUploadRequest  `json:"video,omitempty"`
	Poll   *CreatePollRequest   `json:"poll,omitempty"`
}
type ImageUploadRequest struct {
	URL string `json:"url" validate:"required"`
}

type VideoUploadRequest struct {
	URL       string `json:"url" validate:"required"`
	Thumbnail string `json:"thumbnail,omitempty"`
}
type CreatePollRequest struct {
	Question      string     `json:"question" validate:"required,min=10,max=500"`
	Options       []string   `json:"options" validate:"required,min=2,dive,min=1,max=200"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	AllowMultiple bool       `json:"allow_multiple,omitempty"`
}

type UpdatePostRequest struct {
	Title string `json:"title" validate:"required,min=3,max=300"`
	Text  string `json:"text" validate:"required,min=3,max=300"`
}
type AddImageRequest struct {
	Images []ImageUploadRequest `json:"images" validate:"required,min=1,dive,min=1"`
}

type RemoveImageRequest struct {
	ImageIDs []string `json:"image_ids" validate:"required"`
}
type UpdatePollRequest struct {
	Question      string     `json:"question" validate:"required,min=10,max=500"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	AllowMultiple bool       `json:"allow_multiple,omitempty"`
}
type AddPollOptionRequest struct {
	Options []string `json:"options" validate:"required"`
}
type UpdatePollOptionRequest struct {
	Text string `json:"text" validate:"required,min=3,max=300"`
}
type RemovePollOptionRequest struct {
	OptionIDs []string `json:"option_ids" validate:"required"`
}

// Response DTO
type PostResponse struct {
	ID             string               `json:"id"`
	AuthorID       string               `json:"author_id"`
	AuthorUsername string               `json:"author_username"`
	AuthorAvatar   string               `json:"author_avatar"`
	CommunityID    string               `json:"community_id"`
	CommunityName  string               `json:"community_name"`
	Title          string               `json:"title"`
	Type           string               `json:"type"`
	Content        *PostContentResponse `json:"content,omitempty"`
	VotesCount     *VotesCountResponse  `json:"votes_count"`
	//Vote của Uer
	UserVote      string     `json:"user_vote,omitempty"`
	CommentsCount int        `json:"comments_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type PostContentResponse struct {
	Text   string          `json:"text"`
	Images []ImageResponse `json:"images,omitempty"`
	Poll   *PollResponse   `json:"poll,omitempty"`
	Video  *VideoResponse  `json:"video,omitempty"`
}
type ImageResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
type PollResponse struct {
	Question      string               `json:"question"`
	Options       []PollOptionResponse `json:"options"`
	TotalVotes    int                  `json:"total_votes"`
	UserVoteIDs   []string             `json:"user_vote_ids"`
	ExpiresAt     *time.Time           `json:"expires_at,omitempty"`
	AllowMultiple bool                 `json:"allow_multiple,omitempty"`
}
type PollOptionResponse struct {
	ID         string  `json:"id"`
	Text       string  `json:"text"`
	Votes      int     `json:"votes"`
	Percentage float64 `json:"percentage"`
}
type VotesCountResponse struct {
	Up    int `json:"up"`
	Down  int `json:"down"`
	Score int `json:"score"`
}
type VideoResponse struct {
	ID        string `json:"id"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
}

type GetPostsQuery struct {
	CommunityID string `query:"community_id"`
	AuthorID    string `query:"author_id"`
	Type        string `query:"type" validate:"oneof= text image video poll"`
	Sort        string `query:"sort" validate:"oneof=hot new top controversial"`
	TimeFrame   string `query:"time" validate:"oneof= hour day week month year all"`
	Page        int    `query:"page" validate:"min=1"`
	Limit       int    `query:"limit" validate:"min=1,max=50"`
}
