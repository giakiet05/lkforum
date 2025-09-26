package service

import (
	"fmt"
	"time"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommunityService interface {
	CreateCommunity(req *dto.CreateCommunityRequest) (*model.Community, error)
	GetAllCommunities() ([]*model.Community, error)
	GetCommunityById(id string) (*model.Community, error)
	GetCommunitiesByModeratorId(moderatorId string) ([]*model.Community, error)
	UpdateOneCommunity(req *dto.UpdateCommunityRequest) (*model.Community, error)
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

func (c *communityService) GetCommunitiesByModeratorId(moderatorId string) ([]*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(moderatorId)
	if err != nil {
		return nil, err
	}

	return c.communityRepo.GetByModeratorId(ctx, objectID)
}

func (c *communityService) UpdateOneCommunity(req *dto.UpdateCommunityRequest) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communityID, err := primitive.ObjectIDFromHex(req.CommunityID)
	if err != nil {
		return nil, err
	}

	updates := bson.M{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Banner != nil {
		updates["banner"] = *req.Banner
	}
	if req.Setting != nil {
		updates["setting"] = *req.Setting
	}
	if req.Moderators != nil {
		updates["moderators"] = *req.Moderators
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	return c.communityRepo.UpdateOne(ctx, communityID, updates)
}

func NewCommunityService(communityRepo repo.CommunityRepo) CommunityService {
	return &communityService{communityRepo: communityRepo}
}
