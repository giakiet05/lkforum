package service

import (
	"fmt"
	"log"
	"strings"
	"time"

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

	GetMembersCount(communityID string) (int64, error)
	increaseMembersCount(communityID string) error
	decreaseMembersCount(communityID string) error
	ensureMembersCountExists(communityID string) (string, error)

	StartRedisToMongoMembershipSync()
	syncMemberCounts() error
}

type membershipService struct {
	membershipRepo repo.MembershipRepo
	redisClient    *redis.Client
}

func NewMembershipService(membershipRepo repo.MembershipRepo, redisClient *redis.Client) MembershipService {
	svc := &membershipService{membershipRepo: membershipRepo, redisClient: redisClient}
	svc.StartRedisToMongoMembershipSync()
	return svc
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

	membership, err = m.membershipRepo.Create(ctx, membership)
	if err != nil {
		return nil, err
	}

	err = m.increaseMembersCount(userID)
	if err != nil {
		return nil, err
	}

	return membership, nil
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

	err := m.membershipRepo.Delete(ctx, req.CommunityID)
	if err != nil {
		return err
	}

	err = m.decreaseMembersCount(userID)
	if err != nil {
		return err
	}

	return nil
}

func (m *membershipService) increaseMembersCount(communityID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	key, err := m.ensureMembersCountExists(communityID)
	if err != nil {
		return err
	}

	if err := m.redisClient.Incr(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func (m *membershipService) decreaseMembersCount(communityID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	key, err := m.ensureMembersCountExists(communityID)
	if err != nil {
		return err
	}

	if err := m.redisClient.Decr(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func (m *membershipService) GetMembersCount(communityID string) (int64, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	key, err := m.ensureMembersCountExists(communityID)
	if err != nil {
		return 0, err
	}
	count, err := m.redisClient.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *membershipService) ensureMembersCountExists(communityID string) (string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	key := fmt.Sprintf("community:%s:member_count", communityID)

	exists, err := m.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return "", err
	}

	if exists == 0 {
		dbCount, err := m.membershipRepo.CountMembersByCommunityID(ctx, communityID)
		if err != nil {
			return "", err
		}

		if err := m.redisClient.Set(ctx, key, dbCount, 0).Err(); err != nil {
			return "", err
		}
	}

	return key, nil
}

func (m *membershipService) StartRedisToMongoMembershipSync() {
	// Tạm thời set cứng 1 min
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker.C {
			if err := m.syncMemberCounts(); err != nil {
				log.Printf("⚠️ Redis→Mongo membership sync failed: %v", err)
			} else {
				log.Println("✅ Redis→Mongo membership sync completed successfully")
			}
		}
	}()
}

func (m *membershipService) syncMemberCounts() error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	iter := m.redisClient.Scan(ctx, 0, "community:*:member_count", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		// Key format: community:<id>:member_count
		parts := strings.Split(key, ":")
		if len(parts) < 3 {
			continue
		}
		communityID := parts[1]

		val, err := m.redisClient.Get(ctx, key).Int64()
		if err != nil {
			log.Printf("failed to read %s: %v", key, err)
			continue
		}

		// Update MongoDB
		err = m.membershipRepo.UpdateCommunityMemberCount(ctx, communityID, val)
		if err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("redis scan failed: %w", err)
	}

	return nil
}
