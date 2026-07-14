package seeder

import (
	"time"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	RoleAdminID = "4d8e7b9a-1f6c-4b5d-9e3a-8f7d6c5b4a3e"
	RoleUserID  = "2c5a9d8f-3b7e-4c1f-a6d9-8e7c6b5a4d3f"
)

type roleSeeder struct {
	enabled bool
}

func NewRoleSeeder() Seeder {
	return &roleSeeder{enabled: true}
}

func (s *roleSeeder) Name() string  { return "role_seeder" }
func (s *roleSeeder) Enabled() bool { return s.enabled }

func (s *roleSeeder) Seed(db *gorm.DB) (int, error) {
	now := time.Now()
	roles := []model.Role{
		{ID: uuid.MustParse(RoleAdminID), Name: "admin", Status: true, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.MustParse(RoleUserID), Name: "user", Status: true, CreatedAt: now, UpdatedAt: now},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "status", "updatedAt"}),
	}).Create(&roles)

	return int(result.RowsAffected), result.Error
}
