package models

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		Wallet      struct {
			Address string `json:"address"`
		} `json:"wallet"`
	} `json:"user"`
}

type RegisterResponse struct {
	Message       string `json:"message"`
	UserID        string `json:"user_id"`
	WalletAddress string `json:"wallet_address"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
