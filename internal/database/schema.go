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

const createRefreshTokenTable = `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token TEXT NOT NULL UNIQUE,
		expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		last_used_at TIMESTAMP WITH TIME ZONE,
		revoked BOOLEAN DEFAULT FALSE
	);
`
const createMigrationsTable = `
CREATE TABLE IF NOT EXISTS migrations (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

type Migration struct {
	Name  string
	Query string
}

var migrations = []Migration{
	{
		Name:  "add_avatar_url_to_users",
		Query: `ALTER TABLE users ADD COLUMN avatar_url TEXT;`,
	},
	{
		Name:  "rename_token_to_token_hash_at_refresh_tokens",
		Query: `ALTER TABLE refresh_tokens RENAME COLUMN token TO token_hash;`,
	},
}

func RunMigrations(db *sql.DB) error {
	for _, m := range migrations {
		var exists bool

		err := db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM migrations WHERE name=$1)`,
			m.Name,
		).Scan(&exists)

		if err != nil {
			return err
		}

		if exists {
			continue
		}

		// Run migration
		_, err = db.Exec(m.Query)
		if err != nil {
			return err
		}

		// Record migration
		_, err = db.Exec(
			`INSERT INTO migrations (name) VALUES ($1)`,
			m.Name,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitializeDB(db *sql.DB) error {
	queries := []string{
		createUsersTable,
		createWalletsTable,
		createTransactionsTable,
		createMigrationsTable,
		createRefreshTokenTable,
	}

	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
	`)

	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
	`)

	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_revoked ON refresh_tokens(token, revoked);
	`)

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
