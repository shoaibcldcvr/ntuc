package handler

import (
	"encoding/json"
	"net/http"

	"ntuc/internal/domain"
	"ntuc/internal/service"
	"ntuc/pkg/errors"
	"ntuc/pkg/logger"
	"ntuc/pkg/response"

	"github.com/go-playground/validator/v10"
)

type StudentHandler struct {
	service   *service.StudentService
	logger    *logger.Logger
	validator *validator.Validate
}

func NewStudentHandler(service *service.StudentService, logger *logger.Logger) *StudentHandler {
	return &StudentHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request: %v", err)
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request payload")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.Errorf("validation error: %v", err)
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	ctx := r.Context()
	result, err := h.service.CreateStudent(ctx, &req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()
	result, err := h.service.GetStudent(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *StudentHandler) ListStudents(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")

	ctx := r.Context()
	var results []*domain.StudentResponse
	var err error

	if class != "" {
		results, err = h.service.ListStudentsByClass(ctx, class)
	} else {
		results, err = h.service.ListStudents(ctx)
	}

	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, results)
}

func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req domain.UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request: %v", err)
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request payload")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.Errorf("validation error: %v", err)
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	ctx := r.Context()
	result, err := h.service.UpdateStudent(ctx, id, &req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()
	err := h.service.DeleteStudent(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StudentHandler) handleServiceError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		response.ErrorFromAppError(w, appErr)
		return
	}

	h.logger.Errorf("unexpected error: %v", err)
	response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}
