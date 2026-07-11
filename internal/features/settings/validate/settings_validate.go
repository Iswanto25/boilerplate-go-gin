package validate

import (
	"strings"

	"github.com/edustack/go-boilerplate/internal/features/settings/model"
	"github.com/edustack/go-boilerplate/pkg/errors"
	"github.com/google/uuid"
)

func CreateModuleRequest(req *model.CreateModuleRequest) *errors.AppError {
	if strings.TrimSpace(req.Name) == "" {
		return errors.NewValidationError("Nama modul wajib diisi")
	}
	return nil
}

func UpdateModuleRequest(req *model.UpdateModuleRequest) *errors.AppError {
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return errors.NewValidationError("Nama modul tidak boleh kosong")
	}
	return nil
}

func CreateResourceRequest(req *model.CreateResourceRequest) *errors.AppError {
	errs := []string{}

	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, "Nama resource wajib diisi")
	}

	if req.ModuleID == uuid.Nil {
		errs = append(errs, "Module ID wajib diisi")
	}

	if len(errs) > 0 {
		return errors.NewValidationError(strings.Join(errs, "; "))
	}
	return nil
}

func UpdateResourceRequest(req *model.UpdateResourceRequest) *errors.AppError {
	errs := []string{}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		errs = append(errs, "Nama resource tidak boleh kosong")
	}

	if req.ModuleID != nil && *req.ModuleID == uuid.Nil {
		errs = append(errs, "Module ID tidak valid")
	}

	if len(errs) > 0 {
		return errors.NewValidationError(strings.Join(errs, "; "))
	}
	return nil
}
