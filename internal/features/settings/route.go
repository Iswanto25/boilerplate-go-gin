package settings

import (
	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/features/settings/handler"
	"github.com/edustack/go-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.SettingsHandler, cfg *config.Config, authDeps *middleware.AuthDeps) {
	settings := router.Group("/settings")
	// settings.Use(middleware.AuthMiddleware(cfg, authDeps))
	{
		settings.POST("/modules", h.CreateModule)
		settings.GET("/modules", h.GetAllModules)
		settings.GET("/modules/:id", h.GetModuleByID)
		settings.PUT("/modules/:id", h.UpdateModule)
		settings.DELETE("/modules/:id", h.DeleteModule)

		settings.POST("/resources", h.CreateResource)
		settings.GET("/resources", h.GetAllResources)
		settings.GET("/resources/:id", h.GetResourceByID)
		settings.PUT("/resources/:id", h.UpdateResource)
		settings.DELETE("/resources/:id", h.DeleteResource)

		settings.POST("/roles", h.CreateRole)
		settings.GET("/roles", h.GetAllRoles)
		settings.GET("/roles/:name", h.GetRoleByName)
		settings.GET("/roles/id/:id", h.GetRoleByID)
		settings.PUT("/roles/:id", h.UpdateRole)
		settings.DELETE("/roles/:id", h.DeleteRole)

		settings.POST("/role-permissions", h.CreateRolePermission)
		settings.GET("/role-permissions", h.GetAllRolePermissions)
		settings.GET("/role-permissions/:id", h.GetRolePermissionByID)
		settings.PUT("/role-permissions/:id", h.UpdateRolePermission)
		settings.DELETE("/role-permissions/:id", h.DeleteRolePermission)
	}
}
