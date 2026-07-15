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
	Name            string           `json:"name" binding:"required"`
	ModuleID        uuid.UUID        `json:"moduleId" binding:"required"`
	AvailableAction []ResourceAction `json:"availableAction"`
}

type UpdateResourceRequest struct {
	Name            *string          `json:"name"`
	ModuleID        *uuid.UUID       `json:"moduleId"`
	AvailableAction []ResourceAction `json:"availableAction"`
}

type ResourceResponse struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	ModuleID        uuid.UUID        `json:"moduleId"`
	Module          *ModuleResponse  `json:"module,omitempty"`
	AvailableAction []ResourceAction `json:"availableAction"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

// --- Role DTOs ---

type CreateRoleRequest struct {
	Name   string `json:"name" binding:"required"`
	Status *bool  `json:"status"`
}

type UpdateRoleRequest struct {
	Name   *string `json:"name"`
	Status *bool   `json:"status"`
}

type RoleResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// --- RolePermission DTOs ---

type CreateRolePermissionRequest struct {
	RoleID         uuid.UUID `json:"roleId" binding:"required"`
	ResourceID     uuid.UUID `json:"resourceId" binding:"required"`
	GrantedActions []string  `json:"grantedActions"`
}

type UpdateRolePermissionRequest struct {
	RoleID         *uuid.UUID `json:"roleId"`
	ResourceID     *uuid.UUID `json:"resourceId"`
	GrantedActions []string   `json:"grantedActions"`
}

type RolePermissionResponse struct {
	ID             uuid.UUID        `json:"id"`
	RoleID         uuid.UUID        `json:"roleId"`
	Role           *RoleResponse    `json:"role,omitempty"`
	ResourceID     uuid.UUID        `json:"resourceId"`
	Resource       *ResourceResponse `json:"resource,omitempty"`
	GrantedActions []string         `json:"grantedActions"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
}

type ResourcePermissionDetail struct {
	ID             uuid.UUID `json:"id"`
	ResourceName   string   `json:"resource"`
	GrantedActions []string `json:"actions"`
}

type ModulePermissionDetail struct {
	ID        uuid.UUID                `json:"id"`
	Module    string                     `json:"module"`
	Resources []ResourcePermissionDetail `json:"resources"`
}

type RoleWithPermissionsResponse struct {
	ID          uuid.UUID                `json:"id"`
	Name        string                   `json:"name"`
	Status      bool                     `json:"status"`
	Permissions []ModulePermissionDetail `json:"permissions"`
}

// --- Converters ---

func ToRoleResponse(r *Role) RoleResponse {
	return RoleResponse{
		ID:        r.ID,
		Name:      r.Name,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func ToRoleResponses(roles []Role) []RoleResponse {
	var res []RoleResponse
	for _, r := range roles {
		res = append(res, ToRoleResponse(&r))
	}
	return res
}

func ToRolePermissionResponse(rp *RolePermission) RolePermissionResponse {
	resp := RolePermissionResponse{
		ID:             rp.ID,
		RoleID:         rp.RoleId,
		ResourceID:     rp.ResourceId,
		GrantedActions: rp.GrantedActions,
		CreatedAt:      rp.CreatedAt,
		UpdatedAt:      rp.UpdatedAt,
	}
	if rp.Role.ID != uuid.Nil {
		role := ToRoleResponse(&rp.Role)
		resp.Role = &role
	}
	if rp.Resource.ID != uuid.Nil {
		resource := ToResourceDetailResponse(&rp.Resource)
		resp.Resource = &resource
	}
	return resp
}

func ToRolePermissionResponses(perms []RolePermission) []RolePermissionResponse {
	var res []RolePermissionResponse
	for _, rp := range perms {
		res = append(res, ToRolePermissionResponse(&rp))
	}
	return res
}

func ToModulePermissionDetails(perms []RolePermission) []ModulePermissionDetail {
	moduleMap := make(map[string]*ModulePermissionDetail)
	var orderedModules []string

	for _, rp := range perms {
		moduleName := rp.Resource.Module.Name
		if moduleName == "" {
			moduleName = "unknown"
		}

		resPerm := ResourcePermissionDetail{
			ID:             rp.Resource.ID,
			ResourceName:   rp.Resource.Name,
			GrantedActions: rp.GrantedActions,
		}

		if modDetail, exists := moduleMap[moduleName]; exists {
			modDetail.Resources = append(modDetail.Resources, resPerm)
		} else {
			modDetail = &ModulePermissionDetail{
				ID:        rp.Resource.Module.ID,
				Module:    moduleName,
				Resources: []ResourcePermissionDetail{resPerm},
			}
			moduleMap[moduleName] = modDetail
			orderedModules = append(orderedModules, moduleName)
		}
	}

	var res []ModulePermissionDetail
	for _, mName := range orderedModules {
		res = append(res, *moduleMap[mName])
	}
	return res
}

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
