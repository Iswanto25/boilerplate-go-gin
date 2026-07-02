package router

import (
	"log/slog"
	"net/http"
	"time"

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
	// Gunakan gin.New() agar log teks default [GIN-debug] tidak ikut tercetak
	router := gin.New()

	// Pasang Recovery middleware agar server tidak crash jika ada panic
	router.Use(gin.Recovery())

	// Kustomisasi Logger Gin agar menggunakan struktur log/slog aplikasi Anda
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				slog.Error("HTTP Request Error", "error", e.Error())
			}
		} else {
			// Mencetak log akses menggunakan format yang sinkron dengan InitLogger (Text/JSON otomatis)
			slog.Info("HTTP Request",
				"status", c.Writer.Status(),
				"method", c.Request.Method,
				"path", path,
				"query", query,
				"ip", c.ClientIP(),
				"latency", latency.String(),
				"user-agent", c.Request.UserAgent(),
			)
		}
	})

	// Setup CORS
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/health")
	})
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