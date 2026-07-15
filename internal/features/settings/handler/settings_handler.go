package handler

import (
	"net/http"
	"strings"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/edustack/go-boilerplate/internal/features/settings/service"
	"github.com/edustack/go-boilerplate/internal/features/settings/validate"
	"github.com/edustack/go-boilerplate/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SettingsHandler struct {
	service *service.SettingsService
}

func NewSettingsHandler(svc *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: svc}
}

// --- Module ---

func (h *SettingsHandler) CreateModule(c *gin.Context) {
	var req model.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.CreateModuleRequest(&req); err != nil {
		pkg.WriteError(c, err)
		return
	}

	resp, err := h.service.CreateModule(c.Request.Context(), &req)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusCreated, "module created successfully", resp)
}

func (h *SettingsHandler) GetAllModules(c *gin.Context) {
	modules, err := h.service.GetAllModules(c.Request.Context())
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "modules retrieved successfully", modules)
}

func (h *SettingsHandler) GetModuleByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid module ID")
		return
	}

	resp, err := h.service.GetModuleByID(c.Request.Context(), id)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "module retrieved successfully", resp)
}

func (h *SettingsHandler) UpdateModule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid module ID")
		return
	}

	var req model.UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.UpdateModuleRequest(&req); err != nil {
		pkg.WriteError(c, err)
		return
	}

	resp, err := h.service.UpdateModule(c.Request.Context(), id, &req)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "module updated successfully", resp)
}

func (h *SettingsHandler) DeleteModule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid module ID")
		return
	}

	if err := h.service.DeleteModule(c.Request.Context(), id); err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "module deleted successfully", nil)
}

// --- Resource ---

func (h *SettingsHandler) CreateResource(c *gin.Context) {
	var req model.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.CreateResourceRequest(&req); err != nil {
		pkg.WriteError(c, err)
		return
	}

	resp, err := h.service.CreateResource(c.Request.Context(), &req)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusCreated, "resource created successfully", resp)
}

func (h *SettingsHandler) GetAllResources(c *gin.Context) {
	resources, err := h.service.GetAllResources(c.Request.Context())
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "resources retrieved successfully", resources)
}

func (h *SettingsHandler) GetResourceByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid resource ID")
		return
	}

	resp, err := h.service.GetResourceByID(c.Request.Context(), id)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "resource retrieved successfully", resp)
}

func (h *SettingsHandler) UpdateResource(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid resource ID")
		return
	}

	var req model.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.UpdateResourceRequest(&req); err != nil {
		pkg.WriteError(c, err)
		return
	}

	resp, err := h.service.UpdateResource(c.Request.Context(), id, &req)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "resource updated successfully", resp)
}

func (h *SettingsHandler) DeleteResource(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid resource ID")
		return
	}

	if err := h.service.DeleteResource(c.Request.Context(), id); err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "resource deleted successfully", nil)
}

// --- Role ---

func (h *SettingsHandler) CreateRole(c *gin.Context) {
	var req model.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.CreateRoleRequest(&req); err != nil {
		pkg.WriteError(c, err)
		return
	}

	resp, err := h.service.CreateRole(c.Request.Context(), &req)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusCreated, "role created successfully", resp)
}

func (h *SettingsHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.service.GetAllRoles(c.Request.Context())
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "roles retrieved successfully", roles)
}

func (h *SettingsHandler) GetRoleByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid role ID")
		return
	}

	resp, err := h.service.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "role retrieved successfully", resp)
}

func (h *SettingsHandler) GetRoleByName(c *gin.Context) {
	name := c.Param("name")
	if strings.TrimSpace(name) == "" {
		pkg.Error(c, http.StatusBadRequest, "role name is required")
		return
	}

	resp, err := h.service.GetRoleByName(c.Request.Context(), name)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "role retrieved successfully", resp)
}

func (h *SettingsHandler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid role ID")
		return
	}

	var req model.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.UpdateRoleRequest(&req); err != nil {
		pkg.WriteError(c, err)
		return
	}

	resp, err := h.service.UpdateRole(c.Request.Context(), id, &req)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "role updated successfully", resp)
}

func (h *SettingsHandler) DeleteRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid role ID")
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "role deleted successfully", nil)
}

// --- RolePermission ---

func (h *SettingsHandler) CreateRolePermission(c *gin.Context) {
	var req model.CreateRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.CreateRolePermissionRequest(&req); err != nil {
		pkg.WriteError(c, err)
		return
	}

	resp, err := h.service.CreateRolePermission(c.Request.Context(), &req)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusCreated, "role permission created successfully", resp)
}

func (h *SettingsHandler) GetAllRolePermissions(c *gin.Context) {
	perms, err := h.service.GetAllRolePermissions(c.Request.Context())
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "role permissions retrieved successfully", perms)
}

func (h *SettingsHandler) GetRolePermissionByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid role permission ID")
		return
	}

	resp, err := h.service.GetRolePermissionByID(c.Request.Context(), id)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "role permission retrieved successfully", resp)
}

func (h *SettingsHandler) UpdateRolePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid role permission ID")
		return
	}

	var req model.UpdateRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.UpdateRolePermissionRequest(&req); err != nil {
		pkg.WriteError(c, err)
		return
	}

	resp, err := h.service.UpdateRolePermission(c.Request.Context(), id, &req)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "role permission updated successfully", resp)
}

func (h *SettingsHandler) DeleteRolePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid role permission ID")
		return
	}

	if err := h.service.DeleteRolePermission(c.Request.Context(), id); err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "role permission deleted successfully", nil)
}
