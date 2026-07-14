package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/edustack/go-boilerplate/internal/audit"
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/features/auth"
	authHandler "github.com/edustack/go-boilerplate/internal/features/auth/handler"
	"github.com/edustack/go-boilerplate/internal/features/settings"
	settingsHandler "github.com/edustack/go-boilerplate/internal/features/settings/handler"
	"github.com/edustack/go-boilerplate/internal/features/user"
	userHandler "github.com/edustack/go-boilerplate/internal/features/user/handler"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/edustack/go-boilerplate/pkg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, uh *userHandler.UserHandler, ah *authHandler.AuthHandler, sh *settingsHandler.SettingsHandler, auditHandler *audit.Handler, authDeps *middleware.AuthDeps) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.Recovery(cfg.AppEnv == "production"))
	router.Use(middleware.ErrorHandler(cfg.AppEnv == "production"))

	router.Use(func(c *gin.Context) {
		c.Set("startTime", time.Now().UnixMilli())
		c.Next()
	})

	router.Use(pkg.CaptureRequestBody())

	corsCfg := cors.DefaultConfig()
	allowedOrigins := cfg.AllowedOrigins
	if allowedOrigins == "" || allowedOrigins == "*" {
		corsCfg.AllowAllOrigins = true
		corsCfg.AllowCredentials = false
	} else {
		origins := strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		corsCfg.AllowOrigins = origins
		corsCfg.AllowCredentials = true
	}
	corsCfg.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsCfg))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/health")
	})
	router.GET("/health", func(c *gin.Context) {
		data := map[string]interface{}{
			"status":      "ok",
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
			"environment": cfg.AppEnv,
		}
		pkg.Success(c, http.StatusOK, "Service is healthy", data)
	})

	api := router.Group("/api/v1")
	{
		user.RegisterRoutes(api, uh, cfg, authDeps)
		auth.RegisterRoutes(api, ah, cfg, authDeps)
		settings.RegisterRoutes(api, sh, cfg, authDeps)
		audit.RegisterRoutes(api, auditHandler, cfg, authDeps)
	}

	router.NoRoute(middleware.NotFound(cfg.AppEnv == "production"))
	router.NoMethod(middleware.NotFound(cfg.AppEnv == "production"))

	return router
}
