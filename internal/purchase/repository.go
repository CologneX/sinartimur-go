package purchase

import (
	"context"
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
	Create(req CreateSupplierRequest) (*GetSupplierResponse, error)
	Update(req UpdateSupplierRequest) (*GetSupplierResponse, error)
	Delete(id string) error
}

// SupplierRepositoryImpl implements SupplierRepository
type SupplierRepositoryImpl struct {
	db *sql.DB
}

// GetAll fetches all suppliers with filtering and pagination
func (r *SupplierRepositoryImpl) GetAll(ctx context.Context, req GetSupplierRequest) ([]GetSupplierResponse, int, error) {
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
	err := r.db.QueryRowContext(ctx, countQuery, countParams...).Scan(&totalItems)
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
	rows, err := r.db.QueryContext(ctx, mainQuery, queryParams...)
	if err != nil {
		return nil, 0, fmt.Errorf("main query error: %w", err)
	}
	defer rows.Close()

	// Process results
	var suppliers []GetSupplierResponse
	for rows.Next() {
		var supplier GetSupplierResponse
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.Address,
			&supplier.Telephone,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
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
func (r *SupplierRepositoryImpl) GetByID(ctx context.Context, id string) (*GetSupplierResponse, error) {
	var supplier GetSupplierResponse
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, `
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
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, err
	}

	supplier.CreatedAt = createdAt.Format(time.RFC3339)
	supplier.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &supplier, nil
}

// GetByName fetches a supplier by name
func (r *SupplierRepositoryImpl) GetByName(ctx context.Context, name string) (*GetSupplierResponse, error) {
	var supplier GetSupplierResponse
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, `
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
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, err
	}

	supplier.CreatedAt = createdAt.Format(time.RFC3339)
	supplier.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &supplier, nil
}

// Create inserts a new supplier
func (r *SupplierRepositoryImpl) Create(ctx context.Context, req CreateSupplierRequest) (*GetSupplierResponse, error) {
	var id string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, `
		Insert Into Supplier (Name, Address, Telephone)
		Values ($1, $2, $3)
		Returning Id, Created_At, Updated_At
	`, req.Name, req.Address, req.Telephone).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("create supplier: %w", err)
	}

	return &GetSupplierResponse{
		ID:        id,
		Name:      req.Name,
		Address:   req.Address,
		Telephone: req.Telephone,
		CreatedAt: createdAt.Format(time.RFC3339),
		UpdatedAt: updatedAt.Format(time.RFC3339),
	}, nil
}

// Update modifies an existing supplier
func (r *SupplierRepositoryImpl) Update(ctx context.Context, req UpdateSupplierRequest) (*GetSupplierResponse, error) {
	// First get the current supplier to keep unchanged fields
	current, err := r.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
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

	var updatedAt time.Time
	err = r.db.QueryRowContext(ctx, `
		Update Supplier
		Set Name = $1, Address = $2, Telephone = $3, Updated_At = Current_Timestamp
		Where Id = $4 And Deleted_At Is Null
		Returning Updated_At
	`, name, address, telephone, req.ID).Scan(&updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, fmt.Errorf("update supplier: %w", err)
	}

	return &GetSupplierResponse{
		ID:        req.ID,
		Name:      name,
		Address:   address,
		Telephone: telephone,
		CreatedAt: current.CreatedAt,
		UpdatedAt: updatedAt.Format(time.RFC3339),
	}, nil
}

// Delete soft deletes a supplier
func (r *SupplierRepositoryImpl) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `
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
		return fmt.Errorf("supplier not found")
	}

	return nil
}

// PurchaseOrderRepository interface defines methods for purchase order operations
type PurchaseOrderRepository interface {
}
