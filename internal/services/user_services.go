package services

import (
	"database/sql"
	"errors"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(user *models.User) (*models.User, error)
	SignIn(req *models.SignInRequest) (*models.User, string, error)
}

type userService struct {
	repo       repository.UserRepository
	walletRepo repository.WalletRepository
}

func NewUserService(repo repository.UserRepository, walletRepo repository.WalletRepository) UserService {
	return &userService{repo: repo, walletRepo: walletRepo}
}

func (s *userService) SignUp(user *models.User) (*models.User, error) {
	if user.Name == "" {
		return nil, errors.New("Name is required")
	}

	if user.Password == "" {
		return nil, errors.New("Password is required")
	}

	if user.Email == "" {
		return nil, errors.New("Email is required")
	}

	existingUser, err := s.repo.FindUserByEmail(user.Email)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.ID = uuid.New().String()

	userData := &models.User{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Password: string(hashedPassword),
	}

	wallet := &models.Wallet{
		ID:      user.ID,
		UserID:  user.ID,
		Balance: 500.00,
	}

	if err := s.repo.CreateUserWithWallet(userData, wallet); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) SignIn(req *models.SignInRequest) (*models.User, string, error) {
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, "" , errors.New("cannot find user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, "", errors.New("invalid credentials pass")
	}

	token := "mock-token"
	
	return user, token, nil
}
