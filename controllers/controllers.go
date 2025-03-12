package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/undersleep7x/cryptowallet-v0.1/models"
	"github.com/undersleep7x/cryptowallet-v0.1/services"
	"github.com/gin-gonic/gin"
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
var Ping = func(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PONG"})
}

// handle /price route call and return latest prices from coingecko 
var FetchPrices = func(c *gin.Context) {
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
	prices, err := DefaultPriceFetcher.FetchCryptoPrice(cryptoList, currency) // call service to fetch pricing
	if err != nil { //return error if service error is thrown
		log.Printf("Internal Server Error when calling FetchCryptoPrice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"prices": prices}) // return prices json
	
}

func CreateInvoice(c *gin.Context) {
	var request models.InvoiceRequest

	if err := c.ShouldBindJSON(&request); err != nil { // validation check for json request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	txn, err := services.CreateInvoice(request)
	if err != nil { //catch for service failure
		log.Printf("Internal Server Error when calling CreateInvoice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_address": txn.PaymentAddr,
		"transaction_id": txn.ID,
		"status": txn.Status,
	})
}

func SendPayment(c *gin.Context) {
	var request models.PayoutRequest

	if err := c.ShouldBindJSON(&request); err != nil { // validation check for json request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	txn, err := services.SendPayment(request)
	if err != nil { //catch for service failure
		log.Printf("Internal Server Error when calling sendPayment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": txn.ID,
		"status": txn.Status,
		"tx_hash": txn.TxHash,
	})
}

