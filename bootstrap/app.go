package bootstrap

import (
	"os"

	"test-wallet/config"
	"test-wallet/db"
	"test-wallet/utils"

	"github.com/gin-gonic/gin"
)

// InitializeApp performs all necessary initializations
func InitializeApp() error {
	// Initialize logger
	utils.InitLogger()
	utils.LogInfo("Starting application", nil)

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		return err
	}

	// Initialize database
	if err := db.InitDB(); err != nil {
		return err
	}

	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return nil
}
