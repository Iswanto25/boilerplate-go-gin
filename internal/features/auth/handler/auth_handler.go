package handler

import (
	"net/http"

	"github.com/edustack/go-boilerplate/internal/features/auth/model"
	"github.com/edustack/go-boilerplate/internal/features/auth/service"
	"github.com/edustack/go-boilerplate/internal/features/auth/validate"
	"github.com/edustack/go-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.RegisterRequest(&req); err != nil {
		response.WriteError(c, err)
		return
	}

	result, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "registration successful", result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.LoginRequest(&req); err != nil {
		response.WriteError(c, err)
		return
	}

	result, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "login successful", result)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.RefreshTokenRequest(&req); err != nil {
		response.WriteError(c, err)
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "token refreshed successfully", result)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.authService.Profile(c.Request.Context(), userID.(string))
	if err != nil {
		response.WriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "profile retrieved successfully", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID.(string)); err != nil {
		response.WriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "logout successful", nil)
}
