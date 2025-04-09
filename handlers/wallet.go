package handlers

import (
	"net/http"

	"test-wallet/models"
	"test-wallet/services"

	"github.com/gin-gonic/gin"
)

func CreateWallet(c *gin.Context) {
	wallet, err := services.CreateWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

func GetBalance(c *gin.Context) {
	address := c.Param("address")
	balance, err := services.GetBalance(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func SendETH(c *gin.Context) {
	var request models.SendETHRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	txHash, err := services.SendETH(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send ETH"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transaction_hash": txHash})
}
