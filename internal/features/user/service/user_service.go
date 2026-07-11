package service

import (
	"context"
	"errors"

	"github.com/edustack/go-boilerplate/internal/features/user/model"
	"github.com/edustack/go-boilerplate/internal/features/user/repository"
	appErr "github.com/edustack/go-boilerplate/pkg/errors"
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
	bcryptRounds int
	salt         string
}

func NewUserService(repo repository.UserRepository, bcryptRounds int, salt string) UserService {
	return &userService{repo: repo, bcryptRounds: bcryptRounds, salt: salt}
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.ErrUserNotFound
		}
		return nil, appErr.ErrInternal
	}
	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, appErr.ErrInternal
	}
	return users, nil
}

func (s *userService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	password := req.Password + s.salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptRounds)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	role := req.Role
	if role == "" {
		role = model.RoleUser
	}
	if !role.IsValid() {
		return nil, appErr.NewValidationError("role must be 'admin' or 'user'")
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, appErr.ErrConflict
		}
		return nil, appErr.ErrInternal
	}

	resp := model.ToUserResponse(user)
	return &resp, nil
}
