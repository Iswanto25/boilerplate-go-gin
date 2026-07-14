package pkg

import "net/http"

// HttpStatus convenience constants (mirrors Express HttpStatus enum).
const (
	StatusOK                  = http.StatusOK
	StatusCreated             = http.StatusCreated
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusForbidden           = http.StatusForbidden
	StatusNotFound            = http.StatusNotFound
	StatusConflict            = http.StatusConflict
	StatusPayloadTooLarge     = http.StatusRequestEntityTooLarge
	StatusUnprocessableEntity = http.StatusUnprocessableEntity
	StatusUnsupportedMedia    = http.StatusUnsupportedMediaType
	StatusTooManyRequests     = http.StatusTooManyRequests
	StatusBadGateway          = http.StatusBadGateway
	StatusServiceUnavailable  = http.StatusServiceUnavailable
	StatusInternalServerError = http.StatusInternalServerError
)

type AppError struct {
	Code       int
	Message    string
	StatusCode int
	Hint       string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(msg string) *AppError {
	return &AppError{Code: 40001, Message: msg, StatusCode: http.StatusBadRequest}
}

func NewAppError(statusCode int, code int, msg string, hint string) *AppError {
	return &AppError{Code: code, Message: msg, StatusCode: statusCode, Hint: hint}
}

var (
	ErrNotFound     = &AppError{Code: 40401, Message: "resource not found", StatusCode: http.StatusNotFound}
	ErrConflict     = &AppError{Code: 40901, Message: "resource already exists", StatusCode: http.StatusConflict}
	ErrUnauthorized = &AppError{Code: 40101, Message: "invalid email", StatusCode: http.StatusUnauthorized}
	ErrInvalidCredentials   = &AppError{Code: 40102, Message: "invalid credentials", StatusCode: http.StatusUnauthorized}
	ErrForbidden    = &AppError{Code: 40301, Message: "forbidden", StatusCode: http.StatusForbidden}
	ErrInternal     = &AppError{Code: 50001, Message: "internal server error", StatusCode: http.StatusInternalServerError}

	ErrModuleNotFound          = &AppError{Code: 40402, Message: "module not found", StatusCode: http.StatusNotFound}
	ErrResourceNotFound        = &AppError{Code: 40403, Message: "resource not found", StatusCode: http.StatusNotFound}
	ErrLogNotFound             = &AppError{Code: 40404, Message: "log not found", StatusCode: http.StatusNotFound}
	ErrUserNotFound            = &AppError{Code: 40405, Message: "user not found", StatusCode: http.StatusNotFound}
	ErrRoleNotFound            = &AppError{Code: 40406, Message: "role not found", StatusCode: http.StatusNotFound}
	ErrRolePermissionNotFound  = &AppError{Code: 40407, Message: "role permission not found", StatusCode: http.StatusNotFound}
	ErrResourceNotFoundForRole = &AppError{Code: 40408, Message: "resource not found", StatusCode: http.StatusNotFound}

	ErrValidation = &AppError{Code: 40001, Message: "validation error", StatusCode: http.StatusBadRequest}
)
