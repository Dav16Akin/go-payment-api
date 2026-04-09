package database

import "database/sql"

const createUsersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY ,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

const createWalletsTable = `
	CREATE TABLE IF NOT EXISTS wallets (
		id UUID PRIMARY KEY,
		user_id TEXT NOT NULL,
		balance NUMERIC NOT NULL DEFAULT 0
	);`

const createTransactionsTable = `
	CREATE TABLE IF NOT EXISTS transactions (
		id UUID PRIMARY KEY,
		sender_id TEXT NOT NULL,
		receiver_id TEXT NOT NULL,
		amount NUMERIC NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

func InitializeDB(db *sql.DB) error {
	queries := []string{
		createUsersTable,
		createWalletsTable,
		createTransactionsTable,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
