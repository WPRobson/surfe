package repository

import (
	"encoding/json"
	"os"
	"sort"
	"surfe/internal/models"
)

type actionRepository struct {
	actions []models.Action
}

func NewActionRepository(filePath string) (ActionRepository, error) {
	repo := &actionRepository{}
	if err := repo.loadData(filePath); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *actionRepository) loadData(filePath string) error {
	file, err := os.Open(filePath) // Adjust the path if the file is in a different location
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&r.actions)
}

func (r *actionRepository) GetByUserID(userID int) ([]models.Action, error) {
	userActions := []models.Action{}
	for _, action := range r.actions {
		if action.UserID == userID {
			userActions = append(userActions, action)
		}
	}
	return userActions, nil
}

func (r *actionRepository) GetAll() ([]models.Action, error) {
	actions := make([]models.Action, len(r.actions))
	copy(actions, r.actions)
	return actions, nil
}

func (r *actionRepository) GetNextActions(actionType string) (map[string]int, int, error) {
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
	referrals := make(map[int][]int)
	for _, action := range r.actions {
		if action.Type == "REFER_USER" {
			referrals[action.UserID] = append(referrals[action.UserID], action.TargetUser)
		}
	}
	return referrals, nil
}
