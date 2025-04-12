package purchase

import (
	"sinartimur-go/utils"
)

// GetSupplierResponse represents the response for a supplier retrieval
type GetSupplierResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetSupplierRequest represents the request for filtering suppliers
type GetSupplierRequest struct {
	Name      string `json:"name,omitempty" validate:"omitempty"`
	Telephone string `json:"telephone,omitempty" validate:"omitempty"`
	utils.PaginationParameter
}

// CreateSupplierRequest represents the request for creating a supplier
type CreateSupplierRequest struct {
	Name      string `json:"name" validate:"required"`
	Address   string `json:"address,omitempty" validate:"omitempty"`
	Telephone string `json:"telephone,omitempty" validate:"omitempty"`
}

// UpdateSupplierRequest represents the request for updating a supplier
type UpdateSupplierRequest struct {
	ID        string `json:"id" validate:"required,uuid"`
	Name      string `json:"name,omitempty" validate:"omitempty"`
	Address   string `json:"address,omitempty" validate:"omitempty"`
	Telephone string `json:"telephone,omitempty" validate:"omitempty"`
}

// DeleteSupplierRequest represents the request for deleting a supplier
type DeleteSupplierRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}
