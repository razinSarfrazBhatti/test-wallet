package services

import (
	"encoding/base64"
	"fmt"
	"test-wallet/models"
	"test-wallet/repository"
	"test-wallet/utils"

	"github.com/skip2/go-qrcode"
)

type QRService struct {
	userRepo *repository.UserRepository
}

func NewQRService() *QRService {
	return &QRService{
		userRepo: repository.NewUserRepository(),
	}
}

// GenerateAndStoreQR generates a QR code for a wallet address and stores it
func (s *QRService) GenerateAndStoreQR(wallet *models.Wallet) error {
	// Generate QR code
	qr, err := qrcode.New(wallet.Address, qrcode.Medium)
	if err != nil {
		utils.LogError(err, "Failed to generate QR code", map[string]interface{}{
			"wallet": wallet.Address,
		})
		return fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Convert QR code to PNG
	png, err := qr.PNG(256)
	if err != nil {
		utils.LogError(err, "Failed to convert QR code to PNG", map[string]interface{}{
			"wallet": wallet.Address,
		})
		return fmt.Errorf("failed to convert QR code to PNG: %w", err)
	}

	// Convert PNG to base64
	base64QR := base64.StdEncoding.EncodeToString(png)
	wallet.QRCode = base64QR

	utils.LogInfo("QR code generated and stored", map[string]interface{}{
		"wallet": wallet.Address,
	})

	return nil
}

// GetWalletQR retrieves the stored QR code for a user's wallet
func (s *QRService) GetWalletQR(userID string) (string, error) {
	// Get user's wallet
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		utils.LogError(err, "Failed to get user wallet", map[string]interface{}{
			"user_id": userID,
		})
		return "", fmt.Errorf("failed to get user wallet: %w", err)
	}

	if user.Wallet.QRCode == "" {
		utils.LogError(nil, "QR code not found for wallet", map[string]interface{}{
			"user_id": userID,
			"wallet":  user.Wallet.Address,
		})
		return "", fmt.Errorf("qr code not found for wallet")
	}

	utils.LogInfo("QR code retrieved successfully", map[string]interface{}{
		"user_id": userID,
		"wallet":  user.Wallet.Address,
	})

	return user.Wallet.QRCode, nil
}
