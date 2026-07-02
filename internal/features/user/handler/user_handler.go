package handler

import (
	"net/http"

	"github.com/edustack/go-boilerplate/internal/features/user/model"
	"github.com/edustack/go-boilerplate/internal/features/user/service"
	"github.com/edustack/go-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userResp, err := h.userService.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", userResp)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user ID (must be a valid UUID)")
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", model.ToUserResponse(user))
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]model.UserResponse, len(users))
	for i, u := range users {
		responses[i] = model.ToUserResponse(u)
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", responses)
}
