package repository

import (
	"database/sql"
	"errors"

	"github.com/Dav16Akin/payment-api/internal/models"
)

type TransactionRepository interface {
	Save(transaction *models.Transaction) error
	GetAll() ([]*models.Transaction, error)
	Transfer(transaction *models.Transaction) error
	GetByUser(id string) ([]*models.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (t *transactionRepository) Transfer(transaction *models.Transaction) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var senderBalance float64
	var receiverBalance float64

	err = tx.QueryRow(`SELECT balance FROM wallets WHERE user_id=$1 FOR UPDATE`, transaction.SenderID).Scan(&senderBalance)
	if err != nil {
		return errors.New("sender wallet not found")
	}

	err = tx.QueryRow(
		`SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE`,
		transaction.ReceiverID,
	).Scan(&receiverBalance)
	if err != nil {
		return errors.New("receiver wallet not found")
	}

	if transaction.Amount <= 0 {
		return errors.New("invalid amount")
	}

	if senderBalance < transaction.Amount {
		return errors.New("insufficient funds")
	}

	_, err = tx.Exec(`UPDATE wallets SET balance = balance - $1 WHERE user_id = $2`, transaction.Amount, transaction.SenderID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE wallets SET balance = balance + $1 WHERE user_id = $2`,
		transaction.Amount,
		transaction.ReceiverID,
	)
	if err != nil {
		return err
	}

	transaction.Status = "completed"

	_, err = tx.Exec(
		`INSERT INTO transactions (id, sender_id, receiver_id, amount, status)
		 VALUES ($1, $2, $3, $4, $5)`,
		transaction.ID,
		transaction.SenderID,
		transaction.ReceiverID,
		transaction.Amount,
		transaction.Status,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (t *transactionRepository) Save(transaction *models.Transaction) error {
	query := `INSERT INTO transactions (id , sender_id, receiver_id, amount, status) VALUES ($1,$2,$3,$4,$5)`

	_, err := t.db.Exec(query, transaction.ID, transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.Status)
	if err != nil {
		return err
	}

	return nil
}

func (t *transactionRepository) GetAll() ([]*models.Transaction, error) {

	query := `SELECT id, sender_id, receiver_id, amount, status FROM transactions`

	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		var tx models.Transaction

		err := rows.Scan(
			&tx.ID,
			&tx.SenderID,
			&tx.ReceiverID,
			&tx.Amount,
			&tx.Status,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

func (t *transactionRepository) GetByUser(id string) ([]*models.Transaction, error) {
	query := `SELECT id, sender_id, receiver_id, amount, status FROM transactions WHERE sender_id=$1 OR receiver_id=$1`

	rows , err := t.db.Query(query, id)
	if err != nil {
		return nil , err
	}

	var transactions []*models.Transaction

	for rows.Next() {
		var tx models.Transaction

		err := rows.Scan(
			&tx.ID,
			&tx.SenderID,
			&tx.ReceiverID,
			&tx.Amount,
			&tx.Status,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &tx)
	}

	return transactions, nil
}
