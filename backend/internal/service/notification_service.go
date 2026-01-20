package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService interface {
	Start()
	GetNotifications(recipientID string, page, pageSize int) (*dto.PaginatedNotificationsResponse, error)
	MarkAllAsRead(recipientID string) (int64, error)
	MarkAsRead(notificationID, recipientID string) error
	DeleteNotification(notificationID, recipientID string) error
}

type notificationService struct {
	notificationRepo repo.NotificationRepo
	userRepo         repo.UserRepo
	postRepo         repo.PostRepo
	commentRepo      repo.CommentRepo
	communityRepo    repo.CommunityRepo
	eventBus         bus.EventBus
	redisClient      *redis.Client
}

func NewNotificationService(
	notificationRepo repo.NotificationRepo,
	userRepo repo.UserRepo,
	postRepo repo.PostRepo,
	commentRepo repo.CommentRepo,
	communityRepo repo.CommunityRepo,
	bus bus.EventBus,
	redis *redis.Client,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		postRepo:         postRepo,
		commentRepo:      commentRepo,
		communityRepo:    communityRepo,
		eventBus:         bus,
		redisClient:      redis,
	}
}

func (s *notificationService) Start() {
	eventChannel := make(bus.EventListener, 100)

	s.eventBus.Subscribe(bus.TopicPostUpvoted, eventChannel)
	s.eventBus.Subscribe(bus.TopicCommentCreated, eventChannel)
	s.eventBus.Subscribe(bus.TopicCommentApproved, eventChannel)
	s.eventBus.Subscribe(bus.TopicCommentUpvoted, eventChannel)
	s.eventBus.Subscribe(bus.TopicBroadcast, eventChannel)
	s.eventBus.Subscribe(bus.TopicModeratorAdded, eventChannel)

	log.Println("NotificationService started and subscribed to events.")

	go s.processEvents(eventChannel)
	go s.processBatchedNotifications()
}

func (s *notificationService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicPostUpvoted:
			s.handlePostUpvoted(event)
		case bus.TopicCommentCreated:
			s.handleCommentCreatedReply(event)
		case bus.TopicCommentApproved:
			s.handleCommentApprovedForPost(event)
		case bus.TopicCommentUpvoted:
			s.handleCommentUpvoted(event)
		case bus.TopicModeratorAdded:
			s.handleModeratorsAdded(event)
		case bus.TopicBroadcast:
			s.handleBroadcast(event)
		}
	}
}

func (s *notificationService) handlePostUpvoted(event bus.Event) {
	payload := event.Payload()
	authorID, _ := payload["author_id"].(string)
	voterID, _ := payload["voter_id"].(string)
	postID, _ := payload["post_id"].(string)

	if authorID == "" || voterID == "" || postID == "" {
		return
	}
	if authorID == voterID {
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if recipient wants upvote notifications
	if !s.shouldNotify(ctx, authorID, model.NotificationTypeLike) {
		return
	}

	// Redis batch key
	batchKey := fmt.Sprintf(postUpvoteBatchKeyPrefix, postID)

	// Check if this is the first upvote
	isFirst, err := s.isFirstInBatch(ctx, batchKey, voterID)
	if err != nil {
		log.Printf("ERROR: Failed to check first upvote: %v", err)
		return
	}

	// Add to batch
	if err := s.addToBatch(ctx, batchKey, voterID, postUpvoteBatchTTL); err != nil {
		log.Printf("ERROR: Failed to add to batch: %v", err)
		return
	}

	// If first upvote, send instant notification
	if isFirst {
		voter, err := s.userRepo.GetByID(ctx, voterID)
		if err != nil {
			return
		}

		post, err := s.postRepo.GetByID(ctx, postID)
		if err != nil {
			return
		}

		recipientObjID, _ := primitive.ObjectIDFromHex(authorID)
		actorObjID, _ := primitive.ObjectIDFromHex(voterID)

		notification := &model.Notification{
			RecipientID: recipientObjID,
			ActorID:     actorObjID,
			Type:        model.NotificationTypeLike,
			Message:     fmt.Sprintf("%s đã thích bài viết của bạn: %s", voter.Username, post.Title),
			Link:        fmt.Sprintf("/posts/%s", postID),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  authorID,
			Notification: dto.FromNotification(createdNotification),
		})
		return
	}

	// Check if reached threshold
	count, err := s.getBatchCount(ctx, batchKey)
	if err != nil {
		return
	}

	if count >= postUpvoteThreshold {
		s.sendPostUpvoteBatchNotification(ctx, batchKey)
	}
}

