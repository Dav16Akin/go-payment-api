package services

import (
	"errors"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
	walletRepo repository.WalletRepository
}

func NewUserService(repo repository.UserRepository, walletRepo repository.WalletRepository) UserService {
	return &userService{repo: repo, walletRepo: walletRepo}
}

func (s *userService) CreateUser(user *models.User) (*models.User, error) {
	if user.Name == "" {
		return nil, errors.New("Name is required")
	}

	if user.Email == "" {
		return nil , errors.New("Email is required")
	}

	existingUser, err := s.repo.FindUserByEmail(user.Email)
	if err != nil {
		return nil , err
	}

	if existingUser != nil {
		return nil , errors.New("email already exists")
	}

	user.ID = uuid.New().String()

	if err := s.repo.CreateUser(user); err != nil{
		return nil , err
	}

	wallet := &models.Wallet{
		ID:      user.ID,
        UserID:  user.ID,
        Balance: 0.0,
	}

	if err := s.walletRepo.CreateWallet(wallet); err != nil {
        return nil, errors.New("failed to create wallet")
    }

	return user, nil
}



