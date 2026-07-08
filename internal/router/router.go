package router

import (
	"net/http"
	"time"

	"github.com/edustack/go-boilerplate/internal/audit"
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/features/auth"
	authHandler "github.com/edustack/go-boilerplate/internal/features/auth/handler"
	"github.com/edustack/go-boilerplate/internal/features/settings"
	settingsHandler "github.com/edustack/go-boilerplate/internal/features/settings/handler"
	"github.com/edustack/go-boilerplate/internal/features/user"
	userHandler "github.com/edustack/go-boilerplate/internal/features/user/handler"
	"github.com/edustack/go-boilerplate/pkg/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, uh *userHandler.UserHandler, ah *authHandler.AuthHandler, sh *settingsHandler.SettingsHandler, auditHandler *audit.Handler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	// Start time tracking for response time calculation
	router.Use(func(c *gin.Context) {
		c.Set("startTime", time.Now().UnixMilli())
		c.Next()
	})

	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/health")
	})
	router.GET("/health", func(c *gin.Context) {
		data := map[string]interface{}{
			"status":      "ok",
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
			"environment": cfg.AppEnv,
		}
		response.Success(c, http.StatusOK, "Service is healthy", data)
	})

	api := router.Group("/api/v1")
	{
		user.RegisterRoutes(api, uh, cfg)
		auth.RegisterRoutes(api, ah, cfg)
		settings.RegisterRoutes(api, sh, cfg)
		audit.RegisterRoutes(api, auditHandler, cfg)
	}

	return router
}
