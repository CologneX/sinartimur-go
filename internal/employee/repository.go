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
	GetAttendance(req GetAttendanceRequest) ([]GetAttendanceResponse, error)
	UpdateAttendance(req UpdateAttendanceRequest) error
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

// GetAttendance retrieves attendance records for employees on a specific date
func (r *employeeRepositoryImpl) GetAttendance(req GetAttendanceRequest) ([]GetAttendanceResponse, error) {
	queryBuilder := utils.NewQueryBuilder(`
        SELECT e.Id, e.Name, a.Attendance_Date, a.Status, a.Description
        FROM Employee e
        LEFT JOIN Attendance a ON e.Id = a.Employee_Id AND a.Attendance_Date = $1
        WHERE e.Deleted_At IS NULL
		AND e.Hired_Date <= $1
    `)

	// Add the parameter for Attendance_Date
	queryBuilder.Params = append(queryBuilder.Params, req.AttendanceDate)

	// Execute the query
	query, params := queryBuilder.Build()
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attendance records: %w", err)
	}
	defer rows.Close()

	var attendances []GetAttendanceResponse
	for rows.Next() {
		var attendance GetAttendanceResponse
		err = rows.Scan(&attendance.EmployeeID, &attendance.EmployeeName, &attendance.AttendanceDate, &attendance.AttendanceStatus, &attendance.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance record: %w", err)
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendance records: %w", err)
	}

	return attendances, nil
}

// UpdateAttendance updates the attendance record for an employee
func (r *employeeRepositoryImpl) UpdateAttendance(req UpdateAttendanceRequest) error {
    query := `
        INSERT INTO Attendance (Employee_Id, Attendance_Date, Status, Description, Created_At, Updated_At)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        ON CONFLICT (Employee_Id, Attendance_Date)
        DO UPDATE SET Status = $3, Description = $4, Updated_At = NOW()
    `

    _, err := r.db.Exec(query, req.EmployeeID, req.AttendanceDate, req.AttendanceStatus, req.Description)
    if err != nil {
        return fmt.Errorf("failed to update attendance record: %w", err)
    }

    return nil
}

// GetEmployeeByID retrieves an employee by their ID
func (r *employeeRepositoryImpl) GetEmployeeByID(employeeID string) (*Employee, error) {
	query := `
		SELECT Id, Name, Position, Phone, Nik, Hired_Date, Created_At, Updated_At, Deleted_At
		FROM Employee
		WHERE Id = $1 AND Deleted_At IS NULL
	`

	row := r.db.QueryRow(query, employeeID)

	var employee Employee
	if err := row.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Phone, &employee.Nik, &employee.HiredDate, &employee.CreatedAt, &employee.UpdatedAt, &employee.DeletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Employee not found
		}
		return nil, fmt.Errorf("failed to fetch employee: %w", err)
	}

	return &employee, nil
}
