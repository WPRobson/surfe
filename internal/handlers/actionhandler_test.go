package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockActionService is a mock implementation of services.ActionService
type MockActionService struct {
	mock.Mock
}

func (m *MockActionService) GetNextActionProbabilities(actionType string) (map[string]float64, error) {
	args := m.Called(actionType)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *MockActionService) GetReferralIndex() (map[int]int, error) {
	args := m.Called()
	return args.Get(0).(map[int]int), args.Error(1)
}

func TestGetNextActionProbabilities(t *testing.T) {
	tests := []struct {
		name           string
		actionType     string
		mockResponse   map[string]float64
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:       "successful response",
			actionType: "LOGIN",
			mockResponse: map[string]float64{
				"VIEW_PROFILE": 0.75,
				"REFER_USER":   0.25,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"VIEW_PROFILE": 0.75,
				"REFER_USER":   0.25,
			},
		},
		{
			name:           "service error",
			actionType:     "INVALID",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/actions/:type/next")
			c.SetParamNames("type")
			c.SetParamValues(tt.actionType)

			// Mock service
			mockService := new(MockActionService)
			mockService.On("GetNextActionProbabilities", tt.actionType).Return(tt.mockResponse, tt.mockError)

			// Create handler
			h := NewActionHandler(mockService)

			// Test
			err := h.GetNextActionProbabilities(c)

			// Assertions
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

func TestGetReferralIndex(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   map[int]int
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful response",
			mockResponse: map[int]int{
				1: 3,
				2: 2,
				3: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"1": float64(3),
				"2": float64(2),
				"3": float64(1),
			},
		},
		{
			name:           "service error",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Internal server error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Mock service
			mockService := new(MockActionService)
			mockService.On("GetReferralIndex").Return(tt.mockResponse, tt.mockError)

			// Create handler
			h := NewActionHandler(mockService)

			// Test
			err := h.GetReferralIndex(c)

			// Assertions
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
