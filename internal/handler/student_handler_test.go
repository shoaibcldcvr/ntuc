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

func TestStudentHandler_CreateStudent(t *testing.T) {
	tests := []struct {
		name           string
		input          *domain.CreateStudentRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid request",
			input: &domain.CreateStudentRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Phone:     "9876543210",
				RollNo:    "STU001",
				Class:     "10A",
				GPA:       3.5,
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "missing first name",
			input: &domain.CreateStudentRequest{
				LastName:  "Doe",
				Email:     "john@example.com",
				Phone:     "9876543210",
				RollNo:    "STU001",
				Class:     "10A",
				GPA:       3.5,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid GPA",
			input: &domain.CreateStudentRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Phone:     "9876543210",
				RollNo:    "STU001",
				Class:     "10A",
				GPA:       5.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewStudentService(nil, nil)
			handler := NewStudentHandler(svc, nil)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.CreateStudent(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestStudentHandler_GetStudent(t *testing.T) {
	tests := []struct {
		name           string
		studentID      string
		expectedStatus int
	}{
		{
			name:           "invalid uuid",
			studentID:      "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "valid uuid format",
			studentID:      uuid.New().String(),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewStudentService(nil, nil)
			handler := NewStudentHandler(svc, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/students/"+tt.studentID, nil)
			req.SetPathValue("id", tt.studentID)
			w := httptest.NewRecorder()

			handler.GetStudent(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
