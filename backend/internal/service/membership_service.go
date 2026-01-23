package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MembershipService interface {
	CreateMembership(req *dto.CreateMembershipRequest, userID string) (*model.Membership, error)
	GetMembershipByID(membershipID string) (*model.Membership, error)
	GetMembershipsByUserID(userID string) ([]*model.Membership, error)
	GetAllMemberships(page int, pageSize int) (*dto.PaginatedMembershipsResponse, error)
	GetMembershipByCommunityID(communityID string, page int, pageSize int) (*dto.PaginatedMembershipsResponse, error)
	DeleteMembership(req *dto.DeleteMembershipRequest, userID string) error
	KickMember(communityID string, targetUserID string, modUserID string) error

	// Pending member approval
	GetPendingMembers(communityID string, requesterID string, page int, pageSize int) (*dto.PaginatedMembershipsResponse, error)
	GetApprovedMembers(communityID string, requesterID string, page int, pageSize int) (*dto.PaginatedMembershipsResponse, error)
	ApproveMember(communityID string, membershipID string, requesterID string) error
	RejectMember(communityID string, membershipID string, requesterID string) error

	GetMembersCount(communityID string) (int64, error)
	increaseMembersCount(communityID string) error
	decreaseMembersCount(communityID string) error
	ensureMembersCountExists(communityID string) (string, error)

	StartRedisToMongoMembershipSync()
	syncMemberCounts() error
}

type membershipService struct {
	membershipRepo repo.MembershipRepo
	communityRepo  repo.CommunityRepo
	redisClient    *redis.Client
}

func NewMembershipService(membershipRepo repo.MembershipRepo, communityRepo repo.CommunityRepo, redisClient *redis.Client) MembershipService {
	svc := &membershipService{membershipRepo: membershipRepo, communityRepo: communityRepo, redisClient: redisClient}
	svc.StartRedisToMongoMembershipSync()
	return svc
}

func (m *membershipService) CreateMembership(req *dto.CreateMembershipRequest, userID string) (*model.Membership, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if userID != req.UserID {
		return nil, apperror.ErrForbidden
	}

	userObjectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, err
	}

	communityObjectID, err := primitive.ObjectIDFromHex(req.CommunityID)
	if err != nil {
		return nil, err
	}

	// Check if user already has a membership (pending or approved)
	existingMembership, err := m.membershipRepo.GetByUserIDAndCommunityID(ctx, req.UserID, req.CommunityID)
	if err != nil {
		return nil, err
	}
	if existingMembership != nil {
		if existingMembership.Status == model.MembershipStatusPending {
			return nil, apperror.NewError(nil, "MEMBERSHIP_PENDING", "Bạn đã gửi yêu cầu tham gia và đang chờ duyệt")
		}
		return nil, apperror.ErrAlreadyMember
	}

	// Check if community exists and get its settings
	community, err := m.communityRepo.GetByID(ctx, req.CommunityID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrCommunityNotFound
		}
		return nil, err
	}

	// Debug log
	log.Printf("🔍 CreateMembership - Community: %s, JoinRequireApproval: %v", community.Name, community.Setting.JoinRequireApproval)

	// Determine membership status based on community settings
	status := model.MembershipStatusApproved
	if community.Setting.JoinRequireApproval {
		status = model.MembershipStatusPending
		log.Printf("✅ Setting status to PENDING for user %s", req.UserID)
	} else {
		log.Printf("✅ Setting status to APPROVED for user %s", req.UserID)
	}

	membership := &model.Membership{
		UserID:      userObjectID,
		CommunityID: communityObjectID,
		Status:      status,
		CreatedAt:   time.Now(),
	}

	membership, err = m.membershipRepo.Create(ctx, membership)
	if err != nil {
		return nil, err
	}

	// Only increase member count if approved immediately
	if status == model.MembershipStatusApproved {
		err = m.increaseMembersCount(req.CommunityID)
		if err != nil {
			return nil, err
		}
	}

	return membership, nil
}

func (m *membershipService) GetMembershipByID(membershipID string) (*model.Membership, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	membership, err := m.membershipRepo.GetByID(ctx, membershipID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrMembershipNotFound
		}
		return nil, err
	}

	return membership, nil
}

func (m *membershipService) GetMembershipsByUserID(userID string) ([]*model.Membership, error) {
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
		return apperror.ErrForbidden
	}

	// Find membership by userID and communityID first
	memberships, err := m.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	var membershipID string
	for _, membership := range memberships {
		if membership.CommunityID.Hex() == req.CommunityID {
			membershipID = membership.ID.Hex()
			break
		}
	}

	if membershipID == "" {
		return fmt.Errorf("membership not found")
	}

	// Delete membership by ID
	err = m.membershipRepo.Delete(ctx, membershipID)
	if err != nil {
		return err
	}

	err = m.decreaseMembersCount(req.CommunityID)
	if err != nil {
		return err
	}

	return nil
}

