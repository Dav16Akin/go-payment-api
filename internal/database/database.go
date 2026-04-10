package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func ConnectToDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_PUBLIC_URL")

	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_PUBLIC_URL is not set")
	}

	if strings.Contains(dsn, "sslmode=require") {
		dsn = strings.Replace(dsn, "sslmode=require", "sslmode=disable", 1)
	} else if !strings.Contains(dsn, "sslmode") {
		dsn += "?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	fmt.Println("Connected to PostgreSQL successfully")

	return db, nil
}
