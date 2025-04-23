package purchase_order

import (
	"database/sql"
	"fmt"
	"sinartimur-go/internal/product"
	"sinartimur-go/utils"
	"strings"
	"time"
)

// Repository interface defines methods for purchase purchase-order operations
type Repository interface {
	// Basic CRUD operations
	GetAll(req GetPurchaseOrderRequest) ([]GetPurchaseOrderResponse, int, error)
	GetByID(id string) (*GetPurchaseOrderDetailResponse, error)

	// Core purchase order operations with transaction support
	Create(req CreatePurchaseOrderRequest, userID string, tx *sql.Tx) (string, error)
	Update(req UpdatePurchaseOrderRequest, tx *sql.Tx) (string, error)
	// CompletePurchaseOrder(id string, items []ReceivedItemRequest, userID string, tx *sql.Tx) error
	CompleteFullPurchaseOrder(id string, storageID string, userID string, tx *sql.Tx) error

	// Status change operations
	UpdateStatus(id, status, userID string, tx *sql.Tx) error
	CheckPurchaseOrder(id string, userID string, tx *sql.Tx) error
	CancelPurchaseOrder(id string, userID string, tx *sql.Tx) error

	// Return operations
	CreateReturn(req CreatePurchaseOrderReturnRequest, userID string, tx *sql.Tx) error
	CancelReturn(id string, userID string, tx *sql.Tx) error
	GetAllReturns(req GetPurchaseOrderReturnRequest) ([]GetPurchaseOrderReturnResponse, int, error)

	// Item operations
	AddPurchaseOrderItem(orderID string, req CreatePurchaseOrderItemRequest, tx *sql.Tx) error
	UpdatePurchaseOrderItem(req UpdatePurchaseOrderItemRequest, tx *sql.Tx) error
	RemovePurchaseOrderItem(id string, tx *sql.Tx) error

	// Logging operations
	LogInventoryChange(batchID, storageID, userID, orderID string, action string, quantity float64, description string, tx *sql.Tx) error
	LogFinancialTransaction(userID string, amount float64, transactionType string, orderID string, description string, tx *sql.Tx) error

	// Batch management
	CreateProductBatch(productID, orderID, sku string, quantity, unitPrice float64, tx *sql.Tx) (*string, error)
	AssignBatchToStorage(batchID, storageID string, quantity float64, tx *sql.Tx) error

	// Utility functions
	GenerateBatchSKU(productName, serialID string, supplierName string, date time.Time) (string, error)
	CheckAllItemsReturned(orderID string, tx *sql.Tx) (bool, error)

	// Get Products
	GetProducts(req product.GetProductRequest) ([]product.GetProductResponse, int, error)
}

type RepositoryImpl struct {
	DB *sql.DB
}

// NewRepository creates a new purchase order repository
func NewPurhaseOrderRepository(db *sql.DB) Repository {
	return &RepositoryImpl{
		DB: db,
	}
}

// Create inserts a new purchase order with transaction support
func (r *RepositoryImpl) Create(req CreatePurchaseOrderRequest, userID string, tx *sql.Tx) (string, error) {
	var executor interface {
		QueryRow(string, ...interface{}) *sql.Row
		Exec(string, ...interface{}) (sql.Result, error)
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Parse order date
	orderDate, err := time.Parse(time.RFC3339, req.OrderDate)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	// Validate payment due date for credit orders
	var paymentDueDate sql.NullTime
	if req.PaymentMethod == "credit" {
		if req.PaymentDueDate == "" {
			return "", fmt.Errorf("payment due date is required for credit orders")
		}

		dueDate, err := time.Parse(time.RFC3339, req.PaymentDueDate)
		if err != nil {
			return "", fmt.Errorf("invalid payment due date format: %w", err)
		}

		paymentDueDate = sql.NullTime{
			Time:  dueDate,
			Valid: true,
		}
	}

	// Generate Serial ID
	var serialID string
	if tx != nil {
		serialID, err = utils.GenerateNextSerialID(tx, "PO")
		if err != nil {
			return "", fmt.Errorf("failed to generate serial ID: %w", err)
		}
	} else {
		// If no transaction provided, create one temporarily just for serial ID generation
		err = utils.WithTransaction(r.DB, func(tempTx *sql.Tx) error {
			var genErr error
			serialID, genErr = utils.GenerateNextSerialID(tempTx, "PO")
			return genErr
		})
		if err != nil {
			return "", fmt.Errorf("failed to generate serial ID: %w", err)
		}
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Quantity * item.Price
	}

	// Insert purchase order
	var orderID string
	err = executor.QueryRow(`
        Insert Into Purchase_Order (
            Serial_Id, Supplier_Id, Order_Date, Status, 
            Total_Amount, Payment_Method, Payment_Due_Date, 
            Created_By
        )
        Values ($1, $2, $3, $4, $5, $6, $7, $8)
        Returning Id
    `, serialID, req.SupplierID, orderDate, "order",
		totalAmount, req.PaymentMethod, paymentDueDate, userID).Scan(&orderID)

	if err != nil {
		return "", err
	}

	// Insert order items
	for _, item := range req.Items {
		_, err = executor.Exec(`
            Insert Into Purchase_Order_Detail (
                Purchase_Order_Id, Product_Id, 
                Requested_Quantity, Unit_Price
            )
            Values ($1, $2, $3, $4)
        `, orderID, item.ProductID, item.Quantity, item.Price)

		if err != nil {
			return "", fmt.Errorf("failed to add order item: %w", err)
		}
	}

	return orderID, nil
}

// GetByID retrieves a purchase order with its details
func (r *RepositoryImpl) GetByID(id string) (*GetPurchaseOrderDetailResponse, error) {
	var po GetPurchaseOrderDetailResponse

	// Get purchase order
	err := r.DB.QueryRow(`
        Select 
            Po.Id, Po.Serial_Id, Po.Supplier_Id, S.Name As Suppliername,
            Po.Order_Date, Po.Status, Po.Total_Amount, Po.Payment_Method,
            Po.Payment_Due_Date, Po.Created_By, U.Username As Createdbyname,
            Po.Checked_By, U2.Username As Checkedbyname,
            Po.Created_At, Po.Updated_At,
			S.Address As Supplieraddress, S.Telephone As Supplierphone,
			Po.Cancelled_At, Po.Cancelled_By, U3.Username As Cancelledbyname
        From Purchase_Order Po
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Left Join Appuser U On Po.Created_By = U.Id
        Left Join Appuser U2 On Po.Checked_By = U2.Id
		Left Join Appuser U3 On Po.Cancelled_By = U3.Id
        Where Po.Id = $1
    `, id).Scan(
		&po.ID, &po.SerialID, &po.SupplierID, &po.SupplierName,
		&po.OrderDate, &po.Status, &po.TotalAmount, &po.PaymentMethod,
		&po.PaymentDueDate, &po.CreatedBy, &po.CreatedByName,
		&po.CheckedBy, &po.CheckedByName,
		&po.CreatedAt, &po.UpdatedAt,
		&po.SupplierAddress, &po.SupplierPhone,
		&po.CancelledAt, &po.CancelledBy, &po.CancelledByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase order not found")
		}
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}

	// Get order items with return and received information
	rows, err := r.DB.Query(`
        Select 
            Pod.Id, Pod.Product_Id, P.Name As Productname,
            Pod.Requested_Quantity, Pod.Unit_Price,
            Coalesce(
                (Select Sum(Por.Return_Quantity) 
                 From Purchase_Order_Return Por 
                 Where Por.Product_Detail_Id = Pod.Id And Por.Status = 'returned'),
                0
            ) As Returnedquantity,
            Coalesce(
                (Select Sum(Pb.Initial_Quantity)
                 From Product_Batch Pb
                 Where Pb.Purchase_Order_Id = Pod.Purchase_Order_Id And Pb.Product_Id = Pod.Product_Id),
                0
            ) As Receivedquantity,
            Case When 
                (Select Count(*) From Purchase_Order_Return Por 
                 Where Por.Product_Detail_Id = Pod.Id And Por.Status = 'returned') > 0 
                Then True Else False 
            End As Isreturned,
            (Select Por.Reason 
             From Purchase_Order_Return Por 
             Where Por.Product_Detail_Id = Pod.Id And Por.Status = 'returned'
             Order By Por.Returned_At Desc
             Limit 1) As Returnreason,
            Pod.Created_At, Pod.Updated_At
        From Purchase_Order_Detail Pod
        Join Product P On Pod.Product_Id = P.Id
        Where Pod.Purchase_Order_Id = $1
    `, id)

	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item PurchaseOrderItem
		var returnReason sql.NullString

		err := rows.Scan(
			&item.ID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.Price,
			&item.ReturnedQuantity, &item.ReceivedQuantity, &item.IsReturned, &returnReason,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		if returnReason.Valid {
			item.ReturnReason = &returnReason.String
		}

		po.Items = append(po.Items, item)
	}

	return &po, nil
}

