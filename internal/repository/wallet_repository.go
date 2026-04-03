package repository

import (
	"github.com/Dav16Akin/payment-api/internal/models"
)

type WalletRepository interface {
	CreateWallet(w *models.Wallet) error
	FindWallet(id string) (*models.Wallet, error)
	FindWalletByUserId(id string) (*models.Wallet, error)
}

type walletRepository struct {
	wallets []*models.Wallet
}

func NewWalletRepository() WalletRepository {
	return &walletRepository{wallets: []*models.Wallet{}}
}

func (r *walletRepository) CreateWallet(w *models.Wallet) error {
	r.wallets = append(r.wallets, w)
	return nil
}

func (r *walletRepository) FindWallet(id string) (*models.Wallet, error) {
	for _, w := range r.wallets {
		if w.ID == id {
			return w, nil
		}
	}

	return nil, nil
}

func (r *walletRepository) FindWalletByUserId(id string) (*models.Wallet, error) {
	for _, w := range r.wallets {
		if w.UserID == id {
			return w, nil
		}
	}

	return nil, nil
}
