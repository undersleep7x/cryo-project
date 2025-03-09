package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis struct {
		Address string `yaml:"address"`
	} `yaml:"redis"`
}

func loadConfig() Config {
	var config Config
	data, err := os.ReadFile("config/local_config.yml")
	if err != nil {
		log.Fatalf("Failed to read config file %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}
	return config
}

var config = loadConfig()
var redisClient = redis.NewClient(&redis.Options{
	Addr: config.Redis.Address,
})

func FetchCryptoPrice(cryptoSymbols []string, currency string) (map[string]float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := resty.New()
	priceData := make(map[string]float64)
	var missingCryptos []string

	for _, crypto := range cryptoSymbols {
		cacheKey := fmt.Sprintf("prices:%s:%s", crypto, currency)
		cachedData, err := redisClient.Get(ctx, cacheKey).Result()

		if err == nil {
			var cachedPrice map[string]float64
			err := json.Unmarshal([]byte(cachedData), &cachedPrice)

			if err != nil {
				log.Printf("Failed to parse Redis JSON for %s: %v", crypto, err)
			} else {
				log.Printf("Successfully retrieved cached price data for %s", crypto)
				priceData[crypto] = cachedPrice[crypto]
				continue
			}
		} else if err == redis.Nil {
			log.Printf("No cache for %s in Redis cache, fetching with API", crypto)
			missingCryptos = append(missingCryptos, crypto)
		} else {
			log.Printf("Redis error for %s: %v", crypto, err)
			missingCryptos = append(missingCryptos, crypto)
		}
	}

	if len(missingCryptos) > 0 {
		apiQuery := strings.Join(missingCryptos, ",")
		url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", apiQuery, currency)
		resp, err := client.R().Get(url)

		if err != nil || !resp.IsSuccess() {
			log.Printf("API failure, some prices may be missing: %v", err)
			for _, crypto := range missingCryptos {
				priceData[crypto] = -1
			}
		} else {
			for _, crypto := range missingCryptos {
				price := gjson.Get(resp.String(), fmt.Sprintf("%s.%s", crypto, currency))
				if price.Exists() {
					priceData[crypto] = price.Float()
					cacheKey := fmt.Sprintf("prices:%s:%s", crypto, currency)
					cachedEntry, _ := json.Marshal(map[string]float64{crypto: price.Float()})
					err := redisClient.Set(ctx, cacheKey, cachedEntry, 30*time.Second).Err()
					if err != nil {
						log.Printf("Failed to cache price for %s: %v", crypto, err)
					}
				} else {
					log.Printf("Price for %s not found in API", crypto)
					priceData[crypto] = -1
				}
			}
		}
	}

	return priceData, nil
}