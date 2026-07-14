package service

import (
	"context"
	"errors"

	settingsRepo "github.com/edustack/go-boilerplate/internal/features/settings/repository"
	"github.com/edustack/go-boilerplate/internal/features/user/model"
	"github.com/edustack/go-boilerplate/internal/features/user/repository"
	"github.com/edustack/go-boilerplate/pkg"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
}

type userService struct {
	repo         repository.UserRepository
	settingsRepo settingsRepo.SettingsRepository
	bcryptRounds int
	salt         string
}

func NewUserService(repo repository.UserRepository, bcryptRounds int, salt string, sr settingsRepo.SettingsRepository) UserService {
	return &userService{repo: repo, bcryptRounds: bcryptRounds, salt: salt, settingsRepo: sr}
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrUserNotFound
		}
		return nil, pkg.ErrInternal
	}
	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, pkg.ErrInternal
	}
	return users, nil
}

func (s *userService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	roleName := req.Role
	if roleName == "" {
		roleName = string(model.RoleUser)
	}

	role, err := s.settingsRepo.FindRoleByName(ctx, roleName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NewValidationError("role not found")
		}
		return nil, pkg.ErrInternal
	}

	existing, findErr := s.repo.FindByEmail(ctx, req.Email)
	if findErr == nil && existing != nil {
		return nil, pkg.NewValidationError("email already exists")
	}
	if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, pkg.ErrInternal
	}

	password := req.Password + s.salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptRounds)
	if err != nil {
		return nil, pkg.ErrInternal
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, pkg.NewValidationError("email already exists")
		}
		return nil, pkg.ErrInternal
	}

	user.Role = *role
	resp := model.ToUserResponse(user)
	return &resp, nil
}
