package main

import (
	"log"                // Import the log package for logging errors
	"test-wallet/routes" // Import the routes package to register routes

	"github.com/gin-gonic/gin" // Import the Gin web framework
	"github.com/joho/godotenv" // Import the godotenv package to load environment variables
)

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		// Log and stop execution if loading .env file fails
		log.Fatal("Error loading .env file")
	}

	// Create a default Gin router
	r := gin.Default()

	// Register the API routes defined in the routes package
	routes.RegisterRoutes(r)

	// Start the Gin server on port 8080
	r.Run(":8080")
}
