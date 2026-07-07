package auth

import (
	"github.com/edustack/go-boilerplate/internal/features/auth/handler"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers auth endpoints onto the router group.
func RegisterRoutes(router *gin.RouterGroup, authHandler *handler.AuthHandler, cfg *config.Config) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(cfg), authHandler.Logout)
	}
}
