package handlers

import (
	"net/http"
	"test-wallet/models"
	"test-wallet/services"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var request models.SendETHRequest
	// Bind the JSON body to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		// Return 400 if the request is invalid
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the service to send ETH
	txHash, err := services.SendETH(request)
	if err != nil {
		// Return 500 if sending ETH fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send ETH"})
		return
	}
	// Return the transaction hash
	c.JSON(http.StatusOK, gin.H{"transaction_hash": txHash})
}