// CompletePurchaseOrder implements the process of completing a purchase order
// func (r *RepositoryImpl) CompletePurchaseOrder(id string, items []ReceivedItemRequest, userID string, tx *sql.Tx) error {
// 	var executor interface {
// 		QueryRow(string, ...interface{}) *sql.Row
// 		Exec(string, ...interface{}) (sql.Result, error)
// 		Query(string, ...interface{}) (*sql.Rows, error)
// 	}

// 	if tx != nil {
// 		executor = tx
// 	} else {
// 		executor = r.DB
// 	}

// 	// Get purchase order info (supplier, serial ID)
// 	var serialID, supplierID, supplierName string
// 	err := executor.QueryRow(`
//         SELECT po.Serial_Id, po.Supplier_Id, s.Name
//         FROM Purchase_Order po
//         JOIN Supplier s ON po.Supplier_Id = s.Id
//         WHERE po.Id = $1
//     `, id).Scan(&serialID, &supplierID, &supplierName)

// 	if err != nil {
// 		return fmt.Errorf("failed to get purchase order info: %w", err)
// 	}

// 	// Update purchase order status
// 	_, err = executor.Exec(`
//         UPDATE Purchase_Order
//         SET Status = 'completed', Checked_By = $1, Updated_At = NOW()
//         WHERE Id = $2
//     `, userID, id)

// 	if err != nil {
// 		return fmt.Errorf("failed to update purchase order status: %w", err)
// 	}

// 	// Process received items
// 	now := time.Now()

// 	for _, item := range items {
// 		// Get product details
// 		var productID, productName string
// 		err := executor.QueryRow(`
//             SELECT p.Id, p.Name
//             FROM Product p
//             JOIN Purchase_Order_Detail pod ON pod.Product_Id = p.Id
//             WHERE pod.Id = $1
//         `, item.DetailID).Scan(&productID, &productName)

// 		if err != nil {
// 			return fmt.Errorf("failed to get product details: %w", err)
// 		}

// 		// Generate SKU
// 		sku, err := r.GenerateBatchSKU(productName, supplierName, now)
// 		if err != nil {
// 			return fmt.Errorf("failed to generate SKU: %w", err)
// 		}

// 		// Create product batch
// 		batchID, err := r.CreateProductBatch(productID, id, sku, item.Quantity, item.UnitPrice, tx)
// 		if err != nil {
// 			return fmt.Errorf("failed to create product batch: %w", err)
// 		}

// 		// Associate batch with storage
// 		err = r.AssignBatchToStorage(batchID, item.StorageID, item.Quantity, tx)
// 		if err != nil {
// 			return fmt.Errorf("failed to assign batch to storage: %w", err)
// 		}

// 		// Log inventory change
// 		description := fmt.Sprintf("Pembelian Barang %s", serialID)
// 		err = r.LogInventoryChange(batchID, item.StorageID, userID, id, "add", item.Quantity, description, tx)
// 		if err != nil {
// 			return fmt.Errorf("failed to log inventory change: %w", err)
// 		}
// 	}

// 	// Get total amount for financial log
// 	var totalAmount float64
// 	err = executor.QueryRow(`
//         SELECT Total_Amount FROM Purchase_Order WHERE Id = $1
//     `, id).Scan(&totalAmount)

// 	if err != nil {
// 		return fmt.Errorf("failed to get total amount: %w", err)
// 	}

