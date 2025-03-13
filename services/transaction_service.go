package services

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/undersleep7x/cryo-project/models"
	"github.com/undersleep7x/cryo-project/repository"
)

var GenerateOneTimeAddress = func(currency string) string {
	var genOta = "STUBOTA12345678"
	log.Printf("One time address successfully generated")
	return genOta
}

var CreateInvoice = func(r models.InvoiceRequest) (models.Transaction, error) {

	recipientHash := "RECIPHASH"

	// generates invoice for new merchant payment and saves to db as new transaction
	txn := models.Transaction{
		ID:            "txn_" + uuid.New().String(),
		RecipientHash: recipientHash,
		Amount:        r.Amount,
		Currency:      r.Currency,
		Status:        "invoice",
		PaymentAddr:   GenerateOneTimeAddress(r.Currency),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := repository.SaveTransaction(txn)
	return txn, err
}

var SendPayment = func(r models.PayoutRequest) (models.Transaction, error) {
	txn := models.Transaction{}
	return txn, nil
}
