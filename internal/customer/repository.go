package customer

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sinartimur-go/utils"
	"time"
)

type CustomerRepository interface {
	GetAll(req GetCustomerRequest) ([]GetCustomerResponse, int, error)
	GetByID(id string) (*GetCustomerResponse, error)
	GetByName(name string) (*GetCustomerResponse, error)
	Create(req CreateCustomerRequest) error
	Update(req UpdateCustomerRequest) error
	Delete(req DeleteCustomerRequest) error
}

type RepositoryImpl struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) GetAll(req GetCustomerRequest) ([]GetCustomerResponse, int, error) {
	// Build the base query for selecting customer
	queryBuilder := utils.NewQueryBuilder("SELECT id, name, address, telephone, created_at, updated_at FROM customer WHERE deleted_at IS NULL")

	// Add filters based on the request parameters
	queryBuilder.AddFilter("name ILIKE ", "%"+req.Name+"%")
	queryBuilder.AddFilter("address ILIKE ", "%"+req.Address+"%")
	queryBuilder.AddFilter("telephone ILIKE ", "%"+req.Telephone+"%")

	// Build count query to get total items
	countQuery, countParams := queryBuilder.Build()
	countQuery = fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", countQuery)

	// Execute count query
	var totalItems int
	err := r.db.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total pelanggan: %w", err)
	}

	// Add pagination and sorting
	if req.SortBy != "" {
		direction := "ASC"
		if req.SortOrder == "desc" {
			direction = "DESC"
		}
		queryBuilder.Query.WriteString(fmt.Sprintf(" ORDER BY %s %s", req.SortBy, direction))
	}

	queryBuilder.AddPagination(req.PageSize, req.Page)

	// Execute the final query
	query, params := queryBuilder.Build()
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data pelanggan: %w", err)
	}
	defer rows.Close()

	// Process the results
	var customer []GetCustomerResponse
	for rows.Next() {
		var c GetCustomerResponse
		var createdAt, updatedAt time.Time

		if errScan := rows.Scan(&c.ID, &c.Name, &c.Address, &c.Telephone, &createdAt, &updatedAt); errScan != nil {
			return nil, 0, fmt.Errorf("gagal membaca data pelanggan: %w", errScan)
		}

		c.CreatedAt = createdAt.Format(time.RFC3339)
		c.UpdatedAt = updatedAt.Format(time.RFC3339)
		customer = append(customer, c)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("terjadi kesalahan saat membaca data pelanggan: %w", err)
	}

	return customer, totalItems, nil
}

func (r *RepositoryImpl) GetByID(id string) (*GetCustomerResponse, error) {
	query := `
		SELECT id, name, address, telephone, created_at, updated_at
		FROM customer
		WHERE id = $1 AND deleted_at IS NULL
	`

	var customer GetCustomerResponse
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Address,
		&customer.Telephone,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pelanggan dengan ID %s tidak ditemukan", id)
		}
		return nil, fmt.Errorf("gagal mengambil data pelanggan: %w", err)
	}

	customer.CreatedAt = createdAt.Format(time.RFC3339)
	customer.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &customer, nil
}

func (r *RepositoryImpl) GetByName(name string) (*GetCustomerResponse, error) {
	query := `
		SELECT id, name, address, telephone, created_at, updated_at
		FROM customer
		WHERE name = $1 AND deleted_at IS NULL
	`

	var customer GetCustomerResponse
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, name).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Address,
		&customer.Telephone,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("pelanggan dengan nama %s tidak ditemukan", name)
		}
		return nil, fmt.Errorf("gagal mengambil data pelanggan: %w", err)
	}

	customer.CreatedAt = createdAt.Format(time.RFC3339)
	customer.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &customer, nil
}

func (r *RepositoryImpl) Create(req CreateCustomerRequest) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO customer (id, name, address, telephone, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

		id := uuid.New()
		now := time.Now()

		_, err := tx.Exec(query, id, req.Name, req.Address, req.Telephone, now, now)
		if err != nil {
			return fmt.Errorf("gagal membuat pelanggan baru: %w", err)
		}

		return nil
	})
}

func (r *RepositoryImpl) Update(req UpdateCustomerRequest) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// First check if the customer exists
		checkQuery := "SELECT id FROM customer WHERE id = $1 AND deleted_at IS NULL"
		var customerID string
		err := tx.QueryRow(checkQuery, req.ID).Scan(&customerID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("pelanggan dengan ID %s tidak ditemukan", req.ID)
			}
			return fmt.Errorf("gagal memeriksa keberadaan pelanggan: %w", err)
		}

		// Perform the update
		updateQuery := `
			UPDATE customer
			SET name = $1, address = $2, telephone = $3, updated_at = $4
			WHERE id = $5 AND deleted_at IS NULL
		`

		_, err = tx.Exec(updateQuery, req.Name, req.Address, req.Telephone, time.Now(), req.ID)
		if err != nil {
			return fmt.Errorf("gagal memperbarui data pelanggan: %w", err)
		}

		return nil
	})
}

func (r *RepositoryImpl) Delete(req DeleteCustomerRequest) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Using soft delete by setting deleted_at
		query := `
			UPDATE customer
			SET deleted_at = $1
			WHERE id = $2 AND deleted_at IS NULL
		`

		result, err := tx.Exec(query, time.Now(), req.ID)
		if err != nil {
			return fmt.Errorf("gagal menghapus pelanggan: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("gagal mendapatkan jumlah baris yang terpengaruh: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("pelanggan dengan ID %s tidak ditemukan", req.ID)
		}

		return nil
	})
}
