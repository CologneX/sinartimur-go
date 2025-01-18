package employee

import "database/sql"

type EmployeeRepository interface {
	Create(request CreateEmployeeRequest) error
	Delete(request DeleteEmployeeRequest) error
	Update(request UpdateEmployeeRequest) error
	GetAll(name string) ([]GetEmployeeResponse, error)
	GetByID(id string) (*GetEmployeeResponse, error)
}

type employeeRepositoryImpl struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepositoryImpl{db: db}
}

// Create creates a new employee
func (r *employeeRepositoryImpl) Create(request CreateEmployeeRequest) error {
	_, err := r.db.Exec("INSERT INTO employees (name, position, hired_date) VALUES ($1, $2, $3)", request.Name, request.Position, request.HiredDate)
	if err != nil {
		return err
	}
	return nil
}

// Delete soft deletes an employee
func (r *employeeRepositoryImpl) Delete(request DeleteEmployeeRequest) error {
	_, err := r.db.Exec("UPDATE employees SET deleted_at = NOW() WHERE id = $1", request.ID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates an employee
func (r *employeeRepositoryImpl) Update(request UpdateEmployeeRequest) error {
	_, err := r.db.Exec("UPDATE employees SET name = $1, position = $2, hired_date = $3, updated_at = NOW() WHERE id = $4", request.Name, request.Position, request.HiredDate, request.ID)
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
		rows, err = r.db.Query("SELECT id, name, position, hired_date, created_at, updated_at FROM employees WHERE deleted_at IS NULL AND name ILIKE $1", "%"+name+"%")
	} else {
		rows, err = r.db.Query("SELECT id, name, position, hired_date, created_at, updated_at FROM employees WHERE deleted_at IS NULL")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []GetEmployeeResponse
	for rows.Next() {
		var employee GetEmployeeResponse
		err := rows.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.HiredDate, &employee.CreatedAt, &employee.UpdatedAt)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

// GetByID fetches an employee by ID
func (r *employeeRepositoryImpl) GetByID(id string) (*GetEmployeeResponse, error) {
	row := r.db.QueryRow("SELECT id, name, position, hired_date, created_at, updated_at FROM employees WHERE id = $1 AND deleted_at IS NULL", id)
	var employee GetEmployeeResponse
	err := row.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.HiredDate, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}
