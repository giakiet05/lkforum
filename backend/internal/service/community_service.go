package service

import (
	"context"
	"errors"
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
	CreateCommunity(req *dto.CreateCommunityRequest, requesterID string) (*model.Community, error)
	GetCommunityByID(communityID string, requesterID *string) (*model.Community, error)
	GetCommunitiesFilter(
		requesterID *string,
		name string,
		description string,
		is18Plus bool,
		createFrom time.Time,
		page int, pageSize int,
	) (*dto.PaginatedCommunitiesResponse, error)
	GetCommunitiesByModeratorIDPaginated(moderatorID string, page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error)
	GetAllCommunitiesPaginated(requesterID *string, page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error)
	UpdateCommunity(req *dto.UpdateCommunityRequest, requesterID string) (*model.Community, error)
	AddModerator(req *dto.AddModeratorRequest, requesterID string) error
	RemoveModerator(req *dto.RemoveModeratorRequest, requesterID string) error
	DeleteCommunityByID(communityID string, requesterID string) error

	BanPost(req *dto.CommunityBanPostRequest, requesterID string) error
	UnbanPost(req *dto.CommunityUnbanPostRequest, requesterID string) error

	GetBannedUsers(communityID string, banTypeStr string, expired bool, requesterID string) ([]*model.User, error)
	BanUser(req *dto.CommunityBanUserRequest, requesterID string) error
	UnmuteUser(userID string, communityID string, requesterID string) error
	UnbanUser(userID string, communityID string, requesterID string) error

	GetPendingPosts(communityID string, moderatorID string, page int, pageSize int) (*dto.PaginatedPostsResponse, error)
	ModeratePost(communityID string, postID string, moderatorID string, approve bool, reason *string) error
}

type communityService struct {
	communityRepo  repo.CommunityRepo
	membershipRepo repo.MembershipRepo
	postRepo       repo.PostRepo
	userRepo       repo.UserRepo
	eventBus       bus.EventBus
}

func NewCommunityService(
	communityRepo repo.CommunityRepo,
	membershipRepo repo.MembershipRepo,
	postRepo repo.PostRepo,
	userRepo repo.UserRepo,
	bus bus.EventBus,
) CommunityService {
	return &communityService{
		communityRepo:  communityRepo,
		membershipRepo: membershipRepo,
		postRepo:       postRepo,
		userRepo:       userRepo,
		eventBus:       bus,
	}
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

func (c *communityService) CreateCommunity(req *dto.CreateCommunityRequest, requesterID string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(requesterID)
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

func (c *communityService) GetCommunityByID(communityID string, requesterID *string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if requesterID != nil {
		ok, err := c.communityRepo.IsUserBanned(ctx, *requesterID, model.Banned, communityID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, apperror.ErrUserIsBannedFromCommunity
		}
	}

	community, err := c.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrCommunityNotFound
		}
		return nil, err
	}

	return community, nil
}

