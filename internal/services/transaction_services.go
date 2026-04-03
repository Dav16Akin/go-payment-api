package services

import (
	"errors"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/google/uuid"
)

type TransactionService interface {
	Transfer(transaction *models.Transaction) error
	GetAll() ([]*models.Transaction, error)
}

type transactionService struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
}

func NewTransactionService(walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{walletRepo: walletRepo, transactionRepo: transactionRepo}
}

func (t *transactionService) Transfer(transaction *models.Transaction) error {
	senderWallet, err := t.walletRepo.FindWallet(transaction.SenderID)
	if err != nil {
		return errors.New("sender wallet not found")
	}

	if senderWallet == nil {
		return errors.New("sender wallet not found")
	}

	receiverWallet, err := t.walletRepo.FindWallet(transaction.ReceiverID)
	if err != nil {
		return errors.New("receiver wallet not found")
	}

	if receiverWallet == nil {
		return errors.New("receiver wallet not found")
	}

	if senderWallet.ID == receiverWallet.ID {
		return errors.New("cannot transfer to the same account")
	}

	if transaction.Amount <= 0 {
		return errors.New("invalid amount")
	}

	if senderWallet.Balance < transaction.Amount {
		return errors.New("insufficient funds")
	}

	transaction.ID = uuid.New().String()

	senderWallet.Balance -= transaction.Amount
	receiverWallet.Balance += transaction.Amount

	transaction.Status = "completed"

	if err := t.transactionRepo.Save(transaction); err != nil {

		senderWallet.Balance += transaction.Amount
		receiverWallet.Balance -= transaction.Amount

		return errors.New("failed to save transaction")
	}

	return nil
}

func (t *transactionService) GetAll() ([]*models.Transaction, error) {
	transactions, err := t.transactionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return []*models.Transaction{}, nil
	}

	return transactions, nil
}
