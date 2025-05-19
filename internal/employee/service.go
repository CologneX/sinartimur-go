package employee

import (
	"sinartimur-go/pkg/dto"
)

// EmployeeService is a service that handles user authentication
type EmployeeService struct {
	repo EmployeeRepository
}

// NewEmployeeService creates a new instance of AuthService
func NewEmployeeService(repo EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// CreateEmployee registers a new employee
func (s *EmployeeService) CreateEmployee(request CreateEmployeeRequest) *dto.APIError {
	// Check if employee with the same NIK or phone number already exists
	_, err := s.repo.GetByNIK(request.Nik)
	if err == nil {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"nik": "NIK sudah terdaftar",
			},
		}
	}
	_, err = s.repo.GetByPhone(request.Phone)
	if err == nil {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"phone": "Nomor telepon sudah terdaftar",
			},
		}
	}

	// Create the employee record
	err = s.repo.Create(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Role tidak ditemukan",
			},
		}
	}
	return nil
}

// UpdateEmployee updates an employee
func (s *EmployeeService) UpdateEmployee(request UpdateEmployeeRequest) *dto.APIError {
	// Check if employee exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Employee tidak ditemukan",
			},
		}
	}

	// Check if NIK is already used by another employee
	existingEmployee, err := s.repo.GetByNIK(request.Nik)
	if err == nil && existingEmployee.ID != request.ID {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"nik": "NIK sudah terdaftar",
			},
		}
	}

	// Check if Phone is already used by another employee
	existingEmployee, err = s.repo.GetByPhone(request.Phone)
	if err == nil && existingEmployee.ID != request.ID {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"phone": "Nomor telepon sudah terdaftar",
			},
		}
	}

	err = s.repo.Update(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}

	return nil
}

// DeleteEmployee soft deletes an employee
func (s *EmployeeService) DeleteEmployee(request DeleteEmployeeRequest) *dto.APIError {
	// Check if employee exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Employee tidak ditemukan",
			},
		}
	}
	err = s.repo.Delete(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return nil
}

// GetAllEmployees fetches all employees
func (s *EmployeeService) GetAllEmployees(req GetAllEmployeeRequest) ([]GetEmployeeResponse, int, *dto.APIError) {
	employees, totalItems, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return employees, totalItems, nil
}

// GetAttendanceService retrieves attendance records for employees on a specific date
func (s *EmployeeService) GetAllAttendance(req GetAttendanceRequest) ([]GetAttendanceResponse, *dto.APIError) {
	// Call repository to fetch attendance records
	attendances, err := s.repo.GetAttendance(req)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return attendances, nil
}

// UpdateAttendanceService updates the attendance record for an employee
func (s *EmployeeService) UpdateAttendance(req UpdateAttendanceRequest) *dto.APIError {
	// Check if the employee exists
	employee, err := s.repo.GetByID(req.EmployeeID)
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Gagal mendapatkan data karyawan: " + err.Error(),
			},
		}
	}
	if employee == nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Karyawan tidak ditemukan",
			},
		}
	}

	// Call repository to update attendance
	if err := s.repo.UpdateAttendance(req); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Gagal memperbarui data kehadiran: " + err.Error(),
			},
		}
	}

	return nil
}
