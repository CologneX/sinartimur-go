package finance

import (
	"sinartimur-go/utils"
)

// GetFinanceTransactionResponse defines the response returned to clients
// when fetching finance transaction details
type GetFinanceTransactionResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	Username        string  `json:"username,omitempty"`
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"`
	PurchaseOrderID string  `json:"purchase_order_id,omitempty"`
	SalesOrderID    string  `json:"sales_order_id,omitempty"`
	Description     string  `json:"description"`
	IsSystem        bool    `json:"is_system"`
	TransactionDate string  `json:"transaction_date"`
	CreatedAt       string  `json:"created_at"`
	EditedAt        *string `json:"edited_at,omitempty"`
}

// GetFinanceTransactionRequest defines the parameters for filtering finance transactions
type GetFinanceTransactionRequest struct {
	UserID          string `json:"user_id,omitempty" validate:"omitempty,uuid"`
	Type            string `json:"type,omitempty" validate:"omitempty"`
	PurchaseOrderID string `json:"purchase_order_id,omitempty" validate:"omitempty,uuid"`
	SalesOrderID    string `json:"sales_order_id,omitempty" validate:"omitempty,uuid"`
	IsSystem        *bool  `json:"is_system,omitempty"`
	StartDate       string `json:"start_date,omitempty"`
	EndDate         string `json:"end_date,omitempty"`
	// Add pagination fields
	utils.PaginationParameter
}

// CreateFinanceTransactionRequest defines fields required to create a new finance transaction
type CreateFinanceTransactionRequest struct {
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	Type            string  `json:"type" validate:"required"`
	PurchaseOrderID string  `json:"purchase_order_id,omitempty" validate:"omitempty,uuid"`
	SalesOrderID    string  `json:"sales_order_id,omitempty" validate:"omitempty,uuid"`
	Description     string  `json:"description" validate:"required"`
	TransactionDate string  `json:"transaction_date" validate:"required,rfc3339"`
}

// CancelFinanceTransactionRequest defines fields required to cancel a finance transaction
type CancelFinanceTransactionRequest struct {
	ID          string  `json:"id" validate:"required,uuid"`
	Description *string `json:"description,omitempty" validate:"omitempty"`
}

// FinanceTransactionSummary represents transaction summary statistics
type FinanceTransactionSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetAmount    float64 `json:"net_amount"`
	Period       string  `json:"period,omitempty"`
}
