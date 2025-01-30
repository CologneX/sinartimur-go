package employee

import "github.com/google/uuid"

type Employee struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	Phone     string    `json:"phone"`
	Nik       string    `json:"nik"`
	HiredDate string    `json:"hired_date"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	DeletedAt string    `json:"deleted_at"`
}

type GetEmployeeResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Nik       string    `json:"nik"`
	Phone     string    `json:"phone"`
	Position  string    `json:"position"`
	HiredDate string    `json:"hired_date"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type UpdateEmployeeRequest struct {
	Name      string    `json:"name" validate:"required"`
	Position  string    `json:"position" validate:"required"`
	Phone     string    `json:"phone" validate:"required,min=10,max=13"`
	Nik       string    `json:"nik" validate:"required,len=16"`
	HiredDate string    `json:"hired_date" validate:"required,rfc3339"`
	ID        uuid.UUID `json:"id" validate:"required,uuid"`
}

type CreateEmployeeRequest struct {
	Name      string `json:"name" validate:"required"`
	Phone     string `json:"phone" validate:"required,min=10,max=13"`
	Nik       string `json:"nik" validate:"required,len=16"`
	Position  string `json:"position" validate:"required"`
	HiredDate string `json:"hired_date" validate:"required,rfc3339"`
}

type DeleteEmployeeRequest struct {
	ID uuid.UUID `json:"id"`
}
