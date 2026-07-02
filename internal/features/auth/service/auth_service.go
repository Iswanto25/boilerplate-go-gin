package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/edustack/go-boilerplate/internal/config"
	authModel "github.com/edustack/go-boilerplate/internal/features/auth/model"
	userModel "github.com/edustack/go-boilerplate/internal/features/user/model"
	userRepo "github.com/edustack/go-boilerplate/internal/features/user/repository"
	userService "github.com/edustack/go-boilerplate/internal/features/user/service"
	appErr "github.com/edustack/go-boilerplate/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req *authModel.RegisterRequest) (*authModel.AuthResponse, error)
	Login(ctx context.Context, req *authModel.LoginRequest) (*authModel.AuthResponse, error)
}

type authService struct {
	userService userService.UserService
	userRepo    userRepo.UserRepository
	cfg         *config.Config
}

func NewAuthService(userService userService.UserService, userRepo userRepo.UserRepository, cfg *config.Config) AuthService {
	return &authService{userService: userService, userRepo: userRepo, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, req *authModel.RegisterRequest) (*authModel.AuthResponse, error) {
	createReq := &userModel.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	userResp, err := s.userService.Create(ctx, createReq)
	if err != nil {
		return nil, err
	}

	ttlHours, _ := strconv.Atoi(s.cfg.JWTTTL)
	token, err := s.generateToken(userResp.ID, userResp.Email, ttlHours)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	return &authModel.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		UserID:      userResp.ID,
		Email:       userResp.Email,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *authModel.LoginRequest) (*authModel.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.ErrUnauthorized
		}
		return nil, appErr.ErrInternal
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, appErr.ErrUnauthorized
	}

	ttlHours, _ := strconv.Atoi(s.cfg.JWTTTL)
	token, err := s.generateToken(user.ID, user.Email, ttlHours)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	return &authModel.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		UserID:      user.ID,
		Email:       user.Email,
	}, nil
}

func (s *authService) generateToken(userID uuid.UUID, email string, ttlHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"exp":     time.Now().Add(time.Duration(ttlHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
