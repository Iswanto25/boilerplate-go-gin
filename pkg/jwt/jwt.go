package jwt

import (
	"context"
	"time"

	"github.com/edustack/go-boilerplate/pkg/tokenstore"
	"github.com/golang-jwt/jwt/v5"
)

type JWTUtils struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	tokenStore    *tokenstore.TokenStore
}

func NewJWTUtils(accessSecret, refreshSecret string, accessTTLHours, refreshTTLHours int, ts *tokenstore.TokenStore) *JWTUtils {
	return &JWTUtils{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     time.Duration(accessTTLHours) * time.Hour,
		refreshTTL:    time.Duration(refreshTTLHours) * time.Hour,
		tokenStore:    ts,
	}
}

func (j *JWTUtils) GenerateAndStoreAccessToken(ctx context.Context, userID string, payload map[string]interface{}) (string, error) {
	token, err := j.GenerateAccessToken(payload)
	if err != nil {
		return "", err
	}
	if j.tokenStore != nil {
		_ = j.tokenStore.StoreAccessToken(ctx, userID, token, j.accessTTL)
	}
	return token, nil
}

func (j *JWTUtils) GenerateAndStoreRefreshToken(ctx context.Context, userID string, payload map[string]interface{}) (string, error) {
	token, err := j.GenerateRefreshToken(payload)
	if err != nil {
		return "", err
	}
	if j.tokenStore != nil {
		_ = j.tokenStore.StoreRefreshToken(ctx, userID, token, j.refreshTTL)
	}
	return token, nil
}

func (j *JWTUtils) ValidateAccessTokenInStore(ctx context.Context, userID, tokenString string) (bool, error) {
	if j.tokenStore == nil {
		return true, nil
	}
	stored, err := j.tokenStore.GetAccessToken(ctx, userID)
	if err != nil {
		return false, err
	}
	return stored == "" || stored == tokenString, nil
}

func (j *JWTUtils) ValidateRefreshTokenInStore(ctx context.Context, userID, tokenString string) (bool, error) {
	if j.tokenStore == nil {
		return false, nil
	}
	stored, err := j.tokenStore.GetRefreshToken(ctx, userID)
	if err != nil {
		return false, err
	}
	return stored == tokenString, nil
}

func (j *JWTUtils) RevokeUserTokens(ctx context.Context, userID string) error {
	if j.tokenStore == nil {
		return nil
	}
	return j.tokenStore.DeleteAllTokens(ctx, userID)
}

func (j *JWTUtils) GetStoredRefreshToken(ctx context.Context, userID string) (string, error) {
	if j.tokenStore == nil {
		return "", nil
	}
	return j.tokenStore.GetRefreshToken(ctx, userID)
}

func (j *JWTUtils) GenerateAccessToken(payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(j.accessTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessSecret))
}

func (j *JWTUtils) GenerateRefreshToken(payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(j.refreshTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshSecret))
}

func (j *JWTUtils) VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func (j *JWTUtils) ParseAccessTokenNoExpiry(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.accessSecret), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func (j *JWTUtils) VerifyRefreshToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.refreshSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
