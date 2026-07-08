package audit

import (
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, cfg *config.Config, authDeps *middleware.AuthDeps) {
	auditGroup := router.Group("/audit")
	auditGroup.Use(middleware.AuthMiddleware(cfg, authDeps))
	{
		auditGroup.GET("/logs", handler.GetLogs)
		auditGroup.GET("/logs/:id", handler.GetLogDetail)
	}
}
