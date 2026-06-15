package errors

import "fmt"

type ErrorCode string

const (
	ErrValidation      ErrorCode = "VALIDATION_ERROR"
	ErrNotFound        ErrorCode = "RESOURCE_NOT_FOUND"
	ErrConflict        ErrorCode = "CONFLICT"
	ErrUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrForbidden       ErrorCode = "FORBIDDEN"
	ErrInternal        ErrorCode = "INTERNAL_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewValidationError(message string) *AppError {
	return &AppError{Code: ErrValidation, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: ErrNotFound, Message: message}
}

func NewConflictError(message string) *AppError {
	return &AppError{Code: ErrConflict, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: ErrUnauthorized, Message: message}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{Code: ErrForbidden, Message: message}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{Code: ErrInternal, Message: message, Err: err}
}