// 	// Log financial transaction
// 	description := fmt.Sprintf("Pembelian Barang %s", serialID)
// 	err = r.LogFinancialTransaction(userID, totalAmount, "purchase", id, description, tx)
// 	if err != nil {
// 		return fmt.Errorf("failed to log financial transaction: %w", err)
// 	}

// 	return nil
// }

// CreateProductBatch creates a new product batch
func (r *RepositoryImpl) CreateProductBatch(productID, orderID, sku string, quantity, unitPrice float64, tx *sql.Tx) (*string, error) {
	var executor interface {
		QueryRow(string, ...interface{}) *sql.Row
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	var batchID string
	err := executor.QueryRow(`
        Insert Into Product_Batch (
            Sku, Product_Id, Purchase_Order_Id,
            Initial_Quantity, Current_Quantity, Unit_Price
        )
        Values ($1, $2, $3, $4, $5, $6)
		Returning Id
    `, sku, productID, orderID, quantity, quantity, unitPrice).Scan(&batchID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no batch id returned")
		}
		return nil, fmt.Errorf("insert returning id: %w", err)
	}

	return &batchID, nil
}

// AssignBatchToStorage assigns a batch to a storage location
func (r *RepositoryImpl) AssignBatchToStorage(batchID, storageID string, quantity float64, tx *sql.Tx) error {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	_, err := executor.Exec(`
        Insert Into Batch_Storage (Batch_Id, Storage_Id, Quantity)
        Values ($1, $2, $3)
        On Conflict (Batch_Id, Storage_Id) 
        Do Update Set Quantity = Batch_Storage.Quantity + $3
    `, batchID, storageID, quantity)

	if err != nil {
		return fmt.Errorf("failed to assign batch to storage: %w", err)
	}

	return nil
}

// LogInventoryChange creates an inventory log entry
func (r *RepositoryImpl) LogInventoryChange(batchID, storageID, userID, orderID string, action string, quantity float64, description string, tx *sql.Tx) error {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	_, err := executor.Exec(`
        Insert Into Inventory_Log (
            Batch_Id, Storage_Id, User_Id, Purchase_Order_Id,
            Action, Quantity, Description, Log_Date
        )
        Values ($1, $2, $3, $4, $5, $6, $7, Now())
    `, batchID, storageID, userID, orderID, action, quantity, description)

	if err != nil {
		return fmt.Errorf("failed to log inventory change: %w", err)
	}

	return nil
}

// LogFinancialTransaction creates a financial transaction log entry
func (r *RepositoryImpl) LogFinancialTransaction(userID string, amount float64, transactionType string, orderID string, description string, tx *sql.Tx) error {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	_, err := executor.Exec(`
        Insert Into Financial_Transaction_Log (
            User_Id, Amount, Type, Purchase_Order_Id,
            Description, Transaction_Date
        )
        Values ($1, $2, $3, $4, $5, Now())
    `, userID, amount, transactionType, orderID, description)

	if err != nil {
		return fmt.Errorf("failed to log financial transaction: %w", err)
	}

	return nil
}

