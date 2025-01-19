package user

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type CreateUserRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
