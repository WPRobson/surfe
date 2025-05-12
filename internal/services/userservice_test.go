package services

import (
	"testing"
	"time"

	"surfe/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func TestGetUserByID(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		userID        int
		mockUser      *models.User
		expected      *models.User
		expectedError bool
	}{
		{
			name:   "user found",
			userID: 1,
			mockUser: &models.User{
				ID:        1,
				Name:      "John Doe",
				CreatedAt: now,
			},
			expected: &models.User{
				ID:        1,
				Name:      "John Doe",
				CreatedAt: now,
			},
			expectedError: false,
		},
		{
			name:          "user not found",
			userID:        999,
			mockUser:      nil,
			expected:      nil,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockActionRepo := new(MockActionRepository)
			mockUserRepo.On("GetByID", tt.userID).Return(tt.mockUser, nil)

			service := NewUserService(mockUserRepo, mockActionRepo)
			result, err := service.GetUserByID(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserActionCount(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		userID        int
		mockActions   []models.Action
		expected      int
		expectedError bool
	}{
		{
			name:   "user has actions",
			userID: 1,
			mockActions: []models.Action{
				{ID: 1, Type: "LOGIN", UserID: 1, CreatedAt: now},
				{ID: 2, Type: "VIEW_PROFILE", UserID: 1, CreatedAt: now},
				{ID: 3, Type: "LOGOUT", UserID: 1, CreatedAt: now},
			},
			expected:      3,
			expectedError: false,
		},
		{
			name:          "user has no actions",
			userID:        2,
			mockActions:   []models.Action{},
			expected:      0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockActionRepo := new(MockActionRepository)
			mockActionRepo.On("GetByUserID", tt.userID).Return(tt.mockActions, nil)

			service := NewUserService(mockUserRepo, mockActionRepo)
			result, err := service.GetUserActionCount(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
			mockActionRepo.AssertExpectations(t)
		})
	}
}