func (s *notificationService) handleCommentUpvoted(event bus.Event) {
	payload := event.Payload()
	authorID, _ := payload["author_id"].(string)
	voterID, _ := payload["voter_id"].(string)
	commentID, _ := payload["comment_id"].(string)

	if authorID == "" || voterID == "" || commentID == "" {
		return
	}
	if authorID == voterID {
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if recipient wants upvote notifications
	if !s.shouldNotify(ctx, authorID, model.NotificationTypeLike) {
		return
	}

	// Redis batch key
	batchKey := fmt.Sprintf(commentUpvoteBatchKeyPrefix, commentID)

	// Check if this is the first upvote
	isFirst, err := s.isFirstInBatch(ctx, batchKey, voterID)
	if err != nil {
		log.Printf("ERROR: Failed to check first upvote: %v", err)
		return
	}

	// Add to batch
	if err := s.addToBatch(ctx, batchKey, voterID, commentUpvoteBatchTTL); err != nil {
		log.Printf("ERROR: Failed to add to batch: %v", err)
		return
	}

	// If first upvote, send instant notification
	if isFirst {
		voter, err := s.userRepo.GetByID(ctx, voterID)
		if err != nil {
			return
		}

		comment, err := s.commentRepo.GetByID(ctx, commentID)
		if err != nil {
			return
		}

		recipientObjID, _ := primitive.ObjectIDFromHex(authorID)
		actorObjID, _ := primitive.ObjectIDFromHex(voterID)

		notification := &model.Notification{
			RecipientID: recipientObjID,
			ActorID:     actorObjID,
			Type:        model.NotificationTypeLike,
			Message:     fmt.Sprintf("%s đã thích bình luận của bạn", voter.Username),
			Link:        fmt.Sprintf("/posts/%s#comment-%s", comment.PostID.Hex(), commentID),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  authorID,
			Notification: dto.FromNotification(createdNotification),
		})
		return
	}

	// Check if reached threshold
	count, err := s.getBatchCount(ctx, batchKey)
	if err != nil {
		return
	}

	if count >= commentUpvoteThreshold {
		s.sendCommentUpvoteBatchNotification(ctx, batchKey)
	}
}

// handleCommentCreatedReply handles instant notifications for comment replies only
func (s *notificationService) handleCommentCreatedReply(event bus.Event) {
	payload := event.Payload()
	authorID, _ := payload["author_id"].(string)
	parentAuthorID, _ := payload["parent_author_id"].(string)
	postID, _ := payload["post_id"].(string)
	commentID, _ := payload["comment_id"].(string)

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Only handle reply notifications here (instant notification)
	if parentAuthorID != "" && authorID != parentAuthorID {
		// Check if recipient wants comment notifications
		if !s.shouldNotify(ctx, parentAuthorID, model.NotificationTypeComment) {
			return
		}
		author, err := s.userRepo.GetByID(ctx, authorID)
		if err != nil {
			return
		}

		recipientObjID, _ := primitive.ObjectIDFromHex(parentAuthorID)
		actorObjID, _ := primitive.ObjectIDFromHex(authorID)

		notification := &model.Notification{
			RecipientID: recipientObjID,
			ActorID:     actorObjID,
			Type:        model.NotificationTypeComment,
			Message:     fmt.Sprintf("%s đã trả lời một bình luận của bạn.", author.Username),
			Link:        fmt.Sprintf("/posts/%s#comment-%s", postID, commentID),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  parentAuthorID,
			Notification: dto.FromNotification(createdNotification),
		})
	}
}

