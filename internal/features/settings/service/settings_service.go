package service

import (
	"context"
	"errors"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/edustack/go-boilerplate/internal/features/settings/repository"
	"github.com/edustack/go-boilerplate/pkg"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettingsService interface {
	// Module
	CreateModule(ctx context.Context, req *model.CreateModuleRequest) (*model.ModuleResponse, error)
	GetAllModules(ctx context.Context) ([]model.ModuleResponse, error)
	GetModuleByID(ctx context.Context, id uuid.UUID) (*model.ModuleResponse, error)
	UpdateModule(ctx context.Context, id uuid.UUID, req *model.UpdateModuleRequest) (*model.ModuleResponse, error)
	DeleteModule(ctx context.Context, id uuid.UUID) error

	// Resource
	CreateResource(ctx context.Context, req *model.CreateResourceRequest) (*model.ResourceResponse, error)
	GetAllResources(ctx context.Context) ([]model.ResourceResponse, error)
	GetResourceByID(ctx context.Context, id uuid.UUID) (*model.ResourceResponse, error)
	UpdateResource(ctx context.Context, id uuid.UUID, req *model.UpdateResourceRequest) (*model.ResourceResponse, error)
	DeleteResource(ctx context.Context, id uuid.UUID) error

	// Role
	CreateRole(ctx context.Context, req *model.CreateRoleRequest) (*model.RoleResponse, error)
	GetAllRoles(ctx context.Context) ([]model.RoleResponse, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*model.RoleResponse, error)
	GetRoleByName(ctx context.Context, name string) (*model.RoleWithPermissionsResponse, error)
	UpdateRole(ctx context.Context, id uuid.UUID, req *model.UpdateRoleRequest) (*model.RoleResponse, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error

	// RolePermission
	CreateRolePermission(ctx context.Context, req *model.CreateRolePermissionRequest) (*model.RolePermissionResponse, error)
	GetAllRolePermissions(ctx context.Context) ([]model.RolePermissionResponse, error)
	GetRolePermissionByID(ctx context.Context, id uuid.UUID) (*model.RolePermissionResponse, error)
	UpdateRolePermission(ctx context.Context, id uuid.UUID, req *model.UpdateRolePermissionRequest) (*model.RolePermissionResponse, error)
	DeleteRolePermission(ctx context.Context, id uuid.UUID) error
}

type settingsService struct {
	repo repository.SettingsRepository
}

func NewSettingsService(repo repository.SettingsRepository) SettingsService {
	return &settingsService{repo: repo}
}

// --- Module ---

func (s *settingsService) CreateModule(ctx context.Context, req *model.CreateModuleRequest) (*model.ModuleResponse, error) {
	m := &model.Module{
		Name: req.Name,
	}
	if err := s.repo.CreateModule(ctx, m); err != nil {
		return nil, pkg.ErrInternal
	}
	resp := model.ToModuleResponse(m)
	return &resp, nil
}

func (s *settingsService) GetAllModules(ctx context.Context) ([]model.ModuleResponse, error) {
	modules, err := s.repo.FindAllModules(ctx, 0, 0)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	return model.ToModuleResponses(modules), nil
}

func (s *settingsService) GetModuleByID(ctx context.Context, id uuid.UUID) (*model.ModuleResponse, error) {
	m, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrModuleNotFound
		}
		return nil, pkg.ErrInternal
	}
	resp := model.ToModuleResponse(m)
	return &resp, nil
}

func (s *settingsService) UpdateModule(ctx context.Context, id uuid.UUID, req *model.UpdateModuleRequest) (*model.ModuleResponse, error) {
	updates := &model.Module{}
	if req.Name != nil {
		updates.Name = *req.Name
	}

	if err := s.repo.UpdateModule(ctx, id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrModuleNotFound
		}
		return nil, pkg.ErrInternal
	}

	m, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	resp := model.ToModuleResponse(m)
	return &resp, nil
}

func (s *settingsService) DeleteModule(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteModule(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.ErrModuleNotFound
		}
		return pkg.ErrInternal
	}
	return nil
}

// --- Resource ---

func (s *settingsService) CreateResource(ctx context.Context, req *model.CreateResourceRequest) (*model.ResourceResponse, error) {
	mod, err := s.repo.FindModuleByID(ctx, req.ModuleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrModuleNotFound
		}
		return nil, pkg.ErrInternal
	}

	r := &model.Resource{
		Name:            req.Name,
		ModuleId:        req.ModuleID,
		AvailableAction: req.AvailableAction,
		Module:          *mod,
	}
	if err := s.repo.CreateResource(ctx, r); err != nil {
		return nil, pkg.ErrInternal
	}

	resp := model.ToResourceDetailResponse(r)
	return &resp, nil
}

func (s *settingsService) GetAllResources(ctx context.Context) ([]model.ResourceResponse, error) {
	resources, err := s.repo.FindAllResources(ctx, 0, 0)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	return model.ToResourceResponses(resources), nil
}

func (s *settingsService) GetResourceByID(ctx context.Context, id uuid.UUID) (*model.ResourceResponse, error) {
	r, err := s.repo.FindResourceByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrResourceNotFound
		}
		return nil, pkg.ErrInternal
	}
	resp := model.ToResourceDetailResponse(r)
	return &resp, nil
}

