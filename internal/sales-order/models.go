package customer

// GetSalesOrderResponse
type GetSalesOrderResponse struct {
	ID           string  `json:"id"`
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	OrderDate    string  `json:"order_date"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	CancelledAt  string  `json:"cancelled_at"`
}

// GetSalesOrderRequest
type GetSalesOrderRequest struct {
	CustomerID string `json:"customer_id,omitempty" validate:"omitempty,uuid"`
	OrderDate  string `json:"order_date,omitempty" validate:"omitempty,rfc3339"`
	Status     string `json:"status,omitempty" validate:"omitempty"`
	CreatedAt  string `json:"created_at,omitempty" validate:"omitempty,rfc3339"`
}

// GetDeliveryNoteResponse
type GetDeliveryNoteResponse struct {
	ID            string `json:"id"`
	SalesOrderID  string `json:"sales_order_id"`
	DeliveryDate  string `json:"delivery_date"`
	RecipientName string `json:"recipient_name"`
	IsReceived    bool   `json:"is_received"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	CancelledAt   string `json:"cancelled_at"`
}

// GetDeliveryNoteRequest
type GetDeliveryNoteRequest struct {
	DeliveryDate string `json:"delivery_date,omitempty" validate:"omitempty,rfc3339"`
	IsReceived   bool   `json:"is_received,omitempty" validate:"omitempty,bool"`
}

// CreateSalesOrderRequest
type CreateSalesOrderRequest struct {
	CustomerID  string  `json:"customer_id" validate:"uuid"`
	OrderDate   string  `json:"order_date" validate:"rfc3339"`
	TotalAmount float64 `json:"total_amount"`
}

// DeleteSalesOrderRequest
type DeleteSalesOrderRequest struct {
	ID string `json:"id" validate:"uuid"`
}

// CreateDeliveryNoteRequest
type CreateDeliveryNoteRequest struct {
	SalesOrderID string `json:"sales_order_id" validate:"uuid"`
}

// DeleteDeliveryNoteRequest
type DeleteDeliveryNoteRequest struct {
	DeliveryNoteID string `json:"delivery_note_id" validate:"uuid"`
}
