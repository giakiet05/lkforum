package dto

import "github.com/giakiet05/lkforum/internal/model"

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type PaginatedUsersResponse struct {
	Users      []*UserResponse `json:"users"`
	Pagination Pagination      `json:"pagination"`
}

type PaginatedCommunitiesResponse struct {
	Communities []*CommunityResponse `json:"communities"`
	Pagination  Pagination           `json:"pagination"`
}

type PaginatedMembershipsResponse struct {
	Memberships []*model.Membership `json:"memberships"`
	Pagination  Pagination          `json:"pagination"`
}

type PaginatedPostsResponse struct {
	Posts      []*PostResponse `json:"posts"`
	Pagination Pagination      `json:"pagination"`
}

type PaginatedCommentsResponse struct {
	Comments   []*CommentResponse `json:"comments"`
	Pagination Pagination         `json:"pagination"`
}

type PaginatedChannelsResponse struct {
	Channels   []*ChannelResponse `json:"channels"`
	Pagination Pagination         `json:"pagination"`
}

type PaginatedMessagesResponse struct {
	Messages   []*MessageResponse `json:"messages"`
	Pagination Pagination         `json:"pagination"`
}

type PaginatedPostHistoryResponse struct {
	PostHistories []*PostHistoryResponse `json:"post_histories"`
	Pagination    Pagination             `json:"pagination"`
}

type PaginatedReportResponse struct {
	Reports    []*ReportResponse `json:"reports"`
	Pagination Pagination        `json:"pagination"`
}
