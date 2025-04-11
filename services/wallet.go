package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"test-wallet/config"
	"test-wallet/models"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// CreateWallet generates a new Ethereum wallet with a private key and corresponding address.
// It returns the wallet's address and private key.
func CreateWallet() (*models.CreateWalletResponse, error) {
	// Generate a new private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err // Return error if private key generation fails
	}

	// Generate the Ethereum address from the public key
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Convert the private key to a hexadecimal string
	privateKeyHex := fmt.Sprintf("%x", privateKey.D)

	// Return the wallet's address and private key
	return &models.CreateWalletResponse{
		Address:    address.Hex(),
		PrivateKey: privateKeyHex,
	}, nil
}

// GetBalance retrieves the balance of an Ethereum wallet by address.
// It returns the balance in ETH.
func GetBalance(address string) (string, error) {
	// Dial the Ethereum client using Infura URL from configuration
	client, err := ethclient.Dial(config.GetInfuraURL())
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
	client, err := ethclient.Dial(config.GetInfuraURL())
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

func SendERC20Token(request models.SendERC20Request) (string, error) {
	client, err := ethclient.Dial(config.GetInfuraURL())
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(request.PrivateKey)
	if err != nil {
		return "", err
	}
	account := crypto.PubkeyToAddress(privateKey.PublicKey)
	toAddress := common.HexToAddress(request.ToAddress)
	usdcAddress := common.HexToAddress("0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238")

	amountInUSD, _ := new(big.Float).SetString(request.AmountInUSD)
	amountInWei := new(big.Int)
	amountInWei, _ = amountInUSD.Mul(amountInUSD, big.NewFloat(1e6)).Int(amountInWei)

	erc20ABI, _ := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	data, _ := erc20ABI.Pack("transfer", toAddress, amountInWei)

	nonce, _ := client.PendingNonceAt(context.Background(), account)
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	gasLimit := uint64(100000)
	tx := types.NewTransaction(nonce, usdcAddress, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, _ := client.NetworkID(context.Background())
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}
