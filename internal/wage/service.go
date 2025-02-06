package wage

import (
	"sinartimur-go/pkg/dto"
)

// WageService is a service that handles user authentication
type WageService struct {
	repo WageRepository
}

// NewWageService creates a new instance of AuthService
func NewWageService(repo WageRepository) *WageService {
	return &WageService{repo: repo}
}

// GetAllWages fetches all wages
func (s *WageService) GetAllWages(req GetWageRequest) ([]GetWageResponse, int, *dto.APIError) {
	wages, totalItems, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return wages, totalItems, nil
}

// DeleteWage soft deletes a wage
func (s *WageService) DeleteWage(request DeleteWageRequest) *dto.APIError {
	// Check if wage exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Gaji tidak ditemukan",
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

// GetWageDetail fetches wage details
func (s *WageService) GetWageDetail(wageID string) (*GetWageDetailResponse, *dto.APIError) {
	// Get wage by ID
	wage, err := s.repo.GetByID(wageID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Gaji tidak ditemukan",
			},
		}
	}

	detail, err := s.repo.GetWageDetailByWageID(wageID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	wageDetails := GetWageDetailResponse{
		ID:           wage.ID,
		EmployeeId:   wage.EmployeeId,
		EmployeeName: wage.EmployeeName,
		TotalAmount:  wage.TotalAmount,
		Month:        wage.Month,
		Year:         wage.Year,
		CreatedAt:    wage.CreatedAt,
		UpdatedAt:    wage.UpdatedAt,
		Detail:       detail,
	}

	return &wageDetails, nil
}

// CreateWage creates a new wage
func (s *WageService) CreateWage(request CreateWageRequest) *dto.APIError {
	// Check if employee exists
	_, err := s.repo.GetEmployeeByID(request.EmployeeId.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"employee_id": "Karyawan tidak ditemukan",
			},
		}
	}

	// Check if the employee already has a wage for the specified month and year
	existingWage, err := s.repo.GetByEmployeeIdAndMonthYear(request.EmployeeId.String(), request.Month, request.Year)
	if err == nil && existingWage != nil {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"date": "Karyawan sudah memiliki gaji untuk bulan dan tahun ini",
			},
		}
	}

	// Create the wage record
	err = s.repo.Create(request)
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

// UpdateWage updates a wage
func (s *WageService) UpdateWage(request UpdateWageDetailRequest) *dto.APIError {
	// Check if wage exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Gaji tidak ditemukan",
			},
		}
	}

	// Update wage detail
	err = s.repo.UpdateDetail(request)
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

//// UpdateEmployee updates an employee
//func (s *WageService) UpdateEmployee(request UpdateEmployeeRequest) *dto.APIError {
//	// Check if employee exists
//	_, err := s.repo.GetByID(request.ID.String())
//	if err != nil {
//		return &dto.APIError{
//			StatusCode: 404,
//			Details: map[string]string{
//				"general": "Employee tidak ditemukan",
//			},
//		}
//	}
//
//	// Check if NIK is already used by another employee
//	existingEmployee, err := s.repo.GetByNIK(request.Nik)
//	if err == nil && existingEmployee.ID != request.ID {
//		return &dto.APIError{
//			StatusCode: 409,
//			Details: map[string]string{
//				"nik": "NIK sudah terdaftar",
//			},
//		}
//	}
//
//	// Check if Phone is already used by another employee
//	existingEmployee, err = s.repo.GetByPhone(request.Phone)
//	if err == nil && existingEmployee.ID != request.ID {
//		return &dto.APIError{
//			StatusCode: 409,
//			Details: map[string]string{
//				"phone": "Nomor telepon sudah terdaftar",
//			},
//		}
//	}
//
//	err = s.repo.UpdateDetail(request)
//	if err != nil {
//		return &dto.APIError{
//			StatusCode: 500,
//			Details: map[string]string{
//				"general": "Kesalahan Server",
//			},
//		}
//	}
//
//	return nil
//}
//
//// DeleteEmployee soft deletes an employee
//func (s *WageService) DeleteEmployee(request DeleteEmployeeRequest) *dto.APIError {
//	// Check if employee exists
//	_, err := s.repo.GetByID(request.ID.String())
//	if err != nil {
//		return &dto.APIError{
//			StatusCode: 404,
//			Details: map[string]string{
//				"general": "Employee tidak ditemukan",
//			},
//		}
//	}
//	err = s.repo.Delete(request)
//	if err != nil {
//		return &dto.APIError{
//			StatusCode: 500,
//			Details: map[string]string{
//				"general": "Kesalahan Server",
//			},
//		}
//	}
//	return nil
//}
//
//// GetAllEmployees fetches all employees
//func (s *WageService) GetAllEmployees(name string) ([]GetEmployeeResponse, *dto.APIError) {
//	employees, err := s.repo.GetAll(name)
//	if err != nil {
//		return nil, &dto.APIError{
//			StatusCode: 500,
//			Details: map[string]string{
//				"general": "Kesalahan Server",
//			},
//		}
//	}
//	return employees, nil
//}
