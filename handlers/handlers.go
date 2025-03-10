package handlers

import (
	"net/http"
	"strings"
	"encoding/json"
	// "log"

	"github.com/undersleep7x/cryptowallet-v0.1/services"
)

type PriceFetcher interface {
	FetchCryptoPrice(cryptoList []string, currency string)(any, error)
}

var DefaultPriceFetcher PriceFetcher = defaultPriceFetcher{}

type defaultPriceFetcher struct{}

func (d defaultPriceFetcher) FetchCryptoPrice(cryptoList []string, currency string) (interface{}, error) {
	return services.FetchCryptoPrice(cryptoList, currency)
}

var Ping = func(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"}) //returns 200 OK when hit for healthcheck
}

var FetchPrices = func(w http.ResponseWriter, r *http.Request) {
	cryptos := r.URL.Query().Get("crypto")
	currency := r.URL.Query().Get("currency")

	if cryptos == "" {
		http.Error(w, "Missing 'cryptos' query parameter", http.StatusBadRequest)
		return
	}

	if currency == "" {
		http.Error(w, "Missing 'currency' query parameter", http.StatusBadRequest)
		return
	}
	
	cryptoList := strings.Split(cryptos, ",")
	prices, err := DefaultPriceFetcher.FetchCryptoPrice(cryptoList, currency)
	if err != nil {
		http.Error(w, "Failed to fetch prices", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(prices)
}