// CreateReturn processes a purchase order return
func (r *RepositoryImpl) CreateReturn(req CreatePurchaseOrderReturnRequest, userID string, tx *sql.Tx) error {
	var executor interface {
		QueryRow(string, ...interface{}) *sql.Row
		Exec(string, ...interface{}) (sql.Result, error)
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Insert return record
	var returnID string
	err := executor.QueryRow(`
        Insert Into Purchase_Order_Return (
            Purchase_Order_Id, Product_Detail_Id, Return_Quantity,
            Reason, Status, Returned_By
        )
        Values ($1, $2, $3, $4, $5, $6)
        Returning Id
    `, req.PurchaseOrderID, req.ProductDetailID, req.ReturnQuantity,
		req.Reason, "returned", userID).Scan(&returnID)

	if err != nil {
		return fmt.Errorf("failed to create return: %w", err)
	}

	// Process batch returns
	for _, batch := range req.Batches {
		_, err = executor.Exec(`
            Insert Into Purchase_Order_Return_Batch (
                Purchase_Return_Id, Batch_Id, Quantity
            )
            Values ($1, $2, $3)
        `, returnID, batch.BatchID, batch.Quantity)

		if err != nil {
			return fmt.Errorf("failed to process batch return: %w", err)
		}

		// Update batch current quantity
		_, err = executor.Exec(`
            Update Product_Batch
            Set Current_Quantity = Current_Quantity - $1, Updated_At = Now()
            Where Id = $2
        `, batch.Quantity, batch.BatchID)

		if err != nil {
			return fmt.Errorf("failed to update batch quantity: %w", err)
		}

		// Update batch storage quantity
		_, err = executor.Exec(`
            Update Batch_Storage
            Set Quantity = Quantity - $1, Updated_At = Now()
            Where Batch_Id = $2 And Storage_Id = $3
        `, batch.Quantity, batch.BatchID, batch.StorageID)

		if err != nil {
			return fmt.Errorf("failed to update storage quantity: %w", err)
		}
	}

	// Check if all items are returned and update PO status
	allReturned, err := r.CheckAllItemsReturned(req.PurchaseOrderID, tx)
	if err != nil {
		return fmt.Errorf("failed to check returned status: %w", err)
	}

	var status string
	if allReturned {
		status = "returned"
	} else {
		status = "partially_returned"
	}

	// Update purchase order status
	_, err = executor.Exec(`
        Update Purchase_Order
        Set Status = $1, Updated_At = Now()
        Where Id = $2
    `, status, req.PurchaseOrderID)

	if err != nil {
		return fmt.Errorf("failed to update purchase order status: %w", err)
	}

	// Get order info for logging
	var serialID string
	var totalAmount float64
	err = executor.QueryRow(`
        Select Serial_Id, Total_Amount * ($1 / (
            Select Sum(Requested_Quantity) 
            From Purchase_Order_Detail 
            Where Purchase_Order_Id = $2
        ))
        From Purchase_Order
        Where Id = $2
    `, req.ReturnQuantity, req.PurchaseOrderID).Scan(&serialID, &totalAmount)

	if err != nil {
		return fmt.Errorf("failed to get order info: %w", err)
	}

	// Log financial transaction (negative amount for return)
	description := fmt.Sprintf("Return Pembelian %s", serialID)
	err = r.LogFinancialTransaction(userID, -totalAmount, "purchase_return", req.PurchaseOrderID, description, tx)
	if err != nil {
		return fmt.Errorf("failed to log financial transaction: %w", err)
	}

	// Log inventory changes for each batch
	for _, batch := range req.Batches {
		err = r.LogInventoryChange(
			batch.BatchID, batch.StorageID, userID, req.PurchaseOrderID,
			"return", batch.Quantity, description, tx,
		)
		if err != nil {
			return fmt.Errorf("failed to log inventory change: %w", err)
		}
	}

	return nil
}

// CancelReturn cancels a purchase order return
func (r *RepositoryImpl) CancelReturn(id string, userID string, tx *sql.Tx) error {
	var executor interface {
		QueryRow(string, ...interface{}) *sql.Row
		Exec(string, ...interface{}) (sql.Result, error)
		Query(string, ...interface{}) (*sql.Rows, error)
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Get return info
	var orderID string
	var returnQuantity float64
	err := executor.QueryRow(`
        Select Purchase_Order_Id, Return_Quantity
        From Purchase_Order_Return
        Where Id = $1 And Status = 'returned'
    `, id).Scan(&orderID, &returnQuantity)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("return not found or already cancelled")
		}
		return fmt.Errorf("failed to get return info: %w", err)
	}

	// Update return status
	_, err = executor.Exec(`
        Update Purchase_Order_Return
        Set Status = 'cancelled', Cancelled_By = $1, Cancelled_At = Now()
        Where Id = $2
    `, userID, id)

	if err != nil {
		return fmt.Errorf("failed to cancel return: %w", err)
	}

	// Get batch info for inventory restoration
	rows, err := executor.Query(`
        Select Prb.Batch_Id, Prb.Quantity, Bs.Storage_Id
        From Purchase_Order_Return_Batch Prb
        Join Batch_Storage Bs On Prb.Batch_Id = Bs.Batch_Id
        Where Prb.Purchase_Return_Id = $1
    `, id)

	if err != nil {
		return fmt.Errorf("failed to get batch info: %w", err)
	}
	defer rows.Close()

	// Restore batch quantities
	for rows.Next() {
		var batchID, storageID string
		var quantity float64

		err := rows.Scan(&batchID, &quantity, &storageID)
		if err != nil {
			return fmt.Errorf("failed to scan batch info: %w", err)
		}

		// Update batch current quantity
		_, err = executor.Exec(`
            Update Product_Batch
            Set Current_Quantity = Current_Quantity + $1, Updated_At = Now()
            Where Id = $2
        `, quantity, batchID)

		if err != nil {
			return fmt.Errorf("failed to update batch quantity: %w", err)
		}

		// Update batch storage quantity
		_, err = executor.Exec(`
            Update Batch_Storage
            Set Quantity = Quantity + $1, Updated_At = Now()
            Where Batch_Id = $2 And Storage_Id = $3
        `, quantity, batchID, storageID)

		if err != nil {
			return fmt.Errorf("failed to update storage quantity: %w", err)
		}
	}

	// Check for any active returns and update PO status
	var activeReturns int
	err = executor.QueryRow(`
        Select Count(*)
        From Purchase_Order_Return
        Where Purchase_Order_Id = $1 And Status = 'returned'
    `, orderID).Scan(&activeReturns)

	if err != nil {
		return fmt.Errorf("failed to check active returns: %w", err)
	}

	var newStatus string
	if activeReturns > 0 {
		newStatus = "partially_returned"
	} else {
		newStatus = "completed"
	}

	// Update purchase order status
	_, err = executor.Exec(`
        Update Purchase_Order
        Set Status = $1, Updated_At = Now()
        Where Id = $2
    `, newStatus, orderID)

	if err != nil {
		return fmt.Errorf("failed to update purchase order status: %w", err)
	}

	// Get order info for logging
	var serialID string
	var totalAmount float64
	err = executor.QueryRow(`
        Select Serial_Id, Total_Amount * ($1 / (
            Select Sum(Requested_Quantity) 
            From Purchase_Order_Detail 
            Where Purchase_Order_Id = $2
        ))
        From Purchase_Order
        Where Id = $2
    `, returnQuantity, orderID).Scan(&serialID, &totalAmount)

	if err != nil {
		return fmt.Errorf("failed to get order info: %w", err)
	}

	// Log financial transaction (positive amount for return cancellation)
	description := fmt.Sprintf("Batal Retur Pembelian %s", serialID)
	err = r.LogFinancialTransaction(userID, totalAmount, "purchase_return_cancel", orderID, description, tx)
	if err != nil {
		return fmt.Errorf("failed to log financial transaction: %w", err)
	}

	// Log inventory changes for each batch
	rows, err = executor.Query(`
        Select Prb.Batch_Id, Prb.Quantity, Bs.Storage_Id
        From Purchase_Order_Return_Batch Prb
        Join Batch_Storage Bs On Prb.Batch_Id = Bs.Batch_Id
        Where Prb.Purchase_Return_Id = $1
    `, id)

	if err != nil {
		return fmt.Errorf("failed to get batch info: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var batchID, storageID string
		var quantity float64

		err := rows.Scan(&batchID, &quantity, &storageID)
		if err != nil {
			return fmt.Errorf("failed to scan batch info: %w", err)
		}

		err = r.LogInventoryChange(
			batchID, storageID, userID, orderID,
			"return_cancel", quantity, description, tx,
		)
		if err != nil {
			return fmt.Errorf("failed to log inventory change: %w", err)
		}
	}

	return nil
}

// CheckAllItemsReturned checks if all items in a purchase order have been returned
func (r *RepositoryImpl) CheckAllItemsReturned(orderID string, tx *sql.Tx) (bool, error) {
	var executor interface {
		QueryRow(string, ...interface{}) *sql.Row
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	var totalOrdered, totalReturned float64

	// Get total ordered quantity
	err := executor.QueryRow(`
        Select Coalesce(Sum(Requested_Quantity), 0)
        From Purchase_Order_Detail
        Where Purchase_Order_Id = $1
    `, orderID).Scan(&totalOrdered)

	if err != nil {
		return false, fmt.Errorf("failed to get total ordered quantity: %w", err)
	}

	// Get total returned quantity (only for active returns)
	err = executor.QueryRow(`
        Select Coalesce(Sum(Return_Quantity), 0)
        From Purchase_Order_Return
        Where Purchase_Order_Id = $1 And Status = 'returned'
    `, orderID).Scan(&totalReturned)

	if err != nil {
		return false, fmt.Errorf("failed to get total returned quantity: %w", err)
	}

	// Check if all items are returned (with small epsilon for floating point comparison)
	return (totalOrdered - totalReturned) < 0.001, nil
}

// Update updates a purchase order
func (r *RepositoryImpl) Update(req UpdatePurchaseOrderRequest, tx *sql.Tx) (string, error) {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
		QueryRow(string, ...interface{}) *sql.Row
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Build update query
	query := `Update Purchase_Order Set Updated_At = Now()`
	params := []interface{}{}
	paramCount := 1

	// Add conditional fields
	if req.SupplierID != "" {
		query += fmt.Sprintf(", Supplier_Id = $%d", paramCount)
		params = append(params, req.SupplierID)
		paramCount++
	}

	if req.OrderDate != "" {
		query += fmt.Sprintf(", Order_Date = $%d", paramCount)
		params = append(params, req.OrderDate)
		paramCount++
	}

	if req.PaymentMethod != "" {
		query += fmt.Sprintf(", Payment_Method = $%d", paramCount)
		params = append(params, req.PaymentMethod)
		paramCount++
	}

	if req.PaymentMethod == "cash" {
		query += ", Payment_Due_Date = NULL"
	} else if req.PaymentMethod == "credit" {
		if req.PaymentDueDate == "" {
			return "", fmt.Errorf("tanggal jatuh tempo tidak boleh kosong untuk metode pembayaran kredit")
		}
		query += fmt.Sprintf(", Payment_Due_Date = $%d", paramCount)
		params = append(params, req.PaymentDueDate)
		paramCount++
	}

	if req.CheckedBy != "" {
		query += fmt.Sprintf(", Checked_By = $%d", paramCount)
		params = append(params, req.CheckedBy)
		paramCount++
	}

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE Id = $%d", paramCount)
	params = append(params, req.ID)

	// Execute query
	result, err := executor.Exec(query, params...)
	if err != nil {
		return "", fmt.Errorf("failed to update purchase order: %w", err)
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return "", fmt.Errorf("purchase order tidak ditemukan")
	}

	return req.ID, nil
}

// UpdateStatus updates the status of a purchase order
func (r *RepositoryImpl) UpdateStatus(id, status, userID string, tx *sql.Tx) error {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Validate status
	validStatuses := map[string]bool{
		"order":              true,
		"completed":          true,
		"partially_returned": true,
		"returned":           true,
		"cancelled":          true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Update query
	updateQuery := `
        Update Purchase_Order
        Set Status = $1, Updated_At = Now()
        Where Id = $2
    `

	// For cancelled status, also update cancelled_by and cancelled_at
	if status == "cancelled" {
		updateQuery = `
            Update Purchase_Order
            Set Status = $1, Updated_At = Now(), 
                Cancelled_By = $3, Cancelled_At = Now()
            Where Id = $2
        `
	}

	result, err := executor.Exec(updateQuery, status, id, userID)
	if err != nil {
		return fmt.Errorf("failed to update purchase order status: %w", err)
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order not found")
	}

	return nil
}

// CheckPurchaseOrder marks a purchase order as checked by the given user
func (r *RepositoryImpl) CheckPurchaseOrder(id string, userID string, tx *sql.Tx) error {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
		QueryRow(string, ...interface{}) *sql.Row
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Check if purchase order exists
	var status string
	err := executor.QueryRow(`
        Select Status From Purchase_Order Where Id = $1
    `, id).Scan(&status)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("purchase order not found")
		}
		return fmt.Errorf("failed to check purchase order: %w", err)
	}

	// Only allow checking of orders in 'order' status
	if status != "order" {
		return fmt.Errorf("only purchase orders with status 'order' can be checked")
	}

	// Update purchase order
	result, err := executor.Exec(`
        Update Purchase_Order
        Set Checked_By = $1, Updated_At = Now()
        Where Id = $2
    `, userID, id)

	if err != nil {
		return fmt.Errorf("failed to update purchase order: %w", err)
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order not found")
	}

	return nil
}

// CancelPurchaseOrder cancels a purchase order
func (r *RepositoryImpl) CancelPurchaseOrder(id string, userID string, tx *sql.Tx) error {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
		QueryRow(string, ...interface{}) *sql.Row
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Check if purchase order exists and get its status
	var status string
	var serialID string
	err := executor.QueryRow(`
        Select Status, Serial_Id From Purchase_Order Where Id = $1
    `, id).Scan(&status, &serialID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("purchase order not found")
		}
		return fmt.Errorf("failed to check purchase order: %w", err)
	}

	// Only allow cancellation of orders in 'order' status
	if status != "order" {
		return fmt.Errorf("only purchase orders with status 'order' can be cancelled")
	}

	// Update purchase order
	result, err := executor.Exec(`
        Update Purchase_Order
        Set Status = 'cancelled', Cancelled_By = $1, 
            Cancelled_At = Now(), Updated_At = Now()
        Where Id = $2
    `, userID, id)

	if err != nil {
		return fmt.Errorf("failed to cancel purchase order: %w", err)
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order not found")
	}

	// Log the cancellation
	description := fmt.Sprintf("Pembatalan Pesanan Pembelian %s", serialID)
	err = r.LogFinancialTransaction(userID, 0, "purchase_cancel", id, description, tx)
	if err != nil {
		return fmt.Errorf("failed to log financial transaction: %w", err)
	}

	return nil
}

