package service

import (
	"fmt"
	"time"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommunityService interface {
	CreateCommunity(req *dto.CreateCommunityRequest, userID string) (*model.Community, error)
	GetCommunityByID(id string) (*model.Community, error)
	GetCommunitiesFilter(name string, description string, createFrom time.Time, page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error)
	GetCommunitiesByModeratorIDPaginated(moderatorID string, page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error)
	GetAllCommunitiesPaginated(page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error)
	UpdateCommunity(req *dto.UpdateCommunityRequest, userID string) (*model.Community, error)
	AddModerator(req *dto.AddModeratorRequest, userID string) error
	RemoveModerator(req *dto.RemoveModeratorRequest, userID string) error
	IsModerator(community *model.Community, userID string) (bool, error)
	DeleteCommunityByID(communityID string, userID string) error
}

type communityService struct {
	communityRepo repo.CommunityRepo
}

func (c *communityService) CreateCommunity(req *dto.CreateCommunityRequest, userID string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	community := &model.Community{
		Name:           req.Name,
		Description:    req.Description,
		Avatar:         req.Avatar,
		Banner:         req.Banner,
		Setting:        req.Setting,
		Moderators:     req.Moderators,
		CreateAt:       time.Now(),
		CreateByID:     userObjectID,
		CreateByName:   req.CreatorName,
		CreateByAvatar: req.CreatorAvatar,
		IsDeleted:      false,
		IsBanned:       false,
	}
	return c.communityRepo.Create(ctx, community)
}

func (c *communityService) GetCommunityByID(id string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return c.communityRepo.GetByID(ctx, id)
}

func (c *communityService) GetCommunitiesFilter(
	name string,
	description string,
	createFrom time.Time,
	page int,
	pageSize int,
) (*dto.PaginatedCommunitiesResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communities, total, err := c.communityRepo.GetFilter(ctx, name, description, createFrom, page, pageSize)
	communitiesResponses := dto.FromCommunities(communities)

	var response = &dto.PaginatedCommunitiesResponse{
		Communities: communitiesResponses,
		Pagination: dto.Pagination{
			Total: total,
			Page:  page,
		},
	}

	return response, err
}

func (c *communityService) GetCommunitiesByModeratorIDPaginated(moderatorID string, page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communities, total, err := c.communityRepo.GetByModeratorIDPaginated(ctx, moderatorID, page, pageSize)
	communitiesResponses := dto.FromCommunities(communities)

	var response = &dto.PaginatedCommunitiesResponse{
		Communities: communitiesResponses,
		Pagination: dto.Pagination{
			Total: total,
			Page:  page,
		},
	}

	return response, err
}

func (c *communityService) GetAllCommunitiesPaginated(page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communities, total, err := c.communityRepo.GetAllPaginated(ctx, page, pageSize)
	communitiesResponses := dto.FromCommunities(communities)

	var response = &dto.PaginatedCommunitiesResponse{
		Communities: communitiesResponses,
		Pagination: dto.Pagination{
			Total: total,
			Page:  page,
		},
	}

	return response, err
}

func (c *communityService) UpdateCommunity(req *dto.UpdateCommunityRequest, userID string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, req.CommunityID)
	if err != nil {
		return nil, err
	}

	ok, err := c.IsModerator(community, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("user is not a moderator of the community")
	}

	var updateCount = 0
	if req.Name != nil {
		community.Name = *req.Name
		updateCount++
	}
	if req.Description != nil {
		community.Description = req.Description
		updateCount++
	}
	if req.Avatar != nil {
		community.Avatar = req.Avatar
		updateCount++
	}
	if req.Banner != nil {
		community.Banner = req.Banner
		updateCount++
	}
	if req.Setting != nil {
		community.Setting = *req.Setting
		updateCount++
	}

	if updateCount == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	return community, c.communityRepo.Replace(ctx, community)
}

func (c *communityService) AddModerator(req *dto.AddModeratorRequest, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, req.CommunityID)
	if err != nil {
		return err
	}

	ok, err := c.IsModerator(community, userID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("user is not a moderator of the community")
	}

	var newModerators []model.Moderator
	for _, modDTO := range req.AddedModerator {
		ok, err := c.IsModerator(community, modDTO.ModeratorID)
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		objectID, err := primitive.ObjectIDFromHex(modDTO.ModeratorID)
		if err != nil {
			return fmt.Errorf("invalid new moderator id: %s", modDTO.ModeratorID)
		}
		newModerators = append(
			newModerators,
			model.Moderator{
				UserID:     objectID,
				Username:   modDTO.Username,
				AssignedAt: time.Now(),
			})
	}

	if len(newModerators) == 0 {
		return fmt.Errorf("no new modertor to add")
	}

	community.Moderators = append(community.Moderators, newModerators...)
	return c.communityRepo.Replace(ctx, community)
}

func (c *communityService) RemoveModerator(req *dto.RemoveModeratorRequest, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, req.CommunityID)
	if err != nil {
		return err
	}

	ok, err := c.IsModerator(community, userID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("user is not a moderator of the community")
	}

	for _, modDTO := range req.RemovedModerator {
		if userID == modDTO.ModeratorID {
			return fmt.Errorf("cannot remove yourself as a moderator")
		}

		for i, mod := range community.Moderators {
			if mod.UserID.Hex() == modDTO.ModeratorID {
				community.Moderators = append(community.Moderators[:i], community.Moderators[i+1:]...)
				break
			}
		}
	}

	return c.communityRepo.Replace(ctx, community)
}

func (c *communityService) DeleteCommunityByID(communityID string, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return err
	}

	ok, err := c.IsModerator(community, userID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("user is not a moderator of the community")
	}

	return c.communityRepo.Delete(ctx, communityID)
}

func (c *communityService) IsModerator(community *model.Community, userID string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user id: %s", userID)
	}

	for _, m := range community.Moderators {
		if m.UserID == objectID {
			return true, nil
		}
	}
	return false, nil
}

func NewCommunityService(communityRepo repo.CommunityRepo) CommunityService {
	return &communityService{communityRepo: communityRepo}
}
