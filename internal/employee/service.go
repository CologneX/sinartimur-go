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
			Message:    "Kesalahan Server",
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
			Message:    "Employee tidak ditemukan",
		}
	}
	// Check if employee with the same NIK or phone number already exists
	_, err = s.repo.GetByNIK(request.Nik)
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
	err = s.repo.Update(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Message:    "Kesalahan Server",
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
			Message:    "Employee tidak ditemukan",
		}
	}
	err = s.repo.Delete(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Message:    "Kesalahan Server",
		}
	}
	return nil
}

// GetAllEmployees fetches all employees
func (s *EmployeeService) GetAllEmployees(name string) ([]GetEmployeeResponse, *dto.APIError) {
	employees, err := s.repo.GetAll(name)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Message:    "Kesalahan Server",
		}
	}
	return employees, nil
}
