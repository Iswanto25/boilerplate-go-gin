package router

import (
	"net/http"

	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/features/auth"
	authHandler "github.com/edustack/go-boilerplate/internal/features/auth/handler"
	"github.com/edustack/go-boilerplate/internal/features/user"
	userHandler "github.com/edustack/go-boilerplate/internal/features/user/handler"
	"github.com/edustack/go-boilerplate/pkg/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, uh *userHandler.UserHandler, ah *authHandler.AuthHandler) *gin.Engine {
	router := gin.Default()

	// Setup CORS
	router.Use(cors.Default())

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Service is healthy", nil)
	})

	// API Routes Group
	api := router.Group("/api/v1")
	{
		// Register feature routes (sub-routers)
		user.RegisterRoutes(api, uh)
		auth.RegisterRoutes(api, ah)
	}

	return router
}
