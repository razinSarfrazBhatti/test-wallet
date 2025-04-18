package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"test-wallet/config"
	"test-wallet/models"
	"test-wallet/repository"
	"test-wallet/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type WalletService struct {
	userRepo *repository.UserRepository
	client   *ethclient.Client
}

func NewWalletService() (*WalletService, error) {
	client, err := ethclient.Dial(config.AppConfig.EthConfig.InfuraURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	return &WalletService{
		userRepo: repository.NewUserRepository(),
		client:   client,
	}, nil
}

// GetUserWallet retrieves a user's wallet information
func (s *WalletService) GetUserWallet(userID string) (*models.User, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		utils.LogError(err, "Failed to get user wallet", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get user wallet: %w", err)
	}

	utils.LogDebug("Retrieved user wallet", map[string]interface{}{
		"user_id": userID,
		"wallet":  user.Wallet.Address,
	})

	return user, nil
}

// GetBalance retrieves the ETH balance of a given address
func (s *WalletService) GetBalance(address string) (string, error) {
	addr := common.HexToAddress(address)
	balance, err := s.client.BalanceAt(nil, addr, nil)
	if err != nil {
		utils.LogError(err, "Failed to get balance", map[string]interface{}{
			"address": address,
		})
		return "", fmt.Errorf("failed to get balance: %w", err)
	}

	// Convert balance from Wei to ETH
	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))

	utils.LogDebug("Retrieved balance", map[string]interface{}{
		"address": address,
		"balance": ethBalance.String(),
	})

	return ethBalance.String(), nil
}

// SendETH sends ETH from one address to another
func (s *WalletService) SendETH(req *models.SendETHRequest) (string, error) {
	// Convert private key from hex to ECDSA
	privateKey, err := crypto.HexToECDSA(req.PrivateKey)
	if err != nil {
		utils.LogError(err, "Failed to convert private key", nil)
		return "", fmt.Errorf("failed to convert private key: %w", err)
	}

	// Get the public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to get public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddress := common.HexToAddress(req.ToAddress)

	// Convert amount from ETH to Wei
	amountInETH, _ := new(big.Float).SetString(req.AmountInETH)

	amountInWei, _ := amountInETH.Mul(amountInETH, big.NewFloat(1e18)).Int(nil)

	// Get nonce
	nonce, err := s.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		utils.LogError(err, "Failed to get nonce", nil)
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get gas price
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		utils.LogError(err, "Failed to get gas price", nil)
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	// Create transaction
	tx := types.NewTransaction(nonce, toAddress, amountInWei, 21000, gasPrice, nil)

	// Get chain ID
	chainID, err := s.client.NetworkID(context.Background())
	if err != nil {
		utils.LogError(err, "Failed to get chain ID", nil)
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		utils.LogError(err, "Failed to sign transaction", nil)
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = s.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		utils.LogError(err, "Failed to send transaction", nil)
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	utils.LogInfo("ETH sent successfully", map[string]interface{}{
		"from":    req.FromAddress,
		"to":      req.ToAddress,
		"amount":  req.AmountInETH,
		"tx_hash": signedTx.Hash().Hex(),
	})

	return signedTx.Hash().Hex(), nil
}

// SendERC20Token sends ERC20 tokens from one address to another
func (s *WalletService) SendERC20Token(req *models.SendERC20Request) (string, error) {
	// Convert private key from hex to ECDSA
	privateKey, err := crypto.HexToECDSA(req.PrivateKey)
	if err != nil {
		utils.LogError(err, "Failed to convert private key", nil)
		return "", fmt.Errorf("failed to convert private key: %w", err)
	}

	// Get the public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to get public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddress := common.HexToAddress(req.ToAddress)
	tokenAddress := common.HexToAddress(config.AppConfig.EthConfig.USDCContractAddr)

	// Convert amount from USD to token units (6 decimals for USDC)
	amount := new(big.Int)
	amount.SetString(req.AmountInUSD, 10)
	amount.Mul(amount, big.NewInt(1e6))

	// Get nonce
	nonce, err := s.client.PendingNonceAt(nil, fromAddress)
	if err != nil {
		utils.LogError(err, "Failed to get nonce", nil)
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get gas price
	gasPrice, err := s.client.SuggestGasPrice(nil)
	if err != nil {
		utils.LogError(err, "Failed to get gas price", nil)
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	// Create transfer function data
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	methodID := hash[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// Create transaction
	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), 100000, gasPrice, data)

	// Get chain ID
	chainID, err := s.client.NetworkID(nil)
	if err != nil {
		utils.LogError(err, "Failed to get chain ID", nil)
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		utils.LogError(err, "Failed to sign transaction", nil)
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = s.client.SendTransaction(nil, signedTx)
	if err != nil {
		utils.LogError(err, "Failed to send transaction", nil)
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	utils.LogInfo("ERC20 token sent successfully", map[string]interface{}{
		"from":    req.FromAddress,
		"to":      req.ToAddress,
		"amount":  req.AmountInUSD,
		"tx_hash": signedTx.Hash().Hex(),
	})

	return signedTx.Hash().Hex(), nil
}
