package routes

import (
	"test-wallet/handlers"
	"test-wallet/middleware"
	"test-wallet/utils"

	"github.com/gin-gonic/gin"
)

func RegisterWalletRoutes(r *gin.Engine) {
	walletHandler, err := handlers.NewWalletHandler()
	if err != nil {
		utils.LogFatal(err, "Failed to create wallet handler", nil)
	}

	wallet := r.Group("/wallet")
	wallet.Use(middleware.AuthMiddleware())
	{
		wallet.GET("/balance/:address", walletHandler.GetBalance)
		wallet.POST("/send-eth", walletHandler.SendETH)
		wallet.POST("/send-erc20", walletHandler.SendERC20Token)
		wallet.POST("/recover", walletHandler.RecoverWalletHandler)
		wallet.GET("/qr", walletHandler.GenerateWalletQR)
	}
}
