package cache

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type clientWrapper struct {
	Client *redis.Client
}

func NewRedisClientWrapper(client *redis.Client) RedisClient {
	return &clientWrapper{Client: client}
}

func (r *clientWrapper) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
func (r *clientWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}
func (r *clientWrapper) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
