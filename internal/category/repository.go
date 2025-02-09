package category

import (
	"database/sql"
)

type CategoryRepository interface {
	GetAll(req GetCategoryRequest) ([]GetCategoryResponse, error)
	GetByID(id string) (*GetCategoryResponse, error)
	GetByName(name string) (*GetCategoryResponse, error)
	Create(req CreateCategoryRequest) (*GetCategoryResponse, error)
	Update(req UpdateCategoryRequest) (*GetCategoryResponse, error)
	Delete(req DeleteCategoryRequest) error
}

type CategoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

// GetAll fetches all categories
func (r *CategoryRepositoryImpl) GetAll(req GetCategoryRequest) ([]GetCategoryResponse, error) {
	var categories []GetCategoryResponse
	rows, err := r.db.Query("SELECT id, name, description, created_at, updated_at FROM category WHERE deleted_at is null AND name ILIKE $1", "%"+req.Name+"%")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category GetCategoryResponse
		err = rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetByID fetches a category by ID
func (r *CategoryRepositoryImpl) GetByID(id string) (*GetCategoryResponse, error) {
	var category GetCategoryResponse
	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM category WHERE id = $1 AND deleted_at is null", id).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetByName fetches a category by name
func (r *CategoryRepositoryImpl) GetByName(name string) (*GetCategoryResponse, error) {
	var category GetCategoryResponse
	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM category WHERE name = $1 AND deleted_at is null", name).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Create creates a new category
func (r *CategoryRepositoryImpl) Create(req CreateCategoryRequest) (*GetCategoryResponse, error) {
	var category GetCategoryResponse
	err := r.db.QueryRow("INSERT INTO category (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at", req.Name, req.Description).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update updates an existing category
func (r *CategoryRepositoryImpl) Update(req UpdateCategoryRequest) (*GetCategoryResponse, error) {
	var category GetCategoryResponse
	err := r.db.QueryRow("UPDATE category SET name = $1, description = $2, updated_at = now() WHERE id = $3 RETURNING id, name, description, created_at, updated_at", req.Name, req.Description, req.ID).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Delete deletes a category
func (r *CategoryRepositoryImpl) Delete(req DeleteCategoryRequest) error {
	_, err := r.db.Exec("UPDATE category SET deleted_at = now() WHERE id = $1", req.ID)
	if err != nil {
		return err
	}
	return nil
}
