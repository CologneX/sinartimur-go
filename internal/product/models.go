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
	Category    string    `json:"category,omitempty"`
	CategoryID  uuid.UUID `json:"category_id,omitempty"`
	Unit        string    `json:"unit,omitempty"`
	UnitID      uuid.UUID `json:"unit_id,omitempty"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

// CreateProductRequest is the payload used for creating a new product.
type CreateProductRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty"`
	CategoryID  string `json:"category_id" validate:"required,uuid"`
	UnitID      string `json:"unit_id" validate:"required,uuid"`
}

// UpdateProductRequest is the payload used for updating an existing product.
type UpdateProductRequest struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description string    `json:"description,omitempty"`
	CategoryID  string    `json:"category_id" validate:"required,uuid"`
	UnitID      string    `json:"unit_id" validate:"required,uuid"`
}

// DeleteProductRequest is the payload used for deleting a product.
type DeleteProductRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

// GetProductBatchesRequest is the request for fetching product batches
type GetProductBatchesRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	utils.PaginationParameter
}

// ProductBatchResponse represents a batch of a product with storage info
type ProductBatchResponse struct {
	BatchID         uuid.UUID               `json:"batch_id"`
	SKU             string                  `json:"sku"`
	PurchaseOrderID uuid.UUID               `json:"purchase_order_id"`
	InitialQuantity float64                 `json:"initial_quantity"`
	CurrentQuantity float64                 `json:"current_quantity"`
	UnitPrice       float64                 `json:"unit_price"`
	CreatedAt       string                  `json:"created_at"`
	StorageDetails  []ProductBatchInStorage `json:"storage_details,omitempty"`
}

// ProductBatchInStorage represents quantity of a batch in a specific storage
type ProductBatchInStorage struct {
	StorageID   uuid.UUID `json:"storage_id"`
	StorageName string    `json:"storage_name"`
	Quantity    float64   `json:"quantity"`
}
