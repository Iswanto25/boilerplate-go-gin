package auth

import (
	"github.com/edustack/go-boilerplate/internal/features/auth/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers auth endpoints onto the router group.
func RegisterRoutes(router *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
