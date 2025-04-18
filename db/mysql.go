package db

import (
	"fmt"
	"test-wallet/config"
	"test-wallet/models"
	"test-wallet/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MySql *gorm.DB

// InitDB initializes the database connection
func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBConfig.User,
		config.AppConfig.DBConfig.Password,
		config.AppConfig.DBConfig.Host,
		config.AppConfig.DBConfig.Port,
		config.AppConfig.DBConfig.Name,
	)

	var err error
	MySql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate models
	if err := MySql.AutoMigrate(&models.User{}, &models.Wallet{}); err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	utils.LogInfo("Database connection established", nil)
	return nil
}

// BeginTransaction starts a new database transaction
func BeginTransaction() (*gorm.DB, error) {
	tx := MySql.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	return tx, nil
}

// EndTransaction commits or rolls back a transaction
func EndTransaction(tx *gorm.DB, shouldCommit bool) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}

	if shouldCommit {
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	} else {
		if err := tx.Rollback().Error; err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return MySql
}
