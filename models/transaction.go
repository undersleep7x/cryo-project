package models

import (
	"time"
)

type Transaction interface {
	GetID() string;
	GetSenderType() string
	GetRecipientHash() string
	GetPaymentAddr() string
	GetTxnHash() string
	GetAmount() float64;
	GetCurrency() string;
	GetStatus() string;
	Created() time.Time
	Updated() time.Time
}

//txn struct for both payments out and invoices coming in


type Payout struct{
	Transaction
	SenderHash string `json:"sender_hash,omitempty" gorm:"index"` //hashed sender id only used for payouts
}

type PayoutRequest struct {
	SenderId string `json:"sender_id"`
	RecipientAddr string `json:"recipient_address"`
	Currency string `json:"currency"`
	Amount float64 `json:"amount"`
}
