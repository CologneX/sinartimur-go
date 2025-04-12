package purchase

import (
	"sinartimur-go/pkg/dto"
)

// SupplierService is the service for the Supplier domain
type SupplierService struct {
	repo SupplierRepository
}

// NewSupplierService creates a new instance of SupplierService
func NewSupplierService(repo SupplierRepository) *SupplierService {
	return &SupplierService{repo: repo}
}

// GetAllSuppliers fetches all suppliers with pagination
func (s *SupplierService) GetAllSuppliers(req GetSupplierRequest) ([]GetSupplierResponse, int, *dto.APIError) {
	suppliers, totalItems, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return suppliers, totalItems, nil
}

// GetSupplierByID fetches a supplier by ID
func (s *SupplierService) GetSupplierByID(id string) (*GetSupplierResponse, *dto.APIError) {
	supplier, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Supplier tidak ditemukan",
			},
		}
	}

	return supplier, nil
}

// CreateSupplier creates a new supplier
func (s *SupplierService) CreateSupplier(req CreateSupplierRequest) *dto.APIError {
	// Check if supplier with same name already exists
	existing, err := s.repo.GetByName(req.Name)
	if err == nil && existing != nil {
		return &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"name": "Supplier dengan nama ini sudah terdaftar",
			},
		}
	}

	err = s.repo.Create(req)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}

// UpdateSupplier updates an existing supplier
func (s *SupplierService) UpdateSupplier(req UpdateSupplierRequest) *dto.APIError {
	// Check if supplier exists
	_, err := s.repo.GetByID(req.ID)
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Supplier tidak ditemukan",
			},
		}
	}

	// If name is changing, check if new name is already taken
	if req.Name != "" {
		existing, err := s.repo.GetByName(req.Name)
		if err == nil && existing != nil && existing.ID != req.ID {
			return &dto.APIError{
				StatusCode: 400,
				Details: map[string]string{
					"name": "Supplier dengan nama ini sudah terdaftar",
				},
			}
		}
	}

	err = s.repo.Update(req)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}

// DeleteSupplier deletes a supplier
func (s *SupplierService) DeleteSupplier(id string) *dto.APIError {
	err := s.repo.Delete(id)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}