// handleCommentApprovedForPost handles batched notifications for approved comments on posts
func (s *notificationService) handleCommentApprovedForPost(event bus.Event) {
	payload := event.Payload()
	commentID, _ := payload["comment_id"].(string)
	postID, _ := payload["post_id"].(string)
	authorID, _ := payload["author_id"].(string)
	parentAuthorID, _ := payload["parent_author_id"].(string)

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get post to find post author
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return
	}

	postAuthorID := post.AuthorID.Hex()

	// Skip if commenter is post author
	if authorID == postAuthorID {
		return
	}

	// Skip if this is a reply and parent author is post author (already handled by reply notification)
	if parentAuthorID != "" && parentAuthorID == postAuthorID {
		return
	}

	// Check if recipient wants comment notifications
	if !s.shouldNotify(ctx, postAuthorID, model.NotificationTypeComment) {
		return
	}

	// Batching logic for post author notification
	batchKey := fmt.Sprintf(postCommentBatchKeyPrefix, postID)

	// Check if this is the first comment
	isFirst, err := s.isFirstInBatch(ctx, batchKey, authorID)
	if err != nil {
		log.Printf("ERROR: Failed to check first comment: %v", err)
		return
	}

	// Add to batch
	if err := s.addToBatch(ctx, batchKey, authorID, postCommentBatchTTL); err != nil {
		log.Printf("ERROR: Failed to add to batch: %v", err)
		return
	}

	// If first comment, send instant notification
	if isFirst {
		author, err := s.userRepo.GetByID(ctx, authorID)
		if err != nil {
			return
		}

		recipientObjID, _ := primitive.ObjectIDFromHex(postAuthorID)
		actorObjID, _ := primitive.ObjectIDFromHex(authorID)

		notification := &model.Notification{
			RecipientID: recipientObjID,
			ActorID:     actorObjID,
			Type:        model.NotificationTypeComment,
			Message:     fmt.Sprintf("%s đã bình luận vào bài viết của bạn", author.Username),
			Link:        fmt.Sprintf("/posts/%s#comment-%s", postID, commentID),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  postAuthorID,
			Notification: dto.FromNotification(createdNotification),
		})
		return
	}

	// Check if reached threshold
	count, err := s.getBatchCount(ctx, batchKey)
	if err != nil {
		return
	}

	if count >= postCommentThreshold {
		s.sendPostCommentBatchNotification(ctx, batchKey)
	}
}

func (s *notificationService) handleModeratorsAdded(event bus.Event) {
	payload := event.Payload()

	communityID, _ := payload["community_id"].(string)
	moderatorIDs, _ := payload["moderator_ids"].([]string)

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	community, err := s.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		log.Printf("ERROR: NotificationService: failed to get community: %v", err)
		return
	}

	communityObjID, _ := primitive.ObjectIDFromHex(communityID)

	for _, moderatorID := range moderatorIDs {

		recipientObjID, err := primitive.ObjectIDFromHex(moderatorID)
		if err != nil {
			log.Printf("ERROR: NotificationService: invalid moderatorID: %v", err)
			continue
		}

		notification := &model.Notification{
			RecipientID: recipientObjID,
			ActorID:     communityObjID,
			Type:        model.NotificationTypeSystem,
			Message:     fmt.Sprintf("Bạn được mời làm moderator của cộng đồng %s", community.Name),
			Link:        fmt.Sprintf("/communities/%s", community.Name),
			IsRead:      false,
			Metadata: map[string]interface{}{
				"community_id": communityID,
				"action_type":  "moderator_invitation",
			},
			CreatedAt: time.Now(),
		}

		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			continue
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  moderatorID,
			Notification: dto.FromNotification(createdNotification),
		})
	}
}

func (s *notificationService) handleBroadcast(event bus.Event) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	payload := event.Payload()
	recipientIDs, _ := payload["recipient_ids"].([]string)
	eventType, _ := payload["event_type"].(bus.BroadcastEventType)
	data := payload["data"]

	switch eventType {
	case bus.BroadcastEventMessageCreated:
		// Send notification to user not in chat
		var messageData dto.MessageResponse
		if err := util.DecodeJson(data, &messageData); err != nil {
			log.Printf("Failed to decode message: %v", err)
			return
		}

		key := fmt.Sprintf(config.RedisActiveUsersKey, messageData.ChannelID)
		isInChat, err := s.redisClient.SIsMember(ctx, key, recipientIDs[0]).Result()
		if err != nil {
			log.Printf("Failed to check active users: %v", err)
			return
		}

		if isInChat {
			// Recipient is in chat, skip notification
			return
		}

		recipientObjectID, err := primitive.ObjectIDFromHex(recipientIDs[0])
		if err != nil {
			return
		}

		actorObjectID, err := primitive.ObjectIDFromHex(messageData.SenderID)
		if err != nil {
			return
		}

		notification := &model.Notification{
			RecipientID: recipientObjectID,
			ActorID:     actorObjectID,
			Type:        model.NotificationTypeNewMessage,
			Message:     fmt.Sprintf("Tin nhắn mới từ %s", messageData.SenderUsername),
			Link:        fmt.Sprintf("/channels/%s", messageData.ChannelID),
			IsRead:      false,
			CreatedAt:   time.Now(),
		}

		createdNotification, err := s.notificationRepo.Create(ctx, notification)
		if err != nil {
			log.Printf("ERROR: NotificationService: failed to create notification: %v", err)
			return
		}

		s.eventBus.Publish(bus.NotificationCreatedEvent{
			RecipientID:  recipientIDs[0],
			Notification: dto.FromNotification(createdNotification),
		})
	}
}

