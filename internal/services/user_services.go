package services

import (
	"database/sql"
	"errors"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/Dav16Akin/payment-api/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(user *models.User) (*models.User, error)
	SignIn(req *models.SignInRequest) (*models.User, string, error)
	UpdateProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error)
	ChangePassword(userID string, req *models.ChangePasswordRequest) (string, error)
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
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
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
		return nil, "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

func (s *userService) UpdateProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	updatedUser, err := s.repo.UpdateProfile(userID, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) ChangePassword(userID string, req *models.ChangePasswordRequest) (string, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	if req.NewPassword == req.OldPassword {
		return "", errors.New("old password and new password are same")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return "", errors.New("invalid old password")
	}

	newPasswordHashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return s.repo.ChangePassword(userID, string(newPasswordHashed))
}
