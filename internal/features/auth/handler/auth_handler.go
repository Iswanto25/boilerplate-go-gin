package handler

import (
	"net/http"

	"github.com/edustack/go-boilerplate/internal/features/auth/model"
	"github.com/edustack/go-boilerplate/internal/features/auth/service"
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
		errorMsg := err.Error()
		if len(req.Password) > 0 && len(req.Password) < 8 {
			errorMsg = "Password harus memiliki panjang minimal 8 karakter dan mengandung karakter huruf, angka dan simbol"
		} else if req.Name == "" || req.Email == "" || req.Password == "" {
			errorMsg = "Semua kolom (name, email, password) wajib diisi"
		}

		response.Error(c, http.StatusBadRequest, errorMsg)
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "registration successful", resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "login successful", resp)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "token refreshed successfully", resp)
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
