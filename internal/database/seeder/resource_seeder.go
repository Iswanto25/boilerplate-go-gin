package seeder

import (
	"time"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	ResourceUsersID     = "a2b7e63d-4c8f-4b35-86f7-1e5f8f553a1a"
	ResourceProfilesID  = "d9f8c6b5-0c7e-4a6f-b2d3-2f4c9a5e8c1b"
	ResourceRolesID     = "5e3a8d9c-1b7e-4c2f-8d9a-3f4b5c6d7e8f"
	ResourceModulesID   = "7c9b5a3d-4f1e-4b6a-9d8c-2e3f4a5b6c7d"
	ResourceResourcesID = "3f4a5b6c-7d8e-9f0a-1b2c-3d4e5f6a7b8c"
	ResourceRolePermsID = "b6c7d8e9-f0a1-4b2c-3d4e-5f6a7b8c9d0e"
	ResourceAuditLogsID = "f8e7d6c5-b4a3-4f2e-9d8c-7b6a5f4e3d2c"
)

var crudActions = []model.ResourceAction{
	model.ActionCreate, model.ActionRead, model.ActionUpdate, model.ActionDelete,
}

type resourceSeeder struct {
	enabled bool
}

func NewResourceSeeder() Seeder {
	return &resourceSeeder{enabled: true}
}

func (s *resourceSeeder) Name() string  { return "resource_seeder" }
func (s *resourceSeeder) Enabled() bool { return s.enabled }

func (s *resourceSeeder) Seed(db *gorm.DB) (int, error) {
	now := time.Now()
	modUserMgmt := uuid.MustParse(ModuleUserMgmtID)
	modSettings := uuid.MustParse(ModuleSettingsID)

	resources := []model.Resource{
		{ID: uuid.MustParse(ResourceUsersID), Name: "users", ModuleId: modUserMgmt, AvailableAction: crudActions, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.MustParse(ResourceProfilesID), Name: "profiles", ModuleId: modUserMgmt, AvailableAction: crudActions, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.MustParse(ResourceRolesID), Name: "roles", ModuleId: modSettings, AvailableAction: crudActions, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.MustParse(ResourceModulesID), Name: "modules", ModuleId: modSettings, AvailableAction: crudActions, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.MustParse(ResourceResourcesID), Name: "resources", ModuleId: modSettings, AvailableAction: crudActions, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.MustParse(ResourceRolePermsID), Name: "role-permissions", ModuleId: modSettings, AvailableAction: crudActions, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.MustParse(ResourceAuditLogsID), Name: "audit-logs", ModuleId: modSettings, AvailableAction: []model.ResourceAction{model.ActionRead}, CreatedAt: now, UpdatedAt: now},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "moduleId", "availableAction", "updatedAt"}),
	}).Create(&resources)

	return int(result.RowsAffected), result.Error
}
