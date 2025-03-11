package handlers

import (
	"net/http"
	"strings"
	"encoding/json"
	"log"

	"github.com/undersleep7x/cryptowallet-v0.1/services"
)

//setup interface for price fetching, setting default struct for method, and implementing the method
//improves mocking/testing of handlers
type PriceFetcher interface {
	FetchCryptoPrice(cryptoList []string, currency string)(any, error)
}
var DefaultPriceFetcher PriceFetcher = defaultPriceFetcher{}
type defaultPriceFetcher struct{}
func (d defaultPriceFetcher) FetchCryptoPrice(cryptoList []string, currency string) (interface{}, error) {
	return services.FetchCryptoPrice(cryptoList, currency)
}

//handle /ping route call and return ok to confirm healthy service
var Ping = func(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}) //returns 200 OK when hit for healthcheck
	if err != nil{
		http.Error(w, "Please try again later", http.StatusInternalServerError)
		log.Printf("JSON encoding error: %v", err)
		return
	}
}

// handle /price route call and return latest prices from coingecko 
var FetchPrices = func(w http.ResponseWriter, r *http.Request) {
	//store query params
	cryptos := r.URL.Query().Get("crypto")
	currency := r.URL.Query().Get("currency")

	// param validation before logic, return error if missing
	if cryptos == "" {
		http.Error(w, "Missing 'cryptos' query parameter", http.StatusBadRequest)
		return
	}
	if currency == "" {
		http.Error(w, "Missing 'currency' query parameter", http.StatusBadRequest)
		return
	}
	
	cryptoList := strings.Split(cryptos, ",") // csv -> array of cryptos
	prices, err := DefaultPriceFetcher.FetchCryptoPrice(cryptoList, currency) // call service to fetch pricing
	if err != nil { //return error if service error is thrown
		http.Error(w, "Failed to fetch prices", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(prices); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Json encoding error: %v", err)
		return
	} // return prices json
	
}