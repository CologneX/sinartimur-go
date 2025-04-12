package purchase_order

import (
	"database/sql"
	"sinartimur-go/utils"
	"time"
)

// Request types
type CreatePurchaseOrderRequest struct {
	SupplierID     string                           `json:"supplier_id" validate:"required,uuid"`
	OrderDate      string                           `json:"order_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	PaymentMethod  string                           `json:"payment_method" validate:"required,oneof=cash credit"`
	PaymentDueDate string                           `json:"payment_due_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Items          []CreatePurchaseOrderItemRequest `json:"items" validate:"required,dive"`
}

type CreatePurchaseOrderItemRequest struct {
	ProductID string  `json:"product_id" validate:"required,uuid"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type UpdatePurchaseOrderRequest struct {
	ID             string `json:"id" validate:"required,uuid"`
	SupplierID     string `json:"supplier_id" validate:"omitempty,uuid"`
	OrderDate      string `json:"order_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	PaymentMethod  string `json:"payment_method" validate:"omitempty,oneof=cash credit"`
	PaymentDueDate string `json:"payment_due_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	CheckedBy      string `json:"checked_by" validate:"omitempty,uuid"`
}

type UpdatePurchaseOrderItemRequest struct {
	ID       string  `json:"id" validate:"required,uuid"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
	Price    float64 `json:"price" validate:"required,gt=0"`
}

type ReceivedItemRequest struct {
	DetailID  string  `json:"detail_id" validate:"required,uuid"`
	ProductID string  `json:"product_id" validate:"required,uuid"`
	StorageID string  `json:"storage_id" validate:"required,uuid"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" validate:"required,gt=0"`
}

type CreatePurchaseOrderReturnRequest struct {
	PurchaseOrderID string               `json:"purchase_order_id" validate:"required,uuid"`
	ProductDetailID string               `json:"product_detail_id" validate:"required,uuid"`
	ReturnQuantity  float64              `json:"return_quantity" validate:"required,gt=0"`
	Reason          string               `json:"reason"`
	Batches         []ReturnBatchRequest `json:"batches" validate:"required,dive"`
}

type ReturnBatchRequest struct {
	BatchID   string  `json:"batch_id" validate:"required,uuid"`
	StorageID string  `json:"storage_id" validate:"required,uuid"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
}

type GetPurchaseOrderRequest struct {
	SupplierID string `json:"supplier_id" validate:"omitempty,uuid"`
	Status     string `json:"status" validate:"omitempty,oneof=ordered completed partially_returned returned cancelled"`
	FromDate   string `json:"from_date" validate:"omitempty,datetime=2006-01-02"`
	ToDate     string `json:"to_date" validate:"omitempty,datetime=2006-01-02"`
	utils.PaginationParameter
}

type GetPurchaseOrderReturnRequest struct {
	FromDate string `json:"from_date" validate:"omitempty,datetime=2006-01-02"`
	ToDate   string `json:"to_date" validate:"omitempty,datetime=2006-01-02"`
	utils.PaginationParameter
}

// Response types
type PurchaseOrderDetailResponse struct {
	ID             string              `json:"id"`
	SerialID       string              `json:"serial_id"`
	SupplierID     string              `json:"supplier_id"`
	SupplierName   string              `json:"supplier_name"`
	OrderDate      time.Time           `json:"order_date"`
	Status         string              `json:"status"`
	TotalAmount    float64             `json:"total_amount"`
	PaymentMethod  string              `json:"payment_method"`
	PaymentDueDate sql.NullTime        `json:"payment_due_date"`
	CreatedBy      string              `json:"created_by"`
	CreatedByName  string              `json:"created_by_name"`
	CheckedBy      sql.NullString      `json:"checked_by"`
	CheckedByName  sql.NullString      `json:"checked_by_name"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	Items          []PurchaseOrderItem `json:"items"`
}

type PurchaseOrderItem struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    float64   `json:"quantity"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetPurchaseOrderResponse struct {
	ID           string    `json:"id"`
	SerialID     string    `json:"serial_id"`
	SupplierID   string    `json:"supplier_id"`
	SupplierName string    `json:"supplier_name"`
	OrderDate    time.Time `json:"order_date"`
	Status       string    `json:"status"`
	TotalAmount  float64   `json:"total_amount"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ItemCount    int       `json:"item_count"`
}

type GetPurchaseOrderReturnResponse struct {
	ID              string    `json:"id"`
	PurchaseOrderID string    `json:"purchase_order_id"`
	SerialID        string    `json:"serial_id"`
	ProductID       string    `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ReturnQuantity  float64   `json:"return_quantity"`
	Reason          string    `json:"reason"`
	Status          string    `json:"status"`
	ReturnedBy      string    `json:"returned_by"`
	ReturnedByName  string    `json:"returned_by_name"`
	ReturnedAt      time.Time `json:"returned_at"`
}
