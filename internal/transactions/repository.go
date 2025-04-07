package transactions

import (
	// "log"
)

type TxnRepository interface {
	SaveTransaction(txn Transaction) error
	FindTransactionById(txnId string) Invoice
	UpdateTransactionById(txx Transaction) error
}

type txnRepository struct {

}

func NewTxnRepository() TxnRepository {
	return &txnRepository{}
}

func (r *txnRepository) SaveTransaction(txn Transaction) error {
	return nil
}

func (r *txnRepository) FindTransactionById(txnId string) Invoice {
	return Invoice{}
}

func (r *txnRepository) UpdateTransactionById(txx Transaction) error {
	return nil
}
