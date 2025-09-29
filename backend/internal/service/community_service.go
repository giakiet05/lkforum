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
	GetAllCommunities() ([]*model.Community, error)
	GetCommunityByID(id string) (*model.Community, error)
	GetCommunitiesByModeratorID(moderatorID string) ([]*model.Community, error)
	UpdateCommunity(req *dto.UpdateCommunityRequest, userID string) (*model.Community, error)
	AddModerator(req *dto.AddModeratorRequest, userID string) error
	RemoveModerator(req *dto.RemoveModeratorRequest, userID string) error
	IsModerator(userID string, community *model.Community) (bool, error)
	DeleteCommunityByID(communityID string) error
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
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
		Banner:      req.Banner,
		Setting:     req.Setting,
		Moderators:  req.Moderators,
		CreateAt:    time.Now(),
		CreateBy:    userObjectID,
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

func (c *communityService) GetCommunityByID(id string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communityObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return c.communityRepo.GetById(ctx, communityObjectID)
}

func (c *communityService) GetCommunitiesByModeratorID(moderatorID string) ([]*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	modObjectID, err := primitive.ObjectIDFromHex(moderatorID)
	if err != nil {
		return nil, err
	}

	return c.communityRepo.GetByModeratorId(ctx, modObjectID)
}

func (c *communityService) UpdateCommunity(req *dto.UpdateCommunityRequest, userID string) (*model.Community, error) {
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

func (c *communityService) AddModerator(req *dto.AddModeratorRequest, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communityID, err := primitive.ObjectIDFromHex(req.CommunityID)
	if err != nil {
		return err
	}
	community, err := c.communityRepo.GetById(ctx, communityID)
	if err != nil {
		return err
	}

	ok, err := c.IsModerator(userID, community)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("user is not a moderator of the community")
	}

	var newModerators []model.Moderator
	for _, id := range req.AddedModeratorID {
		ok, err := c.IsModerator(id, community)
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid new moderator id: %s", id)
		}
		newModerators = append(newModerators, model.Moderator{UserID: objectID, AssignedAt: time.Now()})
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

	communityID, err := primitive.ObjectIDFromHex(req.CommunityID)
	if err != nil {
		return err
	}
	community, err := c.communityRepo.GetById(ctx, communityID)
	if err != nil {
		return err
	}

	ok, err := c.IsModerator(userID, community)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("user is not a moderator of the community")
	}

	for _, id := range req.RemovedModeratorID {
		if userID == id {
			return fmt.Errorf("cannot remove yourself as a moderator")
		}

		for i, mod := range community.Moderators {
			if mod.UserID.Hex() == id {
				community.Moderators = append(community.Moderators[:i], community.Moderators[i+1:]...)
				break
			}
		}
	}

	return c.communityRepo.Replace(ctx, community)
}

func (c *communityService) DeleteCommunityByID(communityID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return err
	}

	return c.communityRepo.Delete(ctx, communityObjectID)
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
