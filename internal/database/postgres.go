package database

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/edustack/go-boilerplate/internal/audit"
	"github.com/edustack/go-boilerplate/internal/config"
	settingsModel "github.com/edustack/go-boilerplate/internal/features/settings/model"
	userModel "github.com/edustack/go-boilerplate/internal/features/user/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresConnection(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	gormLogLevel := logger.Warn
	if cfg.AppEnv != "production" {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if cfg.AppEnv != "production" {
		if err := db.AutoMigrate(
			&userModel.User{},
			&settingsModel.Module{},
			&settingsModel.Resource{},
			&settingsModel.Role{},
			&settingsModel.RolePermission{},
			&audit.Logs{},
		); err != nil {
			log.Fatalf("Failed to auto migrate: %v", err)
		}
	} else {
		slog.Warn("production mode — AutoMigrate disabled, apply SQL from migrations/ manually")
	}

	log.Println("Connected to PostgreSQL successfully")
	return db
}