func (s *notificationService) GetNotifications(recipientID string, page, pageSize int) (*dto.PaginatedNotificationsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	notifications, total, err := s.notificationRepo.GetByRecipientID(ctx, recipientID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedNotificationsResponse{
		Notifications: dto.FromNotifications(notifications),
		Pagination: dto.Pagination{
			Total: total,
			Page:  page,
		},
	}, nil
}

func (s *notificationService) MarkAllAsRead(recipientID string) (int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.notificationRepo.MarkAllAsRead(ctx, recipientID)
}

func (s *notificationService) MarkAsRead(notificationID, recipientID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return fmt.Errorf("invalid notification ID: %w", err)
	}

	return s.notificationRepo.MarkAsRead(ctx, objID, recipientID)
}

func (s *notificationService) DeleteNotification(notificationID, recipientID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return fmt.Errorf("invalid notification ID: %w", err)
	}

	return s.notificationRepo.DeleteNotification(ctx, objID, recipientID)
}

// ========== Redis Batching Helpers ==========

const (
	postUpvoteBatchKeyPrefix    = "notification:batch:post:%s:upvotes"
	commentUpvoteBatchKeyPrefix = "notification:batch:comment:%s:upvotes"
	postCommentBatchKeyPrefix   = "notification:batch:post:%s:comments"

	postUpvoteBatchTTL    = 3600 // 1 hour
	commentUpvoteBatchTTL = 3600 // 1 hour
	postCommentBatchTTL   = 1800 // 30 minutes

	postUpvoteThreshold    = 10 // Gửi ngay khi đủ 10 upvotes
	commentUpvoteThreshold = 5  // Gửi ngay khi đủ 5 upvotes
	postCommentThreshold   = 3  // Gửi ngay khi đủ 3 comments

	batchProcessInterval = 300 // Check mỗi 5 phút
)

// addToBatch adds an actor to a batched notification
func (s *notificationService) addToBatch(ctx context.Context, key string, actorID string, ttl int) error {
	score := float64(time.Now().Unix())
	if err := s.redisClient.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: actorID,
	}).Err(); err != nil {
		return err
	}

	// Set TTL
	s.redisClient.Expire(ctx, key, time.Duration(ttl)*time.Second)
	return nil
}

// getBatchMembers gets all members from a batch
func (s *notificationService) getBatchMembers(ctx context.Context, key string) ([]string, error) {
	members, err := s.redisClient.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}

// getBatchCount gets the count of members in a batch
func (s *notificationService) getBatchCount(ctx context.Context, key string) (int64, error) {
	count, err := s.redisClient.ZCard(ctx, key).Result()
	return count, err
}

// clearBatch removes a batch key
func (s *notificationService) clearBatch(ctx context.Context, key string) error {
	return s.redisClient.Del(ctx, key).Err()
}

// isFirstInBatch checks if this is the first item in batch
func (s *notificationService) isFirstInBatch(ctx context.Context, key string, actorID string) (bool, error) {
	// Check if key exists
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// If key doesn't exist, this is the first
	if exists == 0 {
		return true, nil
	}

	// Check if this actor already exists in the batch
	score, err := s.redisClient.ZScore(ctx, key, actorID).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	// If actor already in batch (score exists), this is NOT first
	if err != redis.Nil && score > 0 {
		return false, nil
	}

	// Actor not in batch yet, check current count
	count, err := s.getBatchCount(ctx, key)
	if err != nil {
		return false, err
	}

	// If batch is empty, this is the first
	return count == 0, nil
}

// processBatchedNotifications runs periodically to send batched notifications
func (s *notificationService) processBatchedNotifications() {
	ticker := time.NewTicker(time.Duration(batchProcessInterval) * time.Second)
	defer ticker.Stop()

	log.Println("Batched notification processor started")

	for range ticker.C {
		ctx, cancel := util.NewDefaultDBContext()

		// Process post upvote batches
		s.processPostUpvoteBatches(ctx)

		// Process comment upvote batches
		s.processCommentUpvoteBatches(ctx)

		// Process post comment batches
		s.processPostCommentBatches(ctx)

		cancel()
	}
}

