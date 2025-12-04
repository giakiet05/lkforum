package service

import (
	"context"
	"time"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
)

type AdminStatsService interface {
	GetPlatformOverview() (*dto.PlatformStatsResponse, error)
	GetUserStats(query *dto.GetUserStatsQuery) (*dto.UserStatsResponse, error)
	GetContentStats(query *dto.GetContentStatsQuery) (*dto.ContentStatsResponse, error)
}

type adminStatsService struct {
	userRepo      repo.UserRepo
	communityRepo repo.CommunityRepo
	postRepo      repo.PostRepo
	commentRepo   repo.CommentRepo
	reportRepo    repo.ReportRepo
}

func NewAdminStatsService(
	userRepo repo.UserRepo,
	communityRepo repo.CommunityRepo,
	postRepo repo.PostRepo,
	commentRepo repo.CommentRepo,
	reportRepo repo.ReportRepo,
) AdminStatsService {
	return &adminStatsService{
		userRepo:      userRepo,
		communityRepo: communityRepo,
		postRepo:      postRepo,
		commentRepo:   commentRepo,
		reportRepo:    reportRepo,
	}
}

func (s *adminStatsService) GetPlatformOverview() (*dto.PlatformStatsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := now.AddDate(0, 0, -7)
	startOfMonth := now.AddDate(0, -1, 0)

	// User stats
	userStats, err := s.getUserStatsData(ctx, startOfToday, startOfWeek, startOfMonth)
	if err != nil {
		return nil, err
	}

	// Community stats
	communityStats, err := s.getCommunityStatsData(ctx, startOfWeek, startOfMonth)
	if err != nil {
		return nil, err
	}

	// Content stats
	contentStats, err := s.getContentStatsData(ctx, startOfToday, startOfWeek)
	if err != nil {
		return nil, err
	}

	// Moderation stats
	moderationStats, err := s.getModerationStatsData(ctx, startOfToday)
	if err != nil {
		return nil, err
	}

	return &dto.PlatformStatsResponse{
		Users:       *userStats,
		Communities: *communityStats,
		Content:     *contentStats,
		Moderation:  *moderationStats,
		UpdatedAt:   now,
	}, nil
}

func (s *adminStatsService) GetUserStats(query *dto.GetUserStatsQuery) (*dto.UserStatsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := now.AddDate(0, 0, -7)
	startOfMonth := now.AddDate(0, -1, 0)

	// Overview stats
	userStats, err := s.getUserStatsData(ctx, startOfToday, startOfWeek, startOfMonth)
	if err != nil {
		return nil, err
	}

	// Growth data based on period
	var growth []dto.UserGrowthData
	switch query.Period {
	case "week":
		growth, err = s.getUserGrowthData(ctx, startOfWeek, "daily")
	case "month":
		growth, err = s.getUserGrowthData(ctx, startOfMonth, "daily")
	default:
		growth, err = s.getUserGrowthData(ctx, startOfWeek, "daily")
	}

	if err != nil {
		return nil, err
	}

	return &dto.UserStatsResponse{
		Overview: *userStats,
		Growth:   growth,
	}, nil
}

func (s *adminStatsService) GetContentStats(query *dto.GetContentStatsQuery) (*dto.ContentStatsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := now.AddDate(0, 0, -7)

	// Overview stats
	contentStats, err := s.getContentStatsData(ctx, startOfToday, startOfWeek)
	if err != nil {
		return nil, err
	}

	// Growth data
	var growth []dto.ContentGrowthData
	switch query.Period {
	case "week":
		growth, err = s.getContentGrowthData(ctx, startOfWeek, "daily")
	case "month":
		growth, err = s.getContentGrowthData(ctx, now.AddDate(0, -1, 0), "daily")
	default:
		growth, err = s.getContentGrowthData(ctx, startOfWeek, "daily")
	}

	if err != nil {
		return nil, err
	}

	return &dto.ContentStatsResponse{
		Overview: *contentStats,
		Growth:   growth,
	}, nil
}

