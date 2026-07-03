package settings

import (
	"github.com/edustack/go-boilerplate/internal/audit"
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua endpoint settings ke router group.
// Termasuk audit logs yang dikelola oleh internal/audit package.
func RegisterRoutes(router *gin.RouterGroup, auditHandler *audit.Handler, cfg *config.Config) {
	settings := router.Group("/settings")
	settings.Use(middleware.AuthMiddleware(cfg))
	{
		// Audit Logs
		settings.GET("/logs", auditHandler.GetLogs)
		settings.GET("/logs/:id", auditHandler.GetLogDetail)
	}
}
