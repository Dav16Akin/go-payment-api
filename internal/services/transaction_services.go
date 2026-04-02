package services

import (
	"errors"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/google/uuid"
)

type TransactionService interface {
	Transfer(transaction *models.Transaction) error
}

type transactionService struct {
	walletRepo repository.WalletRepository
}

func NewTransactionService(walletRepo repository.WalletRepository) TransactionService {
	return &transactionService{walletRepo: walletRepo}
}

func (t *transactionService) Transfer(transaction *models.Transaction) error {
	senderWallet, err := t.walletRepo.FindWallet(transaction.SenderID)
	if err != nil {
		return errors.New("sender wallet not found")
	}

	receiverWallet, err := t.walletRepo.FindWallet(transaction.RecieverID)
	if err != nil {
		return errors.New("receiver wallet not found")
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

	transaction.Status = "Completed"

	t.walletRepo.ListAllWallets()

	return nil
}
