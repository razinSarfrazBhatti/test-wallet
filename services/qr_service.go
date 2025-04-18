package services

import (
	"encoding/base64"
	"fmt"
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

// GenerateWalletQR generates a QR code for the user's wallet address
func (s *QRService) GenerateWalletQR(userID string) (string, error) {
	// Get user's wallet address
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		utils.LogError(err, "Failed to get user wallet", map[string]interface{}{
			"user_id": userID,
		})
		return "", fmt.Errorf("failed to get user wallet: %w", err)
	}

	// Generate QR code
	qr, err := qrcode.New(user.Wallet.Address, qrcode.Medium)
	if err != nil {
		utils.LogError(err, "Failed to generate QR code", map[string]interface{}{
			"user_id": userID,
			"wallet":  user.Wallet.Address,
		})
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Convert QR code to PNG
	png, err := qr.PNG(256)
	if err != nil {
		utils.LogError(err, "Failed to convert QR code to PNG", map[string]interface{}{
			"user_id": userID,
			"wallet":  user.Wallet.Address,
		})
		return "", fmt.Errorf("failed to convert QR code to PNG: %w", err)
	}

	// Convert PNG to base64
	base64QR := base64.StdEncoding.EncodeToString(png)

	utils.LogInfo("QR code generated successfully", map[string]interface{}{
		"user_id": userID,
		"wallet":  user.Wallet.Address,
	})

	return base64QR, nil
}
