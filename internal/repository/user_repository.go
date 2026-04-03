package repository

import (
	"github.com/Dav16Akin/payment-api/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db []*models.User
}

func NewUserRepository() UserRepository {
	return &userRepository{db: []*models.User{}}
}

func (r *userRepository) CreateUser(user *models.User) error {
	r.db = append(r.db, user)
	return nil
}

func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	for _, v := range r.db {
		if v.Email == email {
			return v, nil
		}
	}

	return nil, nil
}
