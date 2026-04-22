package repository

import (
	"database/sql"
	"net/mail"

	"github.com/Dav16Akin/payment-api/internal/models"
)

type UserRepository interface {
	CreateUserWithWallet(user *models.User, wallet *models.Wallet) error

	FindUserByID(id string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)

	UpdateProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error)
	ChangePassword(userID string, newPassword string) (string, error)
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

	userQuery := `INSERT INTO users (id, name , email, password) VALUES ($1,$2,$3,$4)`
	_, err = tx.Exec(userQuery, user.ID, user.Name, user.Email, user.Password)
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

func (r *userRepository) FindUserByID(id string) (*models.User, error) {
	var user models.User

	query := `SELECT id, name, email, password, avatar_url FROM users WHERE id=$1`

	var avatar sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&avatar,
	)

	if avatar.Valid {
		user.AvatarURL = &avatar.String
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User

	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, email, password, avatar_url FROM users WHERE email=$1`

	var avatar sql.NullString

	err = r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&avatar,
	)

	if avatar.Valid {
		user.AvatarURL = &avatar.String
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	existingUser, err := r.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, sql.ErrNoRows
	}

	if req.Name != nil {
		existingUser.Name = *req.Name
	}

	if req.AvatarURL != nil {
		existingUser.AvatarURL = req.AvatarURL
	}

	query := `UPDATE users SET name=$1, avatar_url=$2 WHERE id=$3 RETURNING id, name, email, avatar_url`
	var updatedUser models.User
	err = r.db.QueryRow(
		query,
		existingUser.Name,
		existingUser.AvatarURL,
		existingUser.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.AvatarURL,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &updatedUser, nil
}

func (r *userRepository) ChangePassword(userID string, newPassword string) (string, error) {
	query := `UPDATE users SET password=$1 WHERE id=$2`

	result, err := r.db.Exec(query, newPassword, userID)

	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}

	return "password changed succesfully", nil
}
