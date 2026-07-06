package user

import (
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/features/user/handler"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler, cfg *config.Config) {
	users := router.Group("/users")
	{
		users.POST("/", userHandler.Register)
		users.GET("/", middleware.AuthMiddleware(cfg), userHandler.GetAllUsers)
		users.GET("/:id", middleware.AuthMiddleware(cfg), userHandler.GetUser)
	}
}
