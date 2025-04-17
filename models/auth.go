package models

import "time"

type RegisterUserRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Pin         string `json:"pin"`
}

type User struct {
	Id          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	PhoneNumber string    `gorm:"type:text;not null" json:"phone_number"`
	Pin         string    `gorm:"type:text;not null" json:"pin"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	Salt        string    `gorm:"type:text;not null" json:"salt"`
	// Has One relationship (no foreignKey tag here)
	Wallet Wallet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"wallet"`
}
