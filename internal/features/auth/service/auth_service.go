package service

import (
	"context"
	"errors"
	"time"

	"github.com/edustack/go-boilerplate/internal/config"
	authModel "github.com/edustack/go-boilerplate/internal/features/auth/model"
	userModel "github.com/edustack/go-boilerplate/internal/features/user/model"
	userRepo "github.com/edustack/go-boilerplate/internal/features/user/repository"
	userService "github.com/edustack/go-boilerplate/internal/features/user/service"
	appErr "github.com/edustack/go-boilerplate/pkg/errors"
	"github.com/edustack/go-boilerplate/pkg/jwt"
	"github.com/edustack/go-boilerplate/pkg/tokenstore"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req *authModel.RegisterRequest) (*authModel.AuthResponse, error)
	Login(ctx context.Context, req *authModel.LoginRequest) (*authModel.AuthResponse, error)
	RefreshToken(ctx context.Context, req *authModel.RefreshTokenRequest) (*authModel.AuthResponse, error)
	Logout(ctx context.Context, userID string) error
}

type authService struct {
	userService userService.UserService
	userRepo    userRepo.UserRepository
	jwtUtils    *jwt.JWTUtils
	tokenStore  *tokenstore.TokenStore
	cfg         *config.Config
}

func NewAuthService(userService userService.UserService, userRepo userRepo.UserRepository, jwtUtils *jwt.JWTUtils, tokenStore *tokenstore.TokenStore, cfg *config.Config) AuthService {
	return &authService{userService: userService, userRepo: userRepo, jwtUtils: jwtUtils, tokenStore: tokenStore, cfg: cfg}
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

	return s.generateAuthResponse(ctx, userResp.ID, userResp.Email)
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

	return s.generateAuthResponse(ctx, user.ID, user.Email)
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

	// Verify token exists in Redis
	storedToken, _ := s.tokenStore.GetRefreshToken(ctx, userID.String())
	if storedToken == "" || storedToken != req.RefreshToken {
		return nil, appErr.ErrUnauthorized
	}

	// Delete old tokens
	_ = s.tokenStore.DeleteAllTokens(ctx, userID.String())

	// Get user email
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	return s.generateAuthResponse(ctx, user.ID, user.Email)
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	return s.tokenStore.DeleteAllTokens(ctx, userID)
}

func (s *authService) generateAuthResponse(ctx context.Context, userID uuid.UUID, email string) (*authModel.AuthResponse, error) {
	payload := map[string]interface{}{
		"user_id": userID.String(),
		"email":   email,
	}

	accessToken, err := s.jwtUtils.GenerateAccessToken(payload)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	refreshToken, err := s.jwtUtils.GenerateRefreshToken(payload)
	if err != nil {
		return nil, appErr.ErrInternal
	}

	// Store tokens in Redis
	accessTTL := time.Duration(s.cfg.JWTTTL) * time.Hour
	refreshTTL := time.Duration(s.cfg.JWTRefreshTTL) * time.Hour

	_ = s.tokenStore.StoreAccessToken(ctx, userID.String(), accessToken, accessTTL)
	_ = s.tokenStore.StoreRefreshToken(ctx, userID.String(), refreshToken, refreshTTL)

	return &authModel.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userID,
		Email:        email,
	}, nil
}
