package inventory

import (
	"sinartimur-go/utils"

	"github.com/google/uuid"
)

// Storage represents a storage location in the system
type Storage struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

// GetStorageResponse is used when returning storage data to clients
type GetStorageResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
	ID       string `json:"id" validate:"required,uuid"`
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Location string `json:"location" validate:"required"`
}

// DeleteStorageRequest holds data needed to delete a storage
type DeleteStorageRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

// ProductBatch represents a batch of products with inventory tracking
type ProductBatch struct {
	ID              string  `json:"id"`
	SKU             string  `json:"sku"`
	ProductID       string  `json:"product_id"`
	PurchaseOrderID string  `json:"purchase_order_id"`
	InitialQuantity float64 `json:"initial_quantity"`
	CurrentQuantity float64 `json:"current_quantity"`
	UnitPrice       float64 `json:"unit_price"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// GetAllBatchesRequest holds query parameters for batch search
type GetAllBatchesRequest struct {
	ProductID string `json:"product_id" validate:"omitempty,uuid"`
	SKU       string `json:"sku" validate:"omitempty"`
	StorageID string `json:"storage_id" validate:"omitempty,uuid"`
	utils.PaginationParameter
}

// GetAllBatchResponse is used when returning batch data to clients
type GetAllBatchResponse struct {
	ID              string  `json:"id"`
	SKU             string  `json:"sku"`
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	PurchaseOrderID *string `json:"purchase_order_id,omitempty"`
	InitialQuantity float64 `json:"initial_quantity"`
	CurrentQuantity float64 `json:"current_quantity"`
	UnitPrice       float64 `json:"unit_price"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// BatchStorage represents the quantity of a product batch in a specific storage
type BatchStorage struct {
	ID        string  `json:"id"`
	BatchID   string  `json:"batch_id"`
	StorageID string  `json:"storage_id"`
	Quantity  float64 `json:"quantity"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// MoveBatchRequest holds data needed to move products between storages
type MoveBatchRequest struct {
	BatchID         string  `json:"batch_id" validate:"required,uuid"`
	SourceStorageID string  `json:"source_storage_id" validate:"required,uuid"`
	TargetStorageID string  `json:"target_storage_id" validate:"required,uuid,nefield=SourceStorageID"`
	Quantity        float64 `json:"quantity" validate:"required,gt=0"`
	Description     string  `json:"description" validate:"omitempty"`
}

// InventoryLog represents a record of inventory movement or change
type InventoryLog struct {
	ID              string  `json:"id"`
	BatchID         string  `json:"batch_id"`
	StorageID       string  `json:"storage_id"`
	TargetStorageID *string `json:"target_storage_id,omitempty"`
	UserID          string  `json:"user_id"`
	Action          string  `json:"action"` // add, remove, transfer
	Quantity        float64 `json:"quantity"`
	LogDate         string  `json:"log_date"`
	Description     string  `json:"description"`
	CreatedAt       string  `json:"created_at"`
}

// GetInventoryLogsRequest defines filters for querying inventory logs
type GetInventoryLogsRequest struct {
	BatchID         string `json:"batch_id" validate:"omitempty,uuid"`
	ProductID       string `json:"product_id" validate:"omitempty,uuid"`
	StorageID       string `json:"storage_id" validate:"omitempty,uuid"`
	TargetStorageID string `json:"target_storage_id" validate:"omitempty,uuid"`
	UserID          string `json:"user_id" validate:"omitempty,uuid"`
	Action          string `json:"action" validate:"omitempty,oneof=add remove transfer return"`
	FromDate        string `json:"from_date" validate:"omitempty,rfc3339"`
	ToDate          string `json:"to_date" validate:"omitempty,rfc3339"`
	utils.PaginationParameter
}

// GetInventoryLogResponse represents the data structure for inventory log responses
type GetInventoryLogResponse struct {
	ID                string  `json:"id"`
	BatchID           string  `json:"batch_id"`
	BatchSKU          string  `json:"batch_sku"`
	ProductID         string  `json:"product_id"`
	ProductName       string  `json:"product_name"`
	StorageID         string  `json:"storage_id"`
	StorageName       string  `json:"storage_name"`
	TargetStorageID   *string `json:"target_storage_id,omitempty"`
	TargetStorageName *string `json:"target_storage_name,omitempty"`
	UserID            string  `json:"user_id"`
	Username          string  `json:"username"`
	PurchaseOrderID   *string `json:"purchase_order_id,omitempty"`
	SalesOrderID      *string `json:"sales_order_id,omitempty"`
	Action            string  `json:"action"`
	Quantity          float64 `json:"quantity"`
	LogDate           string  `json:"log_date"`
	Description       string  `json:"description,omitempty"`
	CreatedAt         string  `json:"created_at"`
}
