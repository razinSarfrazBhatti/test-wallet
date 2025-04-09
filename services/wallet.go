package services

import (
	"context"
	"fmt"
	"math/big"

	"test-wallet/config"
	"test-wallet/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CreateWallet() (*models.CreateWalletResponse, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	privateKeyHex := fmt.Sprintf("%x", privateKey.D)

	return &models.CreateWalletResponse{
		Address:    address.Hex(),
		PrivateKey: privateKeyHex,
	}, nil
}

func GetBalance(address string) (string, error) {
	client, err := ethclient.Dial(config.GetInfuraURL())
	if err != nil {
		return "", err
	}
	addr := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return "", err
	}
	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	return ethBalance.String(), nil
}

func SendETH(request models.SendETHRequest) (string, error) {
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

	amountInETH, _ := new(big.Float).SetString(request.AmountInETH)
	amountInWei, _ := amountInETH.Mul(amountInETH, big.NewFloat(1e18)).Int(nil)

	nonce, _ := client.PendingNonceAt(context.Background(), account)
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	tx := types.NewTransaction(nonce, toAddress, amountInWei, 21000, gasPrice, nil)
	chainID, _ := client.NetworkID(context.Background())
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
