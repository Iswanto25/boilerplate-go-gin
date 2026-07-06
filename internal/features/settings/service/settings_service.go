package service

import (
	"context"
	"errors"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/edustack/go-boilerplate/internal/features/settings/repository"
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
		return nil, err
	}
	resp := model.ToModuleResponse(m)
	return &resp, nil
}

func (s *settingsService) GetAllModules(ctx context.Context) ([]model.ModuleResponse, error) {
	modules, err := s.repo.FindAllModules(ctx)
	if err != nil {
		return nil, err
	}
	return model.ToModuleResponses(modules), nil
}

func (s *settingsService) GetModuleByID(ctx context.Context, id uuid.UUID) (*model.ModuleResponse, error) {
	m, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}
	resp := model.ToModuleResponse(m)
	return &resp, nil
}

func (s *settingsService) UpdateModule(ctx context.Context, id uuid.UUID, req *model.UpdateModuleRequest) (*model.ModuleResponse, error) {
	_, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	updates := &model.Module{}
	if req.Name != nil {
		updates.Name = *req.Name
	}

	if err := s.repo.UpdateModule(ctx, id, updates); err != nil {
		return nil, err
	}

	m, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := model.ToModuleResponse(m)
	return &resp, nil
}

func (s *settingsService) DeleteModule(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.FindModuleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("module not found")
		}
		return err
	}
	return s.repo.DeleteModule(ctx, id)
}

// --- Resource ---

func (s *settingsService) CreateResource(ctx context.Context, req *model.CreateResourceRequest) (*model.ResourceResponse, error) {
	// Validate module exists
	if _, err := s.repo.FindModuleByID(ctx, req.ModuleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, err
	}

	r := &model.Resource{
		Name:            req.Name,
		ModuleId:        req.ModuleID,
		AvailableAction: req.AvailableAction,
	}
	if err := s.repo.CreateResource(ctx, r); err != nil {
		return nil, err
	}
	// Reload with module relation
	r, err := s.repo.FindResourceByID(ctx, r.ID)
	if err != nil {
		return nil, err
	}
	resp := model.ToResourceDetailResponse(r)
	return &resp, nil
}

func (s *settingsService) GetAllResources(ctx context.Context) ([]model.ResourceResponse, error) {
	resources, err := s.repo.FindAllResources(ctx)
	if err != nil {
		return nil, err
	}
	return model.ToResourceResponses(resources), nil
}

func (s *settingsService) GetResourceByID(ctx context.Context, id uuid.UUID) (*model.ResourceResponse, error) {
	r, err := s.repo.FindResourceByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("resource not found")
		}
		return nil, err
	}
	resp := model.ToResourceDetailResponse(r)
	return &resp, nil
}

func (s *settingsService) UpdateResource(ctx context.Context, id uuid.UUID, req *model.UpdateResourceRequest) (*model.ResourceResponse, error) {
	_, err := s.repo.FindResourceByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("resource not found")
		}
		return nil, err
	}

	updates := &model.Resource{}
	if req.Name != nil {
		updates.Name = *req.Name
	}
	if req.ModuleID != nil {
		// Validate module exists
		if _, err := s.repo.FindModuleByID(ctx, *req.ModuleID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("module not found")
			}
			return nil, err
		}
		updates.ModuleId = *req.ModuleID
	}
	if req.AvailableAction != nil {
		updates.AvailableAction = req.AvailableAction
	}

	if err := s.repo.UpdateResource(ctx, id, updates); err != nil {
		return nil, err
	}

	r, err := s.repo.FindResourceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := model.ToResourceDetailResponse(r)
	return &resp, nil
}

func (s *settingsService) DeleteResource(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.FindResourceByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("resource not found")
		}
		return err
	}
	return s.repo.DeleteResource(ctx, id)
}
