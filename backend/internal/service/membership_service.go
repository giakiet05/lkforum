package service

import (
	"fmt"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MembershipService interface {
	CreateMembership(req *dto.CreateMembershipRequest, userID string) (*model.Membership, error)
	GetMembershipByID(membershipID string) (*model.Membership, error)
	GetMembershipByUserID(userID string) ([]model.Membership, error)
	GetAllMemberships(page int, pageSize int) (*dto.PaginatedMembershipsResponse, error)
	GetMembershipByCommunityID(communityID string, page int, pageSize int) (*dto.PaginatedMembershipsResponse, error)
	DeleteMembership(req *dto.DeleteMembershipRequest, userID string) error
}

type membershipService struct {
	membershipRepo repo.MembershipRepo
	redisClient    *redis.Client
}

func NewMembershipService(membershipRepo repo.MembershipRepo, redisClient *redis.Client) MembershipService {
	return &membershipService{membershipRepo: membershipRepo, redisClient: redisClient}
}

func (m *membershipService) CreateMembership(req *dto.CreateMembershipRequest, userID string) (*model.Membership, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if userID != req.UserID {
		return nil, fmt.Errorf("unathorize to create membership for this user id")
	}

	userObjectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, err
	}

	communityObjectID, err := primitive.ObjectIDFromHex(req.CommunityID)
	if err != nil {
		return nil, err
	}

	membership := &model.Membership{
		UserID:      userObjectID,
		CommunityID: communityObjectID,
	}

	return m.membershipRepo.Create(ctx, membership)
}

func (m *membershipService) GetMembershipByID(membershipID string) (*model.Membership, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return m.membershipRepo.GetByID(ctx, membershipID)
}

func (m *membershipService) GetMembershipByUserID(userID string) ([]model.Membership, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return m.membershipRepo.GetByUserID(ctx, userID)
}

func (m *membershipService) GetAllMemberships(page int, pageSize int) (*dto.PaginatedMembershipsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	memberships, total, err := m.membershipRepo.GetAllPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	response := &dto.PaginatedMembershipsResponse{
		Memberships: memberships,
		Pagination: dto.Pagination{
			Total: total,
			Page:  page,
		},
	}

	return response, nil
}

func (m *membershipService) GetMembershipByCommunityID(communityID string, page int, pageSize int) (*dto.PaginatedMembershipsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	memberships, total, err := m.membershipRepo.GetByCommunityIDPaginated(ctx, communityID, page, pageSize)
	if err != nil {
		return nil, err
	}

	response := &dto.PaginatedMembershipsResponse{
		Memberships: memberships,
		Pagination: dto.Pagination{
			Total: total,
			Page:  page,
		},
	}

	return response, nil
}

func (m *membershipService) DeleteMembership(req *dto.DeleteMembershipRequest, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if userID != req.UserID {
		return fmt.Errorf("unathorize user id")
	}

	return m.membershipRepo.Delete(ctx, req.CommunityID)
}
