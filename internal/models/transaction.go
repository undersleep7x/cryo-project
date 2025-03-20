package models

import (
	"time"
)

// transaction interface for invoice and payment structsx
type Transaction interface {
	GetID() string;
	GetSenderType() string
	GetRecipientRef() string
	GetTxnHash() string
	GetAmount() float64;
	GetCurrency() string;
	GetStatus() string;
	Created() time.Time
	Updated() time.Time
}


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
