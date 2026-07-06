package model

import (
	"time"

	"github.com/google/uuid"
)

// --- Module DTOs ---

type CreateModuleRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateModuleRequest struct {
	Name *string `json:"name"`
}

type ModuleResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// --- Resource DTOs ---

type CreateResourceRequest struct {
	Name            string    `json:"name" binding:"required"`
	ModuleID        uuid.UUID `json:"moduleId" binding:"required"`
	AvailableAction []string  `json:"availableAction"`
}

type UpdateResourceRequest struct {
	Name            *string    `json:"name"`
	ModuleID        *uuid.UUID `json:"moduleId"`
	AvailableAction []string   `json:"availableAction"`
}

type ResourceResponse struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	ModuleID        uuid.UUID       `json:"moduleId"`
	Module          *ModuleResponse `json:"module,omitempty"`
	AvailableAction []string        `json:"availableAction"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

// --- Converters ---

func ToModuleResponse(m *Module) ModuleResponse {
	return ModuleResponse{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToModuleResponses(modules []Module) []ModuleResponse {
	var res []ModuleResponse
	for _, m := range modules {
		res = append(res, ToModuleResponse(&m))
	}
	return res
}

func ToResourceDetailResponse(r *Resource) ResourceResponse {
	resp := ResourceResponse{
		ID:              r.ID,
		Name:            r.Name,
		ModuleID:        r.ModuleId,
		AvailableAction: r.AvailableAction,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
	if r.Module.ID != uuid.Nil {
		m := ToModuleResponse(&r.Module)
		resp.Module = &m
	}
	return resp
}

func ToResourceResponses(resources []Resource) []ResourceResponse {
	var res []ResourceResponse
	for _, r := range resources {
		res = append(res, ToResourceDetailResponse(&r))
	}
	return res
}
