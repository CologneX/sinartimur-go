package employee

import (
	"database/sql"
)

type EmployeeRepository interface {
	Create(request CreateEmployeeRequest) error
	Delete(request DeleteEmployeeRequest) error
	Update(request UpdateEmployeeRequest) error
	GetAll(name string) ([]GetEmployeeResponse, error)
	GetByID(id string) (*GetEmployeeResponse, error)
	GetByNIK(nik string) (*GetEmployeeResponse, error)
	GetByPhone(phone string) (*GetEmployeeResponse, error)
}

type employeeRepositoryImpl struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepositoryImpl{db: db}
}

// Create creates a new employee
func (r *employeeRepositoryImpl) Create(request CreateEmployeeRequest) error {
	_, err := r.db.Exec("Insert Into Employee (Name, Position, Hired_Date, Nik, Phone) Values ($1, $2, $3, $4, $5)", request.Name, request.Position, request.HiredDate, request.Nik, request.Phone)
	if err != nil {

		return err
	}
	return nil
}

// Delete soft deletes an employee
func (r *employeeRepositoryImpl) Delete(request DeleteEmployeeRequest) error {
	_, err := r.db.Exec("Update Employee Set Deleted_At = Now() Where Id = $1", request.ID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an employee
func (r *employeeRepositoryImpl) Update(request UpdateEmployeeRequest) error {
	_, err := r.db.Exec("Update Employee Set Name = $1, Position = $2, Hired_Date = $3, Nik = $4, Phone = $5, Updated_At = Now() Where Id = $6", request.Name, request.Position, request.HiredDate, request.Nik, request.Phone, request.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAll fetches all employees
func (r *employeeRepositoryImpl) GetAll(name string) ([]GetEmployeeResponse, error) {
	var rows *sql.Rows
	var err error

	if name != "" {
		rows, err = r.db.Query("Select Id, Name, Position, Nik, Phone, Hired_Date, Created_At, Updated_At From Employee Where Deleted_At Is Null And Name Ilike $1", "%"+name+"%")
	} else {
		rows, err = r.db.Query("Select Id, Name, Position, Nik, Phone, Hired_Date, Created_At, Updated_At From Employee Where Deleted_At Is Null")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []GetEmployeeResponse
	for rows.Next() {
		var employee GetEmployeeResponse
		err = rows.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Nik, &employee.Phone, &employee.HiredDate, &employee.CreatedAt, &employee.UpdatedAt)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

// GetByID fetches an employee by ID
func (r *employeeRepositoryImpl) GetByID(id string) (*GetEmployeeResponse, error) {
	var employee GetEmployeeResponse
	err := r.db.QueryRow("Select Id, Name, Position, Nik, Phone, Hired_Date, Created_At, Updated_At From Employee Where Id = $1 And Deleted_At Is Null", id).Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Nik, &employee.Phone, &employee.HiredDate, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByNIK fetches an employee by NIK
func (r *employeeRepositoryImpl) GetByNIK(nik string) (*GetEmployeeResponse, error) {
	var employee GetEmployeeResponse
	err := r.db.QueryRow("Select Id, Name, Position, Nik, Phone, Hired_Date, Created_At, Updated_At From Employee Where Nik = $1 And Deleted_At Is Null", nik).Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Nik, &employee.Phone, &employee.HiredDate, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByPhone fetches an employee by phone
func (r *employeeRepositoryImpl) GetByPhone(phone string) (*GetEmployeeResponse, error) {
	var employee GetEmployeeResponse
	err := r.db.QueryRow("Select Id, Name, Position, Nik, Phone, Hired_Date, Created_At, Updated_At From Employee Where Phone = $1 And Deleted_At Is Null", phone).Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Nik, &employee.Phone, &employee.HiredDate, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}
