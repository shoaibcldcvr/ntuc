package response

import (
	"encoding/json"
	"net/http"

	"ntuc/pkg/errors"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{Data: data})
}

func Error(w http.ResponseWriter, statusCode int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

func ErrorFromAppError(w http.ResponseWriter, err *errors.AppError) {
	statusCode := getStatusCode(err.Code)
	Error(w, statusCode, string(err.Code), err.Message)
}

func getStatusCode(code errors.ErrorCode) int {
	switch code {
	case errors.ErrValidation:
		return http.StatusBadRequest
	case errors.ErrUnauthorized:
		return http.StatusUnauthorized
	case errors.ErrForbidden:
		return http.StatusForbidden
	case errors.ErrNotFound:
		return http.StatusNotFound
	case errors.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
