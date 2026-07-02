package user

import (
	"github.com/edustack/go-boilerplate/internal/features/user/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := router.Group("/users")
	{
		users.POST("/", userHandler.Register)
		users.GET("/", userHandler.GetAllUsers)
		users.GET("/:id", userHandler.GetUser)
	}
}
