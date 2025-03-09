package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/undersleep7x/cryptowallet-v0.1/handlers"
)

func SetupRoutes(router *gin.Engine) {
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Welcome to CryptoWallet!"}) // default route for home page
    })
	router.GET("/price", handlers.GetCryptoPrice) // route for sourcing pricing data from CoinGecko API
}