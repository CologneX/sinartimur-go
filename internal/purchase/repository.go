package purchase

import (
	"database/sql"
	"fmt"
	"sinartimur-go/internal/category"
	"sinartimur-go/internal/unit"
)

type InventoryRepository interface {
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

func NewProductRepository(db *sql.DB) InventoryRepository {
	return &ProductRepositoryImpl{db: db}
}

// GetAll fetches all product
func (r *ProductRepositoryImpl) GetAll(req GetProductRequest) ([]GetProductResponse, int, error) {
	var products []GetProductResponse
	var totalItems int
	// Query for fetching all products with category and unit name
	query := "Select P.Id, P.Name,P.Description, P.Price, C.Name As Category, U.Name As Unit, P.Created_At, P.Updated_At From Product P Join Category C On P.Category_Id = C.Id Join Unit U On P.Unit_Id = U.Id Where P.Deleted_At Is Null"

	countQuery := "Select Count(Id) From Product Where Deleted_At Is Null"

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
	err := r.db.QueryRow("Select Id, Name, Description, Price, Category_Id, Unit_Id, Created_At, Updated_At From Product Where Id = $1 And Deleted_At Is Null", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.Unit, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetCategoryByID fetches a category by ID
func (r *ProductRepositoryImpl) GetCategoryByID(id string) (*category.GetCategoryResponse, error) {
	var cat category.GetCategoryResponse
	err := r.db.QueryRow("Select Id, Name, Description, Created_At, Updated_At From Category Where Id = $1 And Deleted_At Is Null", id).Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// GetUnitByID fetches a unit by ID
func (r *ProductRepositoryImpl) GetUnitByID(id string) (*unit.GetUnitResponse, error) {
	var un unit.GetUnitResponse
	err := r.db.QueryRow("Select Id, Name, Description, Created_At, Updated_At From Unit Where Id = $1 And Deleted_At Is Null", id).Scan(&un.ID, &un.Name, &un.Description, &un.CreatedAt, &un.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &un, nil
}

// GetByName fetches a product by name
func (r *ProductRepositoryImpl) GetByName(name string) (*GetProductResponse, error) {
	var product GetProductResponse
	err := r.db.QueryRow("Select Id, Name, Description, Created_At, Updated_At From Product Where Name = $1 And Deleted_At Is Null", name).Scan(&product.ID, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
	fmt.Println(err, product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Create inserts a new product
func (r *ProductRepositoryImpl) Create(req CreateProductRequest) (*GetProductResponse, error) {
	var product GetProductResponse
	err := r.db.QueryRow("Insert Into Product (Name, Description ,Price, Category_Id, Unit_Id) Values ($1, $2, $3, $4, $5 ) Returning Id, Name, Description, Created_At, Updated_At", req.Name, req.Description, req.Price, req.CategoryID, req.UnitID).Scan(&product.ID, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update updates an existing product
func (r *ProductRepositoryImpl) Update(req UpdateProductRequest) (*GetProductResponse, error) {
	var product GetProductResponse
	err := r.db.QueryRow("Update Product Set Name = $1, Description = $2, Price = $3, Category_Id = $4, Unit_Id = $5, Updated_At = Now() Where Id = $6 Returning Id, Name, Description, Created_At, Updated_At", req.Name, req.Description, req.Price, req.CategoryID, req.UnitID, req.ID).Scan(&product.ID, &product.Name, &product.Description, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Delete marks a product as deleted
func (r *ProductRepositoryImpl) Delete(req DeleteProductRequest) error {
	_, err := r.db.Exec("Update Product Set Deleted_At = Now() Where Id = $1", req.ID)
	if err != nil {
		return err
	}
	return nil
}
