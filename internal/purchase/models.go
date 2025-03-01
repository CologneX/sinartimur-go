package purchase

import "sinartimur-go/utils"

// PurchaseOrderItem
type PurchaseOrderItem struct {
	ID              string  `json:"id"`
	PurchaseOrderID string  `json:"purchase_order_id"`
	ProductID       string  `json:"product_id"`
	Quantity        float64 `json:"quantity"`
	Price           float64 `json:"price"`
	Subtotal        float64 `json:"subtotal"`
	Description     string  `json:"description"`
}

// GetPurchaseOrderItemResponse
type GetPurchaseOrderItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
	Description string  `json:"description"`
}

// GetPurchaseOrderResponse
type GetPurchaseOrderResponse struct {
	ID           string `json:"id"`
	IssuedBy     string `json:"issued_by"`
	SupplierID   string `json:"supplier_id"`
	SupplierName string `json:"supplier_name"`
	Amount       string `json:"amount"`
	CreatedAt    string `json:"created_at"`
	CancelledAt  string `json:"cancelled_at"`
	utils.PaginationParameter
}

// UpdatePurchaseOrderItemRequest
type UpdatePurchaseOrderItemRequest = GetPurchaseOrderItemResponse

// DeletePurchaseOrderRequest
type DeletePurchaseOrderRequest struct {
	ID string `json:"id"`
}

// UpdateReceivePurchaseOrderRequest
type UpdateReceivePurchaseOrderRequest struct {
	ID               string `json:"id" validate:"uuid"`
	ReceivedQuantity string `json:"received_quantity" validate:"gt=0"`
	Description      string `json:"description,omitempty" validate:"omitempty"`
}
