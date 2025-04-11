package routes

import (
	"test-wallet/handlers" // Import handler functions for route logic

	"github.com/gin-gonic/gin" // Import Gin web framework
)

// RegisterRoutes sets up the API routes for the application.
func RegisterRoutes(r *gin.Engine) {
	// Route to create a new Ethereum wallet
	r.POST("/create-wallet", handlers.CreateWallet)

	// Route to get the ETH balance of a specific wallet address
	r.GET("/get-balance/:address", handlers.GetBalance)

	// Route to send ETH from one wallet to another
	r.POST("/send-eth", handlers.SendETH)

	// Route to send erc20 from one wallet to another
	r.POST("/send-erc20", handlers.SendERC20Token)
}
