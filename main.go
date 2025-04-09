package main

import (
	"log"
	"test-wallet/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")
}
