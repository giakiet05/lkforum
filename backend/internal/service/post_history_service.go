package service

import (
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostHistoryService interface {
	CreatePostHistory(request *dto.CreatePostHistoryRequest, requesterID string) (*model.PostHistory, error)
	CreatePostHistories(requests []*dto.CreatePostHistoryRequest, requesterID string) ([]*model.PostHistory, error)
	GetPostHistoryByID(id string, requesterID string) (*model.PostHistory, error)
	GetPostHistoryByUserID(userID string, requesterID string, page int, pageSize int) (*dto.PaginatedPostHistoryResponse, error)
	DeletePostHistoryByID(id string) error
}

type postHistoryService struct {
	postHistoryRepository repo.PostHistoryRepo
}

func NewPostHistoryService(postHistoryRepository repo.PostHistoryRepo) PostHistoryService {
	return &postHistoryService{postHistoryRepository: postHistoryRepository}
}

func (p *postHistoryService) CreatePostHistory(request *dto.CreatePostHistoryRequest, requesterID string) (*model.PostHistory, error) {
	if requesterID != request.UserID {
		return nil, apperror.ErrForbidden
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		return nil, err
	}

	postObjectID, err := primitive.ObjectIDFromHex(request.PostID)
	if err != nil {
		return nil, err
	}

	postHistory := &model.PostHistory{
		UserID:   userObjectID,
		PostID:   postObjectID,
		ViewedAt: time.Now(),
	}

	return p.postHistoryRepository.Create(ctx, postHistory)
}

func (p *postHistoryService) CreatePostHistories(requests []*dto.CreatePostHistoryRequest, requesterID string) ([]*model.PostHistory, error) {
	if len(requests) == 0 {
		return nil, nil
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	postHistories := make([]*model.PostHistory, 0, len(requests))

	for _, req := range requests {
		if requesterID != req.UserID {
			return nil, apperror.ErrForbidden
		}

		userObjectID, err := primitive.ObjectIDFromHex(req.UserID)
		if err != nil {
			return nil, err
		}

		postObjectID, err := primitive.ObjectIDFromHex(req.PostID)
		if err != nil {
			return nil, err
		}

		postHistories = append(postHistories, &model.PostHistory{
			UserID:   userObjectID,
			PostID:   postObjectID,
			ViewedAt: time.Now(),
		})
	}

	return p.postHistoryRepository.CreateBatch(ctx, postHistories)
}

func (p *postHistoryService) GetPostHistoryByID(id string, requesterID string) (*model.PostHistory, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	postHistory, err := p.postHistoryRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if postHistory.UserID.Hex() != requesterID {
		return nil, apperror.ErrForbidden
	}

	return postHistory, nil
}

func (p *postHistoryService) GetPostHistoryByUserID(
	userID string, requesterID string,
	page int, pageSize int,
) (*dto.PaginatedPostHistoryResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if userID != requesterID {
		return nil, apperror.ErrForbidden
	}

	postHistories, total, err := p.postHistoryRepository.GetByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	response := dto.PaginatedPostHistoryResponse{
		PostHistories: dto.FromPostHistories(postHistories),
		Pagination: dto.Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	}

	return &response, nil
}

func (p *postHistoryService) DeletePostHistoryByID(id string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return p.postHistoryRepository.DeleteByID(ctx, id)
}
