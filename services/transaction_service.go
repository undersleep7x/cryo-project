package services

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/undersleep7x/cryo-project/models"
	"github.com/undersleep7x/cryo-project/repository"
)

// interface for transaction service
type TransactionService interface {
	CreateInvoice(models.InvoiceRequest) (*models.InvoiceResponse, error)
	SendPayment(models.PayoutRequest)
}
type transactionsServiceImpl struct{}
func NewTransactionsService() TransactionService {
	return &transactionsServiceImpl{}
}

// service function for creating new invoice and saving to db
func (s *transactionsServiceImpl) CreateInvoice(r models.InvoiceRequest) (*models.InvoiceResponse, error) {
	recipientHash := r.RecipientId + "RECIPHASH" //TODO implement hashing functionality here
	resp := models.InvoiceResponse{}

	inv := models.Invoice {
		ID: "txn_" + uuid.NewString(),
		SenderType: r.SenderType,
		RecipientRef: recipientHash,
		WalletRef: GenerateOneTimeAddress(r.Currency),
		RefundRef: r.ExternalRef,
		Amount: r.Amount,
		Currency: r.Currency,
		Status: "invoice",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExternalRef: r.ExternalRef,
	}

	err := repository.SaveTransaction(inv)
	if err != nil {
		log.Printf("Error saving new invoice to database: %v", err)
		return nil, err
	}

	//TODO after creation and save to db, there must be logic that allows for tracking of the invoice
	//such as identifying when payment has been made, following blockchain for confirmation, etc
	//just return invoice for now
	resp.ExternalRef = inv.ExternalRef
	resp.TransactionId = inv.ID
	resp.Status = inv.Status
	return &resp, nil

	//TODO other todos to be mindful of
	// client side encryption for sensitive invoice data (invoice id, recipient id, amount, currency, payment address, sender type, external ref)
	// expiration after certain time (up to 30 days??, cron job with db implementation)
	// refund implementation
	// prevent invoice duplication (check for unique invoice hashes and reject dupe amount, recipient, currency, and metadata)
	// api security and rate limiting 

}

func (s *transactionsServiceImpl) SendPayment(r models.PayoutRequest) {

}

var GenerateOneTimeAddress = func(currency string) string {
	var genOta = "STUBOTA12345678"
	log.Printf("One time address successfully generated")
	return genOta
}
