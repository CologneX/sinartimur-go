package inventory

import (
	"github.com/google/uuid"
	"sinartimur-go/utils"
	"time"
)

// Storage represents a storage location in the system
type Storage struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Location  string     `json:"location"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// GetStorageResponse is used when returning storage data to clients
type GetStorageResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetStorageRequest holds query parameters for storage search
type GetStorageRequest struct {
	Name     string `json:"name" validate:"omitempty"`
	Location string `json:"location" validate:"omitempty"`
	// Adding pagination fields
	utils.PaginationParameter
}

// CreateStorageRequest holds data needed to create a new storage
type CreateStorageRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Location string `json:"location" validate:"required"`
}

// UpdateStorageRequest holds data needed to update an existing storage
type UpdateStorageRequest struct {
	ID       uuid.UUID `json:"id" validate:"required,uuid"`
	Name     string    `json:"name" validate:"required,min=2,max=255"`
	Location string    `json:"location" validate:"required"`
}

// DeleteStorageRequest holds data needed to delete a storage
type DeleteStorageRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

// ProductBatch represents a batch of products with inventory tracking
type ProductBatch struct {
	ID              string    `json:"id"`
	SKU             string    `json:"sku"`
	ProductID       string    `json:"product_id"`
	PurchaseOrderID string    `json:"purchase_order_id"`
	InitialQuantity float64   `json:"initial_quantity"`
	CurrentQuantity float64   `json:"current_quantity"`
	UnitPrice       float64   `json:"unit_price"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BatchStorage represents the quantity of a product batch in a specific storage
type BatchStorage struct {
	ID        string    `json:"id"`
	BatchID   string    `json:"batch_id"`
	StorageID string    `json:"storage_id"`
	Quantity  float64   `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MoveBatchRequest holds data needed to move products between storages
type MoveBatchRequest struct {
	BatchID         uuid.UUID `json:"batch_id" validate:"required,uuid"`
	SourceStorageID uuid.UUID `json:"source_storage_id" validate:"required,uuid"`
	TargetStorageID uuid.UUID `json:"target_storage_id" validate:"required,uuid,nefield=SourceStorageID"`
	Quantity        float64   `json:"quantity" validate:"required,gt=0"`
	Description     string    `json:"description" validate:"omitempty"`
}

// InventoryLog represents a record of inventory movement or change
type InventoryLog struct {
	ID              string    `json:"id"`
	BatchID         string    `json:"batch_id"`
	StorageID       string    `json:"storage_id"`
	TargetStorageID *string   `json:"target_storage_id,omitempty"`
	UserID          string    `json:"user_id"`
	Action          string    `json:"action"` // add, remove, transfer
	Quantity        float64   `json:"quantity"`
	LogDate         time.Time `json:"log_date"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}
