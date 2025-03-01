package finance

type GetFinanceTransactionResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	UserName        string  `json:"user_name"`
	Amount          float32 `json:"amount"`
	Type            string  `json:"type"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"string"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type GetFinanceTransactionRequest struct {
	UserID          string `json:"user_id,omitempty" validate:"omitempty,uuid"`
	Type            string `json:"type,omitempty" validate:"omitempty"`
	TransactionDate string `json:"transaction_date,omitempty" validate:"omitempty,rfc3339"`
}
