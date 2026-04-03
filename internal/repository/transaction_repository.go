package repository

import "github.com/Dav16Akin/payment-api/internal/models"

type TransactionRepository interface {
	Save(transaction *models.Transaction) error 
	GetAll() ([]*models.Transaction, error)
}

type transactionRepository struct {
	transactions []*models.Transaction
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{transactions: []*models.Transaction{}}
}

func (t *transactionRepository) Save(transaction *models.Transaction) error {
	t.transactions = append(t.transactions, transaction)
	return nil
}

func (t *transactionRepository) GetAll() ([]*models.Transaction, error) {
	return t.transactions, nil
}