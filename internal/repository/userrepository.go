package repository

import (
	"encoding/json"
	"os"
	"surfe/internal/models"
	"sync"
)

type userRepository struct {
	users []models.User
	mu    sync.RWMutex // No exactly needed for this usage, as there is no data writes.
}

func NewUserRepository() (UserRepository, error) {
	repo := &userRepository{}
	if err := repo.loadData(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *userRepository) loadData() error {
	file, err := os.Open("../../users.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&r.users)
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]models.User, len(r.users))
	copy(users, r.users)
	return users, nil
}
