package errors

import "net/http"

type AppError struct {
	Code       int
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(msg string) *AppError {
	return &AppError{Code: 40001, Message: msg, StatusCode: http.StatusBadRequest}
}

var (
	ErrNotFound     = &AppError{Code: 40401, Message: "resource not found", StatusCode: http.StatusNotFound}
	ErrConflict     = &AppError{Code: 40901, Message: "resource already exists", StatusCode: http.StatusConflict}
	ErrUnauthorized = &AppError{Code: 40101, Message: "invalid email or password", StatusCode: http.StatusUnauthorized}
	ErrForbidden    = &AppError{Code: 40301, Message: "forbidden", StatusCode: http.StatusForbidden}
	ErrInternal     = &AppError{Code: 50001, Message: "internal server error", StatusCode: http.StatusInternalServerError}

	ErrModuleNotFound   = &AppError{Code: 40402, Message: "module not found", StatusCode: http.StatusNotFound}
	ErrResourceNotFound = &AppError{Code: 40403, Message: "resource not found", StatusCode: http.StatusNotFound}
	ErrLogNotFound      = &AppError{Code: 40404, Message: "log not found", StatusCode: http.StatusNotFound}
	ErrUserNotFound     = &AppError{Code: 40405, Message: "user not found", StatusCode: http.StatusNotFound}

	ErrValidation = &AppError{Code: 40001, Message: "validation error", StatusCode: http.StatusBadRequest}
)
