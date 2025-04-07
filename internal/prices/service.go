package prices

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tidwall/gjson"
	"github.com/undersleep7x/cryo-project/internal/infra/cache"
)

type FetchCryptoPriceService interface {
	FetchCryptoPrice(cryptoSymbols []string, currency string) (map[string]float64, error)
}

type fetchCryptoPriceServiceImpl struct{
	Cache PricesCache
	config Config
}

func NewFetchCryptoPriceService(cache *cache.PriceCache, cfg Config) FetchCryptoPriceService {
	return &fetchCryptoPriceServiceImpl{Cache: cache, config: cfg}
}

func(s *fetchCryptoPriceServiceImpl) FetchCryptoPrice(cryptoSymbols []string, currency string) (map[string]float64, error) {
	// kick off redis context and close at the end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	priceData := make(map[string]float64) // init return variable
	var missingCryptos []string           // init missingcrypto variable

	// loop through array and check cache for any saved data
	for _, crypto := range cryptoSymbols {
		cacheKey := fmt.Sprintf("prices:%s:%s", crypto, currency)
		cachedData, err :=s.Cache.GetCachedPrices(ctx, cacheKey)

		if err == nil {
			var cachedPrice map[string]float64
			err := json.Unmarshal([]byte(cachedData), &cachedPrice)
			if err != nil {
				log.Printf("Failed to parse Redis JSON for %s: %v", crypto, err)
			} else {
				log.Printf("Successfully retrieved cached price data for %s", crypto)
				priceData[crypto] = cachedPrice[crypto]
			}
		} else if err.Error() == "redis: nil" { // if any error, add crypto to missing array and move it api
			log.Printf("No cache for %s in Redis cache, fetching with API", crypto)
			missingCryptos = append(missingCryptos, crypto)
		} else {
			log.Printf("Redis error for %s: %v", crypto, err)
			missingCryptos = append(missingCryptos, crypto)
		}
	}

	if len(missingCryptos) > 0 { // if any were not in cache

		pricesCall, err := FetchPrices(missingCryptos, currency, s.config.BaseURL, s.config.Timeout) // make api call for remaining cryptos

		if err != nil { // set fallback prices if api call fails entirely
			log.Printf("API failure, setting fallback prices: %v", err)
			for _, crypto := range missingCryptos {
				priceData[crypto] = -1
			}
		} else {
			for _, crypto := range missingCryptos {
				price := gjson.Get(pricesCall.String(), fmt.Sprintf("%s.%s", crypto, currency))
				if price.Exists() { // add existing prices to return value
					priceData[crypto] = price.Float()
					cacheKey := fmt.Sprintf("prices:%s:%s", crypto, currency)
					cachedEntry, _ := json.Marshal(map[string]float64{crypto: price.Float()}) // after adding, cache value
					err := s.Cache.CachePrices(ctx, cacheKey, cachedEntry, 30*time.Second)
					if err != nil {
						log.Printf("Failed to cache price for %s: %v", crypto, err)
					}
				} else {
					// if not found in api call or cache, set fallback value
					log.Printf("Price for %s not found in API, setting fallback price", crypto)
					priceData[crypto] = -1
				}
			}
		}
	}

	return priceData, nil
}
