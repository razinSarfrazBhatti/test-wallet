package models

type CreateWalletResponse struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

type SendETHRequest struct {
	FromAddress string `json:"from_address"`
	PrivateKey  string `json:"private_key"`
	ToAddress   string `json:"to_address"`
	AmountInETH string `json:"amount_in_eth"`
}
