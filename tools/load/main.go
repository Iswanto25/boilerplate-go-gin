//go:build ignore

package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"

	"github.com/edustack/go-boilerplate/internal/audit"
	settingsModel "github.com/edustack/go-boilerplate/internal/features/settings/model"
	userModel "github.com/edustack/go-boilerplate/internal/features/user/model"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&userModel.User{},
		&settingsModel.Module{},
		&settingsModel.Resource{},
		&settingsModel.Role{},
		&settingsModel.RolePermission{},
		&audit.Logs{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
