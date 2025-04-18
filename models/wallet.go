package models

import (
	"time"
)

type Wallet struct {
	Id        string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserId    string    `gorm:"type:char(36);not null" json:"user_id"` // Must be the same type and unique
	Address   string    `gorm:"type:text;not null" json:"address"`
	Mnemonic  string    `gorm:"type:text;not null" json:"mnemonic"`
	QRCode    string    `gorm:"type:text" json:"qr_code"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// CreateWalletResponse defines the structure of the response when a new wallet is created.
// It contains the wallet's address and the private key.
type CreateWalletResponse struct {
	Mnemonic   string `json:"mnemonic"`    // passphrase
	Address    string `json:"address"`     // The Ethereum address associated with the wallet
	PrivateKey string `json:"private_key"` // The private key of the wallet in hexadecimal format
}

// SendETHRequest defines the structure of the request to send ETH from one address to another.
// It contains details such as the sender's address, private key, recipient's address, and the amount to be sent.
type SendETHRequest struct {
	FromAddress string `json:"from_address"`  // The address of the sender
	PrivateKey  string `json:"private_key"`   // The private key of the sender (used to sign the transaction)
	ToAddress   string `json:"to_address"`    // The recipient's Ethereum address
	AmountInETH string `json:"amount_in_eth"` // The amount of ETH to send, represented as a string
}

type SendERC20Request struct {
	FromAddress string `json:"from_address"`  // The address of the sender
	PrivateKey  string `json:"private_key"`   // The private key of the sender (used to sign the transaction)
	ToAddress   string `json:"to_address"`    // The recipient's Ethereum address
	AmountInUSD string `json:"amount_in_usd"` // The amount of ETH to send, represented as a string
}

type RecoverWalletRequest struct {
	Mnemonic       string `json:"mnemonic" binding:"required"`
	DerivationPath string `json:"derivation_path" binding:"required"`
}
