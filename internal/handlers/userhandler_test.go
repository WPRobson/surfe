package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"surfe/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserActionCount(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func TestGetUserByID(t *testing.T) {
	fixedTime := time.Date(2024, 3, 11, 20, 0, 0, 0, time.UTC)
	tests := []struct {
		name           string
		userID         string
		mockUser       *models.User
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "user found",
			userID: "1",
			mockUser: &models.User{
				ID:        1,
				Name:      "John Doe",
				CreatedAt: fixedTime,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":        float64(1),
				"name":      "John Doe",
				"createdAt": fixedTime.Format(time.RFC3339),
			},
		},
		{
			name:           "user not found",
			userID:         "999",
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "User not found",
			},
		},
		{
			name:           "invalid user ID",
			userID:         "invalid",
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid user ID",
			},
		},
		{
			name:           "service error",
			userID:         "1",
			mockUser:       nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			mockService := new(MockUserService)
			if tt.userID != "invalid" {
				id, _ := strconv.Atoi(tt.userID)
				mockService.On("GetUserByID", id).Return(tt.mockUser, tt.mockError)
			}

			h := NewUserHandler(mockService)

			err := h.GetUserByID(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetUserActionCount(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockCount      int
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "successful response",
			userID:         "1",
			mockCount:      42,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"count": float64(42),
			},
		},
		{
			name:           "invalid user ID",
			userID:         "invalid",
			mockCount:      0,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid user ID",
			},
		},
		{
			name:           "service error",
			userID:         "1",
			mockCount:      0,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users/:id/actions/count")
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			mockService := new(MockUserService)
			if tt.userID != "invalid" {
				id, _ := strconv.Atoi(tt.userID)
				mockService.On("GetUserActionCount", id).Return(tt.mockCount, tt.mockError)
			}

			h := NewUserHandler(mockService)

			err := h.GetUserActionCount(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}
