package models

import (
	"time"
)

//txn struct for both payments out and invoices coming in
type Transaction struct {
	ID string `json:"transaction_id" gorm:"primaryKey"`
	SenderHash string `json:"sender_hash,omitempty" gorm:"index"` //hashed sender id only used for payouts
	SenderType string `json:"sender_type,omitempty"` //user or merchant
	RecipientHash string `json:"recipient_hash,omitempty" gorm:"index"` // hashed recipientid for invoices or walletid for payouts
	PaymentAddr string `json:"payment_address,omitempty"` //ota for invoice payments out
	TxHash string `json:"tx_hash",omitempty` // blockchain txn hash for payouts
	Amount float64 `json:"amount"`
	Currency string `json:"currency" gorm:"index"`
	Status string `json:"status" gorm:"index"` // invoice, pending, confirmed, failed
	ExternalRef string `json:"external_ref,omitempty" gorm:"index"` //optional tracking id for merchants external systems
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type InvoiceRequest struct {
	RecipientId string `json:"recipient_id"`
	Currency string `json:"currency"`
	Amount float64 `json:"amount"`
}

type PayoutRequest struct {
	SenderId string `json:"sender_id"`
	RecipientAddr string `json:"recipient_address"`
	Currency string `json:"currency"`
	Amount float64 `json:"amount"`
}
