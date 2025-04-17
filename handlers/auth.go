package handlers

import (
	"net/http"
	"test-wallet/db"
	"test-wallet/models"
	"test-wallet/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles the registration of a new user.
// It takes user details, encrypts the PIN, creates a wallet,
// encrypts the wallet's private key and mnemonic with the PIN,
// and saves the user and wallet details to the database.
func RegisterUser(c *gin.Context) {
	var req models.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 1. Generate Salt
	salt, err := services.GenerateSalt(16) // Use a 16-byte salt (adjust as needed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate salt"})
		return
	}

	// 2. Hash the PIN with the Salt
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.Pin+salt), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt PIN"})
		return
	}

	// 3. Create a new Ethereum wallet
	walletResponse, err := services.CreateWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	// 4. Encrypt the mnemonic with the user's PIN
	encryptedMnemonic, err := services.Encrypt(req.Pin, walletResponse.Mnemonic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt mnemonic"})
		return
	}

	// 5. Create a new user record
	newUser := models.User{
		Id:          uuid.New().String(),
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Pin:         string(hashedPin), // Store the salted hash
		Salt:        salt,              // Store the salt
		Wallet: models.Wallet{
			Id:       uuid.New().String(),
			Address:  walletResponse.Address,
			Mnemonic: encryptedMnemonic,
		},
	}

	// 6. Save the user and wallet to the database within a transaction
	tx := db.BeginTransaction()
	if tx == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin database transaction"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			db.EndTransaction(tx, false) // Rollback on panic
			panic(r)
		}
	}()

	if err := tx.Create(&newUser).Error; err != nil {
		db.EndTransaction(tx, false)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	if err := db.EndTransaction(tx, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed: " + err.Error()})
		return
	}

	// Respond with success (you might want to return limited user information)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": newUser.Id, "wallet_address": newUser.Wallet.Address})
}
