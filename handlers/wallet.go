package handlers

import (
	"net/http"
	"test-wallet/models"
	"test-wallet/services"
	"test-wallet/utils"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService *services.WalletService
	qrService     *services.QRService
}

func NewWalletHandler() (*WalletHandler, error) {
	walletService, err := services.NewWalletService()
	if err != nil {
		return nil, err
	}

	return &WalletHandler{
		walletService: walletService,
		qrService:     services.NewQRService(),
	}, nil
}

// GetBalance handles fetching the ETH balance of a given address
func (h *WalletHandler) GetBalance(c *gin.Context) {
	address := c.Param("address")
	balance, err := h.walletService.GetBalance(address)
	if err != nil {
		utils.LogError(err, "Failed to get balance", map[string]interface{}{
			"address": address,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get balance"})
		return
	}

	utils.LogInfo("Balance retrieved successfully", map[string]interface{}{
		"address": address,
		"balance": balance,
	})

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// SendETH handles sending ETH from one address to another
func (h *WalletHandler) SendETH(c *gin.Context) {
	var request models.SendETHRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.LogError(err, "Invalid request payload", nil)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	txHash, err := h.walletService.SendETH(&request)
	if err != nil {
		utils.LogError(err, "Failed to send ETH", map[string]interface{}{
			"from":   request.FromAddress,
			"to":     request.ToAddress,
			"amount": request.AmountInETH,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to send ETH"})
		return
	}

	utils.LogInfo("ETH sent successfully", map[string]interface{}{
		"from":    request.FromAddress,
		"to":      request.ToAddress,
		"amount":  request.AmountInETH,
		"tx_hash": txHash,
	})

	c.JSON(http.StatusOK, gin.H{"transaction_hash": txHash})
}

// SendERC20Token handles sending ERC20 tokens
func (h *WalletHandler) SendERC20Token(c *gin.Context) {
	var request models.SendERC20Request
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.LogError(err, "Invalid request payload", nil)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	txHash, err := h.walletService.SendERC20Token(&request)
	if err != nil {
		utils.LogError(err, "Failed to send ERC20 token", map[string]interface{}{
			"from":   request.FromAddress,
			"to":     request.ToAddress,
			"amount": request.AmountInUSD,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to send ERC20 token"})
		return
	}

	utils.LogInfo("ERC20 token sent successfully", map[string]interface{}{
		"from":    request.FromAddress,
		"to":      request.ToAddress,
		"amount":  request.AmountInUSD,
		"tx_hash": txHash,
	})

	c.JSON(http.StatusOK, gin.H{"transaction_hash": txHash})
}

// RecoverWalletHandler handles wallet recovery
func (h *WalletHandler) RecoverWalletHandler(c *gin.Context) {
	var req models.RecoverWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError(err, "Invalid request payload", nil)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	address, privateKey, err := services.RecoverWalletFromMnemonic(req.Mnemonic, req.DerivationPath)
	if err != nil {
		utils.LogError(err, "Failed to recover wallet", nil)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to recover wallet"})
		return
	}

	utils.LogInfo("Wallet recovered successfully", map[string]interface{}{
		"address": address,
	})

	c.JSON(http.StatusOK, gin.H{
		"address":     address,
		"private_key": privateKey,
	})
}

// GetWalletQR generates a QR code for the user's wallet address
func (h *WalletHandler) GetWalletQR(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogError(nil, "User not authenticated", nil)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "user not authenticated"})
		return
	}

	// Get stored QR code
	qrCode, err := h.qrService.GetWalletQR(userID.(string))
	if err != nil {
		utils.LogError(err, "Failed to get QR code", map[string]interface{}{
			"user_id": userID,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get QR code"})
		return
	}

	// Get user's wallet address
	user, err := h.walletService.GetUserWallet(userID.(string))
	if err != nil {
		utils.LogError(err, "Failed to get wallet address", map[string]interface{}{
			"user_id": userID,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get wallet address"})
		return
	}

	utils.LogInfo("QR code retrieved successfully", map[string]interface{}{
		"user_id": userID,
		"wallet":  user.Wallet.Address,
	})

	c.JSON(http.StatusOK, models.QRCodeResponse{
		QRCode:  qrCode,
		Address: user.Wallet.Address,
	})
}
