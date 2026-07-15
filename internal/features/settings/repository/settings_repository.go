package repository

import (
	"context"
	"fmt"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) CreateModule(ctx context.Context, m *model.Module) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create module: %w", err)
	}
	return nil
}

func (r *SettingsRepository) FindAllModules(ctx context.Context, limit, offset int) ([]model.Module, error) {
	var modules []model.Module
	q := r.db.WithContext(ctx).Order(`"createdAt" ASC`)
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Find(&modules).Error; err != nil {
		return nil, fmt.Errorf("find all modules: %w", err)
	}
	return modules, nil
}

func (r *SettingsRepository) FindModuleByID(ctx context.Context, id uuid.UUID) (*model.Module, error) {
	var m model.Module
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("find module by id: %w", err)
	}
	return &m, nil
}

func (r *SettingsRepository) UpdateModule(ctx context.Context, id uuid.UUID, m *model.Module) error {
	result := r.db.WithContext(ctx).Model(&model.Module{}).Where("id = ?", id).Updates(m)
	if result.Error != nil {
		return fmt.Errorf("update module: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SettingsRepository) DeleteModule(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Module{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("delete module: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SettingsRepository) CreateResource(ctx context.Context, res *model.Resource) error {
	if err := r.db.WithContext(ctx).Create(res).Error; err != nil {
		return fmt.Errorf("create resource: %w", err)
	}
	return nil
}

func (r *SettingsRepository) FindAllResources(ctx context.Context, limit, offset int) ([]model.Resource, error) {
	var resources []model.Resource
	q := r.db.WithContext(ctx).Preload("Module").Order(`"createdAt" ASC`)
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Find(&resources).Error; err != nil {
		return nil, fmt.Errorf("find all resources: %w", err)
	}
	return resources, nil
}

func (r *SettingsRepository) FindResourceByID(ctx context.Context, id uuid.UUID) (*model.Resource, error) {
	var res model.Resource
	if err := r.db.WithContext(ctx).Preload("Module").First(&res, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("find resource by id: %w", err)
	}
	return &res, nil
}

func (r *SettingsRepository) UpdateResource(ctx context.Context, id uuid.UUID, res *model.Resource) error {
	result := r.db.WithContext(ctx).Model(&model.Resource{}).Where("id = ?", id).Updates(res)
	if result.Error != nil {
		return fmt.Errorf("update resource: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SettingsRepository) DeleteResource(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Resource{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("delete resource: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// --- Role ---

func (r *SettingsRepository) CreateRole(ctx context.Context, role *model.Role) error {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return fmt.Errorf("create role: %w", err)
	}
	return nil
}

func (r *SettingsRepository) FindRoleByID(ctx context.Context, id uuid.UUID) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("find role by id: %w", err)
	}
	return &role, nil
}

func (r *SettingsRepository) FindRoleByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, fmt.Errorf("find role by name: %w", err)
	}
	return &role, nil
}

func (r *SettingsRepository) FindAllRoles(ctx context.Context, limit, offset int) ([]model.Role, error) {
	var roles []model.Role
	q := r.db.WithContext(ctx).Order(`"createdAt" ASC`)
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("find all roles: %w", err)
	}
	return roles, nil
}

func (r *SettingsRepository) UpdateRole(ctx context.Context, id uuid.UUID, role *model.Role) error {
	result := r.db.WithContext(ctx).Model(&model.Role{}).Where("id = ?", id).Updates(role)
	if result.Error != nil {
		return fmt.Errorf("update role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SettingsRepository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Role{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("delete role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// --- RolePermission ---

func (r *SettingsRepository) CreateRolePermission(ctx context.Context, rp *model.RolePermission) error {
	if err := r.db.WithContext(ctx).Create(rp).Error; err != nil {
		return fmt.Errorf("create role permission: %w", err)
	}
	return nil
}

func (r *SettingsRepository) FindRolePermissionByID(ctx context.Context, id uuid.UUID) (*model.RolePermission, error) {
	var rp model.RolePermission
	if err := r.db.WithContext(ctx).Preload("Role").Preload("Resource.Module").First(&rp, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("find role permission by id: %w", err)
	}
	return &rp, nil
}

func (r *SettingsRepository) FindAllRolePermissions(ctx context.Context, limit, offset int) ([]model.RolePermission, error) {
	var perms []model.RolePermission
	q := r.db.WithContext(ctx).Preload("Role").Preload("Resource.Module").Order(`"createdAt" ASC`)
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	if err := q.Find(&perms).Error; err != nil {
		return nil, fmt.Errorf("find all role permissions: %w", err)
	}
	return perms, nil
}

func (r *SettingsRepository) FindRolePermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]model.RolePermission, error) {
	var perms []model.RolePermission
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Resource.Module").
		Where(`"roleId" = ?`, roleID).
		Order(`"createdAt" ASC`).
		Find(&perms).Error; err != nil {
		return nil, fmt.Errorf("find role permissions by role id: %w", err)
	}
	return perms, nil
}

func (r *SettingsRepository) UpdateRolePermission(ctx context.Context, id uuid.UUID, rp *model.RolePermission) error {
	result := r.db.WithContext(ctx).Model(&model.RolePermission{}).Where("id = ?", id).Updates(rp)
	if result.Error != nil {
		return fmt.Errorf("update role permission: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SettingsRepository) DeleteRolePermission(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.RolePermission{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("delete role permission: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
