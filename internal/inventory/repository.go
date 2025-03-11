package inventory

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sinartimur-go/utils"
	"time"
)

// StorageRepository defines the interface for storage data operations
type StorageRepository interface {
	// Storage CRUD operations
	GetAllStorages(req GetStorageRequest) ([]GetStorageResponse, int, error)
	GetStorageByID(id string) (*GetStorageResponse, error)
	GetStorageByName(name string) (*Storage, error)
	CreateStorage(req CreateStorageRequest) (*GetStorageResponse, error)
	UpdateStorage(req UpdateStorageRequest) (*GetStorageResponse, error)
	DeleteStorage(id string) error

	// Batch movement operations
	MoveBatch(req MoveBatchRequest, userID string) error
	GetBatchInStorage(batchID string, storageID string) (*BatchStorage, error)
	UpdateBatchInStorage(batchStorage BatchStorage) error
	CreateBatchInStorage(batchStorage BatchStorage) error
	LogInventoryMovement(log InventoryLog) error

	// StorageRepository interface
	GetInventoryLogs(req GetInventoryLogsRequest) ([]GetInventoryLogResponse, int, error)
	RefreshInventoryLogView() error
	GetInventoryLogLastRefreshed() (*time.Time, error)
}

// StorageRepositoryImpl implements the StorageRepository interface
type StorageRepositoryImpl struct {
	db *sql.DB
}

// NewStorageRepository creates a new storage repository instance
func NewStorageRepository(db *sql.DB) StorageRepository {
	return &StorageRepositoryImpl{db: db}
}

// GetAllStorages fetches all storage locations with pagination
func (r *StorageRepositoryImpl) GetAllStorages(req GetStorageRequest) ([]GetStorageResponse, int, error) {
	var storages []GetStorageResponse
	var totalItems int

	// Build base query
	qb := utils.NewQueryBuilder("Select Id, Name, Location, Created_At, Updated_At From Storage Where Deleted_At Is Null")

	// Add filters
	if req.Name != "" {
		qb.AddFilter("Name ILIKE ", `%`+req.Name+`%`)
	}
	if req.Location != "" {
		qb.AddFilter("Location ILIKE ", `%`+req.Location+`%`)
	}

	// Get count first
	//countQuery := fmt.Sprintf("Select Count(*) From (%S) As Filtered_Storages", qb.Query.String())
	countQuery := `Select Count(*) From (` + qb.Query.String() + `) As Filtered_Storages`
	countRow := r.db.QueryRow(countQuery, qb.Params...)
	if err := countRow.Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	// Add sorting
	if req.SortBy != "" && req.SortOrder != "" {
		qb.Query.WriteString(" ORDER BY " + req.SortBy + " " + req.SortOrder)
	} else {
		qb.Query.WriteString(" ORDER BY created_at DESC")
	}

	// Add pagination
	qb.AddPagination(req.PageSize, req.Page)

	// Execute final query
	query, params := qb.Build()
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var storage GetStorageResponse
		if errScan := rows.Scan(&storage.ID, &storage.Name, &storage.Location, &storage.CreatedAt, &storage.UpdatedAt); errScan != nil {
			return nil, 0, errScan
		}
		storages = append(storages, storage)
	}

	return storages, totalItems, nil
}