// AddPurchaseOrderItem adds an item to a purchase order
func (r *RepositoryImpl) AddPurchaseOrderItem(orderID string, req CreatePurchaseOrderItemRequest, tx *sql.Tx) error {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
		QueryRow(string, ...interface{}) *sql.Row
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Check if purchase order exists and its status
	var status string
	err := executor.QueryRow(`
        Select Status From Purchase_Order Where Id = $1
    `, orderID).Scan(&status)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("purchase order not found")
		}
		return fmt.Errorf("failed to check purchase order status: %w", err)
	}

	// Only allow adding items to orders in 'order' status
	if status != "order" {
		return fmt.Errorf("can only add items to purchase orders with status 'order'")
	}

	// Add the item
	_, err = executor.Exec(`
        Insert Into Purchase_Order_Detail (
            Purchase_Order_Id, Product_Id, 
            Requested_Quantity, Unit_Price
        )
        Values ($1, $2, $3, $4)
    `, orderID, req.ProductID, req.Quantity, req.Price)

	if err != nil {
		return fmt.Errorf("failed to add order item: %w", err)
	}

	// Update total amount
	_, err = executor.Exec(`
        Update Purchase_Order
        Set Total_Amount = Total_Amount + $1, Updated_At = Now()
        Where Id = $2
    `, req.Quantity*req.Price, orderID)

	if err != nil {
		return fmt.Errorf("failed to update order total amount: %w", err)
	}

	return nil
}

