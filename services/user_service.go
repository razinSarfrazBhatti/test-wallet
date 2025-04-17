package services

import (
	"errors"
	"test-wallet/middleware"
	"test-wallet/models"
	"test-wallet/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserService) RegisterUser(req *models.RegisterUserRequest) (*models.User, error) {
	// Check if phone number already exists
	exists, err := s.userRepo.PhoneNumberExists(req.PhoneNumber)
	if err != nil {
		return nil, errors.New("failed to check phone number availability")
	}
	if exists {
		return nil, errors.New("phone number is already registered")
	}

	// Generate salt
	salt, err := GenerateSalt(16)
	if err != nil {
		return nil, errors.New("failed to generate salt")
	}

	// Hash the PIN with the salt
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.Pin+salt), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash PIN")
	}

	// Create a new Ethereum wallet
	walletResponse, err := CreateWallet()
	if err != nil {
		return nil, errors.New("failed to create wallet")
	}

	// Encrypt the mnemonic with the user's PIN
	encryptedMnemonic, err := Encrypt(req.Pin, walletResponse.Mnemonic)
	if err != nil {
		return nil, errors.New("failed to encrypt wallet mnemonic")
	}

	// Create a new user record
	newUser := &models.User{
		Id:          uuid.New().String(),
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Pin:         string(hashedPin),
		Salt:        salt,
		Wallet: models.Wallet{
			Id:       uuid.New().String(),
			Address:  walletResponse.Address,
			Mnemonic: encryptedMnemonic,
		},
	}

	// Save the user to the database
	if err := s.userRepo.CreateUser(newUser); err != nil {
		if err.Error() == "phone number already registered" {
			return nil, errors.New("phone number is already registered")
		}
		return nil, errors.New("failed to create user account")
	}

	return newUser, nil
}

func (s *UserService) LoginUser(req *models.LoginRequest) (string, *models.User, error) {
	// Find user by phone number
	user, err := s.userRepo.FindUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return "", nil, errors.New("invalid phone number or PIN")
	}

	// Verify PIN
	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin+user.Salt))
	if err != nil {
		return "", nil, errors.New("invalid phone number or PIN")
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.Id)
	if err != nil {
		return "", nil, errors.New("failed to generate authentication token")
	}

	return token, user, nil
}
