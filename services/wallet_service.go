package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"test-wallet/config"
	"test-wallet/models"
	"test-wallet/repository"
	"test-wallet/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"golang.org/x/crypto/bcrypt"
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

func (s *WalletService) CreateWallet() (*models.CreateWalletResponse, error) {
	// Generate a random mnemonic (12-word by default)
	mnemonic, err := hdwallet.NewMnemonic(128)
	if err != nil {
		return nil, err
	}

	// Create the HD wallet from the mnemonic
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	// Derive a path (you can customize this for multiple accounts)
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	// 	1. m
	//     The master node (root of the HD wallet tree).
	// 2. 44' — Purpose
	//     Follows BIP-44, a standard for multi-account hierarchical deterministic wallets.
	//     The ' indicates it's a "hardened" key (more secure and isolated).
	// 3. 60' — Coin Type
	//     60 is the coin type for Ethereum, defined by SLIP-44.
	// 4. 0' — Account
	//     A unique account index — use 0 for the first wallet. You could use 1' or 2' for multiple user wallets.
	// 5. 0 — Change
	//     0 = external chain (for receiving).
	//     1 = internal chain (used for change addresses, not typically used in Ethereum).
	// 6. 0 — Address Index
	//     Index of the address under the chain. You can increment this (0, 1, 2, …) to get more addresses.
	// Get the private key
	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	// Return the wallet info
	return &models.CreateWalletResponse{
		Mnemonic:   mnemonic,
		Address:    account.Address.Hex(),
		PrivateKey: fmt.Sprintf("%x", privateKey.D),
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
func (s *WalletService) GetBalance(c *gin.Context, address string) (string, error) {
	addr := common.HexToAddress(address)
	balance, err := s.client.BalanceAt(c, addr, nil)
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
func (s *WalletService) SendETH(userID string, req *models.SendETHRequest) (string, error) {
	// Get user's wallet
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		utils.LogError(err, "Failed to get user wallet", map[string]interface{}{
			"user_id": userID,
		})
		return "", fmt.Errorf("failed to get user wallet: %w", err)
	}

	// Verify PIN
	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin+user.Salt))
	if err != nil {
		utils.LogError(err, "Invalid PIN", map[string]interface{}{
			"user_id": userID,
		})
		return "", fmt.Errorf("invalid PIN")
	}

	// Decrypt the mnemonic using the provided PIN
	mnemonic, err := Decrypt(req.Pin, user.Wallet.Mnemonic)
	if err != nil {
		utils.LogError(err, "Failed to decrypt mnemonic", nil)
		return "", fmt.Errorf("failed to decrypt mnemonic: %w", err)
	}

	// Recover wallet from mnemonic
	address, privateKey, err := RecoverWalletFromMnemonic(mnemonic, "m/44'/60'/0'/0/0")
	if err != nil {
		utils.LogError(err, "Failed to recover wallet", nil)
		return "", fmt.Errorf("failed to recover wallet: %w", err)
	}

	// Convert private key from hex to ECDSA
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		utils.LogError(err, "Failed to convert private key", nil)
		return "", fmt.Errorf("failed to convert private key: %w", err)
	}

	// Get the public key and address
	publicKey := privKey.Public()
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

	// Apply a multiplier to the gas price (for faster transactions)
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		GasPrice: gasPrice,
		Value:    amountInWei,
		Data:     nil,
	}

	gasLimit, err := s.client.EstimateGas(context.Background(), msg)
	if err != nil {
		utils.LogError(err, "Failed to estimate gas", nil)
		return "", fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Create transaction
	tx := types.NewTransaction(nonce, toAddress, amountInWei, gasLimit, gasPrice, nil)

	// Get chain ID
	chainID, err := s.client.NetworkID(context.Background())
	if err != nil {
		utils.LogError(err, "Failed to get chain ID", nil)
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
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
		"from":    address,
		"to":      req.ToAddress,
		"amount":  req.AmountInETH,
		"tx_hash": signedTx.Hash().Hex(),
	})

	return signedTx.Hash().Hex(), nil
}

