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
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, appErr.ErrInternal
	}

	resp := model.ToUserResponse(user)
	return &resp, nil
}
