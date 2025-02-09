package product

import (
	"database/sql"
	"fmt"
	"sinartimur-go/internal/category"
	"sinartimur-go/internal/unit"
)

type ProductRepository interface {
	GetAll(req GetProductRequest) ([]GetProductResponse, int, error)
	GetByID(id string) (*GetProductResponse, error)
	GetByName(name string) (*GetProductResponse, error)
	Create(req CreateProductRequest) (*GetProductResponse, error)
	Update(req UpdateProductRequest) (*GetProductResponse, error)
	Delete(req DeleteProductRequest) error
	GetCategoryByID(id string) (*category.GetCategoryResponse, error)
	GetUnitByID(id string) (*unit.GetUnitResponse, error)
}

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

// GetAll fetches all product
func (r *ProductRepositoryImpl) GetAll(req GetProductRequest) ([]GetProductResponse, int, error) {
	var products []GetProductResponse
	var totalItems int
	// Query for fetching all products with category and unit name
	query := "SELECT p.id, p.name,p.description, p.price, c.name as category, u.name as unit, p.created_at, p.updated_at FROM products p JOIN category c ON p.category_id = c.id JOIN unit u ON p.unit_id = u.id WHERE p.deleted_at is null"

	countQuery := "SELECT COUNT(id) FROM products WHERE deleted_at is null"

	if req.Name != "" {
		query += " AND p.name ILIKE '%" + req.Name + "%'"
		countQuery += " AND name ILIKE '%" + req.Name + "%'"
	}

	if req.Category != "" {
		query += " AND p.category_id = '" + req.Category + "'"
		countQuery += " AND category_id = '" + req.Category + "'"
	}

	if req.Unit != "" {
		query += " AND p.unit_id = '" + req.Unit + "'"
		countQuery += " AND unit_id = '" + req.Unit + "'"
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}

	countRow := r.db.QueryRow(countQuery)
	err = countRow.Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var product GetProductResponse
		err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Unit, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		products = append(products, product)
	}

	return products, totalItems, nil
}

// GetByID fetches a product by ID
func (r *ProductRepositoryImpl) GetByID(id string) (*GetProductResponse, error) {
	var product GetProductResponse
	err := r.db.QueryRow("SELECT id, name, description, price, category_id, unit_id, created_at, updated_at FROM products WHERE id = $1 AND deleted_at is null", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Unit, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetCategoryByID fetches a category by ID
func (r *ProductRepositoryImpl) GetCategoryByID(id string) (*category.GetCategoryResponse, error) {
	var cat category.GetCategoryResponse
	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM category WHERE id = $1 AND deleted_at is null", id).Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// GetUnitByID fetches a unit by ID
func (r *ProductRepositoryImpl) GetUnitByID(id string) (*unit.GetUnitResponse, error) {
	var un unit.GetUnitResponse
	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM unit WHERE id = $1 AND deleted_at is null", id).Scan(&un.ID, &un.Name, &un.Description, &un.CreatedAt, &un.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &un, nil
}

// GetByName fetches a product by name
func (r *ProductRepositoryImpl) GetByName(name string) (*GetProductResponse, error) {
	var product GetProductResponse
	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM products WHERE name = $1 AND deleted_at is null", name).Scan(&product.ID, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
	fmt.Println(err, product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Create inserts a new product
func (r *ProductRepositoryImpl) Create(req CreateProductRequest) (*GetProductResponse, error) {
	var product GetProductResponse
	err := r.db.QueryRow("INSERT INTO products (name, description ,price, category_id, unit_id) VALUES ($1, $2, $3, $4, $5 ) RETURNING id, name, description, created_at, updated_at", req.Name, req.Description, req.Price, req.CategoryID, req.UnitID).Scan(&product.ID, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update updates an existing product
func (r *ProductRepositoryImpl) Update(req UpdateProductRequest) (*GetProductResponse, error) {
	var product GetProductResponse
	err := r.db.QueryRow("UPDATE products SET name = $1, description = $2, price = $3, category_id = $4, unit_id = $5, updated_at = now() WHERE id = $6 RETURNING id, name, description, created_at, updated_at", req.Name, req.Description, req.Price, req.CategoryID, req.UnitID, req.ID).Scan(&product.ID, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Delete marks a product as deleted
func (r *ProductRepositoryImpl) Delete(req DeleteProductRequest) error {
	_, err := r.db.Exec("UPDATE products SET deleted_at = now() WHERE id = $1", req.ID)
	if err != nil {
		return err
	}
	return nil
}
