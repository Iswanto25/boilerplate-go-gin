package model

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"`
}

type UserResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	RoleID uuid.UUID `json:"roleId"`
	Role   string    `json:"role"`
}

func ToUserResponse(user *User) UserResponse {
	resp := UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		RoleID: user.RoleID,
	}
	if user.Role.ID != uuid.Nil {
		resp.Role = user.Role.Name
	}
	return resp
}
