package validate

import (
	"strings"
	"unicode"

	"github.com/edustack/go-boilerplate/internal/features/user/model"
	"github.com/edustack/go-boilerplate/pkg/errors"
)

func CreateUserRequest(req *model.CreateUserRequest) *errors.AppError {
	errs := []string{}

	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, "Nama wajib diisi")
	} else if len([]rune(req.Name)) < 2 || len([]rune(req.Name)) > 100 {
		errs = append(errs, "Nama harus memiliki panjang 2-100 karakter")
	}

	if strings.TrimSpace(req.Email) == "" {
		errs = append(errs, "Email wajib diisi")
	}

	if req.Password == "" {
		errs = append(errs, "Password wajib diisi")
	} else if len(req.Password) < 8 {
		errs = append(errs, "Password minimal 8 karakter")
	} else if !isStrongPassword(req.Password) {
		errs = append(errs, "Password harus mengandung huruf, angka, dan simbol")
	}

	if len(errs) > 0 {
		return errors.NewValidationError(strings.Join(errs, "; "))
	}
	return nil
}

func isStrongPassword(password string) bool {
	hasLetter := false
	hasDigit := false
	hasSymbol := false

	for _, ch := range password {
		switch {
		case unicode.IsLetter(ch):
			hasLetter = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsSymbol(ch) || unicode.IsPunct(ch):
			hasSymbol = true
		}
	}

	return hasLetter && hasDigit && hasSymbol
}
