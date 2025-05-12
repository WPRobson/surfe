package services

import (
	"testing"
	"time"

	"surfe/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockActionRepository is a mock implementation of repository.ActionRepository
type MockActionRepository struct {
	mock.Mock
}

func (m *MockActionRepository) GetByUserID(userID int) ([]models.Action, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Action), args.Error(1)
}

func (m *MockActionRepository) GetAll() ([]models.Action, error) {
	args := m.Called()
	return args.Get(0).([]models.Action), args.Error(1)
}

func (m *MockActionRepository) GetNextActions(actionType string) (map[string]int, int, error) {
	args := m.Called(actionType)
	return args.Get(0).(map[string]int), args.Get(1).(int), args.Error(2)
}

func (m *MockActionRepository) GetReferrals() (map[int][]int, error) {
	args := m.Called()
	return args.Get(0).(map[int][]int), args.Error(1)
}

func TestGetNextActionProbabilities(t *testing.T) {
	tests := []struct {
		name          string
		actionType    string
		nextActions   map[string]int
		total         int
		expected      map[string]float64
		expectedError bool
	}{
		{
			name:       "successful probability calculation",
			actionType: "LOGIN",
			nextActions: map[string]int{
				"VIEW_PROFILE": 3,
				"LOGOUT":       1,
			},
			total: 4,
			expected: map[string]float64{
				"VIEW_PROFILE": 0.75,
				"LOGOUT":       0.25,
			},
			expectedError: false,
		},
		{
			name:          "empty next actions",
			actionType:    "LOGOUT",
			nextActions:   map[string]int{},
			total:         0,
			expected:      map[string]float64{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockActionRepository)
			mockRepo.On("GetNextActions", tt.actionType).Return(tt.nextActions, tt.total, nil)

			service := NewActionService(mockRepo)
			result, err := service.GetNextActionProbabilities(tt.actionType)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetReferralIndex(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		actions       []models.Action
		expected      map[int]int
		expectedError bool
	}{
		{
			name: "simple referral chain",
			actions: []models.Action{
				{ID: 1, Type: "REFER_USER", UserID: 1, TargetUser: 2, CreatedAt: now},
				{ID: 2, Type: "REFER_USER", UserID: 2, TargetUser: 3, CreatedAt: now},
				{ID: 3, Type: "REFER_USER", UserID: 3, TargetUser: 4, CreatedAt: now},
			},
			expected: map[int]int{
				1: 3,
				2: 2,
				3: 1,
				4: 0,
			},
			expectedError: false,
		},
		{
			name: "multiple referrals from same user",
			actions: []models.Action{
				{ID: 1, Type: "REFER_USER", UserID: 1, TargetUser: 2, CreatedAt: now},
				{ID: 2, Type: "REFER_USER", UserID: 1, TargetUser: 3, CreatedAt: now},
				{ID: 3, Type: "REFER_USER", UserID: 1, TargetUser: 4, CreatedAt: now},
			},
			expected: map[int]int{
				1: 3,
				2: 0,
				3: 0,
				4: 0,
			},
			expectedError: false,
		},
		{
			name:          "no referrals",
			actions:       []models.Action{},
			expected:      map[int]int{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockActionRepository)
			mockRepo.On("GetAll").Return(tt.actions, nil)

			service := NewActionService(mockRepo)
			result, err := service.GetReferralIndex()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
