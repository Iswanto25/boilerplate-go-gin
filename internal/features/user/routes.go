package user

import (
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/features/user/handler"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler, cfg *config.Config, authDeps *middleware.AuthDeps) {
	users := router.Group("/users")
	users.Use(middleware.AuthMiddleware(cfg, authDeps))
	{
		users.GET("/", userHandler.GetAllUsers)
		users.GET("/:id", userHandler.GetUser)
	}
}
