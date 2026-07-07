package tokenstore

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	accessTokenPrefix  = "access_token:"
	refreshTokenPrefix = "refresh_token:"
)

type TokenStore struct {
	client      *redis.Client
	isAvailable bool
}

func NewTokenStore(client *redis.Client, isAvailable bool) *TokenStore {
	return &TokenStore{client: client, isAvailable: isAvailable}
}

func (s *TokenStore) StoreAccessToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	if !s.isAvailable || s.client == nil {
		slog.Warn("Token storage skipped - Redis not available")
		return nil
	}

	key := accessTokenPrefix + userID
	if err := s.client.Set(ctx, key, token, ttl).Err(); err != nil {
		slog.Warn("Failed to store access token in Redis", "error", err)
		return nil
	}
	return nil
}

func (s *TokenStore) StoreRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	if !s.isAvailable || s.client == nil {
		slog.Warn("Token storage skipped - Redis not available")
		return nil
	}

	key := refreshTokenPrefix + userID
	if err := s.client.Set(ctx, key, token, ttl).Err(); err != nil {
		slog.Warn("Failed to store refresh token in Redis", "error", err)
		return nil
	}
	return nil
}

func (s *TokenStore) GetAccessToken(ctx context.Context, userID string) (string, error) {
	if !s.isAvailable || s.client == nil {
		slog.Warn("Token retrieval skipped - Redis not available")
		return "", nil
	}

	key := accessTokenPrefix + userID
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		slog.Warn("Failed to get access token from Redis", "error", err)
		return "", nil
	}
	return val, nil
}

func (s *TokenStore) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	if !s.isAvailable || s.client == nil {
		slog.Warn("Token retrieval skipped - Redis not available")
		return "", nil
	}

	key := refreshTokenPrefix + userID
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		slog.Warn("Failed to get refresh token from Redis", "error", err)
		return "", nil
	}
	return val, nil
}

func (s *TokenStore) DeleteAccessToken(ctx context.Context, userID string) error {
	if !s.isAvailable || s.client == nil {
		slog.Warn("Token deletion skipped - Redis not available")
		return nil
	}

	key := accessTokenPrefix + userID
	if err := s.client.Del(ctx, key).Err(); err != nil {
		slog.Warn("Failed to delete access token from Redis", "error", err)
	}
	return nil
}

func (s *TokenStore) DeleteRefreshToken(ctx context.Context, userID string) error {
	if !s.isAvailable || s.client == nil {
		slog.Warn("Token deletion skipped - Redis not available")
		return nil
	}

	key := refreshTokenPrefix + userID
	if err := s.client.Del(ctx, key).Err(); err != nil {
		slog.Warn("Failed to delete refresh token from Redis", "error", err)
	}
	return nil
}

func (s *TokenStore) DeleteAllTokens(ctx context.Context, userID string) error {
	if !s.isAvailable || s.client == nil {
		slog.Warn("Token deletion skipped - Redis not available")
		return nil
	}

	if err := s.DeleteAccessToken(ctx, userID); err != nil {
		return err
	}
	return s.DeleteRefreshToken(ctx, userID)
}