// UpdatePurchaseOrderItem updates a purchase order item
func (r *RepositoryImpl) UpdatePurchaseOrderItem(req UpdatePurchaseOrderItemRequest, tx *sql.Tx) error {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
		QueryRow(string, ...interface{}) *sql.Row
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Get current item details
	var orderID string
	var oldQuantity, oldPrice float64
	err := executor.QueryRow(`
        Select Purchase_Order_Id, Requested_Quantity, Unit_Price
        From Purchase_Order_Detail
        Where Id = $1
    `, req.ID).Scan(&orderID, &oldQuantity, &oldPrice)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("purchase order item not found")
		}
		return fmt.Errorf("failed to get current item details: %w", err)
	}

	// Check order status
	var status string
	err = executor.QueryRow(`
        Select Status From Purchase_Order Where Id = $1
    `, orderID).Scan(&status)

	if err != nil {
		return fmt.Errorf("failed to check purchase order status: %w", err)
	}

	// Only allow updating items for orders in 'order' status
	if status != "order" {
		return fmt.Errorf("can only update items for purchase orders with status 'order'")
	}

	// Update the item
	result, err := executor.Exec(`
        Update Purchase_Order_Detail
        Set Requested_Quantity = $1, Unit_Price = $2, Updated_At = Now()
        Where Id = $3
    `, req.Quantity, req.Price, req.ID)

	if err != nil {
		return fmt.Errorf("failed to update order item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order item not found")
	}

	// Calculate the change in total amount
	oldTotal := oldQuantity * oldPrice
	newTotal := req.Quantity * req.Price
	amountDiff := newTotal - oldTotal

	// Update the order's total amount
	_, err = executor.Exec(`
        Update Purchase_Order
        Set Total_Amount = Total_Amount + $1, Updated_At = Now()
        Where Id = $2
    `, amountDiff, orderID)

	if err != nil {
		return fmt.Errorf("failed to update order total amount: %w", err)
	}

	return nil
}

// RemovePurchaseOrderItem removes an item from a purchase order
func (r *RepositoryImpl) RemovePurchaseOrderItem(id string, tx *sql.Tx) error {
	var executor interface {
		Exec(string, ...interface{}) (sql.Result, error)
		QueryRow(string, ...interface{}) *sql.Row
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Get current item details
	var orderID string
	var quantity, price float64
	err := executor.QueryRow(`
        Select Purchase_Order_Id, Requested_Quantity, Unit_Price
        From Purchase_Order_Detail
        Where Id = $1
    `, id).Scan(&orderID, &quantity, &price)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("purchase order item not found")
		}
		return fmt.Errorf("failed to get current item details: %w", err)
	}

	// Check order status
	var status string
	err = executor.QueryRow(`
        Select Status From Purchase_Order Where Id = $1
    `, orderID).Scan(&status)

	if err != nil {
		return fmt.Errorf("failed to check purchase order status: %w", err)
	}

	// Only allow removing items from orders in 'order' status
	if status != "order" {
		return fmt.Errorf("can only remove items from purchase orders with status 'order'")
	}

	// Remove the item
	result, err := executor.Exec(`
        Delete From Purchase_Order_Detail
        Where Id = $1
    `, id)

	if err != nil {
		return fmt.Errorf("failed to remove order item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order item not found")
	}

	// Update the order's total amount
	total := quantity * price
	_, err = executor.Exec(`
        Update Purchase_Order
        Set Total_Amount = Total_Amount - $1, Updated_At = Now()
        Where Id = $2
    `, total, orderID)

	if err != nil {
		return fmt.Errorf("failed to update order total amount: %w", err)
	}

	return nil
}

