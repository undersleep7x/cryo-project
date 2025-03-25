package prices

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// setup interface for price fetching
type PriceHandler struct {
	service FetchCryptoPriceService
}
func NewPriceHandler (service FetchCryptoPriceService) *PriceHandler {
	return &PriceHandler{service: service}
}


// handle /price route call and return latest prices from coingecko
func (f *PriceHandler) FetchPrices (c *gin.Context) {
	//store query params
	cryptos := c.Query("crypto")
	currency := c.Query("currency")

	// param validation before logic, return error if missing
	if cryptos == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'crypto' query parameter"})
		return
	}
	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'currency' query parameter"})
		return
	}

	cryptoList := strings.Split(cryptos, ",") // csv -> array of cryptos
	prices, err := f.service.FetchCryptoPrice(cryptoList, currency) // call service to fetch pricing
	if err != nil {   //return error if service error is thrown
		log.Printf("Internal Server Error when calling FetchCryptoPrice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"prices": prices}) // return prices json

}