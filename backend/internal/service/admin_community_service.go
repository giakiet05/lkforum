package service

import (
	"errors"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminCommunityService interface {
	// Community management
	GetCommunitiesAdmin(query *dto.GetCommunitiesAdminQuery) (*dto.PaginatedCommunitiesResponse, error)
	BanCommunity(communityID string, req *dto.AdminBanCommunityRequest) error
	UnbanCommunity(communityID string) error
}

type adminCommunityService struct {
	communityRepo repo.CommunityRepo
}

func NewAdminCommunityService(communityRepo repo.CommunityRepo) AdminCommunityService {
	return &adminCommunityService{
		communityRepo: communityRepo,
	}
}

func (s *adminCommunityService) GetCommunitiesAdmin(query *dto.GetCommunitiesAdminQuery) (*dto.PaginatedCommunitiesResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Build filter based on status
	filter := repo.Filter{}

	switch query.Status {
	case "active":
		filter["is_banned"] = false
		filter["is_deleted"] = false
	case "banned":
		filter["is_banned"] = true
	case "deleted":
		filter["is_deleted"] = true
	case "all":
		// No filter - get all communities
	default:
		// Default: active communities only
		filter["is_banned"] = false
		filter["is_deleted"] = false
	}

	// Add name search if provided
	if query.Name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: query.Name, Options: "i"}}
	}

	// Pagination
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	findOptions := &repo.FindOptions{
		Skip:  int64((page - 1) * pageSize),
		Limit: int64(pageSize),
		Sort:  map[string]int{"create_at": -1},
	}

	// Use repo Find method to query communities
	communities, total, err := s.communityRepo.FindCommunities(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	communityResponses := dto.FromCommunities(communities)

	return &dto.PaginatedCommunitiesResponse{
		Communities: communityResponses,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (s *adminCommunityService) BanCommunity(communityID string, req *dto.AdminBanCommunityRequest) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get community
	community, err := s.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrCommunityNotFound
		}
		return err
	}

	// Update ban fields
	community.IsBanned = true
	community.BanReason = &req.Reason

	err = s.communityRepo.Replace(ctx, community)
	return err
}

func (s *adminCommunityService) UnbanCommunity(communityID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get community
	community, err := s.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrCommunityNotFound
		}
		return err
	}

	// Unban community
	community.IsBanned = false
	community.BanReason = nil

	err = s.communityRepo.Replace(ctx, community)
	return err
}
