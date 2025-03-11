package customer

import (
	"sinartimur-go/pkg/dto"
)

// CustomerService is the service for the Customer domain
type CustomerService struct {
	repo CustomerRepository
}

// NewCustomerService creates a new instance of CustomerService
func NewCustomerService(repo CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// CreateCustomer creates a new customer
func (s *CustomerService) CreateCustomer(request CreateCustomerRequest) *dto.APIError {
	// Check if customer with the same name already exists
	_, err := s.repo.GetByName(request.Name)
	if err == nil {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"name": "Nama pelanggan sudah terdaftar",
			},
		}
	}

	// Create the customer record
	err = s.repo.Create(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Gagal membuat data pelanggan",
			},
		}
	}
	return nil
}

// UpdateCustomer updates an existing customer
func (s *CustomerService) UpdateCustomer(request UpdateCustomerRequest) *dto.APIError {
	// Check if customer exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Pelanggan tidak ditemukan",
			},
		}
	}

	// Check if name is already used by another customer
	existingCustomer, err := s.repo.GetByName(request.Name)
	if err == nil && existingCustomer.ID != request.ID.String() {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"name": "Nama pelanggan sudah terdaftar",
			},
		}
	}

	err = s.repo.Update(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Gagal memperbarui data pelanggan",
			},
		}
	}

	return nil
}

// DeleteCustomer soft deletes a customer
func (s *CustomerService) DeleteCustomer(request DeleteCustomerRequest) *dto.APIError {
	// Check if customer exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Pelanggan tidak ditemukan",
			},
		}
	}

	err = s.repo.Delete(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Gagal menghapus data pelanggan",
			},
		}
	}

	return nil
}

// GetAllCustomers fetches all customers with pagination
func (s *CustomerService) GetAllCustomers(request GetCustomerRequest) ([]GetCustomerResponse, int, *dto.APIError) {
	customers, totalItems, err := s.repo.GetAll(request)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Gagal mengambil daftar pelanggan",
			},
		}
	}

	return customers, totalItems, nil
}

// GetCustomerByID fetches a customer by ID
func (s *CustomerService) GetCustomerByID(id string) (*GetCustomerResponse, *dto.APIError) {
	customer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Pelanggan tidak ditemukan",
			},
		}
	}

	return customer, nil
}

// GetCustomerByName fetches a customer by name
func (s *CustomerService) GetCustomerByName(name string) (*GetCustomerResponse, *dto.APIError) {
	customer, err := s.repo.GetByName(name)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Pelanggan tidak ditemukan",
			},
		}
	}

	return customer, nil
}
