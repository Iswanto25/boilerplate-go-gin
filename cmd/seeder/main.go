package main

import (
	"log/slog"
	"os"

	"github.com/edustack/go-boilerplate/internal/config"
	"github.com/edustack/go-boilerplate/internal/database"
	"github.com/edustack/go-boilerplate/internal/database/seeder"
	"github.com/edustack/go-boilerplate/pkg"
)

func main() {
	cfg := config.LoadConfig()

	pkg.InitLogger(cfg.AppEnv)
	slog.Info("running seeders", "env", cfg.AppEnv)

	db := database.NewPostgresConnection(cfg)

	runner := seeder.NewRunner(db)
	runner.Register(
		seeder.NewRoleSeeder(),
		seeder.NewModuleSeeder(),
		seeder.NewResourceSeeder(),
		seeder.NewRolePermissionSeeder(),
	)
	runner.RunAll()

	os.Exit(0)
}