// processPostUpvoteBatches processes all pending post upvote batches
func (s *notificationService) processPostUpvoteBatches(ctx context.Context) {
	pattern := fmt.Sprintf(postUpvoteBatchKeyPrefix, "*")
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		log.Printf("ERROR: Failed to get post upvote batch keys: %v", err)
		return
	}

	for _, key := range keys {
		count, err := s.getBatchCount(ctx, key)
		if err != nil {
			continue
		}

		if count == 0 {
			s.clearBatch(ctx, key)
			continue
		}

		s.sendPostUpvoteBatchNotification(ctx, key)
	}
}

// processCommentUpvoteBatches processes all pending comment upvote batches
func (s *notificationService) processCommentUpvoteBatches(ctx context.Context) {
	pattern := fmt.Sprintf(commentUpvoteBatchKeyPrefix, "*")
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		log.Printf("ERROR: Failed to get comment upvote batch keys: %v", err)
		return
	}

	for _, key := range keys {
		// Check if batch has expired or reached threshold
		count, err := s.getBatchCount(ctx, key)
		if err != nil {
			continue
		}

		// Skip if empty
		if count == 0 {
			s.clearBatch(ctx, key)
			continue
		}

		// Send batched notification
		s.sendCommentUpvoteBatchNotification(ctx, key)
	}
}

// processPostCommentBatches processes all pending post comment batches
func (s *notificationService) processPostCommentBatches(ctx context.Context) {
	pattern := fmt.Sprintf(postCommentBatchKeyPrefix, "*")
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		log.Printf("ERROR: Failed to get post comment batch keys: %v", err)
		return
	}

	for _, key := range keys {
		count, err := s.getBatchCount(ctx, key)
		if err != nil {
			continue
		}

		if count == 0 {
			s.clearBatch(ctx, key)
			continue
		}

		s.sendPostCommentBatchNotification(ctx, key)
	}
}

// sendPostUpvoteBatchNotification sends batched post upvote notification
func (s *notificationService) sendPostUpvoteBatchNotification(ctx context.Context, batchKey string) {
	// Extract postID from key: "notification:batch:post:{postID}:upvotes"
	parts := splitBatchKey(batchKey)
	if len(parts) < 4 {
		return
	}
	postID := parts[3]

	// Get all voters
	voterIDs, err := s.getBatchMembers(ctx, batchKey)
	if err != nil || len(voterIDs) == 0 {
		return
	}

	// Get post to find author
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		log.Printf("ERROR: Failed to get post %s: %v", postID, err)
		return
	}

	authorID := post.AuthorID.Hex()

	// Get first voter info
	firstVoter, err := s.userRepo.GetByID(ctx, voterIDs[0])
	if err != nil {
		return
	}

	// Build message
	var message string
	if len(voterIDs) == 1 {
		message = fmt.Sprintf("%s đã thích bài viết của bạn: %s", firstVoter.Username, post.Title)
	} else {
		message = fmt.Sprintf("%s và %d người khác đã thích bài viết của bạn: %s", firstVoter.Username, len(voterIDs)-1, post.Title)
	}

	recipientObjID, _ := primitive.ObjectIDFromHex(authorID)
	actorObjID, _ := primitive.ObjectIDFromHex(voterIDs[0])

	notification := &model.Notification{
		RecipientID: recipientObjID,
		ActorID:     actorObjID,
		Type:        model.NotificationTypeLike,
		Message:     message,
		Link:        fmt.Sprintf("/posts/%s", postID),
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	createdNotification, err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		log.Printf("ERROR: Failed to create batched notification: %v", err)
		return
	}

	s.eventBus.Publish(bus.NotificationCreatedEvent{
		RecipientID:  authorID,
		Notification: dto.FromNotification(createdNotification),
	})

	// Clear batch after sending
	s.clearBatch(ctx, batchKey)
}