// GetAll fetches all purchase orders based on search criteria
func (r *RepositoryImpl) GetAll(req GetPurchaseOrderRequest) ([]GetPurchaseOrderResponse, int, error) {
	// Base query for counting total items
	countQuery := `
        Select Count(*)
        From Purchase_Order Po
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Where 1=1
    `

	// Base query for fetching items
	fetchQuery := `
        Select 
            Po.Id, Po.Serial_Id, Po.Supplier_Id, S.Name As Suppliername,
            Po.Order_Date, Po.Status, Po.Total_Amount,
            U.Username, Po.Created_At, Po.Updated_At,
            (Select Count(*) From Purchase_Order_Detail Pod Where Pod.Purchase_Order_Id = Po.Id) As Itemcount
        From Purchase_Order Po
		Left Join Appuser U On Po.Created_By = U.Id
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Where 1=1
    `

	// Build query with filters
	qb := utils.NewQueryBuilder(fetchQuery)
	countQb := utils.NewQueryBuilder(countQuery)

	if req.SupplierID != "" {
		qb.AddFilter("po.Supplier_Id =", req.SupplierID)
		countQb.AddFilter("po.Supplier_Id =", req.SupplierID)
	}

	if req.Status != "" {
		qb.AddFilter("po.Status =", req.Status)
		countQb.AddFilter("po.Status =", req.Status)
	}

	if req.FromDate != "" {
		fromDate, err := time.Parse("2006-01-02", req.FromDate)
		if err == nil {
			qb.AddFilter("po.Order_Date >=", fromDate)
			countQb.AddFilter("po.Order_Date >=", fromDate)
		}
	}

	if req.ToDate != "" {
		toDate, err := time.Parse("2006-01-02", req.ToDate)
		if err == nil {
			// Add one day to include the end date
			toDate = toDate.Add(24 * time.Hour)
			qb.AddFilter("po.Order_Date <", toDate)
			countQb.AddFilter("po.Order_Date <", toDate)
		}
	}

	// Add order by
	sortField := "po.Created_At"
	if req.SortBy != "" {
		// Map frontend field names to database column names
		sortFieldMap := map[string]string{
			"serialId":     "po.Serial_Id",
			"supplierName": "s.Name",
			"orderDate":    "po.Order_Date",
			"status":       "po.Status",
			"totalAmount":  "po.Total_Amount",
			"createdAt":    "po.Created_At",
		}

		if mappedField, ok := sortFieldMap[req.SortBy]; ok {
			sortField = mappedField
		}
	}

	sortOrder := "DESC"
	if strings.ToUpper(req.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	qb.Query.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder))

	// Add pagination
	qb.AddPagination(req.PageSize, req.Page)

	// Execute count query
	var totalItems int
	countQuery, countParams := countQb.Build()
	err := r.DB.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count purchase orders: %w", err)
	}

	// Execute fetch query
	query, params := qb.Build()
	rows, err := r.DB.Query(query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch purchase orders: %w", err)
	}
	defer rows.Close()

	// Parse results
	var orders []GetPurchaseOrderResponse
	for rows.Next() {
		var order GetPurchaseOrderResponse
		var supplierID, supplierName sql.NullString

		err := rows.Scan(
			&order.ID, &order.SerialID, &supplierID, &supplierName,
			&order.OrderDate, &order.Status, &order.TotalAmount,
			&order.CreatedBy, &order.CreatedAt, &order.UpdatedAt,
			&order.ItemCount,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan purchase order: %w", err)
		}

		if supplierID.Valid {
			order.SupplierID = supplierID.String
		}

		if supplierName.Valid {
			order.SupplierName = supplierName.String
		}

		orders = append(orders, order)
	}

	return orders, totalItems, nil
}

// GetAllReturns fetches all purchase order returns
func (r *RepositoryImpl) GetAllReturns(req GetPurchaseOrderReturnRequest) ([]GetPurchaseOrderReturnResponse, int, error) {
	// Base query for counting total items
	countQuery := `
        Select Count(*)
        From Purchase_Order_Return Por
        Join Purchase_Order Po On Por.Purchase_Order_Id = Po.Id
        Join Purchase_Order_Detail Pod On Por.Product_Detail_Id = Pod.Id
        Join Product P On Pod.Product_Id = P.Id
        Join Appuser U On Por.Returned_By = U.Id
        Where 1=1
    `

	// Base query for fetching items
	fetchQuery := `
        Select 
            Por.Id, Por.Purchase_Order_Id, Po.Serial_Id,
            Pod.Product_Id, P.Name As Productname,
            Por.Return_Quantity, Por.Reason, Por.Status,
            Por.Returned_By, U.Username As Returnedbyname,
            Por.Returned_At
        From Purchase_Order_Return Por
        Join Purchase_Order Po On Por.Purchase_Order_Id = Po.Id
        Join Purchase_Order_Detail Pod On Por.Product_Detail_Id = Pod.Id
        Join Product P On Pod.Product_Id = P.Id
        Join Appuser U On Por.Returned_By = U.Id
        Where 1=1
    `

	// Build query with filters
	qb := utils.NewQueryBuilder(fetchQuery)
	countQb := utils.NewQueryBuilder(countQuery)

	if req.FromDate != "" {
		fromDate, err := time.Parse("2006-01-02", req.FromDate)
		if err == nil {
			qb.AddFilter("por.Returned_At >=", fromDate)
			countQb.AddFilter("por.Returned_At >=", fromDate)
		}
	}

	if req.ToDate != "" {
		toDate, err := time.Parse("2006-01-02", req.ToDate)
		if err == nil {
			// Add one day to include the end date
			toDate = toDate.Add(24 * time.Hour)
			qb.AddFilter("por.Returned_At <", toDate)
			countQb.AddFilter("por.Returned_At <", toDate)
		}
	}

	// Add order by
	qb.Query.WriteString(" ORDER BY por.Returned_At DESC")

	// Add pagination
	qb.AddPagination(req.PageSize, req.Page)

	// Execute count query
	var totalItems int
	countQuery, countParams := countQb.Build()
	err := r.DB.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count purchase order returns: %w", err)
	}

	// Execute fetch query
	query, params := qb.Build()
	rows, err := r.DB.Query(query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch purchase order returns: %w", err)
	}
	defer rows.Close()

	// Parse results
	var returns []GetPurchaseOrderReturnResponse
	for rows.Next() {
		var ret GetPurchaseOrderReturnResponse

		err := rows.Scan(
			&ret.ID, &ret.PurchaseOrderID, &ret.SerialID,
			&ret.ProductID, &ret.ProductName,
			&ret.ReturnQuantity, &ret.Reason, &ret.Status,
			&ret.ReturnedBy, &ret.ReturnedByName,
			&ret.ReturnedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan purchase order return: %w", err)
		}

		returns = append(returns, ret)
	}

	return returns, totalItems, nil
}

// GenerateBatchSKU generates a SKU for product batch
func (r *RepositoryImpl) GenerateBatchSKU(productName, supplierName, serialID string, date time.Time) (string, error) {
	// Product abbreviation (first 3 chars)
	productAbbr := utils.GetAbbreviation(productName, 3)

	// Supplier abbreviation
	supplierAbbr := utils.GetSupplierAbbreviation(supplierName)

	// Format date: DDMMYY
	dateStr := fmt.Sprintf("%02d%02d%02d", date.Day(), date.Month(), date.Year()%100)

	// Extract iteration number from serial ID (PO-20250417-XXXX)
	iterNumber := "0000"
	parts := strings.Split(serialID, "-")
	if len(parts) == 3 {
		iterNumber = parts[2]
	}

	// Format: PROD-SUPPDDMMYY-XXXX (XXXX is the PO iteration number)
	sku := fmt.Sprintf("%s-%s%s-%s", productAbbr, supplierAbbr, dateStr, iterNumber)

	return strings.ToUpper(sku), nil
}