func (s *settingsService) UpdateResource(ctx context.Context, id uuid.UUID, req *model.UpdateResourceRequest) (*model.ResourceResponse, error) {
	updates := &model.Resource{}
	if req.Name != nil {
		updates.Name = *req.Name
	}
	if req.ModuleID != nil {
		if _, err := s.repo.FindModuleByID(ctx, *req.ModuleID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg.ErrModuleNotFound
			}
			return nil, pkg.ErrInternal
		}
		updates.ModuleId = *req.ModuleID
	}
	if req.AvailableAction != nil {
		updates.AvailableAction = req.AvailableAction
	}

	if err := s.repo.UpdateResource(ctx, id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrResourceNotFound
		}
		return nil, pkg.ErrInternal
	}

	r, err := s.repo.FindResourceByID(ctx, id)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	resp := model.ToResourceDetailResponse(r)
	return &resp, nil
}

func (s *settingsService) DeleteResource(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteResource(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.ErrResourceNotFound
		}
		return pkg.ErrInternal
	}
	return nil
}

// --- Role ---

func (s *settingsService) CreateRole(ctx context.Context, req *model.CreateRoleRequest) (*model.RoleResponse, error) {
	status := true
	if req.Status != nil {
		status = *req.Status
	}
	role := &model.Role{
		Name:   req.Name,
		Status: status,
	}
	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, pkg.ErrInternal
	}
	resp := model.ToRoleResponse(role)
	return &resp, nil
}

func (s *settingsService) GetAllRoles(ctx context.Context) ([]model.RoleResponse, error) {
	roles, err := s.repo.FindAllRoles(ctx, 0, 0)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	return model.ToRoleResponses(roles), nil
}

func (s *settingsService) GetRoleByID(ctx context.Context, id uuid.UUID) (*model.RoleResponse, error) {
	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrRoleNotFound
		}
		return nil, pkg.ErrInternal
	}
	resp := model.ToRoleResponse(role)
	return &resp, nil
}

func (s *settingsService) GetRoleByName(ctx context.Context, name string) (*model.RoleWithPermissionsResponse, error) {
	role, err := s.repo.FindRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrRoleNotFound
		}
		return nil, pkg.ErrInternal
	}

	perms, err := s.repo.FindRolePermissionsByRoleID(ctx, role.ID)
	if err != nil {
		return nil, pkg.ErrInternal
	}

	return &model.RoleWithPermissionsResponse{
		ID:          role.ID,
		Name:        role.Name,
		Status:      role.Status,
		Permissions: model.ToRolePermissionResponses(perms),
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

func (s *settingsService) UpdateRole(ctx context.Context, id uuid.UUID, req *model.UpdateRoleRequest) (*model.RoleResponse, error) {
	updates := &model.Role{}
	if req.Name != nil {
		updates.Name = *req.Name
	}
	if req.Status != nil {
		updates.Status = *req.Status
	}

	if err := s.repo.UpdateRole(ctx, id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrRoleNotFound
		}
		return nil, pkg.ErrInternal
	}

	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	resp := model.ToRoleResponse(role)
	return &resp, nil
}

func (s *settingsService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteRole(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.ErrRoleNotFound
		}
		return pkg.ErrInternal
	}
	return nil
}

// --- RolePermission ---

func (s *settingsService) CreateRolePermission(ctx context.Context, req *model.CreateRolePermissionRequest) (*model.RolePermissionResponse, error) {
	if _, err := s.repo.FindRoleByID(ctx, req.RoleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrRoleNotFound
		}
		return nil, pkg.ErrInternal
	}
	if _, err := s.repo.FindResourceByID(ctx, req.ResourceID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrResourceNotFound
		}
		return nil, pkg.ErrInternal
	}

	rp := &model.RolePermission{
		RoleId:         req.RoleID,
		ResourceId:     req.ResourceID,
		GrantedActions: req.GrantedActions,
	}
	if err := s.repo.CreateRolePermission(ctx, rp); err != nil {
		return nil, pkg.ErrInternal
	}

	created, err := s.repo.FindRolePermissionByID(ctx, rp.ID)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	resp := model.ToRolePermissionResponse(created)
	return &resp, nil
}

func (s *settingsService) GetAllRolePermissions(ctx context.Context) ([]model.RolePermissionResponse, error) {
	perms, err := s.repo.FindAllRolePermissions(ctx, 0, 0)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	return model.ToRolePermissionResponses(perms), nil
}

func (s *settingsService) GetRolePermissionByID(ctx context.Context, id uuid.UUID) (*model.RolePermissionResponse, error) {
	rp, err := s.repo.FindRolePermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrRolePermissionNotFound
		}
		return nil, pkg.ErrInternal
	}
	resp := model.ToRolePermissionResponse(rp)
	return &resp, nil
}

func (s *settingsService) UpdateRolePermission(ctx context.Context, id uuid.UUID, req *model.UpdateRolePermissionRequest) (*model.RolePermissionResponse, error) {
	updates := &model.RolePermission{}
	if req.RoleID != nil {
		if _, err := s.repo.FindRoleByID(ctx, *req.RoleID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg.ErrRoleNotFound
			}
			return nil, pkg.ErrInternal
		}
		updates.RoleId = *req.RoleID
	}
	if req.ResourceID != nil {
		if _, err := s.repo.FindResourceByID(ctx, *req.ResourceID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg.ErrResourceNotFound
			}
			return nil, pkg.ErrInternal
		}
		updates.ResourceId = *req.ResourceID
	}
	if req.GrantedActions != nil {
		updates.GrantedActions = req.GrantedActions
	}

	if err := s.repo.UpdateRolePermission(ctx, id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrRolePermissionNotFound
		}
		return nil, pkg.ErrInternal
	}

	rp, err := s.repo.FindRolePermissionByID(ctx, id)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	resp := model.ToRolePermissionResponse(rp)
	return &resp, nil
}

func (s *settingsService) DeleteRolePermission(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteRolePermission(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.ErrRolePermissionNotFound
		}
		return pkg.ErrInternal
	}
	return nil
}
