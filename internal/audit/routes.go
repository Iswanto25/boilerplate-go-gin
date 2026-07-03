package audit

import (
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua endpoint audit ke router group.
// Semua endpoint dilindungi oleh AuthMiddleware.
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, cfg *config.Config) {
	auditGroup := router.Group("/audit")
	auditGroup.Use(middleware.AuthMiddleware(cfg))
	{
		auditGroup.GET("/logs", handler.GetLogs)
		auditGroup.GET("/logs/:id", handler.GetLogDetail)
	}
}
