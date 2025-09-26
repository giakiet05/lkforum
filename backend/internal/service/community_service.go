package service

import "github.com/giakiet05/lkforum/internal/repo"

type CommunityService interface {
}

type communityService struct {
	communityRepo *repo.CommunityRepo
}

func NewCommunityService(communityRepo *repo.CommunityRepo) CommunityService {
	return &communityService{communityRepo: communityRepo}
}
