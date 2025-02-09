package category

import "sinartimur-go/pkg/dto"

// CategoryService is a service that handles category
type CategoryService struct {
	repo CategoryRepository
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAllCategory fetches all categories
func (s *CategoryService) GetAllCategory(req GetCategoryRequest) ([]GetCategoryResponse, *dto.APIError) {
	categories, err := s.repo.GetAll(req)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return categories, nil
}

// DeleteCategory soft deletes a category
func (s *CategoryService) DeleteCategory(request DeleteCategoryRequest) *dto.APIError {
	// Check if category exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Category tidak ditemukan",
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

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(request CreateCategoryRequest) (*GetCategoryResponse, *dto.APIError) {
	// Check if category name is already used
	_, err := s.repo.GetByName(request.Name)
	if err == nil {
		return nil, &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"general": "Kategori sudah ada",
			},
		}
	}

	category, err := s.repo.Create(request)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return category, nil
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(request UpdateCategoryRequest) (*GetCategoryResponse, *dto.APIError) {
	// Check if category exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Kategori tidak ditemukan",
			},
		}
	}

	// Check if category name is already used
	cat, err := s.repo.GetByName(request.Name)
	if err == nil && cat.ID.String() != request.ID.String() {
		return nil, &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"general": "Kategori sudah ada",
			},
		}
	}

	category, err := s.repo.Update(request)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return category, nil
}
