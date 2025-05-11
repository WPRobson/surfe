package services

import (
	"surfe/internal/models"
	"surfe/internal/repository"
)

type userService struct {
	userRepo   repository.UserRepository
	actionRepo repository.ActionRepository
}

func NewUserService(userRepo repository.UserRepository, actionRepo repository.ActionRepository) UserService {
	return &userService{
		userRepo:   userRepo,
		actionRepo: actionRepo,
	}
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetUserActionCount(userID int) (int, error) {
	actions, err := s.actionRepo.GetByUserID(userID)
	if err != nil {
		return 0, err
	}
	return len(actions), nil
}

