package repository

import (
	"encoding/json"
	"os"
	"sort"
	"surfe/internal/models"
	"sync"
)

type actionRepository struct {
	actions []models.Action
	mu      sync.RWMutex // No exactly needed for this usage, as there is no data writes.
}

func NewActionRepository() (ActionRepository, error) {
	repo := &actionRepository{}
	if err := repo.loadData(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *actionRepository) loadData() error {
	file, err := os.Open("../../actions.json") // Adjust the path if the file is in a different location
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&r.actions)
}

func (r *actionRepository) GetByUserID(userID int) ([]models.Action, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userActions []models.Action
	for _, action := range r.actions {
		if action.UserID == userID {
			userActions = append(userActions, action)
		}
	}
	return userActions, nil
}

func (r *actionRepository) GetAll() ([]models.Action, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	actions := make([]models.Action, len(r.actions))
	copy(actions, r.actions)
	return actions, nil
}

func (r *actionRepository) GetNextActions(actionType string) (map[string]int, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userActions := make(map[int][]models.Action)
	for _, a := range r.actions {
		userActions[a.UserID] = append(userActions[a.UserID], a)
	}

	counts := make(map[string]int)
	total := 0

	for _, actions := range userActions {

		sort.Slice(actions, func(i, j int) bool {
			return actions[i].CreatedAt.Before(actions[j].CreatedAt)
		})

		for i := 0; i < len(actions)-1; i++ {
			if actions[i].Type == actionType {
				next := actions[i+1].Type
				counts[next]++
				total++
			}
		}
	}
	return counts, total, nil
}

func (r *actionRepository) GetReferrals() (map[int][]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	referrals := make(map[int][]int)
	for _, action := range r.actions {
		if action.Type == "REFER_USER" {
			referrals[action.UserID] = append(referrals[action.UserID], action.TargetUser)
		}
	}
	return referrals, nil
}
