package transactions

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionsHandler struct {
	service TransactionService
}
func NewTransactionsHandler(service TransactionService) *TransactionsHandler {
	return &TransactionsHandler{service: service}
}

//handle /invoice call
func (f *TransactionsHandler) CreateInvoice(c *gin.Context) {
	var request InvoiceRequest // create request object for json

	if err := c.ShouldBindJSON(&request); err != nil { // validation check for json request after parsing to request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	inv, err := f.service.CreateInvoice(request) // call service for invoices
	if err != nil { //catch for service failure
		log.Printf("Internal Server Error when calling CreateInvoice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}

	c.JSON(http.StatusOK, inv)
}

func (f *TransactionsHandler) SendPayment(c *gin.Context) {
	var request PaymentRequest // create request object for json

	if err := c.ShouldBindJSON(&request); err != nil { // validation check for json request after parsing to request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	txn, err := f.service.SendPayment(request) // call service for invoices
	if err != nil { //catch for service failure
		log.Printf("Internal Server Error when calling SendPayment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send payment"})
		return
	}

	c.JSON(http.StatusOK, txn)
}