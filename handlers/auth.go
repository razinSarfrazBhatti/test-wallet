package handlers

import (
	"net/http"
	"test-wallet/models"
	"test-wallet/services"
	"test-wallet/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService: services.NewUserService(),
	}
}

// RegisterUser handles the registration of a new user
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req models.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError(err, "Invalid request payload", nil)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	user, err := h.userService.RegisterUser(&req)
	if err != nil {
		utils.LogError(err, "Failed to register user", map[string]interface{}{
			"phone_number": req.PhoneNumber,
		})
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	utils.LogInfo("User registered successfully", map[string]interface{}{
		"user_id": user.Id,
		"wallet":  user.Wallet.Address,
	})

	c.JSON(http.StatusCreated, models.RegisterResponse{
		Message:       "User registered successfully",
		UserID:        user.Id,
		WalletAddress: user.Wallet.Address,
	})
}

// LoginUser handles user authentication and returns a JWT token
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError(err, "Invalid request payload", nil)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	token, user, err := h.userService.LoginUser(&req)
	if err != nil {
		utils.LogError(err, "Failed to login user", map[string]interface{}{
			"phone_number": req.PhoneNumber,
		})
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}

	utils.LogInfo("User logged in successfully", map[string]interface{}{
		"user_id": user.Id,
	})

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User: struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			Wallet      struct {
				Address string `json:"address"`
			} `json:"wallet"`
		}{
			ID:          user.Id,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Wallet: struct {
				Address string `json:"address"`
			}{
				Address: user.Wallet.Address,
			},
		},
	})
}
