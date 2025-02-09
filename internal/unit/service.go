package unit

import (
	"sinartimur-go/pkg/dto"
)

// UnitService is a service that handles user authentication
type UnitService struct {
	repo UnitRepository
}

// NewUnitService creates a new instance of AuthService
func NewUnitService(repo UnitRepository) *UnitService {
	return &UnitService{repo: repo}
}

// GetAllUnit fetches all units
func (s *UnitService) GetAllUnit(req GetUnitRequest) ([]GetUnitResponse, *dto.APIError) {
	units, err := s.repo.GetAll(req)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return units, nil

}

// DeleteUnit soft deletes a unit
func (s *UnitService) DeleteUnit(request DeleteUnitRequest) *dto.APIError {
	// Check if unit exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Unit tidak ditemukan",
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

// CreateUnit creates a new unit
func (s *UnitService) CreateUnit(request CreateUnitRequest) (*GetUnitResponse, *dto.APIError) {
	// Check if unit name is already used
	_, err := s.repo.GetByName(request.Name)
	if err == nil {
		return nil, &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"name": "Unit sudah terdaftar",
			},
		}
	}

	unit, err := s.repo.Create(request)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return unit, nil
}

// UpdateUnit updates a unit
func (s *UnitService) UpdateUnit(request UpdateUnitRequest) (*GetUnitResponse, *dto.APIError) {
	// Check if unit exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Unit tidak ditemukan",
			},
		}
	}

	// Check if unit name is already used
	existingUnit, err := s.repo.GetByName(request.Name)
	if err == nil && existingUnit.ID != request.ID {
		return nil, &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"name": "Unit sudah terdaftar",
			},
		}
	}

	unit, err := s.repo.Update(request)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return unit, nil
}

//// UpdateEmployee updates an employee
//func (s *UnitService) UpdateEmployee(request UpdateEmployeeRequest) *dto.APIError {
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
//func (s *UnitService) DeleteEmployee(request DeleteEmployeeRequest) *dto.APIError {
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
//func (s *UnitService) GetAllEmployees(name string) ([]GetEmployeeResponse, *dto.APIError) {
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
