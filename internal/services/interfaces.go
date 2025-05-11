package services

import "surfe/internal/models"

type UserService interface {
	GetUserByID(id int) (*models.User, error)
	GetUserActionCount(userID int) (int, error)
}

type ActionService interface {
	GetNextActionProbabilities(actionType string) (map[string]float64, error)
	GetReferralIndex() (map[int]int, error)
}
