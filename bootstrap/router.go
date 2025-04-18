package bootstrap

import (
	"os"
	"test-wallet/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures the Gin router
func SetupRouter() *gin.Engine {
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	if os.Getenv("ENV") == "development" {
		router.Use(gin.Logger())
	}

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	routes.RegisterRoutes(router)

	return router
}
