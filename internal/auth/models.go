package auth

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	IsActive     bool      `json:"is_active"`
	IsAdmin      bool      `json:"is_admin"`
	IsHr         bool      `json:"is_hr"`
	IsFinance    bool      `json:"is_finance"`
	IsInventory  bool      `json:"is_inventory"`
	IsSales      bool      `json:"is_sales"`
	IsPurchase   bool      `json:"is_purchase"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	Id       string    `json:"id"`
	Username string    `json:"username"`
	Roles    []*string `json:"roles"`
}
