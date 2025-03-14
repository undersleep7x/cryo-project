package models

import (
	"time"
)


type InvoiceRequest struct {
	RecipientId string `json:"recipient_id"`
	Currency string `json:"currency"`
	Amount float64 `json:"amount"`
	ExternalRef *string `json:"external_ref,omitempty"`
	SenderType string `json:"sender_type"`
}

type InvoiceResponse struct {
	TransactionId string `json:"transaction_id"`
	PaymentAddr string `json:"payment_address,omitempty"` //ota for invoice payments out
	Status string `json:"status"` // invoice, pending, confirmed, failed
	ExternalRef *string `json:"external_ref,omitempty"`
}

type Invoice struct {
	ID string `json:"transaction_id" gorm:"primaryKey"`
	SenderType string `json:"sender_type,omitempty"` //user or merchant
	RecipientHash string `json:"recipient_hash,omitempty" gorm:"index"` // hashed recipientid for invoices or walletid for payouts
	PaymentAddr string `json:"payment_address,omitempty"` //ota for invoice payments out
	TxnHash string `json:"tx_hash,omitempty"` // blockchain txn hash for payouts, potentially updated when invoices shift to pending status
	Amount float64 `json:"amount"`
	Currency string `json:"currency" gorm:"index"`
	Status string `json:"status" gorm:"index"` // invoice, pending, confirmed, failed
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	ExternalRef *string `json:"external_ref,omitempty" gorm:"index"` //optional tracking id for merchants external systems
}

// invoice -> transaction implementation func's
func (i Invoice) GetID() string {
	return i.ID
}
func (i Invoice) GetSenderType() string {
	return i.SenderType
}
func (i Invoice) GetRecipientHash() string {
	return i.RecipientHash
}
func (i Invoice) GetPaymentAddr() string {
	return i.PaymentAddr
}
func (i Invoice) GetTxnHash() string {
	return i.TxnHash
}
func (i Invoice) GetAmount() float64 {
	return i.Amount
}
func (i Invoice) GetCurrency() string {
	return i.Currency
}
func (i Invoice) GetStatus() string {
	return i.Status
}
func (i Invoice) Created() time.Time {
	return i.CreatedAt
}
func (i Invoice) Updated() time.Time {
	return i.UpdatedAt
}
func (i Invoice) GetExternalRef() *string {
	if i.ExternalRef != nil{
		return i.ExternalRef
	} else {
		return nil
	}
}