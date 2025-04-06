package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/undersleep7x/cryo-project/internal/prices"
	"github.com/undersleep7x/cryo-project/internal/transactions"
)

func SetupRoutes(router *gin.Engine, priceHandler *prices.PriceHandler, txnHandler *transactions.TransactionsHandler) {
    router.GET("/", Ping) // ping route
	router.GET("/price", priceHandler.FetchPrices)// route for sourcing pricing data from CoinGecko API
	router.POST("/invoice", txnHandler.CreateInvoice) // create a new transaction (p2p payment, invoice, refund, etc)
	router.POST("/send-payment", txnHandler.SendPayment)
}
