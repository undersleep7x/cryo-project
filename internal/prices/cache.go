package prices

import (
	"time"
	"context"
)

type PricesCache interface {
	GetCachedPrices(ctx context.Context, cacheKey string) (string, error)
	CachePrices(ctx context.Context, cacheKey string, value interface{}, ttl time.Duration) error
}