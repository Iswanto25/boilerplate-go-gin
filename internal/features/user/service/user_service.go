package service

import (
	"context"

	"github.com/edustack/go-boilerplate/internal/features/user/model"
	"github.com/edustack/go-boilerplate/internal/features/user/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	return s.repo.FindByID(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *userService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := model.ToUserResponse(user)
	return &resp, nil
}
