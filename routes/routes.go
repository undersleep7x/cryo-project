package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/undersleep7x/cryo-project/controllers"
	"github.com/undersleep7x/cryo-project/services"
)

func SetupRoutes(router *gin.Engine) {
	priceService := services.NewFetchCryptoPriceService()
	priceFetcher := controllers.NewPriceFetcher(priceService)

    router.GET("/", controllers.Ping) // ping route
	router.GET("/price", priceFetcher.FetchPrices)// route for sourcing pricing data from CoinGecko API
	router.POST("/transaction", controllers.CreateInvoice) // create a new transaction (p2p payment, invoice, refund, etc)
}
