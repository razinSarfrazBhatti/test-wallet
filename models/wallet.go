package models

// CreateWalletResponse defines the structure of the response when a new wallet is created.
// It contains the wallet's address and the private key.
type CreateWalletResponse struct {
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
	FromAddress string `json:"from_address"`
	PrivateKey  string `json:"private_key"`
	ToAddress   string `json:"to_address"`
	AmountInUSD string `json:"amount_in_usd"`
}
