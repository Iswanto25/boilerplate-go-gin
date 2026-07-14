package validate

import (
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/edustack/go-boilerplate/internal/features/auth/model"
	"github.com/edustack/go-boilerplate/pkg"
)

func RegisterRequest(req *model.RegisterRequest) *pkg.AppError {
	errs := []string{}

	if strings.TrimSpace(req.Name) == "" {
		errs = append(errs, "Nama wajib diisi")
	} else if len([]rune(req.Name)) < 2 || len([]rune(req.Name)) > 100 {
		errs = append(errs, "Nama harus memiliki panjang 2-100 karakter")
	}

	if strings.TrimSpace(req.Email) == "" {
		errs = append(errs, "Email wajib diisi")
	} else if !isValidEmail(req.Email) {
		errs = append(errs, "Format email tidak valid")
	}

	if req.Password == "" {
		errs = append(errs, "Password wajib diisi")
	} else if len(req.Password) < 8 {
		errs = append(errs, "Password minimal 8 karakter")
	} else if !isStrongPassword(req.Password) {
		errs = append(errs, "Password harus mengandung huruf, angka, dan simbol")
	}

	if req.Role != "" && !req.Role.IsValid() {
		errs = append(errs, "Role harus 'admin' atau 'user'")
	}

	if len(errs) > 0 {
		return &pkg.AppError{
			Code:       40001,
			Message:    strings.Join(errs, "; "),
			StatusCode: http.StatusBadRequest,
		}
	}
	return nil
}

func LoginRequest(req *model.LoginRequest) *pkg.AppError {
	errs := []string{}

	if strings.TrimSpace(req.Email) == "" {
		errs = append(errs, "Email wajib diisi")
	}

	if req.Password == "" {
		errs = append(errs, "Password wajib diisi")
	}

	if len(errs) > 0 {
		return &pkg.AppError{
			Code:       40001,
			Message:    strings.Join(errs, "; "),
			StatusCode: http.StatusBadRequest,
		}
	}
	return nil
}

func RefreshTokenRequest(req *model.RefreshTokenRequest) *pkg.AppError {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return &pkg.AppError{
			Code:       40001,
			Message:    "Refresh token wajib diisi",
			StatusCode: http.StatusBadRequest,
		}
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
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
