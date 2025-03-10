package handlers

import (
	"net/http"
	"strings"
	"encoding/json"
	// "log"

	"github.com/undersleep7x/cryptowallet-v0.1/services"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"}) //returns 200 OK when hit for healthcheck
}

func FetchPrices(w http.ResponseWriter, r *http.Request) {
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
	prices, err := services.FetchCryptoPrice(cryptoList, currency)
	if err != nil {
		http.Error(w, "Failed to fetch prices", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(prices)
}