// Helper methods
func (s *adminStatsService) getUserStatsData(ctx context.Context, startOfToday, startOfWeek, startOfMonth time.Time) (*dto.UserStatsData, error) {
	// These methods need to be implemented in the respective repos
	total, _ := s.userRepo.CountTotal(ctx)
	activeToday, _ := s.userRepo.CountActiveAfter(ctx, startOfToday)
	newThisWeek, _ := s.userRepo.CountCreatedAfter(ctx, startOfWeek)
	newThisMonth, _ := s.userRepo.CountCreatedAfter(ctx, startOfMonth)
	banned, _ := s.userRepo.CountBanned(ctx)
	verified, _ := s.userRepo.CountVerified(ctx)

	return &dto.UserStatsData{
		Total:         total,
		ActiveToday:   activeToday,
		NewThisWeek:   newThisWeek,
		NewThisMonth:  newThisMonth,
		Banned:        banned,
		Verified:      verified,
	}, nil
}

func (s *adminStatsService) getCommunityStatsData(ctx context.Context, startOfWeek, startOfMonth time.Time) (*dto.CommunityStatsData, error) {
	total, _ := s.communityRepo.CountTotal(ctx)
	active, _ := s.communityRepo.CountActive(ctx)
	newThisWeek, _ := s.communityRepo.CountCreatedAfter(ctx, startOfWeek)
	newThisMonth, _ := s.communityRepo.CountCreatedAfter(ctx, startOfMonth)
	banned, _ := s.communityRepo.CountBanned(ctx)
	private, _ := s.communityRepo.CountPrivate(ctx)

	return &dto.CommunityStatsData{
		Total:         total,
		Active:        active,
		NewThisWeek:   newThisWeek,
		NewThisMonth:  newThisMonth,
		Banned:        banned,
		Private:       private,
	}, nil
}

func (s *adminStatsService) getContentStatsData(ctx context.Context, startOfToday, startOfWeek time.Time) (*dto.ContentStatsData, error) {
	totalPosts, _ := s.postRepo.CountTotal(ctx)
	totalComments, _ := s.commentRepo.CountTotal(ctx)
	postsToday, _ := s.postRepo.CountCreatedAfter(ctx, startOfToday)
	commentsToday, _ := s.commentRepo.CountCreatedAfter(ctx, startOfToday)
	postsThisWeek, _ := s.postRepo.CountCreatedAfter(ctx, startOfWeek)
	commentsThisWeek, _ := s.commentRepo.CountCreatedAfter(ctx, startOfWeek)

	return &dto.ContentStatsData{
		TotalPosts:       totalPosts,
		TotalComments:    totalComments,
		PostsToday:       postsToday,
		CommentsToday:    commentsToday,
		PostsThisWeek:    postsThisWeek,
		CommentsThisWeek: commentsThisWeek,
	}, nil
}

func (s *adminStatsService) getModerationStatsData(ctx context.Context, startOfToday time.Time) (*dto.ModerationStatsData, error) {
	pendingReports, _ := s.reportRepo.CountPending(ctx)
	reportsToday, _ := s.reportRepo.CountCreatedAfter(ctx, startOfToday)
	postsNeedApproval, _ := s.postRepo.CountPendingApproval(ctx)
	bannedUsers, _ := s.userRepo.CountBanned(ctx)
	bannedCommunities, _ := s.communityRepo.CountBanned(ctx)

	return &dto.ModerationStatsData{
		PendingReports:      pendingReports,
		ReportsToday:        reportsToday,
		PostsNeedApproval:   postsNeedApproval,
		BannedUsers:         bannedUsers,
		BannedCommunities:   bannedCommunities,
	}, nil
}

func (s *adminStatsService) getUserGrowthData(ctx context.Context, startDate time.Time, interval string) ([]dto.UserGrowthData, error) {
	// This would need aggregation pipeline in MongoDB
	// For now, return empty slice
	return []dto.UserGrowthData{}, nil
}

func (s *adminStatsService) getContentGrowthData(ctx context.Context, startDate time.Time, interval string) ([]dto.ContentGrowthData, error) {
	// This would need aggregation pipeline in MongoDB
	// For now, return empty slice
	return []dto.ContentGrowthData{}, nil
}