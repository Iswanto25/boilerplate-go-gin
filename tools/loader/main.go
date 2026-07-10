package main

import (
	"fmt"
	"log"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/edustack/go-boilerplate/internal/audit"
	settingsModel "github.com/edustack/go-boilerplate/internal/features/settings/model"
	userModel "github.com/edustack/go-boilerplate/internal/features/user/model"
)

func main() {
	l := gormschema.New("postgres")
	s, err := l.Load(
		&userModel.User{},
		&userModel.Profile{},
		&settingsModel.Module{},
		&settingsModel.Resource{},
		&settingsModel.Role{},
		&settingsModel.RolePermission{},
		&audit.Logs{},
	)
	if err != nil {
		log.Fatalf("failed to load schema: %v", err)
	}
	fmt.Println(s)
}
