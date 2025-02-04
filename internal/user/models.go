package user

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

type CreateUserRequest struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
	IsAdmin         bool   `json:"is_admin" validate:"boolean"`
	IsHr            bool   `json:"is_hr" validate:"boolean"`
	IsFinance       bool   `json:"is_finance" validate:"boolean"`
	IsInventory     bool   `json:"is_inventory" validate:"boolean"`
	IsSales         bool   `json:"is_sales" validate:"boolean"`
	IsPurchase      bool   `json:"is_purchase" validate:"boolean"`
}

type UpdateUserRequest struct {
	ID          uuid.UUID `json:"id"`
	IsAdmin     bool      `json:"is_admin" validate:"boolean"`
	IsHr        bool      `json:"is_hr" validate:"boolean"`
	IsFinance   bool      `json:"is_finance" validate:"boolean"`
	IsInventory bool      `json:"is_inventory" validate:"boolean"`
	IsSales     bool      `json:"is_sales" validate:"boolean"`
	IsPurchase  bool      `json:"is_purchase" validate:"boolean"`
	Password    string    `json:"password" validate:"required"`
	Username    string    `json:"username" validate:"required"`
	IsActive    bool      `json:"is_active" validate:"boolean"`
}

//type UserRole struct {
//	ID   uuid.UUID `json:"id"`
//	Name string    `json:"name"`
//}

type GetUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      *[]string `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
