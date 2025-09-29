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
	CreateCommunity(req *dto.CreateCommunityRequest) (*model.Community, error)
	GetAllCommunities() ([]*model.Community, error)
	GetCommunityById(id string) (*model.Community, error)
	GetCommunitiesByModeratorId(moderatorID string) ([]*model.Community, error)
	UpdateOneCommunity(req *dto.UpdateCommunityRequest, userID string) (*model.Community, error)
	IsModerator(userID string, community *model.Community) (bool, error)
}

type communityService struct {
	communityRepo repo.CommunityRepo
}

/////////////////// Thiếu bảo mật moderator ////////////////

func (c *communityService) CreateCommunity(req *dto.CreateCommunityRequest) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community := &model.Community{
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
		Banner:      req.Banner,
		Setting:     req.Setting,
		Moderators:  req.Moderators,
		CreateAt:    time.Now(),
		IsDeleted:   false,
		IsBanned:    false,
	}
	return c.communityRepo.Create(ctx, community)
}

func (c *communityService) GetAllCommunities() ([]*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return c.communityRepo.GetAll(ctx)
}

func (c *communityService) GetCommunityById(id string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return c.communityRepo.GetById(ctx, objectID)
}

func (c *communityService) GetCommunitiesByModeratorId(moderatorID string) ([]*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(moderatorID)
	if err != nil {
		return nil, err
	}

	return c.communityRepo.GetByModeratorId(ctx, objectID)
}

func (c *communityService) UpdateOneCommunity(req *dto.UpdateCommunityRequest, userID string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communityID, err := primitive.ObjectIDFromHex(req.CommunityID)
	if err != nil {
		return nil, err
	}

	community, err := c.communityRepo.GetById(ctx, communityID)
	if err != nil {
		return nil, err
	}

	ok, err := c.IsModerator(userID, community)
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

func (c *communityService) IsModerator(userID string, community *model.Community) (bool, error) {
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
