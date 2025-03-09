package handlers

import (
	"net/http"
	"strings"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/undersleep7x/cryptowallet-v0.1/services"
)

func GetCryptoPrice(c *gin.Context) {
	// parse query params
	cryptos := c.Query("crypto")
	currency := c.Query("currency")

	if cryptos == "" || currency == "" {
		//missing params error handling
		log.Println("Query params cryptos and/or currency are missing from request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	cryptoList := strings.Split(cryptos, ",") // split query param csv into array for flexible processing

	priceData, err := services.FetchCryptoPrice(cryptoList, currency) //process api call w/ error handling
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, priceData)
}