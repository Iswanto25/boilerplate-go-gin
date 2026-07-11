package model

import (
	userModel "github.com/edustack/go-boilerplate/internal/features/user/model"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name     string       `json:"name" binding:"required,min=2,max=100"`
	Email    string       `json:"email" binding:"required,email"`
	Password string       `json:"password" binding:"required,min=8"`
	Role     userModel.Role `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string          `json:"accessToken"`
	RefreshToken string          `json:"refreshToken"`
	UserID       uuid.UUID       `json:"userId"`
	Email        string          `json:"email"`
	Name         string          `json:"name"`
	Role         userModel.Role  `json:"role"`
}
