package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageService interface {
	Start()
	GetMessageByID(channelID string, messageID string, requesterID string) (*model.Message, error)
	GetMessageFilter(query *dto.GetMessageFilterQuery, requesterID string) (*dto.PaginatedMessagesResponse, error)
	DeleteMessage(channelID string, messageID string, requesterID string) error
	MarkChannelAsRead(channelID string, userID string) error
}

type messageService struct {
	messageRepository repo.MessageRepo
	channelRepository repo.ChannelRepo
	userRepository    repo.UserRepo
	eventBus          bus.EventBus
	redisClient       *redis.Client
}

func NewMessageService(messageRepo repo.MessageRepo, channelRepo repo.ChannelRepo, userRepo repo.UserRepo, bus bus.EventBus, redis *redis.Client) MessageService {
	return &messageService{
		messageRepository: messageRepo,
		channelRepository: channelRepo,
		userRepository:    userRepo,
		eventBus:          bus,
		redisClient:       redis,
	}
}

func (m *messageService) Start() {
	eventChannel := make(bus.EventListener, 100)

	m.eventBus.Subscribe(bus.TopicNewMessage, eventChannel)
	m.eventBus.Subscribe(bus.TopicTypingMessage, eventChannel)
	m.eventBus.Subscribe(bus.TopicInChatMessage, eventChannel)

	log.Println("MessageService started and subscribed to events.")

	go m.processEvents(eventChannel)
}

func (m *messageService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicNewMessage:
			m.handleNewMessage(event)
		case bus.TopicTypingMessage:
			m.handleTypingEvent(event)
		case bus.TopicInChatMessage:
			m.handleInChatEvent(event)
		default:
			log.Println("Unhandled event topic:", event.Topic())
		}
	}
}

func (m *messageService) handleNewMessage(event bus.Event) {
	payload := event.Payload()

	tempMessageID, _ := payload["temp_message_id"].(string)
	channelID, _ := payload["channel_id"].(string)
	senderID, _ := payload["sender_id"].(string)
	senderUsername, _ := payload["sender_username"].(string)
	content, _ := payload["content"].(string)

	var msgType model.MessageType
	if t, ok := payload["type"].(string); ok {
		msgType = model.MessageType(t)
	}

	if tempMessageID == "" || channelID == "" || senderID == "" || senderUsername == "" || content == "" {
		m.publishMessageError(senderID, channelID, tempMessageID, apperror.ErrBadRequest)
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		m.publishMessageError(senderID, channelID, tempMessageID, apperror.ErrBadRequest)
		return
	}

	senderObjectID, err := primitive.ObjectIDFromHex(senderID)
	if err != nil {
		m.publishMessageError(senderID, channelID, tempMessageID, apperror.ErrBadRequest)
		return
	}

	channel, err := m.channelRepository.GetByID(ctx, channelID)
	if err != nil {
		m.publishMessageError(senderID, channelID, tempMessageID, apperror.ErrChannelNotFound)
		return
	}

	isMember := false
	var recipientIDs []string
	for _, member := range channel.Members {
		if member.UserID == senderObjectID {
			isMember = true
		} else {
			// Check if recipient allows direct messages
			recipient, err := m.userRepository.GetByID(ctx, member.UserID.Hex())
			if err == nil && recipient.Settings != nil && !recipient.Settings.Privacy.AllowDirectMessages {
				// Recipient has disabled DMs - reject message
				m.publishMessageError(senderID, channelID, tempMessageID, *apperror.NewError(
					nil,
					apperror.ErrForbidden.Code,
					fmt.Sprintf("User %s has disabled direct messages", recipient.Username),
				))
				return
			}
			recipientIDs = append(recipientIDs, member.UserID.Hex())
		}
	}
	if !isMember {
		m.publishMessageError(senderID, channelID, tempMessageID, apperror.ErrForbidden)
		return
	}

	message := &model.Message{
		ChannelID:      channelObjectID,
		SenderID:       &senderObjectID,
		SenderUsername: senderUsername,
		Type:           msgType,
		Content:        content,
		IsSend:         false,
		IsRead:         false,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
	}

	message, err = m.messageRepository.Create(ctx, message)
	if err != nil {
		m.publishMessageError(senderID, channelID, tempMessageID, apperror.ErrInternal)
		return
	}

	broadcastEvent := bus.BroadcastEvent{
		RecipientIDs: recipientIDs,
		EventType:    bus.BroadcastEventMessageCreated,
		TempID:       tempMessageID,
		Data:         dto.FromMessage(message),
	}

	m.eventBus.Publish(broadcastEvent)
}

