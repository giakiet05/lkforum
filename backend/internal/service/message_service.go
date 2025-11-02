package service

import (
	"errors"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageService interface {
	Start()

	//CreateMessage(req *dto.CreateMessageRequest, requesterID string) (*model.Message, error)
	GetMessageByID(channelID string, messageID string, requesterID string) (*model.Message, error)
	GetMessageFilter(query *dto.GetMessageFilterQuery, requesterID string) (*dto.PaginatedMessagesResponse, error)
	DeleteMessage(channelID string, messageID string, requesterID string) error
}

type messageService struct {
	messageRepository repo.MessageRepo
	channelRepository repo.ChannelRepo
	eventBus          *bus.EventBus
}

func NewMessageService(messageRepo repo.MessageRepo, channelRepo repo.ChannelRepo, bus *bus.EventBus) MessageService {
	return &messageService{
		messageRepository: messageRepo,
		channelRepository: channelRepo,
		eventBus:          bus,
	}
}

func (m *messageService) Start() {
	eventChannel := make(bus.EventListener, 100)

	m.eventBus.Subscribe(bus.TopicNewMessage, eventChannel)

	log.Println("MessageService started and subscribed to events.")

	go m.processEvents(eventChannel)
}

func (m *messageService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicNewMessage:
			m.handleNewMessage(event)
		}
	}
}

func (m *messageService) handleNewMessage(event bus.Event) {
	payload := event.Payload()

	tempMessageID, _ := payload["temp_message_id"].(string)
	channelID, _ := payload["channel_id"].(string)
	senderID, _ := payload["sender_id"].(string)
	content, _ := payload["content"].(string)

	var msgType model.MessageType
	if t, ok := payload["type"].(string); ok {
		msgType = model.MessageType(t)
	}

	if tempMessageID == "" || channelID == "" || senderID == "" || content == "" {
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
			recipientIDs = append(recipientIDs, member.UserID.Hex())
		}
	}
	if !isMember {
		m.publishMessageError(senderID, channelID, tempMessageID, apperror.ErrForbidden)
		return
	}

	message := &model.Message{
		ChannelID: channelObjectID,
		SenderID:  &senderObjectID,
		Type:      msgType,
		Content:   content,
		IsSend:    false,
		IsRead:    false,
		IsDeleted: false,
		CreatedAt: time.Now(),
	}

	message, err = m.messageRepository.Create(ctx, message)
	if err != nil {
		m.publishMessageError(senderID, channelID, tempMessageID, apperror.ErrInternal)
		return
	}

	m.eventBus.Publish(bus.MessageCreatedEvent{
		RecipientIDs:  recipientIDs,
		TempMessageID: tempMessageID,
		Message:       *dto.FromMessage(message),
	})
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
