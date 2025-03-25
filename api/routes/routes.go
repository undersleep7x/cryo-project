package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/undersleep7x/cryo-project/internal/prices"
	"github.com/undersleep7x/cryo-project/internal/transactions"
	"github.com/undersleep7x/cryo-project/internal/app"
)

func SetupRoutes(router *gin.Engine) {
	//price interface interactions
	priceService := prices.NewFetchCryptoPriceService()
	priceHandler := prices.NewPriceHandler(priceService)
	//transactions interface interactions
	transactionService := transactions.NewTransactionsService()
	transactionsHandler := transactions.NewTransactionsHandler(transactionService)

    router.GET("/", app.Ping) // ping route
	router.GET("/price", priceHandler.FetchPrices)// route for sourcing pricing data from CoinGecko API
	router.POST("/invoice", transactionsHandler.CreateInvoice) // create a new transaction (p2p payment, invoice, refund, etc)
	router.POST("/send-payment", transactionsHandler.SendPayment)
}
