package product

import (
	"github.com/google/uuid"
	"sinartimur-go/utils"
)

// GetProductRequest is the payload for fetching product.
type GetProductRequest struct {
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty" validate:"omitempty,uuid"`
	Unit     string `json:"unit,omitempty" validate:"omitempty,uuid"`
	// Include PaginationParameter
	utils.PaginationParameter
}

// Product is the model for the Product table.
type Product struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price" validate:"required,numeric,gt=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	UnitID      uuid.UUID `json:"unit_id" validate:"required"`
	CreatedAt   string    `json:"created_at" validate:"rfc3339"`
	UpdatedAt   string    `json:"updated_at" validate:"rfc3339"`
	DeletedAt   *string   `json:"deleted_at,omitempty" validate:"omitempty,rfc3339"`
}

// GetProductResponse is the response payload for the GetProduct endpoint.
type GetProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	Category    string    `json:"category,omitempty"`
	Unit        string    `json:"unit,omitempty"`
	CreatedAt   string    `json:"created_at" validate:"rfc3339"`
	UpdatedAt   string    `json:"updated_at" validate:"rfc3339"`
}

// CreateProductRequest is the payload used for creating a new product.
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Description string  `json:"description" validate:"omitempty"`
	Price       float64 `json:"price" validate:"required,numeric,gt=0"`
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	UnitID      string  `json:"unit_id" validate:"required,uuid"`
}

// UpdateProductRequest is the payload used for updating an existing product.
type UpdateProductRequest struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price" validate:"required,numeric,gt=0"`
	CategoryID  string    `json:"category_id" validate:"required,uuid"`
	UnitID      string    `json:"unit_id" validate:"required,uuid"`
}

// DeleteProductRequest is the payload used for deleting a product.
type DeleteProductRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}
