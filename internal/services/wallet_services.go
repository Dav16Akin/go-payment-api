package services

import (
	"errors"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/repository"
)

type WalletService interface {
	GetWallet(id string) (*models.Wallet, error)
}

type walletService struct {
	walletRepo repository.WalletRepository
}

func NewWalletService(walletRepo repository.WalletRepository) WalletService {
	return &walletService{walletRepo: walletRepo}
}

func (s *walletService) GetWallet(id string) (*models.Wallet, error) {
	wallet , err := s.walletRepo.FindWalletByUserId(id);
	if err != nil {
		return nil , errors.New("wallet not found")
	}

	if wallet == nil {
		return nil , errors.New("wallet not found")
	}

	return wallet , nil
}