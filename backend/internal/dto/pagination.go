package dto

import "github.com/giakiet05/lkforum/internal/model"

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type PaginatedUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination Pagination     `json:"pagination"`
}

type PaginatedCommunitiesResponse struct {
	Communities []CommunityResponse `json:"communities"`
	Pagination  Pagination          `json:"pagination"`
}

type PaginatedMembershipsResponse struct {
	Memberships []model.Membership `json:"memberships"`
	Pagination  Pagination         `json:"pagination"`
}