func (m *messageService) handleTypingEvent(event bus.Event) {
	payload := event.Payload()

	channelID, _ := payload["channel_id"].(string)
	senderID, _ := payload["sender_id"].(string)
	isTyping, _ := payload["is_typing"].(bool)

	if channelID == "" || senderID == "" {
		m.publishMessageError(senderID, channelID, "", apperror.ErrBadRequest)
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	channel, err := m.channelRepository.GetByID(ctx, channelID)
	if err != nil {
		m.publishMessageError(senderID, channelID, "", apperror.ErrChannelNotFound)
		return
	}

	// Gather recipients except the sender
	var recipientIDs []string
	for _, member := range channel.Members {
		if member.UserID.Hex() != senderID {
			recipientIDs = append(recipientIDs, member.UserID.Hex())
		}
	}

	eventType := bus.BroadcastEventTypingStop
	if isTyping {
		eventType = bus.BroadcastEventTypingStart
	}

	data := map[string]interface{}{
		"channel_id": channelID,
		"sender_id":  senderID,
		"is_typing":  isTyping,
	}

	// Publish a broadcast event so the WS Hub can send it to clients
	m.eventBus.Publish(bus.BroadcastEvent{
		RecipientIDs: recipientIDs,
		EventType:    eventType,
		Data:         data,
	})
}

func (m *messageService) handleInChatEvent(event bus.Event) {
	payload := event.Payload()

	channelID, _ := payload["channel_id"].(string)
	userID, _ := payload["user_id"].(string)
	isInChat, _ := payload["is_in_chat"].(bool)

	if channelID == "" || userID == "" {
		m.publishMessageError(userID, channelID, "", apperror.ErrBadRequest)
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	_, err := m.channelRepository.GetByID(ctx, channelID)
	if err != nil {
		m.publishMessageError(userID, channelID, "", apperror.ErrChannelNotFound)
		return
	}

	key := fmt.Sprintf(config.RedisActiveUsersKey, channelID)
	if isInChat {
		// User joins chat
		if err := m.redisClient.SAdd(ctx, key, userID).Err(); err != nil {
			m.publishMessageError(userID, channelID, "", apperror.ErrInternal)
			return
		}
	} else {
		// User leaves chat
		if err := m.redisClient.SRem(ctx, key, userID).Err(); err != nil {
			m.publishMessageError(userID, channelID, "", apperror.ErrInternal)
			return
		}
	}
}

func (m *messageService) publishMessageError(senderID string, channelID string, tempMessageID string, err apperror.AppError) {
	m.eventBus.Publish(bus.MessageErrorEvent{
		SenderID:      senderID,
		ChannelID:     channelID,
		TempMessageID: tempMessageID,
		ErrorCode:     err.Code,
		ErrorMsg:      err.Message,
	})
}

func (m *messageService) GetMessageByID(channelID string, messageID string, requesterID string) (*model.Message, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := m.channelRepository.IsMember(ctx, channelID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	message, err := m.messageRepository.GetByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrNoMessageFound
		}
		return nil, err
	}

	return message, nil
}

func (m *messageService) GetMessageFilter(query *dto.GetMessageFilterQuery, requesterID string) (*dto.PaginatedMessagesResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := m.channelRepository.IsMember(ctx, query.ChannelID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	messages, total, err := m.messageRepository.GetFilter(ctx,
		query.ChannelID, query.SenderID,
		query.SearchContent,
		query.IsRead, query.IsSend, query.IsMedia,
		query.Page, query.PageSize,
	)
	if err != nil {
		return nil, err
	}
	messageResponses := dto.FromMessages(messages)

	var response = dto.PaginatedMessagesResponse{
		Messages: messageResponses,
		Pagination: dto.Pagination{
			Page:     query.Page,
			PageSize: query.PageSize,
			Total:    total,
		},
	}

	return &response, nil
}

func (m *messageService) DeleteMessage(channelID string, messageID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := m.channelRepository.IsMember(ctx, channelID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	ok, err = m.messageRepository.IsSendByUser(ctx, messageID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	return m.messageRepository.Delete(ctx, messageID)
}

func (m *messageService) MarkChannelAsRead(channelID string, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := m.channelRepository.IsMember(ctx, channelID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return apperror.ErrForbidden
	}

	return m.messageRepository.MarkChannelAsRead(ctx, channelID, userID)
}
