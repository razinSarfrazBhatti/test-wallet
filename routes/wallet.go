package routes

import (
	"test-wallet/handlers"
	"test-wallet/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWalletRoutes(r *gin.Engine) {
	wallet := r.Group("/wallet")
	wallet.Use(middleware.AuthMiddleware())
	{
		wallet.GET("/balance/:address", handlers.GetBalance)
		wallet.POST("/send-eth", handlers.SendETH)
		wallet.POST("/send-erc20", handlers.SendERC20Token)
		wallet.POST("/recover", handlers.RecoverWalletHandler)
	}
}
