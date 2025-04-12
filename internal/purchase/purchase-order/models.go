package purchase_order

import (
	"sinartimur-go/utils"
)

// Purchase Order models
type GetPurchaseOrderRequest struct {
	SupplierID string `json:"supplier_id,omitempty" validate:"omitempty,uuid"`
	Status     string `json:"status,omitempty" validate:"omitempty,oneof=ordered received checked completed partially_returned returned cancelled"`
	FromDate   string `json:"from_date,omitempty" validate:"omitempty,rfc3339"`
	ToDate     string `json:"to_date,omitempty" validate:"omitempty,rfc3339"`
	utils.PaginationParameter
}

type GetPurchaseOrderResponse struct {
	ID           string  `json:"id"`
	SupplierID   string  `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	OrderDate    string  `json:"order_date"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
	CreatedBy    string  `json:"created_by"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type StorageAllocationResponse struct {
	StorageID   string  `json:"storage_id"`
	StorageName string  `json:"storage_name"`
	Quantity    float64 `json:"quantity"`
}
type CreatePurchaseOrderRequest struct {
	SupplierID string                           `json:"supplier_id" validate:"required,uuid"`
	OrderDate  string                           `json:"order_date" validate:"required,rfc3339"`
	Status     string                           `json:"status" validate:"required,oneof=ordered received checked completed partially_returned returned cancelled"`
	Items      []CreatePurchaseOrderItemRequest `json:"items" validate:"required,dive"`
}

type CreatePurchaseOrderResponse struct {
	ID           string  `json:"id"`
	SupplierID   string  `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	OrderDate    string  `json:"order_date"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
	CreatedBy    string  `json:"created_by"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type UpdatePurchaseOrderRequest struct {
	ID         string `json:"id" validate:"required,uuid"`
	SupplierID string `json:"supplier_id,omitempty" validate:"omitempty,uuid"`
	OrderDate  string `json:"order_date,omitempty" validate:"omitempty,rfc3339"`
	Status     string `json:"status,omitempty" validate:"omitempty,oneof=ordered received checked completed partially_returned returned cancelled"`
}

type UpdatePurchaseOrderResponse struct {
	ID           string  `json:"id"`
	SupplierID   string  `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	OrderDate    string  `json:"order_date"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
	CreatedBy    string  `json:"created_by"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type GetPurchaseOrderItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

type CreatePurchaseOrderItemRequest struct {
	ProductID string  `json:"product_id" validate:"required,uuid"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type CreatePurchaseOrderItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
	StorageID   string  `json:"storage_id"`
	StorageName string  `json:"storage_name"`
}

type UpdatePurchaseOrderItemRequest struct {
	ID       string  `json:"id" validate:"required,uuid"`
	Quantity float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Price    float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
}

type UpdatePurchaseOrderItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
	StorageID   string  `json:"storage_id"`
	StorageName string  `json:"storage_name"`
}

// Combined response for purchase purchase-order with details
type PurchaseOrderDetailResponse struct {
	GetPurchaseOrderResponse
	Items []GetPurchaseOrderItemResponse `json:"items"`
}

// Purchase Order Return models
type GetPurchaseOrderReturnRequest struct {
	PurchaseOrderID string `json:"purchase_order_id,omitempty" validate:"omitempty,uuid"`
	Status          string `json:"status,omitempty" validate:"omitempty,oneof=pending completed cancelled"`
	FromDate        string `json:"from_date,omitempty" validate:"omitempty,rfc3339"`
	ToDate          string `json:"to_date,omitempty" validate:"omitempty,rfc3339"`
	utils.PaginationParameter
}

type CreatePurchaseOrderReturnRequest struct {
	PurchaseOrderID      string                            `json:"purchase_order_id" validate:"required,uuid"`
	ProductDetailID      string                            `json:"product_detail_id" validate:"required,uuid"`
	ReturnQuantity       float64                           `json:"return_quantity" validate:"required,gt=0"`
	ReturnReason         string                            `json:"return_reason"`
	BatchDetailsToReturn []PurchaseOrderReturnBatchRequest `json:"batch_details,omitempty" validate:"omitempty,dive"`
}

type PurchaseOrderReturnBatchRequest struct {
	BatchID        string  `json:"batch_id" validate:"required,uuid"`
	ReturnQuantity float64 `json:"return_quantity" validate:"required,gt=0"`
}

type GetPurchaseOrderReturnResponse struct {
	ID                string  `json:"id"`
	PurchaseOrderID   string  `json:"purchase_order_id"`
	ProductDetailID   string  `json:"product_detail_id"`
	ReturnQuantity    float64 `json:"return_quantity"`
	RemainingQuantity float64 `json:"remaining_quantity"`
	ReturnReason      string  `json:"return_reason"`
	ReturnStatus      string  `json:"return_status"`
	ReturnedBy        string  `json:"returned_by"`
	ReturnedAt        string  `json:"returned_at"`
	CancelledAt       string  `json:"cancelled_at,omitempty"`
	CancelledBy       string  `json:"cancelled_by,omitempty"`
}

type BatchReturnResponse struct {
	BatchID        string  `json:"batch_id"`
	ReturnQuantity float64 `json:"return_quantity"`
}

type PurchaseOrderReturnDetailResponse struct {
	GetPurchaseOrderReturnResponse
	BatchDetails []BatchReturnResponse `json:"batch_details,omitempty"`
}

type ReceivedItemRequest struct {
	DetailID           string                     `json:"detail_id" validate:"required,uuid"`
	ReceivedQuantity   float64                    `json:"received_quantity" validate:"required,gt=0"`
	UnitPrice          float64                    `json:"unit_price" validate:"required,gt=0"`
	PaymentDueDate     string                     `json:"payment_due_date,omitempty" validate:"omitempty,rfc3339"`
	StorageAllocations []StorageAllocationRequest `json:"storage_allocations" validate:"required,dive"`
}

type StorageAllocationRequest struct {
	StorageID string  `json:"storage_id" validate:"required,uuid"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
}
