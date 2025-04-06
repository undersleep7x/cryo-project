package prices

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// fetch prices from coingecko api
var FetchPrices = func(cryptos []string, currency string, baseURL string, timeoutVal int) (*resty.Response, error) {
	// setup resty client for api call
	client := resty.New()
	timeout := time.Duration(timeoutVal) * time.Second
	client.SetTimeout(timeout)

	// turn passed in array into csv string
	apiQuery := strings.Join(cryptos, ",")
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", baseURL, apiQuery, currency)
	log.Printf("Making API call to %s", url)
	resp, err := client.R().Get(url) // make call to api and return resp
	return resp, err
}
