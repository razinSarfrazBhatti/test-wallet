package routes

import (
	"test-wallet/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/create-wallet", handlers.CreateWallet)
	r.GET("/get-balance/:address", handlers.GetBalance)
	r.POST("/send-eth", handlers.SendETH)
}
