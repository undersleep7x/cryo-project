package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/undersleep7x/cryo-project/internal/controllers"
	"github.com/undersleep7x/cryo-project/internal/services"
)

func SetupRoutes(router *gin.Engine) {
	//price interface interactions
	priceService := services.NewFetchCryptoPriceService()
	priceFetcher := controllers.NewPriceFetcher(priceService)
	//transactions interface interactions
	transactionService := services.NewTransactionsService()
	transactions := controllers.NewTransactionsService(transactionService)

    router.GET("/", controllers.Ping) // ping route
	router.GET("/price", priceFetcher.FetchPrices)// route for sourcing pricing data from CoinGecko API
	router.POST("/invoice", transactions.CreateInvoice) // create a new transaction (p2p payment, invoice, refund, etc)
	router.POST("/send-payment", transactions.SendPayment)
}
