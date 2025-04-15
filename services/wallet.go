package services

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"test-wallet/models"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// CreateWallet generates a new Ethereum wallet with a private key and corresponding address.
// It returns the wallet's address and private key.
func CreateWallet() (*models.CreateWalletResponse, error) {
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
		PrivateKey: fmt.Sprintf("%x", privateKey.D), //converts private key to hex
	}, nil
}

// GetBalance retrieves the balance of an Ethereum wallet by address.
// It returns the balance in ETH.
func GetBalance(address string) (string, error) {
	// Dial the Ethereum client using Infura URL from configuration
	client, err := ethclient.Dial(os.Getenv("INFURA_URL"))
	if err != nil {
		return "", err // Return error if the connection to Ethereum client fails
	}

	// Convert the address string to an Ethereum address
	addr := common.HexToAddress(address)

	// Fetch the balance of the address
	balance, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return "", err // Return error if fetching the balance fails
	}

	// Convert the balance from Wei to ETH (1 ETH = 1e18 Wei)
	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))

	// Return the balance as a string
	return ethBalance.String(), nil
}

// SendETH sends ETH from one address to another.
// It takes a request containing the sender's private key, recipient address, and amount to send.
func SendETH(request models.SendETHRequest) (string, error) {
	// Dial the Ethereum client using Infura URL from configuration
	client, err := ethclient.Dial(os.Getenv("INFURA_URL"))
	if err != nil {
		return "", err // Return error if the connection to Ethereum client fails
	}

	// Convert the sender's private key from hex to ECDSA format
	privateKey, err := crypto.HexToECDSA(request.PrivateKey)
	if err != nil {
		return "", err // Return error if private key conversion fails
	}

	// Derive the sender's Ethereum address from the private key
	account := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Convert the recipient's address to Ethereum address format
	toAddress := common.HexToAddress(request.ToAddress)

	// Convert the amount in ETH to Wei
	amountInETH, _ := new(big.Float).SetString(request.AmountInETH)
	amountInWei, _ := amountInETH.Mul(amountInETH, big.NewFloat(1e18)).Int(nil)

	// Get the nonce (transaction count) for the sender's account
	nonce, _ := client.PendingNonceAt(context.Background(), account)

	// Suggest a gas price for the transaction
	gasPrice, _ := client.SuggestGasPrice(context.Background())

	// Apply a multiplier to the gas price (for faster transactions)
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	// Create a new transaction
	tx := types.NewTransaction(nonce, toAddress, amountInWei, 21000, gasPrice, nil)

	// Get the network ID (chain ID) for signing the transaction
	chainID, _ := client.NetworkID(context.Background())

	// Sign the transaction with the sender's private key
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	// Send the signed transaction to the Ethereum network
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err // Return error if sending the transaction fails
	}

	// Return the transaction hash as a hex string
	return signedTx.Hash().Hex(), nil
}

// SendERC20Token handles the logic for sending ERC20 tokens (e.g., USDC) from one address to another.
// It builds and signs a token transfer transaction and broadcasts it to the Ethereum network.
func SendERC20Token(request models.SendERC20Request) (string, error) {
	// Connect to the Ethereum network using Infura
	client, err := ethclient.Dial(os.Getenv("INFURA_URL"))
	if err != nil {
		return "", err
	}

	// Convert the hex-encoded private key into an ECDSA private key object
	privateKey, err := crypto.HexToECDSA(request.PrivateKey)
	if err != nil {
		return "", err
	}

	// Derive the sender's Ethereum address from the private key
	account := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Convert the recipient's address from string to Ethereum address format
	toAddress := common.HexToAddress(request.ToAddress)

	// USDC token contract address on Ethereum (change if using different token or network)
	usdcAddress := common.HexToAddress(os.Getenv("USDC_CONTRACT_ADDRESS"))

	// Convert the USD amount to USDC token amount in smallest units (USDC has 6 decimal places)
	amountInUSD, _ := new(big.Float).SetString(request.AmountInUSD)
	amountInWei := new(big.Int)
	amountInWei, _ = amountInUSD.Mul(amountInUSD, big.NewFloat(1e6)).Int(amountInWei)

	// Define the ERC20 ABI and pack the `transfer` method call with recipient and amount
	erc20ABI, _ := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	data, _ := erc20ABI.Pack("transfer", toAddress, amountInWei)

	// Get the current nonce for the sender account
	nonce, _ := client.PendingNonceAt(context.Background(), account)

	// Suggest a gas price for the transaction
	gasPrice, _ := client.SuggestGasPrice(context.Background())

	// Set a gas limit for the token transfer transaction
	gasLimit := uint64(100000) // typical for ERC20 token transfers

	// Construct the raw transaction (value is 0 since we’re not sending ETH)
	tx := types.NewTransaction(nonce, usdcAddress, big.NewInt(0), gasLimit, gasPrice, data)

	// Get the chain ID (required for signing the transaction)
	chainID, _ := client.NetworkID(context.Background())

	// Sign the transaction using the sender's private key
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	// Send the signed transaction to the network
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	// Return the transaction hash as confirmation
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

	// Convert the private key to hex
	privateKeyHex := fmt.Sprintf("0x%x", crypto.FromECDSA(privKey))

	return account.Address.Hex(), privateKeyHex, nil
}
