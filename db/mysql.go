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

	err = MySql.AutoMigrate(&models.Wallet{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Auto Migrate Tables
	err = MySql.AutoMigrate(&models.Wallet{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	log.Println("✅ Connected to MySQL!")
}