// sendCommentUpvoteBatchNotification sends batched comment upvote notification
func (s *notificationService) sendCommentUpvoteBatchNotification(ctx context.Context, batchKey string) {
	// Extract commentID from key: "notification:batch:comment:{commentID}:upvotes"
	parts := splitBatchKey(batchKey)
	if len(parts) < 4 {
		return
	}
	commentID := parts[3]

	// Get all voters
	voterIDs, err := s.getBatchMembers(ctx, batchKey)
	if err != nil || len(voterIDs) == 0 {
		return
	}

	// Get comment to find author and postID
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		log.Printf("ERROR: Failed to get comment %s: %v", commentID, err)
		return
	}

	authorID := comment.Author.ID.Hex()

	// Get first voter info
	firstVoter, err := s.userRepo.GetByID(ctx, voterIDs[0])
	if err != nil {
		return
	}

	// Build message
	var message string
	if len(voterIDs) == 1 {
		message = fmt.Sprintf("%s đã thích bình luận của bạn", firstVoter.Username)
	} else {
		message = fmt.Sprintf("%s và %d người khác đã thích bình luận của bạn", firstVoter.Username, len(voterIDs)-1)
	}

	recipientObjID, _ := primitive.ObjectIDFromHex(authorID)
	actorObjID, _ := primitive.ObjectIDFromHex(voterIDs[0])

	notification := &model.Notification{
		RecipientID: recipientObjID,
		ActorID:     actorObjID,
		Type:        model.NotificationTypeLike,
		Message:     message,
		Link:        fmt.Sprintf("/posts/%s#comment-%s", comment.PostID.Hex(), commentID),
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	createdNotification, err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		log.Printf("ERROR: Failed to create batched notification: %v", err)
		return
	}

	s.eventBus.Publish(bus.NotificationCreatedEvent{
		RecipientID:  authorID,
		Notification: dto.FromNotification(createdNotification),
	})

	// Clear batch after sending
	s.clearBatch(ctx, batchKey)
}

// sendPostCommentBatchNotification sends batched post comment notification
func (s *notificationService) sendPostCommentBatchNotification(ctx context.Context, batchKey string) {
	// Extract postID from key: "notification:batch:post:{postID}:comments"
	parts := splitBatchKey(batchKey)
	if len(parts) < 4 {
		return
	}
	postID := parts[3]

	// Get all commenters
	commenterIDs, err := s.getBatchMembers(ctx, batchKey)
	if err != nil || len(commenterIDs) == 0 {
		return
	}

	// Get post to find author
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		log.Printf("ERROR: Failed to get post %s: %v", postID, err)
		return
	}

	authorID := post.AuthorID.Hex()

	// Get first commenter info
	firstCommenter, err := s.userRepo.GetByID(ctx, commenterIDs[0])
	if err != nil {
		return
	}

	// Build message
	var message string
	if len(commenterIDs) == 1 {
		message = fmt.Sprintf("%s đã bình luận vào bài viết của bạn", firstCommenter.Username)
	} else {
		message = fmt.Sprintf("%s và %d người khác đã bình luận vào bài viết của bạn", firstCommenter.Username, len(commenterIDs)-1)
	}

	recipientObjID, _ := primitive.ObjectIDFromHex(authorID)
	actorObjID, _ := primitive.ObjectIDFromHex(commenterIDs[0])

	notification := &model.Notification{
		RecipientID: recipientObjID,
		ActorID:     actorObjID,
		Type:        model.NotificationTypeComment,
		Message:     message,
		Link:        fmt.Sprintf("/posts/%s", postID),
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	createdNotification, err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		log.Printf("ERROR: Failed to create batched notification: %v", err)
		return
	}

	s.eventBus.Publish(bus.NotificationCreatedEvent{
		RecipientID:  authorID,
		Notification: dto.FromNotification(createdNotification),
	})

	// Clear batch after sending
	s.clearBatch(ctx, batchKey)
}

// splitBatchKey splits Redis key by colon
func splitBatchKey(key string) []string {
	return strings.Split(key, ":")
}

// shouldNotify checks if recipient wants to receive this notification type
func (s *notificationService) shouldNotify(ctx context.Context, recipientID string, notifType model.NotificationType) bool {
	recipient, err := s.userRepo.GetByID(ctx, recipientID)
	if err != nil {
		// If can't get user settings, allow notification (fail-open)
		return true
	}

	// No settings = use defaults (all enabled)
	if recipient.Settings == nil {
		return true
	}

	settings := recipient.Settings.Notifications

	// Check in-app enabled first
	if !settings.InAppEnabled {
		return false
	}

	// Check specific notification type preferences
	switch notifType {
	case model.NotificationTypeComment:
		return settings.NotifyOnComment
	case model.NotificationTypeLike:
		return settings.NotifyOnUpvote
	case model.NotificationTypeMention:
		return settings.NotifyOnMention
	case model.NotificationTypeNewMessage:
		return settings.NotifyOnMessage
	default:
		return true // Unknown types allowed by default
	}
}
