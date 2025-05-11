package services

import (
	"math"
	"surfe/internal/models"
	"surfe/internal/repository"
)

type actionService struct {
	actionRepo repository.ActionRepository
}

type ReferralGraph map[int][]int

func NewActionService(actionRepo repository.ActionRepository) ActionService {
	return &actionService{
		actionRepo: actionRepo,
	}
}

func (s *actionService) GetNextActionProbabilities(actionType string) (map[string]float64, error) {
	nextActions, total, err := s.actionRepo.GetNextActions(actionType)
	if err != nil {
		return nil, err
	}

	nextActionProbabilities := make(map[string]float64)

	for actionType, count := range nextActions {
		probability := float64(count) / float64(total)
		nextActionProbabilities[actionType] = math.Round(probability*100) / 100
	}

	return nextActionProbabilities, nil
}

func (s *actionService) GetReferralIndex() (map[int]int, error) {
	actions, err := s.actionRepo.GetAll()
	if err != nil {
		return nil, err
	}
	graph := buildReferralGraph(actions)
	referralIndex := make(map[int]int)
	processed := make(map[int]bool)

	var dfs func(userID int) int
	dfs = func(userID int) int {
		if count, found := referralIndex[userID]; found {
			return count
		}

		count := 0
		processed[userID] = true

		for _, referredUser := range graph[userID] {
			if !processed[referredUser] {
				count += 1 + dfs(referredUser)
				continue
			}
			count += 1 + referralIndex[referredUser]
		}

		referralIndex[userID] = count
		return count
	}
	for userID := range graph {
		if !processed[userID] {
			dfs(userID)
		}
	}

	return referralIndex, nil
}

func buildReferralGraph(actions []models.Action) ReferralGraph {
	graph := make(ReferralGraph)
	for _, action := range actions {
		if action.Type == "REFER_USER" {
			graph[action.UserID] = append(graph[action.UserID], action.TargetUser)
		}
	}
	return graph
}

