package services

import (
	"fmt"
	"log"
	"strings"
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func FetchCryptoPrice(cryptoSymbols []string, currency string) (map[string]float64, error) {
	client := resty.New()
	cryptoQuery := strings.Join(cryptoSymbols, ",")

	cacheKey := fmt.Sprintf("prices:%s:%s", cryptoQuery, currency)
	cachedData, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		var cachedPrices map[string]float64
		json.Unmarshal([]byte(cachedData), &cachedPrices)
		log.Println("Returning cached price data")
		return cachedPrices, nil
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", cryptoQuery, currency)

	resp, err := client.R().Get(url)
	if err != nil {
		log.Println("API Error, returning fallback data")
		return map[string]float64{"bitcoin": -1, "ethereum": -1, "monero": -1}, nil
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("API Request Failed %s", resp.Status())
	}

	priceData := make(map[string]float64)
	for _, crypto := range cryptoSymbols {
		price := gjson.Get(resp.String(), fmt.Sprintf("%s.%s", crypto, currency))
		if price.Exists() {
			priceData[crypto] = price.Float()
		} else {
			log.Printf("Price for %s not found", crypto)
		}
	}

	cacheData, _ := json.Marshal(priceData)
	redisClient.Set(ctx, cacheKey, cacheData, 30*time.Second)

	return priceData, nil
}