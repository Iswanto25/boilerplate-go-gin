package handler

import (
	"net/http"

	"github.com/edustack/go-boilerplate/internal/features/user/model"
	"github.com/edustack/go-boilerplate/internal/features/user/service"
	"github.com/edustack/go-boilerplate/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid user ID (must be a valid UUID)")
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "User retrieved successfully", model.ToUserResponse(user))
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	responses := make([]model.UserResponse, len(users))
	for i, u := range users {
		responses[i] = model.ToUserResponse(u)
	}

	pkg.Success(c, http.StatusOK, "Users retrieved successfully", responses)
}
