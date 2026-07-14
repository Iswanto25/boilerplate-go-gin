package service

import (
	"context"
	"errors"

	"github.com/edustack/go-boilerplate/internal/config"
	authModel "github.com/edustack/go-boilerplate/internal/features/auth/model"
	userModel "github.com/edustack/go-boilerplate/internal/features/user/model"
	userRepo "github.com/edustack/go-boilerplate/internal/features/user/repository"
	userService "github.com/edustack/go-boilerplate/internal/features/user/service"
	"github.com/edustack/go-boilerplate/pkg"
	"github.com/edustack/go-boilerplate/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req *authModel.RegisterRequest) (*authModel.AuthResponse, error)
	Login(ctx context.Context, req *authModel.LoginRequest) (*authModel.AuthResponse, error)
	RefreshToken(ctx context.Context, userID string) (*authModel.AuthResponse, error)
	Logout(ctx context.Context, userID string) error
	Profile(ctx context.Context, userID string) (*userModel.UserResponse, error)
}

type authService struct {
	userService userService.UserService
	userRepo    userRepo.UserRepository
	jwtUtils    *jwt.JWTUtils
	cfg         *config.Config
}

func NewAuthService(userService userService.UserService, userRepo userRepo.UserRepository, jwtUtils *jwt.JWTUtils, cfg *config.Config) AuthService {
	return &authService{userService: userService, userRepo: userRepo, jwtUtils: jwtUtils, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, req *authModel.RegisterRequest) (*authModel.AuthResponse, error) {
	createReq := &userModel.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     string(req.Role),
	}

	userResp, err := s.userService.Create(ctx, createReq)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"userId": userResp.ID.String(),
		"email":  userResp.Email,
		"name":   userResp.Name,
		"role":   userResp.Role,
	}

	accessToken, err := s.jwtUtils.GenerateAndStoreAccessToken(ctx, userResp.ID.String(), payload)
	if err != nil {
		return nil, pkg.ErrInternal.WithCause(err)
	}

	refreshToken, err := s.jwtUtils.GenerateAndStoreRefreshToken(ctx, userResp.ID.String(), payload)
	if err != nil {
		return nil, pkg.ErrInternal.WithCause(err)
	}

	return &authModel.AuthResponse{
		UserID:       userResp.ID,
		Email:        userResp.Email,
		Name:         userResp.Name,
		Role:         userModel.Role(userResp.Role),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *authModel.LoginRequest) (*authModel.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrUnauthorized
		}
		return nil, pkg.ErrInternal.WithCause(err)
	}

	password := req.Password + s.cfg.Salt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, pkg.ErrInvalidCredentials
	}

	payload := map[string]any{
		"userId": user.ID.String(),
		"email":  user.Email,
		"name":   user.Name,
		"role":   user.Role.Name,
	}

	accessToken, err := s.jwtUtils.GenerateAndStoreAccessToken(ctx, user.ID.String(), payload)
	if err != nil {
		return nil, pkg.ErrInternal.WithCause(err)
	}

	refreshToken, err := s.jwtUtils.GenerateAndStoreRefreshToken(ctx, user.ID.String(), payload)
	if err != nil {
		return nil, pkg.ErrInternal.WithCause(err)
	}

	result := &authModel.AuthResponse{
		UserID:       user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         userModel.Role(user.Role.Name),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return result, nil
}

func (s *authService) RefreshToken(ctx context.Context, userID string) (*authModel.AuthResponse, error) {
	refreshToken, err := s.jwtUtils.GetStoredRefreshToken(ctx, userID)
	if err != nil || refreshToken == "" {
		return nil, pkg.ErrUnauthorized
	}

	if valid, _ := s.jwtUtils.ValidateRefreshTokenInStore(ctx, userID, refreshToken); !valid {
		return nil, pkg.ErrUnauthorized
	}

	_ = s.jwtUtils.RevokeUserTokens(ctx, userID)

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	user, err := s.userRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, pkg.ErrInternal.WithCause(err)
	}

	payload := map[string]any{
		"userId": user.ID.String(),
		"email":  user.Email,
		"name":   user.Name,
		"role":   user.Role.Name,
	}

	accessToken, err := s.jwtUtils.GenerateAndStoreAccessToken(ctx, user.ID.String(), payload)
	if err != nil {
		return nil, pkg.ErrInternal.WithCause(err)
	}

	newRefreshToken, err := s.jwtUtils.GenerateAndStoreRefreshToken(ctx, user.ID.String(), payload)
	if err != nil {
		return nil, pkg.ErrInternal.WithCause(err)
	}

	result := &authModel.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	return result, nil
}

func (s *authService) Profile(ctx context.Context, userID string) (*userModel.UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	user, err := s.userService.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	result := userModel.ToUserResponse(user)
	return &result, nil
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	return s.jwtUtils.RevokeUserTokens(ctx, userID)
}
