package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a Redis client with graceful degradation.
// If Redis is unavailable, the app still runs without cache/token storage.
// Returns the client and a boolean indicating availability.
func NewRedisClient(host, port, password string, db int) (*redis.Client, bool) {
	if host == "" || port == "" {
		slog.Warn("Redis configuration not found (REDIS_HOST/REDIS_PORT) - running without Redis")
		return nil, false
	}

	addr := host + ":" + port
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   3,
	})

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("Redis not available - running without Redis cache", "error", err)
		return nil, false
	}

	slog.Info("Redis connected", "addr", addr)
	return client, true
}
