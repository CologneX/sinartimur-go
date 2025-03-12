package employee

import (
	"database/sql"
	"fmt"
	"sinartimur-go/utils"
)

type EmployeeRepository interface {
	Create(request CreateEmployeeRequest) error
	Delete(request DeleteEmployeeRequest) error
	Update(request UpdateEmployeeRequest) error
	GetAll(req GetAllEmployeeRequest) ([]GetEmployeeResponse, int, error)
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
func (r *employeeRepositoryImpl) GetAll(req GetAllEmployeeRequest) ([]GetEmployeeResponse, int, error) {
	// Build the base query
	queryBuilder := utils.NewQueryBuilder(`
		SELECT Id, Name, Position, Nik, Phone, Hired_Date, Created_At, Updated_At 
		FROM Employee 
		WHERE Deleted_At IS NULL
	`)

	// Add filters based on request parameters
	if req.Name != "" {
		queryBuilder.AddFilter("Name ILIKE", "%"+req.Name+"%")
	}

	//if req.Position != "" {
	//	queryBuilder.AddFilter("Position ILIKE", "%"+req.Position+"%")
	//}

	// Build count query to get total items
	countQuery, countParams := queryBuilder.Build()
	countQuery = fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", countQuery)

	// Execute count query
	var totalItems int
	err := r.db.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total karyawan: %w", err)
	}

	// Add sorting if provided
	if req.SortBy != "" {
		direction := "ASC"
		if req.SortOrder == "desc" {
			direction = "DESC"
		}
		queryBuilder.Query.WriteString(fmt.Sprintf(" ORDER BY %s %s", req.SortBy, direction))
	}

	// Add pagination
	queryBuilder.AddPagination(req.PageSize, req.Page)

	// Execute final query
	query, params := queryBuilder.Build()
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data karyawan: %w", err)
	}
	defer rows.Close()

	var employees []GetEmployeeResponse
	for rows.Next() {
		var employee GetEmployeeResponse
		err = rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Position,
			&employee.Nik,
			&employee.Phone,
			&employee.HiredDate,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal membaca data karyawan: %w", err)
		}
		employees = append(employees, employee)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("terjadi kesalahan saat membaca data karyawan: %w", err)
	}

	return employees, totalItems, nil
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
