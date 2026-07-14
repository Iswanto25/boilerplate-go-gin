package seeder

import (
	"time"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type rolePermissionSeeder struct {
	enabled bool
}

func NewRolePermissionSeeder() Seeder {
	return &rolePermissionSeeder{enabled: true}
}

func (s *rolePermissionSeeder) Name() string    { return "role_permission_seeder" }
func (s *rolePermissionSeeder) Enabled() bool   { return s.enabled }

func (s *rolePermissionSeeder) Seed(db *gorm.DB) (int, error) {
	now := time.Now()
	adminID := uuid.MustParse(RoleAdminID)
	userID := uuid.MustParse(RoleUserID)

	crud := []string{"create", "read", "update", "delete"}

	adminPerms := []model.RolePermission{
		{ID: uuid.New(), RoleId: adminID, ResourceId: uuid.MustParse(ResourceUsersID), GrantedActions: crud, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.New(), RoleId: adminID, ResourceId: uuid.MustParse(ResourceProfilesID), GrantedActions: crud, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.New(), RoleId: adminID, ResourceId: uuid.MustParse(ResourceRolesID), GrantedActions: crud, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.New(), RoleId: adminID, ResourceId: uuid.MustParse(ResourceModulesID), GrantedActions: crud, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.New(), RoleId: adminID, ResourceId: uuid.MustParse(ResourceResourcesID), GrantedActions: crud, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.New(), RoleId: adminID, ResourceId: uuid.MustParse(ResourceRolePermsID), GrantedActions: crud, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.New(), RoleId: adminID, ResourceId: uuid.MustParse(ResourceAuditLogsID), GrantedActions: []string{"read"}, CreatedAt: now, UpdatedAt: now},
	}

	userPerms := []model.RolePermission{
		{ID: uuid.New(), RoleId: userID, ResourceId: uuid.MustParse(ResourceProfilesID), GrantedActions: []string{"read", "update"}, CreatedAt: now, UpdatedAt: now},
	}

	all := append(adminPerms, userPerms...)

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "roleId"}, {Name: "resourceId"}},
		DoUpdates: clause.AssignmentColumns([]string{"grantedActions", "updatedAt"}),
	}).Create(&all)

	return int(result.RowsAffected), result.Error
}
