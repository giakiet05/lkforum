package service

import (
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChannelService interface {
	Start()
	CreateChannel(req *dto.CreateChannelRequest, requesterID string) (*model.Channel, int64, error)
	GetChannelByID(channelID string, requesterID string) (*model.Channel, int64, error)
	GetChannelsByUserID(userID string, requesterID string, page int, pageSize int) (*dto.PaginatedChannelsResponse, error)
	GetChannelByBothUserID(user1ID string, user2ID string, requesterID string) (*model.Channel, int64, error)
	UpdateChannel(req *dto.UpdateChannelRequest, requesterID string) (*model.Channel, error)
	DeleteChannel(channelID string, userID string) error
}

type channelService struct {
	channelRepository repo.ChannelRepo
	messageRepository repo.MessageRepo
	eventBus          *bus.EventBus
}

func NewChannelService(channelRepo repo.ChannelRepo, messageRepo repo.MessageRepo, bus *bus.EventBus) ChannelService {
	return &channelService{
		channelRepository: channelRepo,
		messageRepository: messageRepo,
		eventBus:          bus,
	}
}

func (s *channelService) Start() {
	eventChannel := make(bus.EventListener, 100)

	s.eventBus.Subscribe(bus.TopicUserChangeAvatar, eventChannel)

	log.Println("ChannelService started and subscribed to events.")

	go s.processEvents(eventChannel)
}

func (s *channelService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicUserChangeAvatar:
			s.handleNewAvatar(event)
		default:
			log.Println("Unhandled event topic:", event.Topic())
		}
	}
}

func (s *channelService) handleNewAvatar(event bus.Event) {
	payload := event.Payload()
	userID, _ := payload["user_id"].(string)
	newAvatar, _ := payload["new_avatar"].(string)

	if userID == "" || newAvatar == "" {
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	err := s.channelRepository.UpdateUserAvatar(ctx, userID, newAvatar)
	if err != nil {
		log.Printf("Failed to update avatar: %v", err)
	}
}

func (s *channelService) CreateChannel(req *dto.CreateChannelRequest, requesterID string) (*model.Channel, int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if requesterID != req.Member1 && requesterID != req.Member2 {
		return nil, 0, apperror.ErrForbidden
	}

	member1ID, err := primitive.ObjectIDFromHex(req.Member1)
	if err != nil {
		return nil, 0, err
	}
	member2ID, err := primitive.ObjectIDFromHex(req.Member2)
	if err != nil {
		return nil, 0, err
	}

	channel := &model.Channel{
		ID: primitive.NewObjectID(),
		Members: []model.ChannelMember{
			{
				ID:       member1ID,
				Username: req.Member1Username,
				Avatar:   req.Member1Avatar,
			},
			{
				ID:       member2ID,
				Username: req.Member2Username,
				Avatar:   req.Member2Avatar,
			},
		},
		Settings: []model.ChannelSetting{
			{
				UserID:          member1ID,
				Notification:    true,
				TypingIndicator: true,
				IsDeleted:       false,
			},
			{
				UserID:          member2ID,
				Notification:    true,
				TypingIndicator: true,
				IsDeleted:       false,
			},
		},
		Status:    model.ChannelStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	channel, err = s.channelRepository.Create(ctx, channel)
	if err != nil {
		return nil, 0, err
	}

	unreadMessageCount, err := s.messageRepository.CountUnreadMessages(ctx, channel.ID.Hex(), requesterID)
	if err != nil {
		return nil, 0, err
	}

	return channel, unreadMessageCount, nil
}

func (s *channelService) GetChannelByID(channelID string, requesterID string) (*model.Channel, int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	channel, err := s.channelRepository.GetByID(ctx, channelID)
	if err != nil {
		return nil, 0, err
	}

	isMember, err := s.channelRepository.IsMember(ctx, channelID, requesterID)
	if err != nil {
		return nil, 0, err
	}
	if !isMember {
		return nil, 0, apperror.ErrForbidden
	}

	unreadMessageCount, err := s.messageRepository.CountUnreadMessages(ctx, channel.ID.Hex(), requesterID)
	if err != nil {
		return nil, 0, err
	}

	return channel, unreadMessageCount, nil
}

func (s *channelService) GetChannelsByUserID(userID string, requesterID string, page int, pageSize int) (*dto.PaginatedChannelsResponse, error) {
	if userID != requesterID {
		return nil, apperror.ErrForbidden
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	channels, total, err := s.channelRepository.GetByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	unreadCounts := make([]*int64, len(channels))
	for i, ch := range channels {
		count, err := s.messageRepository.CountUnreadMessages(ctx, ch.ID.Hex(), userID)
		if err != nil {
			return nil, err
		}
		unreadCounts[i] = &count
	}
	channelResponses, err := dto.FromChannels(channels, unreadCounts)
	if err != nil {
		return nil, err
	}

	var response = dto.PaginatedChannelsResponse{
		Channels: channelResponses,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}

	return &response, nil
}

func (s *channelService) GetChannelByBothUserID(user1ID string, user2ID string, requesterID string) (*model.Channel, int64, error) {
	if requesterID != user1ID && requesterID != user2ID {
		return nil, 0, apperror.ErrForbidden
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	channel, err := s.channelRepository.GetByBothUserID(ctx, user1ID, user2ID)
	if err != nil {
		return nil, 0, err
	}

	unreadMessageCount, err := s.messageRepository.CountUnreadMessages(ctx, channel.ID.Hex(), requesterID)
	if err != nil {
		return nil, 0, err
	}

	return channel, unreadMessageCount, nil
}

func (s *channelService) UpdateChannel(req *dto.UpdateChannelRequest, requesterID string) (*model.Channel, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	channel, err := s.channelRepository.GetByID(ctx, req.ChannelID)
	if err != nil {
		return nil, err
	}

	isMember := false
	for _, member := range channel.Members {
		if member.ID.Hex() == requesterID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	requesterObjectID, err := primitive.ObjectIDFromHex(requesterID)
	if err != nil {
		return nil, err
	}

	updated := false
	for i := range channel.Settings {
		if channel.Settings[i].UserID == requesterObjectID {
			if req.Nickname != nil {
				channel.Settings[i].Nickname = req.Nickname
			}
			if req.Notification != nil {
				channel.Settings[i].Notification = *req.Notification
			}
			if req.TypingIndicator != nil {
				channel.Settings[i].TypingIndicator = *req.TypingIndicator
			}
			updated = true
			break
		}
	}

	if !updated {
		return nil, apperror.ErrNoFieldsToUpdate
	}

	channel.UpdatedAt = time.Now()

	return s.channelRepository.Update(ctx, channel)
}

func (s *channelService) DeleteChannel(channelID string, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.channelRepository.Delete(ctx, channelID, userID)
}
