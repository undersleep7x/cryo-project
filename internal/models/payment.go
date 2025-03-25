package models

import (
	"time"
)

type PaymentRequest struct {
	SenderId string `json:"sender_id"` //invoices can only be sent by merchants, this would essentially be a ref to merchant id
	Currency string `json:"currency"`
	Amount float64 `json:"amount"`
	PaymentAddr string `json:"payment_address"`
	SenderType string `json:"sender_type"`
	WalletRef string `json:"wallet_ref"`
	InvoiceId string `json:"invoice_id"`
	// RefundRef *string `json:"refund_ref,omitempty"` // may add refund functionality later but for now not needed
}

type PaymentResponse struct {
	TransactionId string `json:"transaction_id"`
	PaymentAddr string `json:"payment_address,omitempty"` //ota for invoice payments out
	Status string `json:"status"` // invoice, pending, confirmed, failed
	ExternalRef *string `json:"external_ref,omitempty"`
}

type Payment struct {
	ID string `json:"transaction_id" gorm:"primaryKey"`
	SenderType string `json:"sender_type,omitempty"` //user or merchant
	RecipientRef string `json:"recipient_ref,omitempty" gorm:"index"` // hashed recipientid for invoices or walletid for payouts
	SenderRef string `json:"sender_ref" gorm:"index"`
	PaymentAddr string `json:"payment_address,omitempty"` //ota for invoice payments out
	TxnRef string `json:"tx_ref,omitempty"` // blockchain txn hash for payouts, potentially updated when invoices shift to pending status
	Amount float64 `json:"amount"`
	// RefundRef *string `json:"refund_ref,omitempty"` // ref to refund table
	Currency string `json:"currency" gorm:"index"`
	Status string `json:"status" gorm:"index"` // invoice, pending, confirmed, failed
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// invoice -> transaction implementation func's
func (p Payment) GetID() string {
	return p.ID
}
func (p Payment) GetSenderType() string {
	return p.SenderType
}
func (p Payment) GetSenderRef() string {
	return p.SenderRef
}
func (p Payment) GetRecipientRef() string {
	return p.RecipientRef
}
func (p Payment) GetPaymentAddr() string {
	return p.PaymentAddr
}
func (p Payment) GetTxnHash() string {
	return p.TxnRef
}
func (p Payment) GetAmount() float64 {
	return p.Amount
}
func (p Payment) GetCurrency() string {
	return p.Currency
}
func (p Payment) GetStatus() string {
	return p.Status
}
func (p Payment) Created() time.Time {
	return p.CreatedAt
}
func (p Payment) Updated() time.Time {
	return p.UpdatedAt
}
// func (p Payment) GetExternalRef() *string {
// 	if p.ExternalRef != nil{
// 		return p.ExternalRef
// 	} else {
// 		return nil
// 	}
// }
// func (p Payment) GetRefundRef() *string {
// 	if p.RefundRef != nil {
// 		return p.ExternalRef
// 	} else {
// 		return nil
// 	}
// }