// KickMember allows moderator/creator to remove a member from community
func (m *membershipService) KickMember(communityID string, targetUserID string, modUserID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if mod user has permission (is creator or moderator)
	community, err := m.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return err
	}

	isCreator := community.CreateByID.Hex() == modUserID
	isModerator := false
	for _, mod := range community.Moderators {
		if mod.UserID.Hex() == modUserID {
			isModerator = true
			break
		}
	}

	if !isCreator && !isModerator {
		return apperror.ErrForbidden
	}

	// Cannot kick creator
	if community.CreateByID.Hex() == targetUserID {
		return apperror.NewError(nil, "CANNOT_KICK_CREATOR", "Cannot kick the community creator")
	}

	// Cannot kick other moderators unless you're creator
	if !isCreator {
		for _, mod := range community.Moderators {
			if mod.UserID.Hex() == targetUserID {
				return apperror.NewError(nil, "CANNOT_KICK_MODERATOR", "Only creator can kick moderators")
			}
		}
	}

	// Find and delete membership
	memberships, err := m.membershipRepo.GetByUserID(ctx, targetUserID)
	if err != nil {
		return err
	}

	var membershipID string
	for _, membership := range memberships {
		if membership.CommunityID.Hex() == communityID {
			membershipID = membership.ID.Hex()
			break
		}
	}

	if membershipID == "" {
		return apperror.NewError(nil, "MEMBERSHIP_NOT_FOUND", "User is not a member of this community")
	}

	err = m.membershipRepo.Delete(ctx, membershipID)
	if err != nil {
		return err
	}

	return m.decreaseMembersCount(communityID)
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

	key := fmt.Sprintf(config.RedisMembersCountKey, communityID)

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
	// Tạm thời set cứng 5 min
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for range ticker.C {
			if err := m.syncMemberCounts(); err != nil {
				log.Printf("Redis→Mongo membership sync failed: %v", err)
			} else {
				log.Println("Redis→Mongo membership sync completed successfully")
			}
		}
	}()
}

func (m *membershipService) syncMemberCounts() error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	iter := m.redisClient.Scan(ctx, 0, fmt.Sprintf(config.RedisMembersCountKey, "*"), 0).Iterator()
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

// GetPendingMembers returns all pending membership requests for a community
func (m *membershipService) GetPendingMembers(communityID string, requesterID string, page int, pageSize int) (*dto.PaginatedMembershipsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is creator or moderator
	community, err := m.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return nil, err
	}

	isAuthorized := community.CreateByID.Hex() == requesterID
	if !isAuthorized {
		for _, mod := range community.Moderators {
			if mod.UserID.Hex() == requesterID {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		return nil, apperror.ErrForbidden
	}

	memberships, total, err := m.membershipRepo.GetPendingByCommunityID(ctx, communityID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedMembershipsResponse{
		Memberships: memberships,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

// GetApprovedMembers returns all approved members for a community
func (m *membershipService) GetApprovedMembers(communityID string, requesterID string, page int, pageSize int) (*dto.PaginatedMembershipsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is creator or moderator
	community, err := m.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return nil, err
	}

	isAuthorized := community.CreateByID.Hex() == requesterID
	if !isAuthorized {
		for _, mod := range community.Moderators {
			if mod.UserID.Hex() == requesterID {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		return nil, apperror.ErrForbidden
	}

	memberships, total, err := m.membershipRepo.GetApprovedByCommunityID(ctx, communityID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedMembershipsResponse{
		Memberships: memberships,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

// ApproveMember approves a pending membership request
func (m *membershipService) ApproveMember(communityID string, membershipID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is creator or moderator
	community, err := m.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return err
	}

	isAuthorized := community.CreateByID.Hex() == requesterID
	if !isAuthorized {
		for _, mod := range community.Moderators {
			if mod.UserID.Hex() == requesterID {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		return apperror.ErrForbidden
	}

	// Get membership and verify it belongs to this community
	membership, err := m.membershipRepo.GetByID(ctx, membershipID)
	if err != nil {
		return err
	}
	if membership == nil {
		return apperror.ErrMembershipNotFound
	}

	if membership.CommunityID.Hex() != communityID {
		return apperror.ErrForbidden
	}

	if membership.Status != model.MembershipStatusPending {
		return apperror.ErrMembershipAlreadyProcessed
	}

	// Update status to approved
	err = m.membershipRepo.UpdateStatus(ctx, membershipID, model.MembershipStatusApproved)
	if err != nil {
		return err
	}

	// Increase member count
	return m.increaseMembersCount(communityID)
}

// RejectMember rejects a pending membership request
func (m *membershipService) RejectMember(communityID string, membershipID string, requesterID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check if requester is creator or moderator
	community, err := m.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return err
	}

	isAuthorized := community.CreateByID.Hex() == requesterID
	if !isAuthorized {
		for _, mod := range community.Moderators {
			if mod.UserID.Hex() == requesterID {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		return apperror.ErrForbidden
	}

	// Get membership and verify it belongs to this community
	membership, err := m.membershipRepo.GetByID(ctx, membershipID)
	if err != nil {
		return err
	}
	if membership == nil {
		return apperror.ErrMembershipNotFound
	}

	if membership.CommunityID.Hex() != communityID {
		return apperror.ErrForbidden
	}

	if membership.Status != model.MembershipStatusPending {
		return apperror.ErrMembershipAlreadyProcessed
	}

	// Delete the membership request
	return m.membershipRepo.Delete(ctx, membershipID)
}
