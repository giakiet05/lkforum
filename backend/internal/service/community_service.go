package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	model "github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommunityService interface {
	Start()
	CreateCommunity(req *dto.CreateCommunityRequest, userID string) (*model.Community, error)
	GetCommunityByID(id string) (*model.Community, error)
	GetCommunitiesFilter(name string, description string, is18Plus bool, createFrom time.Time, page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error)
	GetCommunitiesByModeratorIDPaginated(moderatorID string, page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error)
	GetAllCommunitiesPaginated(page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error)
	UpdateCommunity(req *dto.UpdateCommunityRequest, userID string) (*model.Community, error)
	AddModerator(req *dto.AddModeratorRequest, userID string) error
	RemoveModerator(req *dto.RemoveModeratorRequest, userID string) error
	IsModerator(community *model.Community, userID string) (bool, error)
	DeleteCommunityByID(communityID string, userID string) error

	GetBannedUsers(communityID string, banTypeStr string, requesterID string) ([]*model.User, error)
	BanUser(req *dto.BanUserRequest, requesterID string) error
	UnmuteUser(userID string, communityID string, requesterID string) error
	UnbanUser(userID string, communityID string, requesterID string) error
}

type communityService struct {
	communityRepo  repo.CommunityRepo
	membershipRepo repo.MembershipRepo
	eventBus       bus.EventBus
}

func NewCommunityService(communityRepo repo.CommunityRepo, membershipRepo repo.MembershipRepo, bus bus.EventBus) CommunityService {
	return &communityService{communityRepo: communityRepo, membershipRepo: membershipRepo, eventBus: bus}
}

func (c *communityService) Start() {
	eventChannel := make(bus.EventListener, 100)

	c.eventBus.Subscribe(bus.TopicUserChangeAvatar, eventChannel)

	log.Println("ChannelService started and subscribed to events.")

	go c.processEvents(eventChannel)
}

func (c *communityService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicUserChangeAvatar:
			c.handleNewAvatar(event)
		default:
			log.Println("Unhandled event topic:", event.Topic())
		}
	}
}

func (c *communityService) handleNewAvatar(event bus.Event) {
	payload := event.Payload()
	userID, _ := payload["user_id"].(string)
	newAvatar, _ := payload["new_avatar"].(string)

	if userID == "" || newAvatar == "" {
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	err := c.communityRepo.UpdateUserAvatar(ctx, userID, newAvatar)
	if err != nil {
		log.Printf("Failed to update avatar: %v", err)
	}
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
		Rules:          req.Rules,
		Moderators:     req.Moderators,
		CreateAt:       time.Now(),
		CreateByID:     userObjectID,
		CreateByName:   req.CreatorName,
		CreateByAvatar: req.CreatorAvatar,
		IsDeleted:      false,
		IsBanned:       false,
	}
	community, err = c.communityRepo.Create(ctx, community)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, apperror.ErrCommunityNameExists
		}
		return nil, err
	}
	membership := &model.Membership{
		UserID:      userObjectID,
		CommunityID: community.ID,
	}
	_, err = c.membershipRepo.Create(ctx, membership)
	if err != nil {
		return nil, err
	}

	return community, nil
}

func (c *communityService) GetCommunityByID(id string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrCommunityNotFound
		}
		return nil, err
	}

	return community, nil
}

func (c *communityService) GetCommunitiesFilter(
	name string,
	description string,
	is18Plus bool,
	createFrom time.Time,
	page int,
	pageSize int,
) (*dto.PaginatedCommunitiesResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communities, total, err := c.communityRepo.GetFilter(ctx, name, description, is18Plus, createFrom, page, pageSize)
	if err != nil {
		return nil, err
	}
	communitiesResponses := dto.FromCommunities(communities)

	var response = &dto.PaginatedCommunitiesResponse{
		Communities: communitiesResponses,
		Pagination: dto.Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	}

	return response, nil
}

func (c *communityService) GetCommunitiesByModeratorIDPaginated(moderatorID string, page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communities, total, err := c.communityRepo.GetByModeratorIDPaginated(ctx, moderatorID, page, pageSize)
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return nil, err
	}
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
		return nil, apperror.ErrForbidden
	}

	var updateCount = 0
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
	if req.Rules != nil {
		community.Rules = *req.Rules
		updateCount++
	}

	if updateCount == 0 {
		return nil, apperror.ErrNoFieldsToUpdate
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
		return apperror.ErrForbidden
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
			return apperror.ErrInvalidID
		}

		existed, err := c.communityRepo.IsUserExist(ctx, modDTO.ModeratorID)
		if err != nil {
			return err
		}
		if !existed {
			return apperror.ErrInvalidID
		}

		ok, err = c.membershipRepo.IsMember(ctx, modDTO.ModeratorID, userID)
		if err != nil {
			return err
		}
		if !ok {
			return apperror.ErrUserNotMember
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
		return apperror.ErrNoFieldsToUpdate
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
		return apperror.ErrForbidden
	}

	for _, modID := range req.RemovedModerator {
		if userID == modID {
			return apperror.ErrCannotRemoveModerator
		}

		for i, mod := range community.Moderators {
			if mod.UserID.Hex() == modID {
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

func (c *communityService) GetBannedUsers(communityID string, banTypeStr string, requesterID string) ([]*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return nil, err
	}

	ok, err := c.IsModerator(community, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	if banTypeStr == "" {
		return nil, apperror.ErrBadRequest
	}

	banType := model.CommunityBanType(banTypeStr)
	if banType == model.Banned {
		return c.communityRepo.GetBannedUsers(ctx, communityID)
	}

	if banType == model.Muted {
		return c.communityRepo.GetBannedUsers(ctx, communityID)
	}

	return nil, apperror.ErrBadRequest
}

func (c *communityService) BanUser(req *dto.BanUserRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, req.CommunityID)
	if err != nil {
		return err
	}

	ok, err := c.IsModerator(community, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	communityObjectID, err := primitive.ObjectIDFromHex(req.CommunityID)
	if err != nil {
		return apperror.ErrInternal
	}

	userObjectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return apperror.ErrInternal
	}

	requesterObjectID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return apperror.ErrInternal
	}

	// Validate ban type
	banType := model.CommunityBanType(req.Type)
	if banType != model.Banned && banType != model.Muted {
		return apperror.ErrBadRequest
	}

	// Calculate ban expiration
	expiresAt := time.Now().Add(time.Hour * 24 * time.Duration(req.LengthDays))

	// Create ban record
	ban := model.CommunityBan{
		CommunityID: communityObjectID,
		UserID:      userObjectID,
		Type:        banType,
		Reason:      req.Reason,
		BannedAt:    time.Now(),
		BannedBy:    requesterObjectID,
		ExpiresAt:   expiresAt,
	}

	// Call repo method to save ban
	err = c.communityRepo.BanUser(ctx, &ban)

	if err != nil {
		return apperror.ErrInternal
	}

	return nil
}

func (c *communityService) UnmuteUser(userID string, communityID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return err
	}

	ok, err := c.IsModerator(community, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	return c.communityRepo.UnmuteUser(ctx, userID, communityID)
}

func (c *communityService) UnbanUser(userID string, communityID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return err
	}

	ok, err := c.IsModerator(community, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	return c.communityRepo.UnbanUser(ctx, userID, communityID)
}
