package repository

import (
	"context"
	"fmt"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettingsRepository interface {
	CreateModule(ctx context.Context, m *model.Module) error
	FindAllModules(ctx context.Context, limit, offset int) ([]model.Module, error)
	FindModuleByID(ctx context.Context, id uuid.UUID) (*model.Module, error)
	UpdateModule(ctx context.Context, id uuid.UUID, m *model.Module) error
	DeleteModule(ctx context.Context, id uuid.UUID) error

	CreateResource(ctx context.Context, r *model.Resource) error
	FindAllResources(ctx context.Context, limit, offset int) ([]model.Resource, error)
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

func (r *settingsRepository) CreateModule(ctx context.Context, m *model.Module) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create module: %w", err)
	}
	return nil
}

func (r *settingsRepository) FindAllModules(ctx context.Context, limit, offset int) ([]model.Module, error) {
	var modules []model.Module
	q := r.db.WithContext(ctx).Order("created_at ASC")
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Find(&modules).Error; err != nil {
		return nil, fmt.Errorf("find all modules: %w", err)
	}
	return modules, nil
}

func (r *settingsRepository) FindModuleByID(ctx context.Context, id uuid.UUID) (*model.Module, error) {
	var m model.Module
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("find module by id: %w", err)
	}
	return &m, nil
}

func (r *settingsRepository) UpdateModule(ctx context.Context, id uuid.UUID, m *model.Module) error {
	result := r.db.WithContext(ctx).Model(&model.Module{}).Where("id = ?", id).Updates(m)
	if result.Error != nil {
		return fmt.Errorf("update module: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *settingsRepository) DeleteModule(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Module{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("delete module: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *settingsRepository) CreateResource(ctx context.Context, res *model.Resource) error {
	if err := r.db.WithContext(ctx).Create(res).Error; err != nil {
		return fmt.Errorf("create resource: %w", err)
	}
	return nil
}

func (r *settingsRepository) FindAllResources(ctx context.Context, limit, offset int) ([]model.Resource, error) {
	var resources []model.Resource
	q := r.db.WithContext(ctx).Preload("Module").Order("created_at ASC")
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Find(&resources).Error; err != nil {
		return nil, fmt.Errorf("find all resources: %w", err)
	}
	return resources, nil
}

func (r *settingsRepository) FindResourceByID(ctx context.Context, id uuid.UUID) (*model.Resource, error) {
	var res model.Resource
	if err := r.db.WithContext(ctx).Preload("Module").First(&res, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("find resource by id: %w", err)
	}
	return &res, nil
}

func (r *settingsRepository) UpdateResource(ctx context.Context, id uuid.UUID, res *model.Resource) error {
	result := r.db.WithContext(ctx).Model(&model.Resource{}).Where("id = ?", id).Updates(res)
	if result.Error != nil {
		return fmt.Errorf("update resource: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *settingsRepository) DeleteResource(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Resource{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("delete resource: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
