package employee

import "github.com/google/uuid"

type Employee struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	HiredDate string    `json:"hired_date"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	DeletedAt string    `json:"deleted_at"`
}

type GetEmployeeResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	HiredDate string    `json:"hired_date"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type UpdateEmployeeRequest struct {
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	HiredDate string    `json:"hired_date"`
	ID        uuid.UUID `json:"id"`
}

type CreateEmployeeRequest struct {
	Name      string `json:"name"`
	Position  string `json:"position"`
	HiredDate string `json:"hired_date"`
}

type DeleteEmployeeRequest struct {
	ID uuid.UUID `json:"id"`
}
