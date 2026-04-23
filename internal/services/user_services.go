package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/Dav16Akin/payment-api/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(user *models.User) (*models.User, error)
	SignIn(req *models.SignInRequest) (string, string, error)

	RefreshToken(token string) (string, string, error)
	Logout(token string) error

	GetUserProfile(userID string) (*models.User, error)
	UpdateUserProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error)
	ChangeUserPassword(userID string, req *models.ChangePasswordRequest) (string, error)
}

type userService struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
	tokenRepo  repository.RefreshTokenRepository
}

func NewUserService(userRepo repository.UserRepository, walletRepo repository.WalletRepository, tokenRepo repository.RefreshTokenRepository) UserService {
	return &userService{userRepo: userRepo, walletRepo: walletRepo, tokenRepo: tokenRepo}
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

	existingUser, err := s.userRepo.FindUserByEmail(user.Email)

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

	if err := s.userRepo.CreateUserWithWallet(userData, wallet); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) SignIn(req *models.SignInRequest) (string, string, error) {
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", "", errors.New("invalid password")
	}

	accessToken, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	refreshToken, err := utils.GenerateRandomToken()
	if err != nil {
		return "", "", err
	}

	hashed := utils.HashToken(refreshToken)

	rt := models.RefreshToken{
		ID:         uuid.New().String(),
		UserID:     user.ID,
		TokenHash:  hashed,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour),
		LastUsedAt: time.Now(),
		Revoked:    false,
	}

	if err := s.tokenRepo.Create(&rt); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) RefreshToken(token string) (string, string, error) {
	hashed := utils.HashToken(token)

	stored, err := s.tokenRepo.FindByTokenHash(hashed)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	if stored.Revoked {
		return "", "", errors.New("token revoked")
	}

	if stored.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("token expired")
	}

	if err := s.tokenRepo.Revoke(stored.ID); err != nil {
		return "", "", err
	}

	newAccess, err := utils.GenerateJWT(stored.UserID)
	if err != nil {
		return "", "", err
	}

	newRefresh, err := utils.GenerateRandomToken()
	if err != nil {
		return "", "", err
	}

	newRT := &models.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    stored.UserID,
		TokenHash: utils.HashToken(newRefresh),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.tokenRepo.Create(newRT); err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

func (s *userService) Logout(token string) error {
	hashed := utils.HashToken(token)

	stored, err := s.tokenRepo.FindByTokenHash(hashed)
	if err != nil {
		return errors.New("invalid token")
	}

	return s.tokenRepo.Revoke(stored.ID)
}

func (s *userService) GetUserProfile(userID string) (*models.User, error) {

	user, err := s.userRepo.FindUserByID(userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUserProfile(userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	updatedUser, err := s.userRepo.UpdateUserProfile(userID, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) ChangeUserPassword(userID string, req *models.ChangePasswordRequest) (string, error) {
	user, err := s.userRepo.FindUserByID(userID)
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

	str, err := s.userRepo.ChangeUserPassword(userID, string(newPasswordHashed))
	if err != nil {
		return "", err
	}

	err = s.tokenRepo.RevokeAllByUserID(userID)
	if err != nil {
		return "", err
	}

	return str, nil
}
