package repository

import (
	"errors"
	"fmt"
	"test-wallet/db"
	"test-wallet/models"
	"test-wallet/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: db.GetDB(),
	}
}

// PhoneNumberExists checks if a phone number is already registered
func (r *UserRepository) PhoneNumberExists(phoneNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("phone_number = ?", phoneNumber).Count(&count).Error
	if err != nil {
		utils.LogError(err, "Failed to check phone number existence", map[string]interface{}{
			"phone_number": phoneNumber,
		})
		return false, fmt.Errorf("failed to check phone number existence: %w", err)
	}

	exists := count > 0
	utils.LogDebug("Phone number existence checked", map[string]interface{}{
		"phone_number": phoneNumber,
		"exists":       exists,
	})

	return exists, nil
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	// Check if phone number already exists
	exists, err := r.PhoneNumberExists(user.PhoneNumber)
	if err != nil {
		return err
	}
	if exists {
		utils.LogInfo("Phone number already registered", map[string]interface{}{
			"phone_number": user.PhoneNumber,
		})
		return errors.New("phone number already registered")
	}

	// Start transaction
	tx, err := db.BeginTransaction()
	if err != nil {
		utils.LogError(err, "Failed to begin transaction", nil)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create user
	if err := tx.Create(user).Error; err != nil {
		utils.LogError(err, "Failed to create user", map[string]interface{}{
			"user_id": user.Id,
		})
		if err := db.EndTransaction(tx, false); err != nil {
			utils.LogError(err, "Failed to rollback transaction", nil)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Commit transaction
	if err := db.EndTransaction(tx, true); err != nil {
		utils.LogError(err, "Failed to commit transaction", nil)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	utils.LogInfo("User created successfully", map[string]interface{}{
		"user_id": user.Id,
	})

	return nil
}

// FindUserByPhoneNumber finds a user by their phone number
func (r *UserRepository) FindUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Wallet").Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogInfo("User not found", map[string]interface{}{
				"phone_number": phoneNumber,
			})
			return nil, fmt.Errorf("user not found")
		}
		utils.LogError(err, "Failed to find user by phone number", map[string]interface{}{
			"phone_number": phoneNumber,
		})
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	utils.LogDebug("User found by phone number", map[string]interface{}{
		"user_id":      user.Id,
		"phone_number": phoneNumber,
	})

	return &user, nil
}

// FindUserByID finds a user by their ID
func (r *UserRepository) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Wallet").First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogInfo("User not found", map[string]interface{}{
				"user_id": userID,
			})
			return nil, fmt.Errorf("user not found")
		}
		utils.LogError(err, "Failed to find user by ID", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	utils.LogDebug("User found by ID", map[string]interface{}{
		"user_id": userID,
	})

	return &user, nil
}
