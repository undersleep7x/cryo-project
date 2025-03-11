package api

import (
	"log"
	"fmt"
	"time"
	"strings"

	"github.com/undersleep7x/cryptowallet-v0.1/app"
	"github.com/go-resty/resty/v2"
)

var FetchPrices = func(cryptos []string, currency string) (*resty.Response, error){
	client := resty.New()
	timeout := time.Duration(app.Config.CoinGecko.Timeout)*time.Second
	client.SetTimeout(timeout)

	apiQuery := strings.Join(cryptos, ",")
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", app.Config.CoinGecko.BaseURL, apiQuery, currency)
	log.Printf("Making API call to %s", url)
	resp, err := client.R().Get(url)
	return resp, err
}