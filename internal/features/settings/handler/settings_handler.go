package handler

import (
	"net/http"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/edustack/go-boilerplate/internal/features/settings/service"
	"github.com/edustack/go-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SettingsHandler struct {
	service service.SettingsService
}

func NewSettingsHandler(svc service.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: svc}
}

// --- Module ---

func (h *SettingsHandler) CreateModule(c *gin.Context) {
	var req model.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.CreateModule(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "module created successfully", resp)
}

func (h *SettingsHandler) GetAllModules(c *gin.Context) {
	modules, err := h.service.GetAllModules(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "modules retrieved successfully", modules)
}

func (h *SettingsHandler) GetModuleByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid module ID")
		return
	}

	resp, err := h.service.GetModuleByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "module not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "module retrieved successfully", resp)
}

func (h *SettingsHandler) UpdateModule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid module ID")
		return
	}

	var req model.UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.UpdateModule(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "module not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "module updated successfully", resp)
}

func (h *SettingsHandler) DeleteModule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid module ID")
		return
	}

	if err := h.service.DeleteModule(c.Request.Context(), id); err != nil {
		if err.Error() == "module not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "module deleted successfully", nil)
}

// --- Resource ---

func (h *SettingsHandler) CreateResource(c *gin.Context) {
	var req model.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.CreateResource(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "module not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "resource created successfully", resp)
}

func (h *SettingsHandler) GetAllResources(c *gin.Context) {
	resources, err := h.service.GetAllResources(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "resources retrieved successfully", resources)
}

func (h *SettingsHandler) GetResourceByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid resource ID")
		return
	}

	resp, err := h.service.GetResourceByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "resource not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "resource retrieved successfully", resp)
}

func (h *SettingsHandler) UpdateResource(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid resource ID")
		return
	}

	var req model.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.service.UpdateResource(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "resource not found" || err.Error() == "module not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "resource updated successfully", resp)
}

func (h *SettingsHandler) DeleteResource(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid resource ID")
		return
	}

	if err := h.service.DeleteResource(c.Request.Context(), id); err != nil {
		if err.Error() == "resource not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "resource deleted successfully", nil)
}
