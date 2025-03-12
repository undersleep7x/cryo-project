package routes

import (
	"github.com/undersleep7x/cryptowallet-v0.1/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
    router.GET("/", controllers.Ping) // ping route
	router.GET("/price", controllers.FetchPrices)// route for sourcing pricing data from CoinGecko API
	router.POST("/transaction", controllers.CreateInvoice) // create a new transaction (p2p payment, invoice, refund, etc)
}