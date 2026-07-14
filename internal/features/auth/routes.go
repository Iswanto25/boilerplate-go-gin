package auth

import (
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/features/auth/handler"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, authHandler *handler.AuthHandler, cfg *config.Config, authDeps *middleware.AuthDeps) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.GET("/profile", middleware.AuthMiddleware(cfg, authDeps), authHandler.Profile)
		auth.POST("/logout", middleware.AuthMiddleware(cfg, authDeps), authHandler.Logout)
	}
}
