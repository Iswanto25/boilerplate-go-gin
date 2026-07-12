package database

import (
	"fmt"
	"time"

	"github.com/edustack/go-boilerplate/internal/audit"
	"github.com/edustack/go-boilerplate/internal/config"
	settingsModel "github.com/edustack/go-boilerplate/internal/features/settings/model"
	userModel "github.com/edustack/go-boilerplate/internal/features/user/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"log/slog"
)

func NewPostgresConnection(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get sql.DB", "error", err)
		panic(err)
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
			slog.Error("AutoMigrate failed", "error", err)
			panic(err)
		}
		slog.Info("AutoMigrate completed")
	} else {
		slog.Warn("production mode — AutoMigrate disabled, apply SQL from migrations/ manually")
	}

	slog.Info("Connected to PostgreSQL successfully")
	return db
}
