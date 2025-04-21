package cache

import (
	"context"
	"time"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Ping(ctx context.Context) error
}