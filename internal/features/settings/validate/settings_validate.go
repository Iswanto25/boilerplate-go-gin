package validate

import (
	"strings"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/edustack/go-boilerplate/pkg"
	"github.com/google/uuid"
)

func CreateModuleRequest(req *model.CreateModuleRequest) *pkg.AppError {
	if strings.TrimSpace(req.Name) == "" {
		return pkg.NewValidationError("Nama modul wajib diisi")
	}
	return nil
}

func UpdateModuleRequest(req *model.UpdateModuleRequest) *pkg.AppError {
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return pkg.NewValidationError("Nama modul tidak boleh kosong")
	}
	return nil
}

func CreateResourceRequest(req *model.CreateResourceRequest) *pkg.AppError {
	errs := []string{}

	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, "Nama resource wajib diisi")
	}

	if req.ModuleID == uuid.Nil {
		errs = append(errs, "Module ID wajib diisi")
	}

	if len(errs) > 0 {
		return pkg.NewValidationError(strings.Join(errs, "; "))
	}
	return nil
}

func UpdateResourceRequest(req *model.UpdateResourceRequest) *pkg.AppError {
	errs := []string{}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		errs = append(errs, "Nama resource tidak boleh kosong")
	}

	if req.ModuleID != nil && *req.ModuleID == uuid.Nil {
		errs = append(errs, "Module ID tidak valid")
	}

	if len(errs) > 0 {
		return pkg.NewValidationError(strings.Join(errs, "; "))
	}
	return nil
}

func CreateRoleRequest(req *model.CreateRoleRequest) *pkg.AppError {
	if strings.TrimSpace(req.Name) == "" {
		return pkg.NewValidationError("Nama role wajib diisi")
	}
	return nil
}

func UpdateRoleRequest(req *model.UpdateRoleRequest) *pkg.AppError {
	errs := []string{}
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		errs = append(errs, "Nama role tidak boleh kosong")
	}
	if len(errs) > 0 {
		return pkg.NewValidationError(strings.Join(errs, "; "))
	}
	return nil
}

func CreateRolePermissionRequest(req *model.CreateRolePermissionRequest) *pkg.AppError {
	errs := []string{}
	if req.RoleID == uuid.Nil {
		errs = append(errs, "Role ID wajib diisi")
	}
	if req.ResourceID == uuid.Nil {
		errs = append(errs, "Resource ID wajib diisi")
	}
	if len(errs) > 0 {
		return pkg.NewValidationError(strings.Join(errs, "; "))
	}
	return nil
}

func UpdateRolePermissionRequest(req *model.UpdateRolePermissionRequest) *pkg.AppError {
	errs := []string{}
	if req.RoleID != nil && *req.RoleID == uuid.Nil {
		errs = append(errs, "Role ID tidak valid")
	}
	if req.ResourceID != nil && *req.ResourceID == uuid.Nil {
		errs = append(errs, "Resource ID tidak valid")
	}
	if len(errs) > 0 {
		return pkg.NewValidationError(strings.Join(errs, "; "))
	}
	return nil
}
