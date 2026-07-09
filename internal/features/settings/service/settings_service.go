package service

import (
	"context"
	"errors"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/edustack/go-boilerplate/internal/features/settings/repository"
	appErr "github.com/edustack/go-boilerplate/pkg/errors"
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
		return nil, appErr.ErrInternal
	}
	resp := model.ToModuleResponse(m)
	return &resp, nil
}

func (s *settingsService) GetAllModules(ctx context.Context) ([]model.ModuleResponse, error) {
	modules, err := s.repo.FindAllModules(ctx, 0, 0)
	if err != nil {
		return nil, appErr.ErrInternal
	}
	return model.ToModuleResponses(modules), nil
}

func (s *settingsService) GetModuleByID(ctx context.Context, id uuid.UUID) (*model.ModuleResponse, error) {
	m, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.ErrModuleNotFound
		}
		return nil, appErr.ErrInternal
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
			return nil, appErr.ErrModuleNotFound
		}
		return nil, appErr.ErrInternal
	}

	m, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		return nil, appErr.ErrInternal
	}
	resp := model.ToModuleResponse(m)
	return &resp, nil
}

func (s *settingsService) DeleteModule(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteModule(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErr.ErrModuleNotFound
		}
		return appErr.ErrInternal
	}
	return nil
}

// --- Resource ---

func (s *settingsService) CreateResource(ctx context.Context, req *model.CreateResourceRequest) (*model.ResourceResponse, error) {
	mod, err := s.repo.FindModuleByID(ctx, req.ModuleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.ErrModuleNotFound
		}
		return nil, appErr.ErrInternal
	}

	r := &model.Resource{
		Name:            req.Name,
		ModuleId:        req.ModuleID,
		AvailableAction: req.AvailableAction,
		Module:          *mod,
	}
	if err := s.repo.CreateResource(ctx, r); err != nil {
		return nil, appErr.ErrInternal
	}

	resp := model.ToResourceDetailResponse(r)
	return &resp, nil
}

func (s *settingsService) GetAllResources(ctx context.Context) ([]model.ResourceResponse, error) {
	resources, err := s.repo.FindAllResources(ctx, 0, 0)
	if err != nil {
		return nil, appErr.ErrInternal
	}
	return model.ToResourceResponses(resources), nil
}

func (s *settingsService) GetResourceByID(ctx context.Context, id uuid.UUID) (*model.ResourceResponse, error) {
	r, err := s.repo.FindResourceByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.ErrResourceNotFound
		}
		return nil, appErr.ErrInternal
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
				return nil, appErr.ErrModuleNotFound
			}
			return nil, appErr.ErrInternal
		}
		updates.ModuleId = *req.ModuleID
	}
	if req.AvailableAction != nil {
		updates.AvailableAction = req.AvailableAction
	}

	if err := s.repo.UpdateResource(ctx, id, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.ErrResourceNotFound
		}
		return nil, appErr.ErrInternal
	}

	r, err := s.repo.FindResourceByID(ctx, id)
	if err != nil {
		return nil, appErr.ErrInternal
	}
	resp := model.ToResourceDetailResponse(r)
	return &resp, nil
}

func (s *settingsService) DeleteResource(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteResource(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErr.ErrResourceNotFound
		}
		return appErr.ErrInternal
	}
	return nil
}
