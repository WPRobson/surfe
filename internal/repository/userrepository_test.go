package repository

import (
	"io"
	"os"
	"testing"
	"time"

	"surfe/internal/models"

	"github.com/stretchr/testify/assert"
)

func setupUserTestFile(t *testing.T) (string, func()) {
	tmpFile, err := os.CreateTemp("", "users-*.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(tmpFile.Name(), []byte(`[
		{"id": 1, "name": "John Doe", "createdAt": "2024-03-11T20:00:00Z"},
		{"id": 2, "name": "Jane Smith", "createdAt": "2024-03-11T20:00:00Z"}
	]`), 0644); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	originalFile := "users.json"
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

func TestUserRepository_GetByID(t *testing.T) {
	filePath, cleanup := setupUserTestFile(t)
	defer cleanup()

	tests := []struct {
		name          string
		userID        int
		expected      *models.User
		expectedError bool
	}{
		{
			name:   "user found",
			userID: 1,
			expected: &models.User{
				ID:        1,
				Name:      "John Doe",
				CreatedAt: time.Date(2024, 3, 11, 20, 0, 0, 0, time.UTC),
			},
			expectedError: false,
		},
		{
			name:          "user not found",
			userID:        999,
			expected:      nil,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewUserRepository(filePath)
			if err != nil {
				t.Fatal(err)
			}

			result, err := repo.GetByID(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUserRepository_GetAll(t *testing.T) {
	filePath, cleanup := setupUserTestFile(t)
	defer cleanup()

	expectedUsers := []models.User{
		{
			ID:        1,
			Name:      "John Doe",
			CreatedAt: time.Date(2024, 3, 11, 20, 0, 0, 0, time.UTC),
		},
		{
			ID:        2,
			Name:      "Jane Smith",
			CreatedAt: time.Date(2024, 3, 11, 20, 0, 0, 0, time.UTC),
		},
	}

	t.Run("get all users", func(t *testing.T) {
		repo, err := NewUserRepository(filePath)
		if err != nil {
			t.Fatal(err)
		}

		result, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, result)
	})
}
