package routes

import (
	"test-wallet/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	authHandler := handlers.NewAuthHandler()

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.RegisterUser)
		auth.POST("/login", authHandler.LoginUser)
		auth.POST("/reset-pin" /*handlers.ResetPassword*/)
	}
}
