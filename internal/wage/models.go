package wage

import "github.com/google/uuid"

type Wage struct {
	ID          uuid.UUID `json:"id"`
	EmployeeId  uuid.UUID `json:"employee_id"`
	TotalAmount float64   `json:"total_amount"`
	Month       int       `json:"month"`
	Year        int       `json:"year"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	DeletedAt   string    `json:"deleted_at"`
}

type WageDetail struct {
	ID            uuid.UUID `json:"id"`
	WageId        uuid.UUID `json:"wage_id"`
	ComponentName string    `json:"component_name"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	DeletedAt     string    `json:"deleted_at"`
}

type GetWageResponse struct {
	ID           uuid.UUID `json:"id"`
	EmployeeId   uuid.UUID `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	TotalAmount  float64   `json:"total_amount"`
	Month        int       `json:"month"`
	Year         int       `json:"year"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type GetWageDetail struct {
	ID            uuid.UUID `json:"id"`
	ComponentName string    `json:"component_name"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

type GetWageDetailResponse struct {
	ID           uuid.UUID        `json:"id"`
	EmployeeId   uuid.UUID        `json:"employee_id"`
	EmployeeName string           `json:"employee_name"`
	TotalAmount  float64          `json:"total_amount"`
	Month        int              `json:"month"`
	Year         int              `json:"year"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
	Detail       []*GetWageDetail `json:"detail"`
}

type CreateWageRequest struct {
	EmployeeId uuid.UUID                 `json:"employee_id" validate:"required"`
	Month      int                       `json:"month" validate:"required,numeric"`
	Year       int                       `json:"year" validate:"required,numeric,len=4"`
	WageDetail []CreateWageDetailRequest `json:"wage_detail" validate:"required,dive"`
}

type CreateWageDetailRequest struct {
	ComponentName string  `json:"component_name" validate:"required,len=100"`
	Description   string  `json:"description" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,numeric,gt=0"`
}

type UpdateWageDetailRequest struct {
	ID         uuid.UUID                 `json:"id"`
	WageDetail []CreateWageDetailRequest `json:"wage_detail"`
}

type DeleteWageRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteWageDetailRequest struct {
	ID uuid.UUID `json:"id"`
}
