package product

import (
	"database/sql"
	"fmt"
	"sinartimur-go/internal/category"
	"sinartimur-go/internal/unit"
	"strings"
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
	GetProductBatches(req GetProductBatchesRequest) ([]ProductBatchResponse, int, error)
}

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

// GetAll fetches all products
func (r *ProductRepositoryImpl) GetAll(req GetProductRequest) ([]GetProductResponse, int, error) {
	var products []GetProductResponse
	var totalItems int

	// Using strings.Builder for query construction
	var queryBuilder strings.Builder
	var countQueryBuilder strings.Builder

	// Base queries
	queryBuilder.WriteString("Select P.Id, P.Name, P.Description, C.Name As Category, Category_Id, U.Name As Unit,Unit_Id, P.Created_At, P.Updated_At From Product P Join Category C On P.Category_Id = C.Id Join Unit U On P.Unit_Id = U.Id Where P.Deleted_At Is Null")
	countQueryBuilder.WriteString("Select Count(Id) From Product Where Deleted_At Is Null")

	// Apply filters
	if req.Name != "" {
		queryBuilder.WriteString(" AND p.name ILIKE '%")
		queryBuilder.WriteString(req.Name)
		queryBuilder.WriteString("%'")

		countQueryBuilder.WriteString(" AND name ILIKE '%")
		countQueryBuilder.WriteString(req.Name)
		countQueryBuilder.WriteString("%'")
	}

	if req.Category != "" {
		queryBuilder.WriteString(" AND p.category_id = '")
		queryBuilder.WriteString(req.Category)
		queryBuilder.WriteString("'")

		countQueryBuilder.WriteString(" AND category_id = '")
		countQueryBuilder.WriteString(req.Category)
		countQueryBuilder.WriteString("'")
	}

	if req.Unit != "" {
		queryBuilder.WriteString(" AND p.unit_id = '")
		queryBuilder.WriteString(req.Unit)
		queryBuilder.WriteString("'")

		countQueryBuilder.WriteString(" AND unit_id = '")
		countQueryBuilder.WriteString(req.Unit)
		countQueryBuilder.WriteString("'")
	}

	// Add sorting
	if req.SortBy != "" {
		queryBuilder.WriteString(" ORDER BY ")
		queryBuilder.WriteString(req.SortBy)
		if req.SortOrder != "" {
			queryBuilder.WriteString(" ")
			queryBuilder.WriteString(req.SortOrder)
		}
	} else {
		queryBuilder.WriteString(" ORDER BY P.Name")
	}

	// Add pagination
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", req.PageSize, (req.Page-1)*req.PageSize))

	// Execute count query first
	err := r.db.QueryRow(countQueryBuilder.String()).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Execute main query
	rows, err := r.db.Query(queryBuilder.String())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var product GetProductResponse
		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Category,
			&product.CategoryID,
			&product.Unit,
			&product.UnitID,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
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
	query := `
		Select P.Id, P.Name, P.Description,
		       C.Name As Category, U.Name As Unit, 
		       P.Created_At, P.Updated_At 
		From Product P
		Join Category C On P.Category_Id = C.Id
		Join Unit U On P.Unit_Id = U.Id
		Where P.Id = $1 And P.Deleted_At Is Null`

	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.Unit,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
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
	query := `
		Select P.Id, P.Name, P.Description,
		       C.Name As Category, U.Name As Unit, 
		       P.Created_At, P.Updated_At 
		From Product P
		Join Category C On P.Category_Id = C.Id
		Join Unit U On P.Unit_Id = U.Id
		Where P.Name = $1 And P.Deleted_At Is Null`

	err := r.db.QueryRow(query, name).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.Unit,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Create inserts a new product
func (r *ProductRepositoryImpl) Create(req CreateProductRequest) (*GetProductResponse, error) {
	var product GetProductResponse
	query := `
		Insert Into Product (Name, Description, Category_Id, Unit_Id) 
		Values ($1, $2, $3, $4) 
		Returning Id, Name, Description,
		(Select Name From Category Where Id = $3) As Category, 
		(Select Name From Unit Where Id = $4) As Unit, 
		Created_At, Updated_At`

	err := r.db.QueryRow(query, req.Name, req.Description, req.CategoryID, req.UnitID).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.Unit,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update updates an existing product
func (r *ProductRepositoryImpl) Update(req UpdateProductRequest) (*GetProductResponse, error) {
	var product GetProductResponse
	query := `
		Update Product 
		Set Name = $1, Description = $2, Category_Id = $3, Unit_Id = $4, Updated_At = Now() 
		Where Id = $5
		Returning Id, Name, Description,
		(Select Name From Category Where Id = $3) As Category, 
		(Select Name From Unit Where Id = $4) As Unit, 
		Created_At, Updated_At`

	err := r.db.QueryRow(query, req.Name, req.Description, req.CategoryID, req.UnitID, req.ID).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.Unit,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
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

// GetProductBatches fetches all batches for a product with storage information
func (r *ProductRepositoryImpl) GetProductBatches(req GetProductBatchesRequest) ([]ProductBatchResponse, int, error) {
	var batches []ProductBatchResponse
	var totalItems int

	// Count total batches for this product
	countQuery := "Select Count(Id) From Product_Batch Where Product_Id = $1"
	err := r.db.QueryRow(countQuery, req.ProductID).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Build query for batches with pagination
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		Select Id, Sku, Purchase_Order_Id, Initial_Quantity, Current_Quantity, Unit_Price, Created_At
		From Product_Batch
		Where Product_Id = $1
		Order By Created_At Desc
		Limit $2 Offset $3`)

	rows, err := r.db.Query(queryBuilder.String(), req.ProductID, req.PageSize, (req.Page-1)*req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var batch ProductBatchResponse
		err = rows.Scan(
			&batch.BatchID,
			&batch.SKU,
			&batch.PurchaseOrderID,
			&batch.InitialQuantity,
			&batch.CurrentQuantity,
			&batch.UnitPrice,
			&batch.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Get storage information for this batch
		storageQuery := `
			Select Bs.Storage_Id, S.Name, Bs.Quantity
			From Batch_Storage Bs
			Join Storage S On Bs.Storage_Id = S.Id
			Where Bs.Batch_Id = $1 And S.Deleted_At Is Null
		`
		storageRows, errQ := r.db.Query(storageQuery, batch.BatchID)
		if errQ != nil {
			return nil, 0, errQ
		}

		for storageRows.Next() {
			var storage ProductBatchInStorage
			err = storageRows.Scan(&storage.StorageID, &storage.StorageName, &storage.Quantity)
			if err != nil {
				storageRows.Close()
				return nil, 0, err
			}
			batch.StorageDetails = append(batch.StorageDetails, storage)
		}
		storageRows.Close()

		batches = append(batches, batch)
	}

	return batches, totalItems, nil
}
