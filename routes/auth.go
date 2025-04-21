package routes

import (
	"test-wallet/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	authHandler, err := handlers.NewAuthHandler()
	if err != nil {
		panic(err) // Handle error appropriately in production
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.RegisterUser)
		auth.POST("/login", authHandler.LoginUser)
		auth.POST("/reset-pin" /*handlers.ResetPin*/)
	}
}
