package repository

import "surfe/internal/models"

type ActionRepository interface {
	GetByUserID(userID int) ([]models.Action, error)
	GetAll() ([]models.Action, error)
	GetNextActions(actionType string) (map[string]int, int, error)
	GetReferrals() (map[int][]int, error)
}

type UserRepository interface {
	GetByID(id int) (*models.User, error)
	GetAll() ([]models.User, error)
}
