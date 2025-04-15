package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register" /*handlers.RegisterUser*/)
		auth.POST("/login" /*handlers.LoginUser*/)
		auth.POST("/reset-password" /*handlers.ResetPassword*/)
	}
}
