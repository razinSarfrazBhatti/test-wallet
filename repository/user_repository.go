package repository

import (
	"errors"
	"test-wallet/db"
	"test-wallet/models"

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

func (r *UserRepository) PhoneNumberExists(phoneNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("phone_number = ?", phoneNumber).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	// Check if phone number already exists
	exists, err := r.PhoneNumberExists(user.PhoneNumber)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("phone number already registered")
	}

	tx := db.BeginTransaction()
	if tx == nil {
		return gorm.ErrInvalidTransaction
	}

	defer func() {
		if r := recover(); r != nil {
			db.EndTransaction(tx, false)
			panic(r)
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		db.EndTransaction(tx, false)
		return err
	}

	return db.EndTransaction(tx, true)
}

func (r *UserRepository) FindUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Wallet").Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Wallet").First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
