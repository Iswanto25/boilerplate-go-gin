package service

import (
	"context"
	"errors"

	"github.com/edustack/go-boilerplate/internal/config"
	authModel "github.com/edustack/go-boilerplate/internal/features/auth/model"
	userModel "github.com/edustack/go-boilerplate/internal/features/user/model"
	userRepo "github.com/edustack/go-boilerplate/internal/features/user/repository"
	userService "github.com/edustack/go-boilerplate/internal/features/user/service"
	appErr "github.com/edustack/go-boilerplate/pkg/errors"
	"github.com/edustack/go-boilerplate/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req *authModel.RegisterRequest) (*authModel.AuthResponse, error)
	Login(ctx context.Context, req *authModel.LoginRequest) (*authModel.AuthResponse, error)
	RefreshToken(ctx context.Context, req *authModel.RefreshTokenRequest) (*authModel.AuthResponse, error)
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
		Role:     req.Role,
	}

	userResp, err := s.userService.Create(ctx, createReq)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"userId": userResp.ID.String(),
		"email":  userResp.Email,
		"name":   userResp.Name,
		"role":   string(userResp.Role),
	}

	accessToken, err := s.jwtUtils.GenerateAndStoreAccessToken(ctx, userResp.ID.String(), payload)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	refreshToken, err := s.jwtUtils.GenerateAndStoreRefreshToken(ctx, userResp.ID.String(), payload)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	return &authModel.AuthResponse{
		UserID:       userResp.ID,
		Email:        userResp.Email,
		Name:         userResp.Name,
		Role:         userResp.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

	password := req.Password + s.cfg.Salt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, appErr.ErrUnauthorized
	}

	payload := map[string]any{
		"userId": user.ID.String(),
		"email":   user.Email,
		"name":    user.Name,
		"role":    string(user.Role),
	}

	accessToken, err := s.jwtUtils.GenerateAndStoreAccessToken(ctx, user.ID.String(), payload)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	refreshToken, err := s.jwtUtils.GenerateAndStoreRefreshToken(ctx, user.ID.String(), payload)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	result := &authModel.AuthResponse{
		UserID:       user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return result, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *authModel.RefreshTokenRequest) (*authModel.AuthResponse, error) {
	claims, err := s.jwtUtils.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, appErr.ErrUnauthorized
	}

	userIDStr, _ := claims["user_id"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, appErr.ErrUnauthorized
	}

	if valid, _ := s.jwtUtils.ValidateRefreshTokenInStore(ctx, userID.String(), req.RefreshToken); !valid {
		return nil, appErr.ErrUnauthorized
	}

	_ = s.jwtUtils.RevokeUserTokens(ctx, userID.String())

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	payload := map[string]any{
		"userId": user.ID.String(),
		"email":   user.Email,
		"name":    user.Name,
		"role":    string(user.Role),
	}

	accessToken, err := s.jwtUtils.GenerateAndStoreAccessToken(ctx, user.ID.String(), payload)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	refreshToken, err := s.jwtUtils.GenerateAndStoreRefreshToken(ctx, user.ID.String(), payload)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	result := &authModel.AuthResponse{
		UserID:       user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return result, nil
}

func (s *authService) Profile(ctx context.Context, userID string) (*userModel.UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, appErr.ErrUnauthorized
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