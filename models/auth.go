package models

type RegisterUserRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Pin         string `json:"pin"`
}
