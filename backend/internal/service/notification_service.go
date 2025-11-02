package service

import (
	"fmt"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService interface {
	Start()
	GetNotifications(recipientID string, page, pageSize int) (*dto.PaginatedNotificationsResponse, error)
	MarkAllAsRead(recipientID string) (int64, error)
}

type notificationService struct {
	notificationRepo repo.NotificationRepo
	userRepo         repo.UserRepo
	postRepo         repo.PostRepo
	commentRepo      repo.CommentRepo
	eventBus         *bus.EventBus
}

func NewNotificationService(notificationRepo repo.NotificationRepo, userRepo repo.UserRepo, postRepo repo.PostRepo, commentRepo repo.CommentRepo, bus *bus.EventBus) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		postRepo:         postRepo,
		commentRepo:      commentRepo,
		eventBus:         bus,
	}
}

func (s *notificationService) Start() {
	eventChannel := make(bus.EventListener, 100)

	s.eventBus.Subscribe(bus.TopicPostUpvoted, eventChannel)
	s.eventBus.Subscribe(bus.TopicCommentCreated, eventChannel)

	log.Println("NotificationService started and subscribed to events.")

	go s.processEvents(eventChannel)
}

func (s *notificationService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicPostUpvoted:
			s.handlePostUpvoted(event)
		case bus.TopicCommentCreated:
			s.handleCommentCreated(event)
		}
	}
}

func (s *notificationService) handlePostUpvoted(event bus.Event) {
	payload := event.Payload()
	authorID, _ := payload["authorId"].(string)
	voterID, _ := payload["voterId"].(string)
	postID, _ := payload["postId"].(string)

	if authorID == "" || voterID == "" || postID == "" {
		return
	}
	if authorID == voterID {
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	voter, err := s.userRepo.GetByID(ctx, voterID)
	if err != nil {
		return
	}

	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return
	}

	post, err := s.postRepo.GetByID(ctx, postObjID)
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
}

func (s *notificationService) handleCommentCreated(event bus.Event) {
	payload := event.Payload()
	authorID, _ := payload["authorId"].(string)
	parentAuthorID, _ := payload["parentAuthorId"].(string)
	postID, _ := payload["postId"].(string)
	commentID, _ := payload["commentId"].(string)

	if parentAuthorID == "" || authorID == parentAuthorID {
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

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
