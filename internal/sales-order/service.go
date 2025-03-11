package customer

import (
	"sinartimur-go/pkg/dto"
)

// ProductService is the service for the Product domain.
type ProductService struct {
	repo InventoryRepository
}

// NewProductService creates a new instance of ProductService
func NewProductService(repo InventoryRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetAllProducts fetches all products
func (s *ProductService) GetAllProducts(search GetProductRequest) ([]GetProductResponse, int, *dto.APIError) {
	products, totalItem, err := s.repo.GetAll(search)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	return products, totalItem, nil
}

// GetProductByID fetches a product by ID
func (s *ProductService) GetProductByID(id string) (*GetProductResponse, *dto.APIError) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Produk tidak ditemukan",
			},
		}
	}
	return product, nil
}

// GetProductByName fetches a product by name
func (s *ProductService) GetProductByName(name string) (*GetProductResponse, *dto.APIError) {
	product, err := s.repo.GetByName(name)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Produk tidak ditemukan",
			},
		}
	}
	return product, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(request CreateProductRequest) (*GetProductResponse, *dto.APIError) {
	// Check if product name is already used
	product, err := s.repo.GetByName(request.Name)
	if err == nil && product != nil {
		return nil, &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"name": "Nama produk sudah digunakan",
			},
		}
	}

	// Check if category exists
	_, err = s.repo.GetCategoryByID(request.CategoryID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"category": "Kategori tidak ditemukan",
			},
		}
	}

	// Check if unit exists
	_, err = s.repo.GetUnitByID(request.UnitID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"unit": "Unit tidak ditemukan",
			},
		}
	}

	product, err = s.repo.Create(request)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	return product, nil
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(request UpdateProductRequest) (*GetProductResponse, *dto.APIError) {
	// Check if product exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Produk tidak ditemukan",
			},
		}
	}

	// Check if product name is already used
	prod, err := s.repo.GetByName(request.Name)
	if err == nil && prod.ID != request.ID {
		return nil, &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"name": "Nama produk sudah digunakan",
			},
		}
	}

	// Check if category exists
	_, err = s.repo.GetCategoryByID(request.CategoryID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"category": "Kategori tidak ditemukan",
			},
		}
	}

	// Check if unit exists
	_, err = s.repo.GetUnitByID(request.UnitID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"unit": "Unit tidak ditemukan",
			},
		}
	}

	product, err := s.repo.Update(request)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return product, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(request DeleteProductRequest) *dto.APIError {
	// Check if product exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Produk tidak ditemukan",
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
