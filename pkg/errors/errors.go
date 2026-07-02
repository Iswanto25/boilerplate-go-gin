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

var (
	ErrNotFound     = &AppError{Code: 40401, Message: "resource not found", StatusCode: http.StatusNotFound}
	ErrConflict     = &AppError{Code: 40901, Message: "resource already exists", StatusCode: http.StatusConflict}
	ErrUnauthorized = &AppError{Code: 40101, Message: "invalid email or password", StatusCode: http.StatusUnauthorized}
	ErrForbidden    = &AppError{Code: 40301, Message: "forbidden", StatusCode: http.StatusForbidden}
	ErrInternal     = &AppError{Code: 50001, Message: "internal server error", StatusCode: http.StatusInternalServerError}
)
