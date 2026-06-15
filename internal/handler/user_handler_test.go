package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ntuc/internal/domain"
	"ntuc/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		input          *domain.CreateUserRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid request",
			input: &domain.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "missing name",
			input: &domain.CreateUserRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "short password",
			input: &domain.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "short",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			svc := service.NewUserService(nil, nil)
			handler := NewUserHandler(svc, nil)

			// Create request
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
			w := httptest.NewRecorder()

			// Execute
			handler.CreateUser(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "invalid uuid",
			userID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "valid uuid format",
			userID:         uuid.New().String(),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewUserService(nil, nil)
			handler := NewUserHandler(svc, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.userID, nil)
			req.SetPathValue("id", tt.userID)
			w := httptest.NewRecorder()

			handler.GetUser(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
