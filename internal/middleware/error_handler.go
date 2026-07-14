package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	appErr "github.com/edustack/go-boilerplate/pkg/errors"
	"github.com/edustack/go-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

// writeErrorResponse writes a unified error response. In production, a 500
// is masked so internal system/database details are not leaked to clients.
// Mirrors the Express `errorHandler` 500-masking behavior.
func writeErrorResponse(c *gin.Context, err error, statusCode int, isProduction bool) {
	msg := err.Error()
	if isProduction && statusCode == http.StatusInternalServerError {
		msg = "Internal server error"
	}
	response.Error(c, statusCode, msg)
}

// logError logs the full error for internal debugging. The stack trace is
// omitted in production to avoid leaking implementation details.
func logError(c *gin.Context, err error, statusCode int, isProduction bool) {
	slog.Error("request failed",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"statusCode", statusCode,
		"error", err.Error(),
	)
	if !isProduction {
		slog.Error("stack trace", "stack", string(debug.Stack()))
	}
}

// Recovery replaces gin.Recovery(). It catches panics and returns a unified
// JSON response via the central response helper instead of Gin's default
// panic output, and masks internal errors in production.
func Recovery(isProduction bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				if !c.Writer.Written() {
					writeErrorResponse(c, err, http.StatusInternalServerError, isProduction)
				}
				logError(c, err, http.StatusInternalServerError, isProduction)
			}
		}()
		c.Next()
	}
}

// ErrorHandler processes errors pushed through c.Error(). It runs after the
// handler chain and writes a unified response only if nothing has been written
// yet. This mirrors Express's `next(err)` propagation to the error middleware.
func ErrorHandler(isProduction bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err

		statusCode := http.StatusInternalServerError
		var appError *appErr.AppError
		if errors.As(err, &appError) {
			statusCode = appError.StatusCode
		}

		logError(c, err, statusCode, isProduction)
		writeErrorResponse(c, err, statusCode, isProduction)
	}
}

// NotFound handles unmatched routes and methods. Mirrors the Express
// `notFoundHandler` that returns "Route METHOD PATH not found".
func NotFound(isProduction bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isProduction {
			response.Error(c, http.StatusNotFound,
				fmt.Sprintf("Route %s %s not found", c.Request.Method, c.Request.URL.Path))
		} else {
			response.Error(c, http.StatusNotFound, "Not found")
		}
	}
}
