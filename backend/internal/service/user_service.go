package service

import (
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
)

type UserService interface {
	GetAllUsers() ([]*model.User, error)
}

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
func (s *userService) GetAllUsers() ([]*model.User, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()
	return s.userRepo.GetAll(ctx)
}
