package repository

import (
	"fmt"

	"github.com/Dav16Akin/payment-api/internal/models"
)

type WalletRepository interface {
	CreateWallet(w *models.Wallet) error
    FindWallet(id string) (*models.Wallet, error)
    ListAllWallets()
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


func (r *walletRepository) ListAllWallets() {
	fmt.Println("\n--- Wallets---")
	for _, u := range r.wallets {
		fmt.Printf("ID: %s | Balance: %g\n", u.ID, u.Balance)
	}
	fmt.Println("--------------------------")
}