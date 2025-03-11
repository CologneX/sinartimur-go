package customer

import (
	"github.com/google/uuid"
	"sinartimur-go/utils"
	"time"
)

// Customer represents a customer entity
type Customer struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Address   string     `json:"address"`
	Telephone string     `json:"telephone"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// CreateCustomerRequest represents the data needed to create a customer
type CreateCustomerRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=255"`
	Address   string `json:"address" validate:"omitempty,max=1000"`
	Telephone string `json:"telephone" validate:"omitempty,max=50"`
}

// UpdateCustomerRequest represents the data needed to update a customer
type UpdateCustomerRequest struct {
	ID        uuid.UUID `json:"id" validate:"required,uuid"`
	Name      string    `json:"name" validate:"required,min=2,max=255"`
	Address   string    `json:"address" validate:"omitempty,max=1000"`
	Telephone string    `json:"telephone" validate:"omitempty,max=50"`
}

// GetCustomerResponse represents the customer data returned from read operations
type GetCustomerResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetCustomerRequest represents the parameters for filtering customer queries
type GetCustomerRequest struct {
	Name      string `json:"name,omitempty" validate:"omitempty,max=255"`
	Address   string `json:"address,omitempty" validate:"omitempty,max=1000"`
	Telephone string `json:"telephone,omitempty" validate:"omitempty,max=50"`
	utils.PaginationParameter
}

// DeleteCustomerRequest represents the request to delete a customer
type DeleteCustomerRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

// GetCustomerByIDRequest represents the request to get a customer by ID
type GetCustomerByIDRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

// CustomerPaginationResponse represents the paginated response for customers
type CustomerPaginationResponse struct {
	Data       []GetCustomerResponse `json:"data"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalItems int                   `json:"total_items"`
}
