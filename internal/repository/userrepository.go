package repository

import (
	"encoding/json"
	"errors"
	"os"
	"surfe/internal/models"
)

type userRepository struct {
	users []models.User
}

func NewUserRepository(filePath string) (UserRepository, error) {
	repo := &userRepository{}
	if err := repo.loadData(filePath); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *userRepository) loadData(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&r.users); err != nil {
		return errors.New("invalid JSON data")
	}
	return nil
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userRepository) GetAll() ([]models.User, error) {
	users := make([]models.User, len(r.users))
	copy(users, r.users)
	return users, nil
}
