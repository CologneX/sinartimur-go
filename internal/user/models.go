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
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password" validate:"required"`
	Username string    `json:"username" validate:"required"`
	IsActive bool      `json:"is_active" validate:"boolean"`
}

type UserRole struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GetUserResponse struct {
	ID        uuid.UUID   `json:"id"`
	Username  string      `json:"username"`
	Role      *[]UserRole `json:"role"`
	IsActive  bool        `json:"is_active"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}
