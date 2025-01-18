package employee

import "errors"

// EmployeeService is a service that handles user authentication
type EmployeeService struct {
	repo EmployeeRepository
}

// NewEmployeeService creates a new instance of AuthService
func NewEmployeeService(repo EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// CreateEmployee registers a new employee
func (s *EmployeeService) CreateEmployee(request CreateEmployeeRequest) error {
	err := s.repo.Create(request)
	if err != nil {
		return err
	}
	return nil
}

// UpdateEmployee updates an employee
func (s *EmployeeService) UpdateEmployee(request UpdateEmployeeRequest) error {
	// Check if employee exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return errors.New("Employee tidak ditemukan")
	}
	err = s.repo.Update(request)
	if err != nil {
		return err
	}
	return nil
}

// DeleteEmployee soft deletes an employee
func (s *EmployeeService) DeleteEmployee(request DeleteEmployeeRequest) error {
	// Check if employee exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return errors.New("Employee tidak ditemukan")

	}
	err = s.repo.Delete(request)
	if err != nil {
		return err
	}
	return nil
}

// GetAllEmployees fetches all employees
func (s *EmployeeService) GetAllEmployees(name string) ([]GetEmployeeResponse, error) {
	employees, err := s.repo.GetAll(name)
	if err != nil {
		return nil, err
	}
	return employees, nil
}