// SendERC20Token sends ERC20 tokens from one address to another
func (s *WalletService) SendERC20Token(userID string, req *models.SendERC20Request) (string, error) {
	// Get user's wallet
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		utils.LogError(err, "Failed to get user wallet", map[string]interface{}{
			"user_id": userID,
		})
		return "", fmt.Errorf("failed to get user wallet: %w", err)
	}

	// Verify PIN
	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin+user.Salt))
	if err != nil {
		utils.LogError(err, "Invalid PIN", map[string]interface{}{
			"user_id": userID,
		})
		return "", fmt.Errorf("invalid PIN")
	}

	// Decrypt the mnemonic using the provided PIN
	mnemonic, err := Decrypt(req.Pin, user.Wallet.Mnemonic)
	if err != nil {
		utils.LogError(err, "Failed to decrypt mnemonic", nil)
		return "", fmt.Errorf("failed to decrypt mnemonic: %w", err)
	}

	// Recover wallet from mnemonic
	address, privateKey, err := RecoverWalletFromMnemonic(mnemonic, "m/44'/60'/0'/0/0")
	if err != nil {
		utils.LogError(err, "Failed to recover wallet", nil)
		return "", fmt.Errorf("failed to recover wallet: %w", err)
	}

	// Convert private key from hex to ECDSA
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		utils.LogError(err, "Failed to convert private key", nil)
		return "", fmt.Errorf("failed to convert private key: %w", err)
	}

	// Get the public key and address
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to get public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddress := common.HexToAddress(req.ToAddress)

	// USDC token contract address on Ethereum (change if using different token or network)
	usdcAddress := common.HexToAddress(config.AppConfig.EthConfig.USDCContractAddr)

	// Convert the USD amount to USDC token amount in smallest units (USDC has 6 decimal places)
	amountInUSD, _ := new(big.Float).SetString(req.AmountInUSD)
	amountInWei := new(big.Int)
	amountInWei, _ = amountInUSD.Mul(amountInUSD, big.NewFloat(1e6)).Int(amountInWei)

	// Define the ERC20 ABI and pack the `transfer` method call with recipient and amount
	erc20ABI, _ := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	data, _ := erc20ABI.Pack("transfer", toAddress, amountInWei)

	// Get the current nonce for the sender account
	nonce, _ := s.client.PendingNonceAt(context.Background(), fromAddress)

	// Suggest a gas price for the transaction
	gasPrice, _ := s.client.SuggestGasPrice(context.Background())

	// Apply a multiplier to the gas price (for faster transactions)
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	// Set a gas limit for the token transfer transaction
	gasLimit := uint64(100000) // typical for ERC20 token transfers

	// Construct the raw transaction (value is 0 since we're not sending ETH)
	tx := types.NewTransaction(nonce, usdcAddress, big.NewInt(0), gasLimit, gasPrice, data)

	// Get the chain ID (required for signing the transaction)
	chainID, _ := s.client.NetworkID(context.Background())

	// Sign the transaction using the sender's private key
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)

	// Send the signed transaction to the network
	err = s.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	utils.LogInfo("ERC20 token sent successfully", map[string]interface{}{
		"from":    address,
		"to":      req.ToAddress,
		"amount":  req.AmountInUSD,
		"tx_hash": signedTx.Hash().Hex(),
	})

	return signedTx.Hash().Hex(), nil
}

// RecoverWalletFromMnemonic recovers a wallet using mnemonic and derivation path
func RecoverWalletFromMnemonic(mnemonic, derivationPath string) (string, string, error) {
	// Create a new wallet from the mnemonic
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", "", fmt.Errorf("failed to create wallet: %w", err)
	}

	// Parse the derivation path (e.g. m/44'/60'/0'/0/0)
	path := hdwallet.MustParseDerivationPath(derivationPath)

	// Derive the account
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", "", fmt.Errorf("failed to derive account: %w", err)
	}

	// Get the private key
	privKey, err := wallet.PrivateKey(account)
	if err != nil {
		return "", "", fmt.Errorf("failed to get private key: %w", err)
	}

	// Convert the private key to hex without the '0x' prefix
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privKey))

	return account.Address.Hex(), privateKeyHex, nil
}
