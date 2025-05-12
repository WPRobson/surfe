package repository

import (
	"io"
	"os"
	"testing"
	"time"

	"surfe/internal/models"

	"github.com/stretchr/testify/assert"
)

func setupActionTestFile(t *testing.T) (string, func()) {
	tmpFile, err := os.CreateTemp("", "actions-*.json")
	if err != nil {
		t.Fatal(err)
	}

	// Create test data with various action types and a referral chain
	testData := `[
		{
			"id": 1,
			"type": "LOGIN",
			"userId": 1,
			"targetUser": 0,
			"createdAt": "2024-03-11T20:00:00Z"
		},
		{
			"id": 2,
			"type": "VIEW_PROFILE",
			"userId": 1,
			"targetUser": 0,
			"createdAt": "2024-03-11T20:01:00Z"
		},
		{
			"id": 3,
			"type": "REFER_USER",
			"userId": 1,
			"targetUser": 2,
			"createdAt": "2024-03-11T20:02:00Z"
		},
		{
			"id": 4,
			"type": "REFER_USER",
			"userId": 2,
			"targetUser": 3,
			"createdAt": "2024-03-11T20:03:00Z"
		},
		{
			"id": 5,
			"type": "LOGIN",
			"userId": 2,
			"targetUser": 0,
			"createdAt": "2024-03-11T20:04:00Z"
		},
		{
			"id": 6,
			"type": "REFER_USER",
			"userId": 2,
			"targetUser": 0,
			"createdAt": "2024-03-11T20:05:00Z"
		}
	]`

	if err := os.WriteFile(tmpFile.Name(), []byte(testData), 0644); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	originalFile := "actions.json"
	if _, err := os.Stat(originalFile); err == nil {
		if err := os.Rename(originalFile, originalFile+".bak"); err != nil {
			t.Fatal(err)
		}
	}

	src, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()
	dst, err := os.Create(originalFile)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		t.Fatal(err)
	}
	dst.Close()

	return originalFile, func() {
		os.Remove(originalFile)
		if _, err := os.Stat(originalFile + ".bak"); err == nil {
			os.Rename(originalFile+".bak", originalFile)
		}
	}
}

func TestActionRepository_GetByUserID(t *testing.T) {
	filePath, cleanup := setupActionTestFile(t)
	defer cleanup()

	tests := []struct {
		name          string
		userID        int
		expected      []models.Action
		expectedError bool
	}{
		{
			name:   "user has actions",
			userID: 1,
			expected: []models.Action{
				{
					ID:         1,
					Type:       "LOGIN",
					UserID:     1,
					TargetUser: 0,
					CreatedAt:  time.Date(2024, 3, 11, 20, 0, 0, 0, time.UTC),
				},
				{
					ID:         2,
					Type:       "VIEW_PROFILE",
					UserID:     1,
					TargetUser: 0,
					CreatedAt:  time.Date(2024, 3, 11, 20, 1, 0, 0, time.UTC),
				},
				{
					ID:         3,
					Type:       "REFER_USER",
					UserID:     1,
					TargetUser: 2,
					CreatedAt:  time.Date(2024, 3, 11, 20, 2, 0, 0, time.UTC),
				},
			},
			expectedError: false,
		},
		{
			name:          "user has no actions",
			userID:        999,
			expected:      []models.Action{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewActionRepository(filePath)
			if err != nil {
				t.Fatal(err)
			}

			result, err := repo.GetByUserID(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestActionRepository_GetAll(t *testing.T) {
	filePath, cleanup := setupActionTestFile(t)
	defer cleanup()

	expectedActions := []models.Action{
		{
			ID:         1,
			Type:       "LOGIN",
			UserID:     1,
			TargetUser: 0,
			CreatedAt:  time.Date(2024, 3, 11, 20, 0, 0, 0, time.UTC),
		},
		{
			ID:         2,
			Type:       "VIEW_PROFILE",
			UserID:     1,
			TargetUser: 0,
			CreatedAt:  time.Date(2024, 3, 11, 20, 1, 0, 0, time.UTC),
		},
		{
			ID:         3,
			Type:       "REFER_USER",
			UserID:     1,
			TargetUser: 2,
			CreatedAt:  time.Date(2024, 3, 11, 20, 2, 0, 0, time.UTC),
		},
		{
			ID:         4,
			Type:       "REFER_USER",
			UserID:     2,
			TargetUser: 3,
			CreatedAt:  time.Date(2024, 3, 11, 20, 3, 0, 0, time.UTC),
		},
		{
			ID:         5,
			Type:       "LOGIN",
			UserID:     2,
			TargetUser: 0,
			CreatedAt:  time.Date(2024, 3, 11, 20, 4, 0, 0, time.UTC),
		},
		{
			ID:         6,
			Type:       "REFER_USER",
			UserID:     2,
			TargetUser: 0,
			CreatedAt:  time.Date(2024, 3, 11, 20, 5, 0, 0, time.UTC),
		},
	}

	t.Run("get all actions", func(t *testing.T) {
		repo, err := NewActionRepository(filePath)
		if err != nil {
			t.Fatal(err)
		}

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Equal(t, expectedActions, result)
	})
}

func TestActionRepository_GetNextActions(t *testing.T) {
	filePath, cleanup := setupActionTestFile(t)
	defer cleanup()

	tests := []struct {
		name          string
		actionType    string
		expected      map[string]int
		expectedTotal int
		expectedError bool
	}{
		{
			name:       "get next actions after login",
			actionType: "LOGIN",
			expected: map[string]int{
				"VIEW_PROFILE": 1,
				"REFER_USER":   1,
			},
			expectedTotal: 2,
			expectedError: false,
		},
		{
			name:       "get next actions after view profile",
			actionType: "VIEW_PROFILE",
			expected: map[string]int{
				"REFER_USER": 1,
			},
			expectedTotal: 1,
			expectedError: false,
		},
		{
			name:          "no next actions",
			actionType:    "LOGOUT",
			expected:      map[string]int{},
			expectedTotal: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewActionRepository(filePath)
			if err != nil {
				t.Fatal(err)
			}

			result, total, err := repo.GetNextActions(tt.actionType)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
				assert.Equal(t, tt.expectedTotal, total)
			}
		})
	}
}

func TestActionRepository_GetReferrals(t *testing.T) {
	filePath, cleanup := setupActionTestFile(t)
	defer cleanup()

	expectedReferrals := map[int][]int{
		1: {2},
		2: {3,0},
	}

	t.Run("get referrals", func(t *testing.T) {
		repo, err := NewActionRepository(filePath)
		if err != nil {
			t.Fatal(err)
		}

		result, err := repo.GetReferrals()
		assert.NoError(t, err)
		assert.Equal(t, expectedReferrals, result)
	})
}
