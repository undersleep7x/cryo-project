package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/undersleep7x/cryo-project/app"
)

// fetch prices from coingecko api
var FetchPrices = func(cryptos []string, currency string) (*resty.Response, error) {
	// setup resty client for api call
	client := resty.New()
	timeout := time.Duration(app.Config.CoinGecko.Timeout) * time.Second
	client.SetTimeout(timeout)

	// turn passed in array into csv string
	apiQuery := strings.Join(cryptos, ",")
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", app.Config.CoinGecko.BaseURL, apiQuery, currency)
	log.Printf("Making API call to %s", url)
	resp, err := client.R().Get(url) // make call to api and return resp
	return resp, err
}
