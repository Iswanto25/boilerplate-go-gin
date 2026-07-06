package settings

import (
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/features/settings/handler"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.SettingsHandler, cfg *config.Config) {
	settings := router.Group("/settings")
	settings.Use(middleware.AuthMiddleware(cfg))
	{
		// Module
		settings.POST("/modules", h.CreateModule)
		settings.GET("/modules", h.GetAllModules)
		settings.GET("/modules/:id", h.GetModuleByID)
		settings.PUT("/modules/:id", h.UpdateModule)
		settings.DELETE("/modules/:id", h.DeleteModule)

		// Resource
		settings.POST("/resources", h.CreateResource)
		settings.GET("/resources", h.GetAllResources)
		settings.GET("/resources/:id", h.GetResourceByID)
		settings.PUT("/resources/:id", h.UpdateResource)
		settings.DELETE("/resources/:id", h.DeleteResource)
	}
}