// GetStorageByID fetches a storage location by ID
func (r *StorageRepositoryImpl) GetStorageByID(id string) (*GetStorageResponse, error) {
	var storage GetStorageResponse
	err := r.db.QueryRow("Select Id, Name, Location, Created_At, Updated_At From Storage Where Id = $1 And Deleted_At Is Null", id).
		Scan(&storage.ID, &storage.Name, &storage.Location, &storage.CreatedAt, &storage.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &storage, nil
}

// GetStorageByName fetches a storage location by name
func (r *StorageRepositoryImpl) GetStorageByName(name string) (*Storage, error) {
	var storage Storage
	err := r.db.QueryRow("Select Id, Name, Location, Created_At, Updated_At, Deleted_At From Storage Where Name = $1 And Deleted_At Is Null", name).
		Scan(&storage.ID, &storage.Name, &storage.Location, &storage.CreatedAt, &storage.UpdatedAt, &storage.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &storage, nil
}

// CreateStorage creates a new storage location
func (r *StorageRepositoryImpl) CreateStorage(req CreateStorageRequest) (*GetStorageResponse, error) {
	var storage GetStorageResponse
	err := r.db.QueryRow("Insert Into Storage (Id, Name, Location) Values ($1, $2, $3) Returning Id, Name, Location, Created_At, Updated_At",
		uuid.New().String(), req.Name, req.Location).
		Scan(&storage.ID, &storage.Name, &storage.Location, &storage.CreatedAt, &storage.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &storage, nil
}

// UpdateStorage updates an existing storage location
func (r *StorageRepositoryImpl) UpdateStorage(req UpdateStorageRequest) (*GetStorageResponse, error) {
	var storage GetStorageResponse
	err := r.db.QueryRow("Update Storage Set Name = $1, Location = $2, Updated_At = Now() Where Id = $3 And Deleted_At Is Null Returning Id, Name, Location, Created_At, Updated_At",
		req.Name, req.Location, req.ID.String()).
		Scan(&storage.ID, &storage.Name, &storage.Location, &storage.CreatedAt, &storage.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &storage, nil
}

// DeleteStorage performs a soft delete on a storage location
func (r *StorageRepositoryImpl) DeleteStorage(id string) error {
	_, err := r.db.Exec("Update Storage Set Deleted_At = Now() Where Id = $1 And Deleted_At Is Null", id)
	return err
}

// GetBatchInStorage fetches a batch in a specific storage
func (r *StorageRepositoryImpl) GetBatchInStorage(batchID string, storageID string) (*BatchStorage, error) {
	var batchStorage BatchStorage
	err := r.db.QueryRow("Select Id, Batch_Id, Storage_Id, Quantity, Created_At, Updated_At From Batch_Storage Where Batch_Id = $1 And Storage_Id = $2",
		batchID, storageID).
		Scan(&batchStorage.ID, &batchStorage.BatchID, &batchStorage.StorageID, &batchStorage.Quantity, &batchStorage.CreatedAt, &batchStorage.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &batchStorage, nil
}

// UpdateBatchInStorage updates the quantity of a batch in storage
func (r *StorageRepositoryImpl) UpdateBatchInStorage(batchStorage BatchStorage) error {
	_, err := r.db.Exec("Update Batch_Storage Set Quantity = $1, Updated_At = Now() Where Id = $2",
		batchStorage.Quantity, batchStorage.ID)
	return err
}

// CreateBatchInStorage creates a new batch in storage entry
func (r *StorageRepositoryImpl) CreateBatchInStorage(batchStorage BatchStorage) error {
	_, err := r.db.Exec("Insert Into Batch_Storage (Id, Batch_Id, Storage_Id, Quantity) Values ($1, $2, $3, $4)",
		uuid.New().String(), batchStorage.BatchID, batchStorage.StorageID, batchStorage.Quantity)
	return err
}

// LogInventoryMovement logs an inventory movement action
func (r *StorageRepositoryImpl) LogInventoryMovement(log InventoryLog) error {
	_, err := r.db.Exec("Insert Into Inventory_Log (Id, Batch_Id, Storage_Id, Target_Storage_Id, User_Id, Action, Quantity, Log_Date, Description) Values ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		uuid.New().String(), log.BatchID, log.StorageID, log.TargetStorageID, log.UserID, log.Action, log.Quantity, log.LogDate, log.Description)
	return err
}

// MoveBatch moves a batch from one storage to another
func (r *StorageRepositoryImpl) MoveBatch(req MoveBatchRequest, userID string) error {
	// Use a transaction to ensure data consistency
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Get batch in source storage
		var sourceBatchStorage BatchStorage
		err := tx.QueryRow("Select Id, Batch_Id, Storage_Id, Quantity, Created_At, Updated_At From Batch_Storage Where Batch_Id = $1 And Storage_Id = $2",
			req.BatchID, req.SourceStorageID).
			Scan(&sourceBatchStorage.ID, &sourceBatchStorage.BatchID, &sourceBatchStorage.StorageID, &sourceBatchStorage.Quantity, &sourceBatchStorage.CreatedAt, &sourceBatchStorage.UpdatedAt)
		if err != nil {
			return err
		}

		// Check if source has enough quantity
		if sourceBatchStorage.Quantity < req.Quantity {
			return fmt.Errorf("kuantitas tidak mencukupi di gudang sumber")
		}

		// Update source storage quantity
		_, err = tx.Exec("Update Batch_Storage Set Quantity = Quantity - $1, Updated_At = Now() Where Id = $2",
			req.Quantity, sourceBatchStorage.ID)
		if err != nil {
			return err
		}

		// Check if batch exists in target storage
		var targetBatchExists bool
		err = tx.QueryRow("Select Exists(Select 1 From Batch_Storage Where Batch_Id = $1 And Storage_Id = $2)",
			req.BatchID, req.TargetStorageID).Scan(&targetBatchExists)
		if err != nil {
			return err
		}

		// If batch exists in target, update quantity, otherwise create new entry
		if targetBatchExists {
			_, err = tx.Exec("Update Batch_Storage Set Quantity = Quantity + $1, Updated_At = Now() Where Batch_Id = $2 And Storage_Id = $3",
				req.Quantity, req.BatchID, req.TargetStorageID)
		} else {
			newID := uuid.New().String()
			_, err = tx.Exec("Insert Into Batch_Storage (Id, Batch_Id, Storage_Id, Quantity) Values ($1, $2, $3, $4)",
				newID, req.BatchID, req.TargetStorageID, req.Quantity)
		}
		if err != nil {
			return err
		}

		// Log the movement
		_, err = tx.Exec("Insert Into Inventory_Log (Id, Batch_Id, Storage_Id, Target_Storage_Id, User_Id, Action, Quantity, Log_Date, Description) Values ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
			uuid.New().String(), req.BatchID, req.SourceStorageID, req.TargetStorageID, userID, "transfer", req.Quantity, time.Now(), req.Description)
		if err != nil {
			return err
		}

		return nil
	})
}

// RefreshInventoryLogView refreshes the materialized view
func (r *StorageRepositoryImpl) RefreshInventoryLogView() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Refresh the materialized view
	_, err = tx.Exec("REFRESH MATERIALIZED VIEW inventory_log_view")
	if err != nil {
		return err
	}

	// Update the refresh timestamp
	_, err = tx.Exec("UPDATE materialized_view_refresh SET last_refreshed = NOW() WHERE view_name = 'inventory_log_view'")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *StorageRepositoryImpl) GetInventoryLogLastRefreshed() (*time.Time, error) {
	var lastRefreshed time.Time
	err := r.db.QueryRow("SELECT last_refreshed FROM materialized_view_refresh WHERE view_name = 'inventory_log_view'").Scan(&lastRefreshed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No refresh record found
		}
		return nil, err
	}
	return &lastRefreshed, nil
}

// GetInventoryLogs retrieves inventory logs based on the provided filters
func (r *StorageRepositoryImpl) GetInventoryLogs(req GetInventoryLogsRequest) ([]GetInventoryLogResponse, int, error) {
	var logs []GetInventoryLogResponse
	var totalItems int

	// Build base query
	qb := utils.NewQueryBuilder("SELECT * FROM inventory_log_view WHERE 1=1")

	// Add filters
	if req.BatchID != "" {
		qb.AddFilter("batch_id =", req.BatchID)
	}

	if req.ProductID != "" {
		qb.AddFilter("product_id =", req.ProductID)
	}

	if req.StorageID != "" {
		qb.AddFilter("storage_id =", req.StorageID)
	}

	if req.TargetStorageID != "" {
		qb.AddFilter("target_storage_id =", req.TargetStorageID)
	}

	if req.UserID != "" {
		qb.AddFilter("user_id =", req.UserID)
	}

	if req.Action != "" {
		qb.AddFilter("action =", req.Action)
	}

	if req.FromDate != "" {
		qb.AddFilter("log_date >", req.FromDate)
	}

	if req.ToDate != "" {
		qb.AddFilter("log_date <", req.ToDate)
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS filtered_logs", qb.Query.String())
	err := r.db.QueryRow(countQuery, qb.Params...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count total logs: %w", err)
	}

	// Add sorting
	if req.SortBy != "" && req.SortOrder != "" {
		qb.Query.WriteString(fmt.Sprintf(" ORDER BY %s %s", req.SortBy, req.SortOrder))
	} else {
		qb.Query.WriteString(" ORDER BY log_date DESC")
	}

	// Add pagination
	qb.AddPagination(req.PageSize, req.Page)

	// Execute final query
	query, params := qb.Build()
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query inventory logs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var log GetInventoryLogResponse
		var targetStorageID, targetStorageName, purchaseOrderID, salesOrderID sql.NullString

		errScan := rows.Scan(
			&log.ID,
			&log.BatchID,
			&log.BatchSKU,
			&log.ProductID,
			&log.ProductName,
			&log.StorageID,
			&log.StorageName,
			&targetStorageID,
			&targetStorageName,
			&log.UserID,
			&log.Username,
			&purchaseOrderID,
			&salesOrderID,
			&log.Action,
			&log.Quantity,
			&log.LogDate,
			&log.Description,
			&log.CreatedAt,
		)
		if errScan != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory log: %w", errScan)
		}

		if targetStorageID.Valid {
			log.TargetStorageID = &targetStorageID.String
		}

		if targetStorageName.Valid {
			log.TargetStorageName = &targetStorageName.String
		}

		if purchaseOrderID.Valid {
			log.PurchaseOrderID = &purchaseOrderID.String
		}

		if salesOrderID.Valid {
			log.SalesOrderID = &salesOrderID.String
		}

		logs = append(logs, log)
	}

	return logs, totalItems, nil
}
