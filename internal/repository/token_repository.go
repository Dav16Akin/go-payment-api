package repository

import (
	"database/sql"

	"github.com/Dav16Akin/payment-api/internal/models"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	FindByTokenHash(hash string) (*models.RefreshToken, error)
	Revoke(id string) error
	RevokeAllByUserID(id string) error
}

type refreshTokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at,last_used_at, revoked) VALUES($1,$2,$3,$4,$5,$6)`

	_, err := r.db.Exec(query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt,token.LastUsedAt, token.Revoked)
	if err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepository) FindByTokenHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	query := `SELECT id , user_id, token_hash, expires_at, revoked FROM refresh_tokens WHERE token_hash=$1 and revoked=false`
	err := r.db.QueryRow(query, hash).Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.Revoked)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &token, nil
}


func (r *refreshTokenRepository) Revoke(id string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE id=$1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *refreshTokenRepository) RevokeAllByUserID(id string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
