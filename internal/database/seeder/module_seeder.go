package seeder

import (
	"time"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	ModuleUserMgmtID = "c7d4e3f2-1a9b-4e8c-a7d6-3b2a1f0e9d8c"
	ModuleSettingsID = "9f8a7b6c-5d4e-4f3a-b2c1-0d9e8f7a6b5c"
)

type moduleSeeder struct {
	enabled bool
}

func NewModuleSeeder() Seeder {
	return &moduleSeeder{enabled: true}
}

func (s *moduleSeeder) Name() string  { return "module_seeder" }
func (s *moduleSeeder) Enabled() bool { return s.enabled }

func (s *moduleSeeder) Seed(db *gorm.DB) (int, error) {
	now := time.Now()
	modules := []model.Module{
		{ID: uuid.MustParse(ModuleUserMgmtID), Name: "User Management", CreatedAt: now, UpdatedAt: now},
		{ID: uuid.MustParse(ModuleSettingsID), Name: "Settings", CreatedAt: now, UpdatedAt: now},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "updatedAt"}),
	}).Create(&modules)

	return int(result.RowsAffected), result.Error
}
