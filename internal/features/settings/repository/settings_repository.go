package repository

import (
	"context"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettingsRepository interface {
	// Module
	CreateModule(ctx context.Context, m *model.Module) error
	FindAllModules(ctx context.Context) ([]model.Module, error)
	FindModuleByID(ctx context.Context, id uuid.UUID) (*model.Module, error)
	UpdateModule(ctx context.Context, id uuid.UUID, m *model.Module) error
	DeleteModule(ctx context.Context, id uuid.UUID) error

	// Resource
	CreateResource(ctx context.Context, r *model.Resource) error
	FindAllResources(ctx context.Context) ([]model.Resource, error)
	FindResourceByID(ctx context.Context, id uuid.UUID) (*model.Resource, error)
	UpdateResource(ctx context.Context, id uuid.UUID, r *model.Resource) error
	DeleteResource(ctx context.Context, id uuid.UUID) error
}

type settingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{db: db}
}

// --- Module ---

func (r *settingsRepository) CreateModule(ctx context.Context, m *model.Module) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *settingsRepository) FindAllModules(ctx context.Context) ([]model.Module, error) {
	var modules []model.Module
	err := r.db.WithContext(ctx).Order("created_at ASC").Find(&modules).Error
	return modules, err
}

func (r *settingsRepository) FindModuleByID(ctx context.Context, id uuid.UUID) (*model.Module, error) {
	var m model.Module
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *settingsRepository) UpdateModule(ctx context.Context, id uuid.UUID, m *model.Module) error {
	return r.db.WithContext(ctx).Model(&model.Module{}).Where("id = ?", id).Updates(m).Error
}

func (r *settingsRepository) DeleteModule(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Module{}, "id = ?", id).Error
}

// --- Resource ---

func (r *settingsRepository) CreateResource(ctx context.Context, res *model.Resource) error {
	return r.db.WithContext(ctx).Create(res).Error
}

func (r *settingsRepository) FindAllResources(ctx context.Context) ([]model.Resource, error) {
	var resources []model.Resource
	err := r.db.WithContext(ctx).Preload("Module").Order("created_at ASC").Find(&resources).Error
	return resources, err
}

func (r *settingsRepository) FindResourceByID(ctx context.Context, id uuid.UUID) (*model.Resource, error) {
	var res model.Resource
	err := r.db.WithContext(ctx).Preload("Module").First(&res, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *settingsRepository) UpdateResource(ctx context.Context, id uuid.UUID, res *model.Resource) error {
	return r.db.WithContext(ctx).Model(&model.Resource{}).Where("id = ?", id).Updates(res).Error
}

func (r *settingsRepository) DeleteResource(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Resource{}, "id = ?", id).Error
}
