package db

import (
	"fmt"
	"log"
	"os"
	"test-wallet/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MySql *gorm.DB

// InitDB initializes the database connection.
func InitDB() {
	// Example: user:password@tcp(localhost:3306)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_DB_USER"),
		os.Getenv("MYSQL_DB_PASS"),
		os.Getenv("MYSQL_DB_HOST"),
		os.Getenv("MYSQL_DB_PORT"),
		os.Getenv("MYSQL_DB_NAME"),
	)

	var err error
	MySql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Auto Migrate Tables
	err = MySql.AutoMigrate(
		&models.User{},
		&models.Wallet{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	log.Println("✅ Connected to MySQL!")
}

// BeginTransaction starts a new database transaction.
func BeginTransaction() *gorm.DB {
	tx := MySql.Begin()
	if tx.Error != nil {
		log.Printf("Failed to begin transaction: %v", tx.Error) // Log the error
		return nil                                              // Return nil to indicate failure
	}
	return tx
}

// EndTransaction commits or rolls back a transaction.
func EndTransaction(tx *gorm.DB, shouldCommit bool) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil") // Handle nil transaction
	}
	if shouldCommit {
		if err := tx.Commit().Error; err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			return err
		}
	} else {
		if err := tx.Rollback().Error; err != nil {
			log.Printf("Failed to rollback transaction: %v", err)
			return err
		}
	}
	return nil
}

// GetDB returns the database connection.
func GetDB() *gorm.DB {
	return MySql
}