func (c *communityService) GetCommunitiesFilter(
	requesterID *string,
	name string,
	description string,
	is18Plus bool,
	createFrom time.Time,
	page int, pageSize int,
) (*dto.PaginatedCommunitiesResponse, error) {

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communities, total, err := c.communityRepo.GetFilter(ctx, name, description, is18Plus, createFrom, page, pageSize)
	if err != nil {
		return nil, err
	}

	if requesterID != nil && len(communities) > 0 {
		communities, err = c.filterOutBannedCommunities(ctx, *requesterID, communities)
		if err != nil {
			return nil, err
		}
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

func (c *communityService) GetAllCommunitiesPaginated(requesterID *string, page int, pageSize int) (*dto.PaginatedCommunitiesResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	communities, total, err := c.communityRepo.GetAllPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	if requesterID != nil && len(communities) > 0 {
		communities, err = c.filterOutBannedCommunities(ctx, *requesterID, communities)
		if err != nil {
			return nil, err
		}
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

func (c *communityService) UpdateCommunity(req *dto.UpdateCommunityRequest, requesterID string) (*model.Community, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, req.CommunityID)
	if err != nil {
		return nil, err
	}

	ok, err := c.communityRepo.IsModerator(ctx, req.CommunityID, requesterID)
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

func (c *communityService) AddModerator(req *dto.AddModeratorRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, req.CommunityID)
	if err != nil {
		return err
	}

	ok, err := c.communityRepo.IsModerator(ctx, req.CommunityID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	var newModerators []model.Moderator
	for _, modDTO := range req.AddedModerator {
		ok, err := c.communityRepo.IsModerator(ctx, req.CommunityID, modDTO.ModeratorID)
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

		ok, err = c.membershipRepo.IsMember(ctx, modDTO.ModeratorID, requesterID)
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

func (c *communityService) RemoveModerator(req *dto.RemoveModeratorRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := c.communityRepo.GetByID(ctx, req.CommunityID)
	if err != nil {
		return err
	}

	ok, err := c.communityRepo.IsModerator(ctx, req.CommunityID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	for _, modID := range req.RemovedModerator {
		if requesterID == modID {
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

func (c *communityService) DeleteCommunityByID(communityID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := c.communityRepo.IsModerator(ctx, communityID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	return c.communityRepo.Delete(ctx, communityID)
}

func (c *communityService) BanPost(req *dto.CommunityBanPostRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := c.communityRepo.IsModerator(ctx, req.CommunityID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	return c.postRepo.BanPost(ctx, req.PostID, req.Reason)
}

func (c *communityService) UnbanPost(req *dto.CommunityUnbanPostRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := c.communityRepo.IsModerator(ctx, req.CommunityID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	return c.postRepo.UnbanPost(ctx, req.PostID)
}

func (c *communityService) GetBannedUsers(communityID string, banTypeStr string, expired bool, requesterID string) ([]*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := c.communityRepo.IsModerator(ctx, communityID, requesterID)
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
		return c.communityRepo.GetBannedUsers(ctx, communityID, expired)
	}

	if banType == model.Muted {
		return c.communityRepo.GetBannedUsers(ctx, communityID, expired)
	}

	return nil, apperror.ErrBadRequest
}

func (c *communityService) BanUser(req *dto.CommunityBanUserRequest, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := c.communityRepo.IsModerator(ctx, req.CommunityID, requesterID)
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

	ok, err := c.communityRepo.IsModerator(ctx, communityID, requesterID)
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

	ok, err := c.communityRepo.IsModerator(ctx, communityID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	return c.communityRepo.UnbanUser(ctx, userID, communityID)
}

// FilterOutBannedCommunities removes communities from the list where the user is banned.
func (c *communityService) filterOutBannedCommunities(
	ctx context.Context,
	requesterID string,
	communities []*model.Community,
) ([]*model.Community, error) {

	if len(communities) == 0 {
		return communities, nil
	}

	// 1. Collect community IDs
	communityIDs := make([]string, 0, len(communities))
	for _, cmt := range communities {
		communityIDs = append(communityIDs, cmt.ID.Hex())
	}

	// 2. Get banned community IDs
	bannedIDs, err := c.communityRepo.GetBannedCommunityIDs(ctx, requesterID, model.Banned, communityIDs)
	if err != nil {
		return nil, err
	}

	// 3. Build banned set for fast lookup
	bannedSet := make(map[string]struct{}, len(bannedIDs))
	for _, id := range bannedIDs {
		bannedSet[id] = struct{}{}
	}

	// 4. Filter out banned communities
	filtered := make([]*model.Community, 0, len(communities))
	for _, cmt := range communities {
		if _, banned := bannedSet[cmt.ID.Hex()]; !banned {
			filtered = append(filtered, cmt)
		}
	}

	return filtered, nil
}

func (c *communityService) GetPendingPosts(communityID string, moderatorID string, page int, pageSize int) (*dto.PaginatedPostsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get community
	community, err := c.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrCommunityNotFound
		}
		return nil, err
	}

	// Check if requester is moderator
	isMod, err := c.communityRepo.IsModerator(ctx, communityID, moderatorID)
	if err != nil {
		return nil, err
	}
	if !isMod {
		return nil, apperror.ErrForbidden
	}

	// Build filter for pending posts in this community
	communityObjID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	filter := repo.Filter{
		"community_id":      communityObjID,
		"moderation_status": model.ModerationPending,
		"is_deleted":        false,
	}

	// Build find options
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	findOptions := &repo.FindOptions{
		Skip:  int64((page - 1) * pageSize),
		Limit: int64(pageSize),
		Sort:  map[string]int{"created_at": -1},
	}

	// Get pending posts
	posts, total, err := c.postRepo.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &dto.PaginatedPostsResponse{
			Posts: []*dto.PostResponse{},
			Pagination: dto.Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    0,
			},
		}, nil
	}

	// Get authors and communities info
	authorIDs := make([]string, 0, len(posts))
	authorIDMap := make(map[string]bool)
	for _, post := range posts {
		if !authorIDMap[post.AuthorID.Hex()] {
			authorIDs = append(authorIDs, post.AuthorID.Hex())
			authorIDMap[post.AuthorID.Hex()] = true
		}
	}

	authors, _ := c.userRepo.GetByIDs(ctx, authorIDs)
	authorsMap := make(map[string]*model.User)
	for _, author := range authors {
		authorsMap[author.ID.Hex()] = author
	}

	// Community map (only one community)
	communitiesMap := map[string]*model.Community{
		communityID: community,
	}

	// Convert to response (with moderation info)
	responses := dto.FromPostsWithModeration(posts, authorsMap, communitiesMap, nil, nil)

	return &dto.PaginatedPostsResponse{
		Posts: responses,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (c *communityService) ModeratePost(communityID string, postID string, moderatorID string, approve bool, reason *string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is moderator
	isMod, err := c.communityRepo.IsModerator(ctx, communityID, moderatorID)
	if err != nil {
		return err
	}
	if !isMod {
		return apperror.ErrForbidden
	}

	// Get post
	post, err := c.postRepo.GetByID(ctx, postID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrPostNotFound
		}
		return err
	}

	// Check if post belongs to this community
	if post.CommunityID.Hex() != communityID {
		return apperror.NewError(nil, apperror.ErrBadRequest.Code, "post does not belong to this community")
	}

	// Check if post is pending moderation
	if post.ModerationStatus != model.ModerationPending {
		return apperror.NewError(nil, apperror.ErrBadRequest.Code, "post is not pending moderation")
	}

	now := time.Now()
	var newStatus model.ModerationStatus

	if approve {
		// Approve post
		newStatus = model.ModerationApproved
		update := repo.UpdateDocument{
			"$set": map[string]interface{}{
				"moderation_status": newStatus,
				"moderated_at":      &now,
			},
		}

		err = c.postRepo.UpdateByID(ctx, postID, update)
		if err != nil {
			return err
		}

		// Publish approved event
		c.eventBus.Publish(&bus.PostApprovedEvent{
			PostID:   postID,
			AuthorID: post.AuthorID.Hex(),
		})

		log.Printf("Post %s approved by moderator %s in community %s", postID, moderatorID, communityID)
	} else {
		// Reject post
		newStatus = model.ModerationRejected
		moderationResult := &model.ModerationResult{
			IsViolation: true,
			Reason:      "Rejected by moderator",
		}
		if reason != nil {
			moderationResult.Reason = *reason
		}

		update := repo.UpdateDocument{
			"$set": map[string]interface{}{
				"moderation_status": newStatus,
				"moderation_result": moderationResult,
				"moderated_at":      &now,
			},
		}

		err = c.postRepo.UpdateByID(ctx, postID, update)
		if err != nil {
			return err
		}

		// Publish rejected event
		c.eventBus.Publish(&bus.PostRejectedEvent{
			PostID:   postID,
			AuthorID: post.AuthorID.Hex(),
			Reason:   moderationResult.Reason,
		})

		log.Printf("Post %s rejected by moderator %s in community %s: %s", postID, moderatorID, communityID, moderationResult.Reason)
	}

	return nil
}
