package cache

import (
	"context"
	"time"
	platformRedis "github.com/undersleep7x/cryo-project/internal/platform/redisstore"
)

type PriceCache struct {
	Redis platformRedis.RedisClient
}

func NewPriceCache (client platformRedis.RedisClient) *PriceCache {
	return &PriceCache{Redis: client}
}

func (c *PriceCache) GetCachedPrices(ctx context.Context, cacheKey string) (string, error) {
	return c.Redis.Get(ctx, cacheKey)
}

func (c *PriceCache) CachePrices(ctx context.Context, cacheKey string, value interface{}, ttl time.Duration) error {
	return c.Redis.Set(ctx, cacheKey, value, ttl)
}