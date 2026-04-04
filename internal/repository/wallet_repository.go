package repository

import (
	"database/sql"
	"errors"

	"github.com/Dav16Akin/payment-api/internal/models"
)

type WalletRepository interface {
	FindWallet(id string) (*models.Wallet, error)
	FindWalletByUserId(id string) (*models.Wallet, error)
}

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepository{db: db}
}

var ErrWalletNotFound = errors.New("wallet not found")


func (r *walletRepository) FindWallet(id string) (*models.Wallet, error) {
	var wallet models.Wallet

	query := `SELECT id, user_id, balance FROM wallets WHERE id=$1`

	err := r.db.QueryRow(query, id).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil , ErrWalletNotFound
		}

		return nil , err
	}

	return &wallet, nil
}

func (r *walletRepository) FindWalletByUserId(id string) (*models.Wallet, error) {
	var wallet models.Wallet
	
	query := `SELECT id, user_id, balance FROM wallets WHERE user_id=$1`

	err := r.db.QueryRow(query, id).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil , ErrWalletNotFound
		}

		return nil , err
	}

	return &wallet, nil
}
