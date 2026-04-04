package repository

import (
	"database/sql"

	"github.com/Dav16Akin/payment-api/internal/models"
)

type UserRepository interface {
	CreateUserWithWallet(user *models.User, wallet *models.Wallet) error
	FindUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUserWithWallet(user *models.User, wallet *models.Wallet) error {
	tx, err := r.db.Begin() // begin transaction
	if err != nil {
		return err
	}
	defer tx.Rollback() // rollback if anything fails

	userQuery := `INSERT INTO users (id, name , email) VALUES ($1,$2,$3)`
	_, err = tx.Exec(userQuery, user.ID, user.Name, user.Email)
	if err != nil {
		return err
	}

	walletQuery := `INSERT INTO wallets (id, user_id, balance) VALUES ($1,$2,$3)`
	_, err = tx.Exec(walletQuery, wallet.ID, wallet.UserID, wallet.Balance)
	if err != nil {
		return err
	}

	return tx.Commit() // commit transaction if both succeed
}

func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User

	query := `SELECT id , name, email FROM users WHERE email=$1`

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email)

	if err == sql.ErrNoRows {
		return nil, nil 
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
