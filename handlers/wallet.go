package handlers

import (
	"net/http"

	"test-wallet/models"
	"test-wallet/services"

	"github.com/gin-gonic/gin"
)

// CreateWallet handles the creation of a new Ethereum wallet.
// It calls the service layer to generate a wallet and returns it in the response.
func CreateWallet(c *gin.Context) {
	wallet, err := services.CreateWallet()
	if err != nil {
		// Return 500 if wallet creation fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}
	// Return the newly created wallet
	c.JSON(http.StatusOK, wallet)
}

// GetBalance handles fetching the ETH balance of a given address.
// The address is extracted from the URL parameter.
func GetBalance(c *gin.Context) {
	address := c.Param("address")
	balance, err := services.GetBalance(address)
	if err != nil {
		// Return 500 if balance retrieval fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
		return
	}
	// Return the balance of the address
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// SendETH handles sending ETH from one address to another.
// It expects a JSON body with fields defined in the SendETHRequest model.
func SendETH(c *gin.Context) {
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

// SendERC20Token handles the HTTP request to send ERC20 tokens.
func SendERC20Token(c *gin.Context) {
	var request models.SendERC20Request

	// Bind the incoming JSON payload to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		// Return a 400 Bad Request response if the request is malformed
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the service layer to perform the token transfer
	txHash, err := services.SendERC20Token(request)
	if err != nil {
		// Return a 500 Internal Server Error if the transfer fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send ERC20 token"})
		return
	}

	// Return the transaction hash in a 200 OK response
	c.JSON(http.StatusOK, gin.H{"transaction_hash": txHash})
}

func RecoverWalletHandler(c *gin.Context) {
	var req models.RecoverWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, privateKey, err := services.RecoverWalletFromMnemonic(req.Mnemonic, req.DerivationPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address":     address,
		"private_key": privateKey,
	})
}
