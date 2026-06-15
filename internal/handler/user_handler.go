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

type UserHandler struct {
	service   *service.UserService
	logger    *logger.Logger
	validator *validator.Validate
}

func NewUserHandler(service *service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
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
	result, err := h.service.CreateUser(ctx, &req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()
	result, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	results, err := h.service.ListUsers(ctx)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, results)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req domain.UpdateUserRequest
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
	result, err := h.service.UpdateUser(ctx, id, &req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx := r.Context()
	err := h.service.DeleteUser(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		response.ErrorFromAppError(w, appErr)
		return
	}

	h.logger.Errorf("unexpected error: %v", err)
	response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}