// GetProducts fetches products based on search criteria
func (r *RepositoryImpl) GetProducts(req product.GetProductRequest) ([]product.GetProductResponse, int, error) {
	// Base query for counting total items
	countQuery := `
		Select Count(*)
		From Product P
		Join Supplier S On P.Supplier_Id = S.Id
		Where 1=1
	`

	// Base query for fetching items
	fetchQuery := `
		Select 
			P.Id, P.Name, P.Category_Id, C.Name As Category,
			P.Unit_Id, U.Name As Unit,
			P.Created_At, P.Updated_At
		From Product P
		Join Category C On P.Category_Id = C.Id
		Join Unit U On P.Unit_Id = U.Id
		Where 1=1
	`

	// Build query with filters
	qb := utils.NewQueryBuilder(fetchQuery)
	countQb := utils.NewQueryBuilder(countQuery)

	if req.Name != "" {
		qb.AddFilter("p.Name ILIKE %?%", req.Name)
		countQb.AddFilter("p.Name ILIKE %?%", req.Name)
	}

	// Query execution
	query, params := qb.Build()
	rows, err := r.DB.Query(query, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Parse results
	var products []product.GetProductResponse
	for rows.Next() {
		var prod product.GetProductResponse
		err := rows.Scan(
			&prod.ID, &prod.Name, &prod.CategoryID, &prod.Category,
			&prod.UnitID, &prod.Unit,
			&prod.CreatedAt, &prod.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		products = append(products, prod)
	}

	// Execute count query
	var totalItems int
	countQuery, countParams := countQb.Build()
	err = r.DB.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}
	return products, totalItems, nil
}

// CompleteFullPurchaseOrder processes all items in a purchase order automatically
func (r *RepositoryImpl) CompleteFullPurchaseOrder(id string, storageID string, userID string, tx *sql.Tx) error {
	var executor interface {
		QueryRow(string, ...interface{}) *sql.Row
		Exec(string, ...interface{}) (sql.Result, error)
		Query(string, ...interface{}) (*sql.Rows, error)
	}

	if tx != nil {
		executor = tx
	} else {
		executor = r.DB
	}

	// Get purchase order info (supplier, serial ID)
	var serialID, supplierID, supplierName string
	err := executor.QueryRow(`
        Select Po.Serial_Id, Po.Supplier_Id, S.Name
        From Purchase_Order Po
        Join Supplier S On Po.Supplier_Id = S.Id
        Where Po.Id = $1
    `, id).Scan(&serialID, &supplierID, &supplierName)

	if err != nil {
		return fmt.Errorf("failed to get purchase order info: %w", err)
	}

	// Update purchase order status
	_, err = executor.Exec(`
        Update Purchase_Order
        Set Status = 'completed', Checked_By = $1, Updated_At = Now()
        Where Id = $2
    `, userID, id)

	if err != nil {
		return fmt.Errorf("failed to update purchase order status: %w", err)
	}

	now := time.Now()

	// First, collect all the items we need to process
	rows, err := executor.Query(`
        Select Pod.Id, Pod.Product_Id, P.Name, Pod.Requested_Quantity, Pod.Unit_Price
        From Purchase_Order_Detail Pod
        Join Product P On Pod.Product_Id = P.Id
        Where Pod.Purchase_Order_Id = $1
    `, id)

	if err != nil {
		return fmt.Errorf("failed to get purchase order details: %w", err)
	}

	// Create a structure to store all the items we need to process
	type orderItem struct {
		detailID    string
		productID   string
		productName string
		quantity    float64
		unitPrice   float64
	}

	var items []orderItem

	// Read all rows into our slice
	for rows.Next() {
		var item orderItem
		err := rows.Scan(&item.detailID, &item.productID, &item.productName, &item.quantity, &item.unitPrice)
		if err != nil {
			rows.Close() // Important to close in case of error
			return fmt.Errorf("failed to scan purchase order detail: %w", err)
		}
		items = append(items, item)
	}

	// Close the rows now that we've read all the data
	rows.Close()

	// Check for errors from iterating rows
	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating through purchase order details: %w", err)
	}

	// Now process each item without the connection being used for row iteration
	for _, item := range items {
		// Generate SKU with the serialID
		sku, err := r.GenerateBatchSKU(item.productName, supplierName, serialID, now)
		if err != nil {
			return fmt.Errorf("failed to generate SKU: %w", err)
		}

		// Create product batch
		batchID, err := r.CreateProductBatch(item.productID, id, sku, item.quantity, item.unitPrice, tx)
		if err != nil {
			return fmt.Errorf("failed to create product batch: %w", err)
		}

		// Associate batch with storage
		err = r.AssignBatchToStorage(*batchID, storageID, item.quantity, tx)
		if err != nil {
			return fmt.Errorf("failed to assign batch to storage: %w", err)
		}

		// Log inventory change
		description := fmt.Sprintf("Pembelian Barang %s", serialID)
		err = r.LogInventoryChange(*batchID, storageID, userID, id, "add", item.quantity, description, tx)
		if err != nil {
			return fmt.Errorf("failed to log inventory change: %w", err)
		}
	}

	// Get total amount for financial log
	var totalAmount float64
	err = executor.QueryRow(`
        Select Total_Amount From Purchase_Order Where Id = $1
    `, id).Scan(&totalAmount)

	if err != nil {
		return fmt.Errorf("failed to get total amount: %w", err)
	}

	// Log financial transaction
	description := fmt.Sprintf("Pembelian Barang %s", serialID)
	err = r.LogFinancialTransaction(userID, totalAmount, "purchase", id, description, tx)
	if err != nil {
		return fmt.Errorf("failed to log financial transaction: %w", err)
	}

	return nil
}
