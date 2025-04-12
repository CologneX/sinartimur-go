package purchase

import (
	"database/sql"
	"errors"
	"fmt"
	"sinartimur-go/utils"
	"time"
)

// SupplierRepository interface defines methods for supplier operations
type SupplierRepository interface {
	GetAll(req GetSupplierRequest) ([]GetSupplierResponse, int, error)
	GetByID(id string) (*GetSupplierResponse, error)
	GetByName(name string) (*GetSupplierResponse, error)
	Create(req CreateSupplierRequest) error
	Update(req UpdateSupplierRequest) error
	Delete(id string) error
}

// SupplierRepositoryImpl implements SupplierRepository
type SupplierRepositoryImpl struct {
	db *sql.DB
}

// NewSupplierRepository creates a new instance of SupplierRepositoryImpl
func NewSupplierRepository(db *sql.DB) SupplierRepository {
	return &SupplierRepositoryImpl{db: db}
}

// GetAll fetches all suppliers with filtering and pagination
func (r *SupplierRepositoryImpl) GetAll(req GetSupplierRequest) ([]GetSupplierResponse, int, error) {
	// Build count query using QueryBuilder
	countBuilder := utils.NewQueryBuilder(`
		Select Count(Id)
		From Supplier
		Where Deleted_At Is Null
	`)

	if req.Name != "" {
		countBuilder.AddFilter("Name ILIKE", "%"+req.Name+"%")
	}

	if req.Telephone != "" {
		countBuilder.AddFilter("Telephone ILIKE", "%"+req.Telephone+"%")
	}

	countQuery, countParams := countBuilder.Build()

	// Execute count query
	var totalItems int
	err := r.db.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("count query error: %w", err)
	}

	// Build main query
	queryBuilder := utils.NewQueryBuilder(`
		Select Id, Name, Address, Telephone, Created_At, Updated_At
		From Supplier
		Where Deleted_At Is Null
	`)

	if req.Name != "" {
		queryBuilder.AddFilter("Name ILIKE", "%"+req.Name+"%")
	}

	if req.Telephone != "" {
		queryBuilder.AddFilter("Telephone ILIKE", "%"+req.Telephone+"%")
	}

	// Add sorting
	queryBuilder.Query.WriteString(" ORDER BY Created_At DESC")

	// Add pagination
	mainQuery, queryParams := queryBuilder.AddPagination(req.PageSize, req.Page).Build()
	// Execute main query
	rows, err := r.db.Query(mainQuery, queryParams...)
	if err != nil {
		return nil, 0, fmt.Errorf("main query error: %w", err)
	}
	defer rows.Close()

	// Process results
	var suppliers []GetSupplierResponse
	for rows.Next() {
		var supplier GetSupplierResponse
		var createdAt, updatedAt time.Time

		errScan := rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.Address,
			&supplier.Telephone,
			&createdAt,
			&updatedAt,
		)
		if errScan != nil {
			return nil, 0, fmt.Errorf("scan error: %w", errScan)
		}

		supplier.CreatedAt = createdAt.Format(time.RFC3339)
		supplier.UpdatedAt = updatedAt.Format(time.RFC3339)
		suppliers = append(suppliers, supplier)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return suppliers, totalItems, nil
}

// GetByID fetches a supplier by ID
func (r *SupplierRepositoryImpl) GetByID(id string) (*GetSupplierResponse, error) {
	var supplier GetSupplierResponse
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(`
		Select Id, Name, Address, Telephone, Created_At, Updated_At
		From Supplier
		Where Id = $1 And Deleted_At Is Null
	`, id).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.Address,
		&supplier.Telephone,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("supplier tidak ditemukan")
		}
		return nil, err
	}

	supplier.CreatedAt = createdAt.Format(time.RFC3339)
	supplier.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &supplier, nil
}

// GetByName fetches a supplier by name
func (r *SupplierRepositoryImpl) GetByName(name string) (*GetSupplierResponse, error) {
	var supplier GetSupplierResponse
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(`
		Select Id, Name, Address, Telephone, Created_At, Updated_At
		From Supplier
		Where Name = $1 And Deleted_At Is Null
	`, name).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.Address,
		&supplier.Telephone,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("supplier tidak ditemukan")
		}
		return nil, err
	}

	supplier.CreatedAt = createdAt.Format(time.RFC3339)
	supplier.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &supplier, nil
}

// Create inserts a new supplier
func (r *SupplierRepositoryImpl) Create(req CreateSupplierRequest) error {
	_, err := r.db.Exec(`
		Insert Into Supplier (Name, Address, Telephone)
		Values ($1, $2, $3)
	`, req.Name, req.Address, req.Telephone)

	if err != nil {
		return fmt.Errorf("create supplier: %w", err)
	}

	return nil
}

// Update modifies an existing supplier
func (r *SupplierRepositoryImpl) Update(req UpdateSupplierRequest) error {
	// First get the current supplier to keep unchanged fields
	current, err := r.GetByID(req.ID)
	if err != nil {
		return err
	}

	// Prepare update values
	name := req.Name
	if name == "" {
		name = current.Name
	}

	address := req.Address
	if address == "" {
		address = current.Address
	}

	telephone := req.Telephone
	if telephone == "" {
		telephone = current.Telephone
	}

	_, err = r.db.Exec(`
		Update Supplier
		Set Name = $1, Address = $2, Telephone = $3, Updated_At = Current_Timestamp
		Where Id = $4 And Deleted_At Is Null
	`, name, address, telephone, req.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("supplier tidak ditemukan")
		}
		return fmt.Errorf("update supplier: %w", err)
	}

	return nil
}

// Delete soft deletes a supplier
func (r *SupplierRepositoryImpl) Delete(id string) error {
	result, err := r.db.Exec(`
		Update Supplier
		Set Deleted_At = Current_Timestamp
		Where Id = $1 And Deleted_At Is Null
	`, id)

	if err != nil {
		return fmt.Errorf("delete supplier: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("supplier tidak ditemukan")
	}

	return nil
}
