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

type GetWageRequest struct {
	EmployeeId string `json:"employee_id" validate:"omitempty,uuid"`
	Month      int    `json:"month" validate:"omitempty,numeric"`
	Year       int    `json:"year" validate:"omitempty,numeric,len=4"`
	Page       int    `json:"page" validate:"omitempty,numeric,min=1"`
	PageSize   int    `json:"page_size" validate:"omitempty,numeric,min=1"`
	SortBy     string `json:"sort_by" validate:"omitempty,oneof=id employee_id total_amount month year created_at updated_at"`
	SortOrder  string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}
type WageDetailRequest struct {
	ComponentName string  `json:"component_name" validate:"required"`
	Description   string  `json:"description"`
	Amount        float64 `json:"amount" validate:"required,numeric,gt=0"`
}

type CreateWageRequest struct {
	EmployeeId uuid.UUID           `json:"employee_id" validate:"required"`
	Month      int                 `json:"month" validate:"required,numeric,min=1,max=12"`
	Year       int                 `json:"year" validate:"required,numeric,min=1000,max=9999"`
	WageDetail []WageDetailRequest `json:"wage_detail" validate:"required,min=1,dive,required"`
}

type UpdateWageDetailRequest struct {
	ID         uuid.UUID           `json:"id"`
	WageDetail []WageDetailRequest `json:"wage_detail" validate:"required,min=1,dive,required"`
}

type DeleteWageRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteWageDetailRequest struct {
	ID uuid.UUID `json:"id"`
}
