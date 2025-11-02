package service

import (
	"errors"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageService interface {
	//CreateMessage(req *dto.CreateMessageRequest, requesterID string) (*model.Message, error)
	GetMessageByID(channelID string, messageID string, requesterID string) (*model.Message, error)
	GetMessageFilter(query *dto.GetMessageFilterQuery, requesterID string) (*dto.PaginatedMessagesResponse, error)
	DeleteMessage(channelID string, messageID string, requesterID string) error
}

type messageService struct {
	messageRepository repo.MessageRepo
	channelRepository repo.ChannelRepo
}

func NewMessageService(messageRepo repo.MessageRepo, channelRepo repo.ChannelRepo) MessageService {
	return &messageService{
		messageRepository: messageRepo,
		channelRepository: channelRepo,
	}
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
			return nil, apperror.ErrCommunityNotFound
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
