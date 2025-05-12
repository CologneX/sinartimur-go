package finance

import (
	"database/sql"
	"errors"
	"fmt"
	"sinartimur-go/utils"
	"time"
)

// FinanceTransactionRepository defines operations for finance transactions
type FinanceTransactionRepository interface {
	Create(req CreateFinanceTransactionRequest, userID string) error
	GetAll(req GetFinanceTransactionRequest) ([]GetFinanceTransactionResponse, int, error)
	GetByID(id string) (*GetFinanceTransactionResponse, error)
	Cancel(req CancelFinanceTransactionRequest, userID string) error
	GetSummary(startDate, endDate time.Time) (*FinanceTransactionSummary, error)
	RefreshFinanceTransactionView() error
	GetFinanceTransactionViewLastRefreshed() (*time.Time, error)
}

type financeTransactionRepositoryImpl struct {
	db *sql.DB
}

// NewFinanceTransactionRepository creates a new finance transaction repository
func NewFinanceTransactionRepository(db *sql.DB) FinanceTransactionRepository {
	return &financeTransactionRepositoryImpl{db: db}
}

// Create adds a new finance transaction to the database
func (r *financeTransactionRepositoryImpl) Create(req CreateFinanceTransactionRequest, userID string) error {
	query := `
		INSERT INTO Financial_Transaction_Log 
		(User_Id, Amount, Type, Purchase_Order_Id, Sales_Order_Id, Description, Is_System, Transaction_Date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// Prepare parameters
	var purchaseOrderID, salesOrderID *string
	if req.PurchaseOrderID != "" {
		purchaseOrderID = &req.PurchaseOrderID
	}
	if req.SalesOrderID != "" {
		salesOrderID = &req.SalesOrderID
	}

	// For manually created transactions, Is_System is always false
	isSystem := false

	_, err := r.db.Exec(
		query,
		userID,
		req.Amount,
		req.Type,
		purchaseOrderID,
		salesOrderID,
		req.Description,
		isSystem,
		req.TransactionDate,
	)

	if err != nil {
		return fmt.Errorf("gagal membuat transaksi keuangan: %w", err)
	}

	return nil
}

// GetAll fetches financial transactions with filtering and pagination using the materialized view
func (r *financeTransactionRepositoryImpl) GetAll(req GetFinanceTransactionRequest) ([]GetFinanceTransactionResponse, int, error) {
	// Build base query using the materialized view
	queryBuilder := utils.NewQueryBuilder(`
        SELECT 
            Id, 
            User_Id, 
            Username, 
            Amount, 
            Type, 
            Purchase_Order_Id, 
            Sales_Order_Id, 
            Description, 
            Is_System, 
            Transaction_Date, 
            Created_At, 
            Edited_At,
			Deleted_At
        FROM 
            finance_transaction_log_view
        WHERE 
            1=1
    `)

	// Add filters
	if req.UserID != "" {
		queryBuilder.AddFilter("User_Id =", req.UserID)
	}

	if req.Type != "" {
		queryBuilder.AddFilter("Type =", req.Type)
	}

	if req.PurchaseOrderID != "" {
		queryBuilder.AddFilter("Purchase_Order_Id =", req.PurchaseOrderID)
	}

	if req.SalesOrderID != "" {
		queryBuilder.AddFilter("Sales_Order_Id =", req.SalesOrderID)
	}

	if req.IsSystem != nil {
		queryBuilder.AddFilter("Is_System =", *req.IsSystem)
	}

	// Date range filters
	if req.StartDate != "" {
		queryBuilder.AddFilter("Transaction_Date >=", req.StartDate)
	}

	if req.EndDate != "" {
		queryBuilder.AddFilter("Transaction_Date <=", req.EndDate)
	}

	// Get total count query
	countQuery, countParams := queryBuilder.Build()
	countQuery = fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", countQuery)

	// Execute count query
	var totalItems int
	err := r.db.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total transaksi: %w", err)
	}

	// Add sorting
	if req.SortBy != "" {
		direction := "ASC"
		if req.SortOrder == "desc" {
			direction = "DESC"
		}
		queryBuilder.Query.WriteString(fmt.Sprintf(" ORDER BY %s %s", req.SortBy, direction))
	} else {
		// Default sort by transaction date descending
		queryBuilder.Query.WriteString(" ORDER BY Transaction_Date DESC")
	}

	// Add pagination
	queryBuilder.AddPagination(req.PageSize, req.Page)

	// Execute final query
	query, params := queryBuilder.Build()
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data transaksi: %w", err)
	}
	defer rows.Close()

	var transactions []GetFinanceTransactionResponse
	for rows.Next() {
		var tx GetFinanceTransactionResponse
		var purchaseOrderID, salesOrderID sql.NullString

		err = rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Username,
			&tx.Amount,
			&tx.Type,
			&purchaseOrderID,
			&salesOrderID,
			&tx.Description,
			&tx.IsSystem,
			&tx.TransactionDate,
			&tx.CreatedAt,
			&tx.EditedAt,
			&tx.DeletedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("gagal membaca data transaksi: %w", err)
		}

		if purchaseOrderID.Valid {
			tx.PurchaseOrderID = purchaseOrderID.String
		}

		if salesOrderID.Valid {
			tx.SalesOrderID = salesOrderID.String
		}

		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("terjadi kesalahan saat membaca data transaksi: %w", err)
	}

	return transactions, totalItems, nil
}

// GetByID fetches a single finance transaction by ID
func (r *financeTransactionRepositoryImpl) GetByID(id string) (*GetFinanceTransactionResponse, error) {
	query := `
		SELECT 
			ft.Id, 
			ft.User_Id, 
			u.Username, 
			ft.Amount, 
			ft.Type, 
			ft.Purchase_Order_Id, 
			ft.Sales_Order_Id, 
			ft.Description, 
			ft.Is_System, 
			ft.Transaction_Date, 
			ft.Created_At, 
			ft.Edited_At
		FROM 
			Financial_Transaction_Log ft
		LEFT JOIN 
			Appuser u ON ft.User_Id = u.Id
		WHERE 
			ft.Id = $1 AND ft.Deleted_At IS NULL
	`

	var tx GetFinanceTransactionResponse
	var purchaseOrderID, salesOrderID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.UserID,
		&tx.Username,
		&tx.Amount,
		&tx.Type,
		&purchaseOrderID,
		&salesOrderID,
		&tx.Description,
		&tx.IsSystem,
		&tx.TransactionDate,
		&tx.CreatedAt,
		&tx.EditedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaksi keuangan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data transaksi: %w", err)
	}

	if purchaseOrderID.Valid {
		tx.PurchaseOrderID = purchaseOrderID.String
	}

	if salesOrderID.Valid {
		tx.SalesOrderID = salesOrderID.String
	}

	return &tx, nil
}

// Cancel soft deletes a finance transaction and adds cancellation info
func (r *financeTransactionRepositoryImpl) Cancel(req CancelFinanceTransactionRequest, userID string) error {
	// Get the original transaction first to verify it exists
	tx, err := r.GetByID(req.ID)
	if err != nil {
		return err
	}

	// Can't cancel system-generated transactions
	if tx.IsSystem {
		return fmt.Errorf("transaksi sistem tidak dapat dibatalkan")
	}

	// Update the description to include cancellation reason and mark as deleted
	query := `
		UPDATE Financial_Transaction_Log 
		SET 
			Description = $1, 
			Edited_At = NOW(), 
			Deleted_At = NOW() 
		WHERE Id = $2
	`

	cancelDescription := fmt.Sprintf("%s [DIBATALKAN: %s]", tx.Description, req.Description)

	_, err = r.db.Exec(query, cancelDescription, req.ID)
	if err != nil {
		return fmt.Errorf("gagal membatalkan transaksi keuangan: %w", err)
	}

	return nil
}

// GetSummary retrieves a summary of financial transactions within a date range
func (r *financeTransactionRepositoryImpl) GetSummary(startDate, endDate time.Time) (*FinanceTransactionSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN Type IN ('income', 'sale') THEN Amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN Type IN ('expense', 'purchase') THEN Amount ELSE 0 END), 0) AS total_expense
		FROM
			Financial_Transaction_Log
		WHERE
			Deleted_At IS NULL
			AND Transaction_Date BETWEEN $1 AND $2
	`

	var summary FinanceTransactionSummary

	err := r.db.QueryRow(query, startDate, endDate).Scan(
		&summary.TotalIncome,
		&summary.TotalExpense,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil ringkasan keuangan: %w", err)
	}

	// Calculate net amount
	summary.NetAmount = summary.TotalIncome - summary.TotalExpense

	// Format period
	if startDate.Year() == endDate.Year() && startDate.Month() == endDate.Month() {
		summary.Period = fmt.Sprintf("%s %d", startDate.Month().String(), startDate.Year())
	} else {
		summary.Period = fmt.Sprintf("%s %d - %s %d",
			startDate.Month().String(), startDate.Year(),
			endDate.Month().String(), endDate.Year())
	}

	return &summary, nil
}

func (r *financeTransactionRepositoryImpl) RefreshFinanceTransactionView() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Refresh the materialized view
	_, err = tx.Exec("REFRESH MATERIALIZED VIEW finance_transaction_log_view")
	if err != nil {
		return err
	}

	// Update the refresh timestamp
	_, err = tx.Exec("UPDATE materialized_view_refresh SET last_refreshed = NOW() WHERE view_name = 'finance_transaction_log_view'")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *financeTransactionRepositoryImpl) GetFinanceTransactionViewLastRefreshed() (*time.Time, error) {
	var lastRefreshed time.Time
	err := r.db.QueryRow("SELECT last_refreshed FROM materialized_view_refresh WHERE view_name = 'finance_transaction_log_view'").Scan(&lastRefreshed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &lastRefreshed, nil
}
