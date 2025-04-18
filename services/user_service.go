package services

import (
	"errors"
	"fmt"
	"test-wallet/middleware"
	"test-wallet/models"
	"test-wallet/repository"
	"test-wallet/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  *repository.UserRepository
	qrService *QRService
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:  repository.NewUserRepository(),
		qrService: NewQRService(),
	}
}

func (s *UserService) RegisterUser(req *models.RegisterUserRequest) (*models.User, error) {
	// Check if phone number already exists
	exists, err := s.userRepo.PhoneNumberExists(req.PhoneNumber)
	if err != nil {
		utils.LogError(err, "Failed to check phone number availability", map[string]interface{}{
			"phone_number": req.PhoneNumber,
		})
		return nil, errors.New("failed to check phone number availability")
	}
	if exists {
		utils.LogInfo("Phone number already registered", map[string]interface{}{
			"phone_number": req.PhoneNumber,
		})
		return nil, errors.New("phone number is already registered")
	}

	// Generate salt
	salt, err := GenerateSalt(16)
	if err != nil {
		utils.LogError(err, "Failed to generate salt", nil)
		return nil, errors.New("failed to generate salt")
	}

	// Hash the PIN with the salt
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.Pin+salt), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError(err, "Failed to hash PIN", nil)
		return nil, errors.New("failed to hash PIN")
	}

	// Create a new Ethereum wallet
	walletResponse, err := CreateWallet()
	if err != nil {
		utils.LogError(err, "Failed to create wallet", nil)
		return nil, errors.New("failed to create wallet")
	}

	// Encrypt the mnemonic with the user's PIN
	encryptedMnemonic, err := Encrypt(req.Pin, walletResponse.Mnemonic)
	if err != nil {
		utils.LogError(err, "Failed to encrypt wallet mnemonic", nil)
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

	// Generate and store QR code
	if err := s.qrService.GenerateAndStoreQR(&newUser.Wallet); err != nil {
		utils.LogError(err, "Failed to generate QR code", map[string]interface{}{
			"user_id": newUser.Id,
		})
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Save the user to the database
	if err := s.userRepo.CreateUser(newUser); err != nil {
		utils.LogError(err, "Failed to create user account", map[string]interface{}{
			"user_id": newUser.Id,
		})
		return nil, fmt.Errorf("failed to create user account: %w", err)
	}

	utils.LogInfo("User registered successfully", map[string]interface{}{
		"user_id": newUser.Id,
		"wallet":  newUser.Wallet.Address,
	})

	return newUser, nil
}

func (s *UserService) LoginUser(req *models.LoginRequest) (string, *models.User, error) {
	// Find user by phone number
	user, err := s.userRepo.FindUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		utils.LogError(err, "User not found", map[string]interface{}{
			"phone_number": req.PhoneNumber,
		})
		return "", nil, errors.New("invalid phone number or PIN")
	}

	// Verify PIN
	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin+user.Salt))
	if err != nil {
		utils.LogError(err, "Invalid PIN", map[string]interface{}{
			"user_id": user.Id,
		})
		return "", nil, errors.New("invalid phone number or PIN")
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.Id)
	if err != nil {
		utils.LogError(err, "Failed to generate JWT token", map[string]interface{}{
			"user_id": user.Id,
		})
		return "", nil, errors.New("failed to generate authentication token")
	}

	utils.LogInfo("User logged in successfully", map[string]interface{}{
		"user_id": user.Id,
	})

	return token, user, nil
}
