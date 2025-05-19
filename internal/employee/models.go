package employee

import (
	"sinartimur-go/utils"

	"github.com/google/uuid"
)

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

// enum of attendance status "present", "absent", "late"
type AttendanceStatus string

const (
	Present AttendanceStatus = "present"
	Absent  AttendanceStatus = "absent"
	Late    AttendanceStatus = "late"
)

type EmployeeAttendance struct {
	ID               string           `json:"id"`
	EmployeeID       string           `json:"employee_id"`
	AttendanceDate   string           `json:"attendance_date"`
	AttendanceStatus AttendanceStatus `json:"attendance_status"`
	Description      string           `json:"description"`
	CreatedAt        string           `json:"created_at"`
	UpdatedAt        string           `json:"updated_at"`
}

type GetAttendanceRequest struct {
	// EmployeeID     string `json:"employee_id,omitempty" validate:"uuid,omitempty"`
	AttendanceDate string `json:"attendance_date" validate:"required,rfc3339"`
	// utils.PaginationParameter
}

type UpdateAttendanceRequest struct {
	EmployeeID       string           `json:"employee_id" validate:"required,uuid"`
	AttendanceDate   string           `json:"attendance_date" validate:"required,rfc3339"`
	AttendanceStatus AttendanceStatus `json:"attendance_status" validate:"required,oneof=present absent late"`
	Description      string           `json:"description,omitempty"`
}

type GetAttendanceResponse struct {
	EmployeeID       string            `json:"employee_id"`
	EmployeeName     string            `json:"employee_name"`
	AttendanceDate   *string           `json:"attendance_date"`
	AttendanceStatus *AttendanceStatus `json:"attendance_status"`
	Description      *string           `json:"description,omitempty"`
}

type GetAllAttendanceResponse struct {
	Items []GetAttendanceResponse `json:"items"`
	// utils.PaginationResponse
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

type GetAllEmployeeRequest struct {
	Name string `json:"name,omitempty"`
	utils.PaginationParameter
}
