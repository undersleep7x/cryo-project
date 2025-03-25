package services

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/undersleep7x/cryo-project/internal/models"
	"github.com/undersleep7x/cryo-project/internal/repository"
)

// interface for transaction service
type TransactionService interface {
	CreateInvoice(models.InvoiceRequest) (*models.InvoiceResponse, error)
	SendPayment(models.PaymentRequest) (*models.PaymentResponse, error)
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
	resp.ExternalRef = inv.GetExternalRef()
	resp.TransactionId = inv.GetID()
	resp.Status = inv.GetStatus()
	return &resp, nil

	//TODO other todos to be mindful of
	// client side encryption for sensitive invoice data (invoice id, recipient id, amount, currency, payment address, sender type, external ref)
	// expiration after certain time (up to 30 days??, cron job with db implementation)
	// refund implementation
	// prevent invoice duplication (check for unique invoice hashes and reject dupe amount, recipient, currency, and metadata)
	// api security and rate limiting 

}

func (s *transactionsServiceImpl) SendPayment(r models.PaymentRequest) (*models.PaymentResponse, error) {
	senderRef := r.SenderId + "hash"
	recipRef := r.PaymentAddr + "hash"
	response := models.PaymentResponse{}

	if r.InvoiceId == "" { // flow for a direct payment

		pay := models.Payment {
			ID: "txn_" + uuid.NewString(),
			SenderType: r.SenderType,
			RecipientRef: recipRef,
			SenderRef: senderRef,
			PaymentAddr: r.PaymentAddr,
			TxnRef: "txnrefhash",
			Amount: r.Amount,
			Currency: r.Currency,
			Status: "Pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repository.SaveTransaction(pay)
		if err != nil {
			log.Printf("Error saving new invoice to database: %v", err)
			return nil, err
		}

		response.PaymentAddr = pay.PaymentAddr
		response.Status = pay.Status
		response.TransactionId = pay.ID

	} else { // flow for invoice payment
		inv := repository.FindTransactionById(r.InvoiceId)
		// stubbed info
		externalRef := "anexternalref"
		inv.ExternalRef = &externalRef
		inv.ID = "txn_" + uuid.NewString()
		inv.Status = "Invoice"
		inv.ID = "1234555x05"
		

		inv.SetTxnHash("txnHashFromBlockChain") //txnhash should be present even if txn is still "otw"
		inv.SetStatus("Pending") //txn is on the way, will next be confirmed or failed
		inv.SetUpdate(time.Now())

		err := repository.UpdateTransactionById(inv)
		if err != nil {
			log.Printf("Error saving new invoice to database: %v", err)
			return nil, err
		}

		response.ExternalRef = inv.ExternalRef
		response.PaymentAddr = r.PaymentAddr
		response.Status = inv.Status
		response.TransactionId = inv.ID
	}

	return &response, nil
}

var GenerateOneTimeAddress = func(currency string) string {
	var genOta = "STUBOTA12345678"
	log.Printf("One time address successfully generated")
	return genOta
}
