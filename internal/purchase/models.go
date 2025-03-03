package purchase

import "sinartimur-go/utils"

// GetPurchaseOrderItemResponse represents the response for a purchase order item retrieval
type GetPurchaseOrderItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

// GetPurchaseOrderResponse represents the response for a purchase order retrieval
type GetPurchaseOrderResponse struct {
	ID           string  `json:"id"`
	CreatedBy    string  `json:"created_by"`
	SupplierID   string  `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	OrderDate    string  `json:"order_date"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// GetPurchaseOrderRequest represents the request for filtering purchase orders
type GetPurchaseOrderRequest struct {
	SupplierName string `json:"supplier_name,omitempty" validate:"omitempty"`
	OrderDate    string `json:"order_date,omitempty" validate:"omitempty,rfc3339"`
	Status       string `json:"status,omitempty" validate:"omitempty,oneof=pending approved completed cancelled"`
	FromDate     string `json:"from_date,omitempty" validate:"omitempty,rfc3339"`
	ToDate       string `json:"to_date,omitempty" validate:"omitempty,rfc3339"`
	utils.PaginationParameter
}

// CreatePurchaseOrderRequest represents the request for creating a purchase order
type CreatePurchaseOrderRequest struct {
	SupplierID string                           `json:"supplier_id" validate:"required,uuid"`
	OrderDate  string                           `json:"order_date" validate:"required,rfc3339"`
	Status     string                           `json:"status" validate:"required,oneof=pending approved completed cancelled"`
	Items      []CreatePurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// CreatePurchaseOrderItemRequest represents the request for creating a purchase order item
type CreatePurchaseOrderItemRequest struct {
	ProductID   string  `json:"product_id" validate:"required,uuid"`
	Quantity    int     `json:"quantity" validate:"required,gt=0"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Description string  `json:"description,omitempty" validate:"omitempty"`
}

// UpdatePurchaseOrderRequest represents the request for updating a purchase order
type UpdatePurchaseOrderRequest struct {
	ID          string  `json:"id" validate:"required,uuid"`
	SupplierID  string  `json:"supplier_id,omitempty" validate:"omitempty,uuid"`
	OrderDate   string  `json:"order_date,omitempty" validate:"omitempty,rfc3339"`
	Status      string  `json:"status,omitempty" validate:"omitempty,oneof=pending approved completed cancelled"`
	TotalAmount float64 `json:"total_amount,omitempty" validate:"omitempty,gt=0"`
}

// UpdatePurchaseOrderItemRequest represents the request for updating a purchase order item
type UpdatePurchaseOrderItemRequest struct {
	ID          string  `json:"id"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

// DeletePurchaseOrderRequest represents the request for deleting a purchase order
type DeletePurchaseOrderRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

// DeletePurchaseOrderItemRequest represents the request for deleting a purchase order item
type DeletePurchaseOrderItemRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

// UpdateReceivePurchaseOrderRequest represents the request for receiving items from a purchase order
type UpdateReceivePurchaseOrderRequest struct {
	ID               string `json:"id" validate:"required,uuid"`
	ReceivedQuantity string `json:"received_quantity" validate:"required,gt=0"`
	Description      string `json:"description,omitempty" validate:"omitempty"`
}

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

// CancelPurchaseOrderRequest represents the request for canceling a purchase order
type CancelPurchaseOrderRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

// PurchaseOrderDetailResponse represents the full purchase order with its items
type PurchaseOrderDetailResponse struct {
	GetPurchaseOrderResponse
	Items []GetPurchaseOrderItemResponse `json:"items"`
}
