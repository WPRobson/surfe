package repository

import (
	"encoding/json"
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

	return json.NewDecoder(file).Decode(&r.users)
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
	users := make([]models.User, len(r.users))
	copy(users, r.users)
	return users, nil
}
