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

// PurchaseOrderRepository interface defines methods for purchase order operations
type PurchaseOrderRepository interface {
	// Basic CRUD operations
	GetAll(req GetPurchaseOrderRequest) ([]GetPurchaseOrderResponse, int, error)
	GetByID(id string) (*PurchaseOrderDetailResponse, error)
	Create(req CreatePurchaseOrderRequest, userID string) error
	Update(req UpdatePurchaseOrderRequest) error
	Delete(id string) error

	// Order processing operations
	ReceiveOrder(orderID string, receivedItems []ReceivedItemRequest, userID string) error
	CheckOrder(orderID string, userID string) error
	CancelOrder(orderID string, userID string) error

	// Return operations
	CreateReturn(returnRequest CreatePurchaseOrderReturnRequest, userID string) error
	GetAllReturns(req GetPurchaseOrderReturnRequest) ([]GetPurchaseOrderReturnResponse, int, error)
	GetReturnByID(returnID string) (*PurchaseOrderReturnDetailResponse, error)
	CancelReturn(returnID string, userID string) error

	// Order item operations
	GetByOrderID(orderID string) ([]GetPurchaseOrderItemResponse, error)
	AddItem(orderID string, req CreatePurchaseOrderItemRequest) error
	UpdateItem(req UpdatePurchaseOrderItemRequest) error
	RemoveItem(itemID string) error
}

// PurchaseOrderRepositoryImpl implements PurchaseOrderRepository
type PurchaseOrderRepositoryImpl struct {
	db *sql.DB
}

// NewPurchaseOrderRepository creates a new instance of PurchaseOrderRepositoryImpl
func NewPurchaseOrderRepository(db *sql.DB) PurchaseOrderRepository {
	return &PurchaseOrderRepositoryImpl{db: db}
}

// GetAll fetches all purchase orders with filtering and pagination
func (r *PurchaseOrderRepositoryImpl) GetAll(req GetPurchaseOrderRequest) ([]GetPurchaseOrderResponse, int, error) {
	// Build count query
	countBuilder := utils.NewQueryBuilder(`
        Select Count(Po.Id)
        From Purchase_Order Po
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Where 1=1
    `)

	if req.SupplierID != "" {
		countBuilder.AddFilter("s.Id =", req.SupplierID)
	}

	if req.Status != "" {
		countBuilder.AddFilter("po.Status =", req.Status)
	}

	if req.FromDate != "" {
		countBuilder.AddFilter("po.Order_Date >=", req.FromDate)
	}

	if req.ToDate != "" {
		countBuilder.AddFilter("po.Order_Date <=", req.ToDate)
	}

	countQuery, countParams := countBuilder.Build()

	// Execute count query
	var totalItems int
	err := r.db.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung jumlah pesanan: %w", err)
	}

	// Build main query
	queryBuilder := utils.NewQueryBuilder(`
        Select Po.Id, Po.Supplier_Id, S.Name As Suppliername,
               Po.Order_Date, Po.Status, Po.Total_Amount, 
               U.Username As Createdby, Po.Created_At, Po.Updated_At
        From Purchase_Order Po
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Left Join Appuser U On Po.Created_By = U.Id    
        Where 1=1`)

	if req.SupplierID != "" {
		queryBuilder.AddFilter("s.Id =", req.SupplierID)
	}

	if req.Status != "" {
		queryBuilder.AddFilter("po.Status =", req.Status)
	}

	if req.FromDate != "" {
		queryBuilder.AddFilter("po.Order_Date >=", req.FromDate)
	}

	if req.ToDate != "" {
		queryBuilder.AddFilter("po.Order_Date <=", req.ToDate)
	}

	// Add sorting
	queryBuilder.Query.WriteString(" ORDER BY po.Order_Date DESC")

	// Add pagination
	mainQuery, queryParams := queryBuilder.AddPagination(req.PageSize, req.Page).Build()

	fmt.Println(queryParams)
	// Execute main query
	rows, err := r.db.Query(mainQuery, queryParams...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar pesanan: %w", err)
	}
	defer rows.Close()

	// Process results
	var orders []GetPurchaseOrderResponse
	for rows.Next() {
		var order GetPurchaseOrderResponse
		err := rows.Scan(
			&order.ID,
			&order.SupplierID,
			&order.SupplierName,
			&order.OrderDate,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedBy,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal memproses data pesanan: %w", err)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("gagal saat iterasi hasil query: %w", err)
	}

	return orders, totalItems, nil
}

// GetByID fetches a purchase order by ID with its details
func (r *PurchaseOrderRepositoryImpl) GetByID(id string) (*PurchaseOrderDetailResponse, error) {
	// Fetch the purchase order header
	var order GetPurchaseOrderResponse

	queryErr := r.db.QueryRow(`
        Select Po.Id, Po.Supplier_Id, Coalesce(S.Name, 'Supplier Dihapus') As Suppliername,
               Po.Order_Date, Po.Status, Po.Total_Amount,
               U.Username As Createdby, Po.Created_At, Po.Updated_At
        From Purchase_Order Po
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Left Join Appuser U On Po.Created_By = U.Id
        Where Po.Id = $1
    `, id).Scan(
		&order.ID,
		&order.SupplierID,
		&order.SupplierName,
		&order.OrderDate,
		&order.Status,
		&order.TotalAmount,
		&order.CreatedBy,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			return nil, fmt.Errorf("pesanan pembelian tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil pesanan pembelian: %w", queryErr)
	}

	// Fetch the purchase order items
	rows, queryErr := r.db.Query(`
        Select Pod.Id, Pod.Product_Id, P.Name, Pod.Requested_Quantity, Pod.Unit_Price
        From Purchase_Order_Detail Pod
        Join Product P On Pod.Product_Id = P.Id
        Where Pod.Purchase_Order_Id = $1
    `, id)

	if queryErr != nil {
		return nil, fmt.Errorf("gagal mengambil detail pesanan: %w", queryErr)
	}
	defer rows.Close()

	var items []GetPurchaseOrderItemResponse
	for rows.Next() {
		var item GetPurchaseOrderItemResponse
		scanErr := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.Price,
		)
		if scanErr != nil {
			return nil, fmt.Errorf("gagal memproses detail pesanan: %w", scanErr)
		}
		item.Subtotal = item.Quantity * item.Price
		items = append(items, item)
	}

	if rowErr := rows.Err(); rowErr != nil {
		return nil, fmt.Errorf("gagal saat iterasi hasil query: %w", rowErr)
	}

	return &PurchaseOrderDetailResponse{
		GetPurchaseOrderResponse: order,
		Items:                    items,
	}, nil
}

// Create creates a new purchase order
func (r *PurchaseOrderRepositoryImpl) Create(req CreatePurchaseOrderRequest, userID string) error {
	var orderID string

	txErr := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Insert purchase order header
		var orderDate time.Time
		parseErr := orderDate.UnmarshalText([]byte(req.OrderDate))
		if parseErr != nil {
			return fmt.Errorf("format tanggal tidak valid: %w", parseErr)
		}

		// Calculate total amount from items
		totalAmount := 0.0
		for _, item := range req.Items {
			totalAmount += item.Quantity * item.Price
		}

		var createdAt time.Time
		insertErr := tx.QueryRow(`
            Insert Into Purchase_Order (Supplier_Id, Order_Date, Status, Total_Amount, Created_By)
            Values ($1, $2, $3, $4, $5)
            Returning Id, Created_At
        `, req.SupplierID, orderDate, req.Status, totalAmount, userID).Scan(&orderID, &createdAt)

		if insertErr != nil {
			return fmt.Errorf("gagal membuat pesanan pembelian: %w", insertErr)
		}

		// Insert purchase order items
		var items []GetPurchaseOrderItemResponse
		for _, item := range req.Items {
			var detailID string
			var productName string

			// Get product name
			productErr := tx.QueryRow(`
                Select Name From Product Where Id = $1
            `, item.ProductID).Scan(&productName)

			if productErr != nil {
				return fmt.Errorf("produk tidak ditemukan: %w", productErr)
			}

			// Insert item
			detailErr := tx.QueryRow(`
                Insert Into Purchase_Order_Detail (Purchase_Order_Id, Product_Id, Requested_Quantity, Unit_Price)
                Values ($1, $2, $3, $4)
                Returning Id
            `, orderID, item.ProductID, item.Quantity, item.Price).Scan(&detailID)

			if detailErr != nil {
				return fmt.Errorf("gagal menambahkan item pesanan: %w", detailErr)
			}

			// Add to response items
			items = append(items, GetPurchaseOrderItemResponse{
				ID:          detailID,
				ProductID:   item.ProductID,
				ProductName: productName,
				Quantity:    item.Quantity,
				Price:       item.Price,
				Subtotal:    item.Quantity * item.Price,
			})
		}

		// Get supplier name
		var supplierName string
		supplierErr := tx.QueryRow(`
            Select Name From Supplier Where Id = $1
        `, req.SupplierID).Scan(&supplierName)

		if supplierErr != nil {
			if errors.Is(supplierErr, sql.ErrNoRows) {
				supplierName = "Supplier Dihapus"
			} else {
				return fmt.Errorf("gagal mengambil informasi supplier: %w", supplierErr)
			}
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

// Update modifies an existing purchase order
func (r *PurchaseOrderRepositoryImpl) Update(req UpdatePurchaseOrderRequest) error {
	// Get current purchase order to keep unchanged fields
	var currentStatus string
	var supplierID sql.NullString
	var orderDate, updatedAt time.Time

	queryErr := r.db.QueryRow(`
        Select Status, Supplier_Id, Order_Date
        From Purchase_Order
        Where Id = $1
    `, req.ID).Scan(&currentStatus, &supplierID, &orderDate)

	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			return fmt.Errorf("pesanan pembelian tidak ditemukan")
		}
		return fmt.Errorf("gagal mengambil pesanan pembelian: %w", queryErr)
	}

	// Check if status change is allowed
	if req.Status != "" && req.Status != currentStatus {
		if !isStatusChangeAllowed(currentStatus, req.Status) {
			return fmt.Errorf("perubahan status dari %s ke %s tidak diizinkan", currentStatus, req.Status)
		}
	}

	// Prepare update values
	status := req.Status
	if status == "" {
		status = currentStatus
	}

	updateSupplierID := supplierID
	if req.SupplierID != "" {
		updateSupplierID = sql.NullString{
			String: req.SupplierID,
			Valid:  true,
		}
	}

	updateOrderDate := orderDate
	if req.OrderDate != "" {
		parseErr := updateOrderDate.UnmarshalText([]byte(req.OrderDate))
		if parseErr != nil {
			return fmt.Errorf("format tanggal tidak valid: %w", parseErr)
		}
	}

	// Update the purchase order
	updateErr := r.db.QueryRow(`
        Update Purchase_Order
        Set Supplier_Id = $1, Order_Date = $2, Status = $3, Updated_At = Current_Timestamp
        Where Id = $4
        Returning Updated_At
    `, updateSupplierID, updateOrderDate, status, req.ID).Scan(&updatedAt)

	if updateErr != nil {
		if errors.Is(updateErr, sql.ErrNoRows) {
			return fmt.Errorf("pesanan pembelian tidak ditemukan")
		}
		return fmt.Errorf("gagal memperbarui pesanan pembelian: %w", updateErr)
	}

	// Get updated purchase order
	var supplierName string
	var totalAmount float64
	var createdBy string
	var createdAt time.Time

	getErr := r.db.QueryRow(`
        Select Coalesce(S.Name, 'Supplier Dihapus') As Suppliername,
               Po.Total_Amount, U.Username, Po.Created_At
        From Purchase_Order Po
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Left Join Appuser U On Po.Created_By = U.Id
        Where Po.Id = $1
    `, req.ID).Scan(&supplierName, &totalAmount, &createdBy, &createdAt)

	if getErr != nil {
		return fmt.Errorf("gagal mengambil data pesanan yang diperbarui: %w", getErr)
	}

	return nil
}

// Delete soft deletes a purchase order
func (r *PurchaseOrderRepositoryImpl) Delete(id string) error {
	// Check if purchase order exists
	var status string
	queryErr := r.db.QueryRow(`
        Select Status From Purchase_Order Where Id = $1
    `, id).Scan(&status)

	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			return fmt.Errorf("pesanan pembelian tidak ditemukan")
		}
		return fmt.Errorf("gagal memeriksa status pesanan: %w", queryErr)
	}

	// Only allow deletion of 'ordered' status orders
	if status != "ordered" {
		return fmt.Errorf("hanya pesanan dengan status 'ordered' yang dapat dihapus, status saat ini: %s", status)
	}

	// Delete the order (using UPDATE for soft delete)
	result, execErr := r.db.Exec(`
        Update Purchase_Order
        Set Status = 'cancelled', Cancelled_At = Current_Timestamp
        Where Id = $1 And Status = 'ordered'
    `, id)

	if execErr != nil {
		return fmt.Errorf("gagal menghapus pesanan pembelian: %w", execErr)
	}

	affected, rowErr := result.RowsAffected()
	if rowErr != nil {
		return fmt.Errorf("gagal mendapatkan jumlah baris yang terpengaruh: %w", rowErr)
	}

	if affected == 0 {
		return fmt.Errorf("pesanan pembelian tidak ditemukan atau tidak dalam status yang dapat dihapus")
	}

	return nil
}

// Helper function to determine if status change is allowed
func isStatusChangeAllowed(currentStatus, newStatus string) bool {
	// Define allowed status transitions
	allowedTransitions := map[string][]string{
		"ordered":            {"received", "cancelled"},
		"received":           {"checked"},
		"checked":            {"completed", "partially_returned"},
		"partially_returned": {"completed", "returned"},
		"completed":          {},
		"returned":           {},
		"cancelled":          {},
	}

	// Check if transition is allowed
	allowed := false
	if transitions, exists := allowedTransitions[currentStatus]; exists {
		for _, allowedStatus := range transitions {
			if allowedStatus == newStatus {
				allowed = true
				break
			}
		}
	}

	return allowed
}

// ReceiveOrder marks a purchase order as received and creates product batches
func (r *PurchaseOrderRepositoryImpl) ReceiveOrder(orderID string, receivedItems []ReceivedItemRequest, userID string) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Check if order exists and is in 'ordered' status
		var status string
		orderErr := tx.QueryRow(`
            Select Status 
            From Purchase_Order 
            Where Id = $1
        `, orderID).Scan(&status)

		if orderErr != nil {
			if errors.Is(orderErr, sql.ErrNoRows) {
				return fmt.Errorf("pesanan pembelian tidak ditemukan")
			}
			return fmt.Errorf("gagal memeriksa status pesanan: %w", orderErr)
		}

		if status != "ordered" {
			return fmt.Errorf("hanya pesanan dengan status 'ordered' yang dapat diterima, status saat ini: %s", status)
		}

		// Set payment due date if provided with the first item
		var paymentDueDate *time.Time
		if len(receivedItems) > 0 && receivedItems[0].PaymentDueDate != "" {
			parsedDate, parseErr := time.Parse(time.RFC3339, receivedItems[0].PaymentDueDate)
			if parseErr != nil {
				return fmt.Errorf("format tanggal jatuh tempo tidak valid: %w", parseErr)
			}
			paymentDueDate = &parsedDate
		}

		// Process each received item
		for _, item := range receivedItems {
			// Verify the detail exists and belongs to this order
			var productID string
			var requestedQty float64
			detailErr := tx.QueryRow(`
                Select Product_Id, Requested_Quantity
                From Purchase_Order_Detail
                Where Id = $1 And Purchase_Order_Id = $2
            `, item.DetailID, orderID).Scan(&productID, &requestedQty)

			if detailErr != nil {
				if errors.Is(detailErr, sql.ErrNoRows) {
					return fmt.Errorf("detail pesanan tidak ditemukan atau tidak termasuk dalam pesanan ini")
				}
				return fmt.Errorf("gagal memeriksa detail pesanan: %w", detailErr)
			}

			// Validate received quantity
			if item.ReceivedQuantity > requestedQty {
				return fmt.Errorf("jumlah diterima (%f) tidak boleh lebih dari jumlah pesanan (%f)", item.ReceivedQuantity, requestedQty)
			}

			// Update detail with received quantity
			_, updateErr := tx.Exec(`
                Update Purchase_Order_Detail
                Set Received_Quantity = $1, Unit_Price = $2
                Where Id = $3
            `, item.ReceivedQuantity, item.UnitPrice, item.DetailID)

			if updateErr != nil {
				return fmt.Errorf("gagal memperbarui detail pesanan: %w", updateErr)
			}

			// Generate SKU for this batch
			var sku string
			skuErr := tx.QueryRow(`
                Select 'B-' || To_Char(Now(), 'YYYYMMDD') || '-' || 
                       Lpad(Coalesce(
                          (Select Count(*) + 1 From Product_Batch 
                           Where Created_At::date = Current_Date), 
                          1)::text, 4, '0')
            `).Scan(&sku)

			if skuErr != nil {
				return fmt.Errorf("gagal generate SKU: %w", skuErr)
			}

			// Create product batch
			var batchID string
			batchErr := tx.QueryRow(`
                Insert Into Product_Batch 
                (Sku, Product_Id, Purchase_Order_Id, Initial_Quantity, Current_Quantity, Unit_Price)
                Values ($1, $2, $3, $4, $5, $6)
                Returning Id
            `, sku, productID, orderID, item.ReceivedQuantity, item.ReceivedQuantity, item.UnitPrice).Scan(&batchID)

			if batchErr != nil {
				return fmt.Errorf("gagal membuat batch produk: %w", batchErr)
			}

			// Allocate to storage
			var totalAllocated float64 = 0
			for _, allocation := range item.StorageAllocations {
				// Check if storage exists
				var exists bool
				storageErr := tx.QueryRow(`
                    Select Exists(Select 1 From Storage Where Id = $1 And Deleted_At Is Null)
                `, allocation.StorageID).Scan(&exists)

				if storageErr != nil {
					return fmt.Errorf("gagal memeriksa penyimpanan: %w", storageErr)
				}

				if !exists {
					return fmt.Errorf("penyimpanan dengan ID %s tidak ditemukan", allocation.StorageID)
				}

				// Create storage allocation
				_, allocErr := tx.Exec(`
                    Insert Into Batch_Storage (Batch_Id, Storage_Id, Quantity)
                    Values ($1, $2, $3)
                `, batchID, allocation.StorageID, allocation.Quantity)

				if allocErr != nil {
					return fmt.Errorf("gagal mengalokasikan batch ke penyimpanan: %w", allocErr)
				}

				// Create inventory log
				_, logErr := tx.Exec(`
                    Insert Into Inventory_Log 
                    (Batch_Id, Storage_Id, User_Id, Purchase_Order_Id, Action, Quantity, Description)
                    Values ($1, $2, $3, $4, $5, $6, $7)
                `, batchID, allocation.StorageID, userID, orderID, "add", allocation.Quantity,
					fmt.Sprintf("Penerimaan barang dari PO #%s", orderID))

				if logErr != nil {
					return fmt.Errorf("gagal mencatat log inventaris: %w", logErr)
				}

				totalAllocated += allocation.Quantity
			}

			// Verify total allocated matches received quantity
			if totalAllocated != item.ReceivedQuantity {
				return fmt.Errorf("jumlah alokasi (%f) tidak sama dengan jumlah diterima (%f)",
					totalAllocated, item.ReceivedQuantity)
			}
		}

		// Update purchase order status and mark as received
		_, updateOrderErr := tx.Exec(`
            Update Purchase_Order
            Set Status = 'received', Received_By = $1, Received_At = Current_Timestamp,
                Payment_Due_Date = $2
            Where Id = $3
        `, userID, paymentDueDate, orderID)

		if updateOrderErr != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", updateOrderErr)
		}

		return nil
	})
}

// CheckOrder marks a received purchase order as checked
func (r *PurchaseOrderRepositoryImpl) CheckOrder(orderID string, userID string) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Check if order exists and is in 'received' status
		var status string
		orderErr := tx.QueryRow(`
            Select Status 
            From Purchase_Order 
            Where Id = $1
        `, orderID).Scan(&status)

		if orderErr != nil {
			if errors.Is(orderErr, sql.ErrNoRows) {
				return fmt.Errorf("pesanan pembelian tidak ditemukan")
			}
			return fmt.Errorf("gagal memeriksa status pesanan: %w", orderErr)
		}

		if status != "received" {
			return fmt.Errorf("hanya pesanan dengan status 'received' yang dapat diperiksa, status saat ini: %s", status)
		}

		// Update purchase order status and mark as checked
		_, updateErr := tx.Exec(`
            Update Purchase_Order
            Set Status = 'checked', Checked_By = $1, Checked_At = Current_Timestamp
            Where Id = $2
        `, userID, orderID)

		if updateErr != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", updateErr)
		}

		return nil
	})
}

// CancelOrder cancels a purchase order
func (r *PurchaseOrderRepositoryImpl) CancelOrder(orderID string, userID string) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Check if order exists and is in 'ordered' status
		var status string
		orderErr := tx.QueryRow(`
            Select Status 
            From Purchase_Order 
            Where Id = $1
        `, orderID).Scan(&status)

		if orderErr != nil {
			if errors.Is(orderErr, sql.ErrNoRows) {
				return fmt.Errorf("pesanan pembelian tidak ditemukan")
			}
			return fmt.Errorf("gagal memeriksa status pesanan: %w", orderErr)
		}

		// Only allow cancellation of orders in 'ordered' status
		if status != "ordered" {
			return fmt.Errorf("hanya pesanan dengan status 'ordered' yang dapat dibatalkan, status saat ini: %s", status)
		}

		// Update purchase order status and mark as cancelled
		_, updateErr := tx.Exec(`
            Update Purchase_Order
            Set Status = 'cancelled', Cancelled_By = $1, Cancelled_At = Current_Timestamp
            Where Id = $2
        `, userID, orderID)

		if updateErr != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", updateErr)
		}

		return nil
	})
}

// CreateReturn creates a new purchase order return record
func (r *PurchaseOrderRepositoryImpl) CreateReturn(returnRequest CreatePurchaseOrderReturnRequest, userID string) error {
	var returnResponse PurchaseOrderReturnDetailResponse

	txErr := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Verify purchase order exists and check status
		var orderStatus string
		orderErr := tx.QueryRow(`
            Select Status
            From Purchase_Order
            Where Id = $1
        `, returnRequest.PurchaseOrderID).Scan(&orderStatus)

		if orderErr != nil {
			if errors.Is(orderErr, sql.ErrNoRows) {
				return fmt.Errorf("pesanan pembelian tidak ditemukan")
			}
			return fmt.Errorf("gagal memeriksa status pesanan: %w", orderErr)
		}

		// Only allow returns for received or checked orders
		if orderStatus != "received" && orderStatus != "checked" {
			return fmt.Errorf("hanya pesanan dengan status 'received' atau 'checked' yang dapat dikembalikan, status saat ini: %s", orderStatus)
		}

		// Verify purchase order detail exists and has sufficient quantity
		var productID string
		var receivedQty, returnedQty float64
		detailErr := tx.QueryRow(`
            Select Pod.Product_Id, Pod.Received_Quantity, Pod.Total_Returned_Quantity
            From Purchase_Order_Detail Pod
            Where Pod.Id = $1 And Pod.Purchase_Order_Id = $2
        `, returnRequest.ProductDetailID, returnRequest.PurchaseOrderID).Scan(&productID, &receivedQty, &returnedQty)

		if detailErr != nil {
			if errors.Is(detailErr, sql.ErrNoRows) {
				return fmt.Errorf("detail produk tidak ditemukan")
			}
			return fmt.Errorf("gagal memverifikasi detail pesanan: %w", detailErr)
		}

		// Calculate remaining quantity that can be returned
		remainingQty := receivedQty - returnedQty

		// Check if return quantity is valid
		if returnRequest.ReturnQuantity > remainingQty {
			return fmt.Errorf("jumlah pengembalian melebihi jumlah yang tersisa (%f)", remainingQty)
		}

		// Create return record
		var returnID string
		var returnedAt time.Time
		createReturnErr := tx.QueryRow(`
            Insert Into Purchase_Order_Return (
                Purchase_Order_Id, Product_Detail_Id, Return_Quantity, 
                Remaining_Quantity, Return_Reason, Return_Status, Returned_By
            ) Values ($1, $2, $3, $4, $5, 'pending', $6)
            Returning Id, Returned_At
        `,
			returnRequest.PurchaseOrderID,
			returnRequest.ProductDetailID,
			returnRequest.ReturnQuantity,
			remainingQty-returnRequest.ReturnQuantity,
			returnRequest.ReturnReason,
			userID,
		).Scan(&returnID, &returnedAt)

		if createReturnErr != nil {
			return fmt.Errorf("gagal membuat catatan pengembalian: %w", createReturnErr)
		}

		// Update returned quantity in purchase order detail
		_, updateDetailErr := tx.Exec(`
            Update Purchase_Order_Detail
            Set Total_Returned_Quantity = Total_Returned_Quantity + $1
            Where Id = $2
        `, returnRequest.ReturnQuantity, returnRequest.ProductDetailID)

		if updateDetailErr != nil {
			return fmt.Errorf("gagal memperbarui jumlah pengembalian pada detail pesanan: %w", updateDetailErr)
		}

		// If all items returned, update purchase order status
		var totalReceived, totalReturned float64
		totalsErr := tx.QueryRow(`
            Select Sum(Received_Quantity), Sum(Total_Returned_Quantity)
            From Purchase_Order_Detail
            Where Purchase_Order_Id = $1
        `, returnRequest.PurchaseOrderID).Scan(&totalReceived, &totalReturned)

		if totalsErr != nil {
			return fmt.Errorf("gagal menghitung total barang: %w", totalsErr)
		}

		// If all received items are now returned, update order status
		if totalReceived == totalReturned {
			_, updateOrderErr := tx.Exec(`
                Update Purchase_Order 
                Set Status = 'returned', Fully_Returned_By = $1, Fully_Returned_At = Current_Timestamp
                Where Id = $2
            `, userID, returnRequest.PurchaseOrderID)

			if updateOrderErr != nil {
				return fmt.Errorf("gagal memperbarui status pesanan: %w", updateOrderErr)
			}
		} else if totalReturned > 0 {
			// If some items returned, update status to partially_returned
			_, updatePartialErr := tx.Exec(`
                Update Purchase_Order 
                Set Status = 'partially_returned'
                Where Id = $1 And Status != 'returned'
            `, returnRequest.PurchaseOrderID)

			if updatePartialErr != nil {
				return fmt.Errorf("gagal memperbarui status pesanan sebagian: %w", updatePartialErr)
			}
		}

		// Process batch returns
		for _, batchReturn := range returnRequest.BatchDetailsToReturn {
			// Verify batch exists and has sufficient quantity
			var batchQty float64
			batchErr := tx.QueryRow(`
                Select Current_Quantity
                From Product_Batch
                Where Id = $1 And Purchase_Order_Id = $2
            `, batchReturn.BatchID, returnRequest.PurchaseOrderID).Scan(&batchQty)

			if batchErr != nil {
				if errors.Is(batchErr, sql.ErrNoRows) {
					return fmt.Errorf("batch dengan ID %s tidak ditemukan", batchReturn.BatchID)
				}
				return fmt.Errorf("gagal memeriksa batch: %w", batchErr)
			}

			if batchReturn.ReturnQuantity > batchQty {
				return fmt.Errorf("jumlah pengembalian untuk batch %s melebihi jumlah yang tersedia (%f)", batchReturn.BatchID, batchQty)
			}

			// Create batch return record
			_, createBatchReturnErr := tx.Exec(`
                Insert Into Purchase_Order_Return_Batch (
                    Purchase_Return_Id, Batch_Id, Return_Quantity
                ) Values ($1, $2, $3)
            `, returnID, batchReturn.BatchID, batchReturn.ReturnQuantity)

			if createBatchReturnErr != nil {
				return fmt.Errorf("gagal membuat catatan pengembalian batch: %w", createBatchReturnErr)
			}

			// Update batch quantity
			_, updateBatchErr := tx.Exec(`
                Update Product_Batch
                Set Current_Quantity = Current_Quantity - $1
                Where Id = $2
            `, batchReturn.ReturnQuantity, batchReturn.BatchID)

			if updateBatchErr != nil {
				return fmt.Errorf("gagal memperbarui jumlah batch: %w", updateBatchErr)
			}
		}

		// Set the return response data
		returnResponse.ID = returnID
		returnResponse.PurchaseOrderID = returnRequest.PurchaseOrderID
		returnResponse.ProductDetailID = returnRequest.ProductDetailID
		returnResponse.ReturnQuantity = returnRequest.ReturnQuantity
		returnResponse.RemainingQuantity = remainingQty - returnRequest.ReturnQuantity
		returnResponse.ReturnReason = returnRequest.ReturnReason
		returnResponse.ReturnStatus = "pending"
		returnResponse.ReturnedBy = userID
		returnResponse.ReturnedAt = returnedAt.Format(time.RFC3339)

		// Set batch details in response
		var batchDetails []BatchReturnResponse
		for _, batch := range returnRequest.BatchDetailsToReturn {
			batchDetails = append(batchDetails, BatchReturnResponse{
				BatchID:        batch.BatchID,
				ReturnQuantity: batch.ReturnQuantity,
			})
		}
		returnResponse.BatchDetails = batchDetails

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

// GetAllReturns fetches all purchase order returns with filtering and pagination
func (r *PurchaseOrderRepositoryImpl) GetAllReturns(req GetPurchaseOrderReturnRequest) ([]GetPurchaseOrderReturnResponse, int, error) {
	// Build count query
	countBuilder := utils.NewQueryBuilder(`
        Select Count(R.Id)
        From Purchase_Order_Return R
        Join Purchase_Order Po On R.Purchase_Order_Id = Po.Id
        Where 1=1
    `)

	if req.PurchaseOrderID != "" {
		countBuilder.AddFilter("r.Purchase_Order_Id =", req.PurchaseOrderID)
	}

	if req.Status != "" {
		countBuilder.AddFilter("r.Return_Status =", req.Status)
	}

	if req.FromDate != "" {
		countBuilder.AddFilter("r.Returned_At >=", req.FromDate)
	}

	if req.ToDate != "" {
		countBuilder.AddFilter("r.Returned_At <=", req.ToDate)
	}

	countQuery, countParams := countBuilder.Build()

	// Execute count query
	var totalItems int
	countErr := r.db.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if countErr != nil {
		return nil, 0, fmt.Errorf("gagal menghitung jumlah pengembalian: %w", countErr)
	}

	// Build main query
	queryBuilder := utils.NewQueryBuilder(`
        Select R.Id, R.Purchase_Order_Id, R.Product_Detail_Id,
               R.Return_Quantity, R.Remaining_Quantity, R.Return_Reason,
               R.Return_Status, Coalesce(U1.Username, '') As Returnedby,
               R.Returned_At, U2.Username As Cancelledby, R.Cancelled_At
        From Purchase_Order_Return R
        Join Purchase_Order Po On R.Purchase_Order_Id = Po.Id
        Left Join Appuser U1 On R.Returned_By = U1.Id
        Left Join Appuser U2 On R.Cancelled_By = U2.Id
        Where 1=1
    `)

	if req.PurchaseOrderID != "" {
		queryBuilder.AddFilter("r.Purchase_Order_Id =", req.PurchaseOrderID)
	}

	if req.Status != "" {
		queryBuilder.AddFilter("r.Return_Status =", req.Status)
	}

	if req.FromDate != "" {
		queryBuilder.AddFilter("r.Returned_At >=", req.FromDate)
	}

	if req.ToDate != "" {
		queryBuilder.AddFilter("r.Returned_At <=", req.ToDate)
	}

	// Add sorting
	queryBuilder.Query.WriteString(" ORDER BY r.Returned_At DESC")

	// Add pagination
	mainQuery, queryParams := queryBuilder.AddPagination(req.PageSize, req.Page).Build()

	// Execute main query
	rows, queryErr := r.db.Query(mainQuery, queryParams...)
	if queryErr != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data pengembalian: %w", queryErr)
	}
	defer rows.Close()

	// Process results
	var returns []GetPurchaseOrderReturnResponse
	for rows.Next() {
		var returnItem GetPurchaseOrderReturnResponse
		var returnedAt time.Time
		var cancelledBy sql.NullString
		var cancelledAt sql.NullTime

		scanErr := rows.Scan(
			&returnItem.ID,
			&returnItem.PurchaseOrderID,
			&returnItem.ProductDetailID,
			&returnItem.ReturnQuantity,
			&returnItem.RemainingQuantity,
			&returnItem.ReturnReason,
			&returnItem.ReturnStatus,
			&returnItem.ReturnedBy,
			&returnedAt,
			&cancelledBy,
			&cancelledAt,
		)
		if scanErr != nil {
			return nil, 0, fmt.Errorf("gagal membaca data pengembalian: %w", scanErr)
		}

		returnItem.ReturnedAt = returnedAt.Format(time.RFC3339)

		if cancelledAt.Valid {
			returnItem.CancelledAt = cancelledAt.Time.Format(time.RFC3339)
		}

		if cancelledBy.Valid {
			returnItem.CancelledBy = cancelledBy.String
		}

		returns = append(returns, returnItem)
	}

	if rowErr := rows.Err(); rowErr != nil {
		return nil, 0, fmt.Errorf("gagal memproses data pengembalian: %w", rowErr)
	}

	return returns, totalItems, nil
}

// GetReturnByID fetches a specific purchase order return by ID
func (r *PurchaseOrderRepositoryImpl) GetReturnByID(returnID string) (*PurchaseOrderReturnDetailResponse, error) {
	// Fetch the return record
	var returnDetail PurchaseOrderReturnDetailResponse
	var returnedAt time.Time
	var cancelledBy sql.NullString
	var cancelledAt sql.NullTime

	queryErr := r.db.QueryRow(`
        Select R.Id, R.Purchase_Order_Id, R.Product_Detail_Id,
               R.Return_Quantity, R.Remaining_Quantity, R.Return_Reason,
               R.Return_Status, Coalesce(U1.Username, '') As Returnedby,
               R.Returned_At, U2.Username As Cancelledby, R.Cancelled_At
        From Purchase_Order_Return R
        Left Join Appuser U1 On R.Returned_By = U1.Id
        Left Join Appuser U2 On R.Cancelled_By = U2.Id
        Where R.Id = $1
    `, returnID).Scan(
		&returnDetail.ID,
		&returnDetail.PurchaseOrderID,
		&returnDetail.ProductDetailID,
		&returnDetail.ReturnQuantity,
		&returnDetail.RemainingQuantity,
		&returnDetail.ReturnReason,
		&returnDetail.ReturnStatus,
		&returnDetail.ReturnedBy,
		&returnedAt,
		&cancelledBy,
		&cancelledAt,
	)

	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			return nil, fmt.Errorf("pengembalian dengan ID %s tidak ditemukan", returnID)
		}
		return nil, fmt.Errorf("gagal mengambil data pengembalian: %w", queryErr)
	}

	returnDetail.ReturnedAt = returnedAt.Format(time.RFC3339)

	if cancelledAt.Valid {
		returnDetail.CancelledAt = cancelledAt.Time.Format(time.RFC3339)
	}

	if cancelledBy.Valid {
		returnDetail.CancelledBy = cancelledBy.String
	}

	// Fetch the batch details
	rows, batchQueryErr := r.db.Query(`
        Select Rb.Batch_Id, Rb.Return_Quantity
        From Purchase_Order_Return_Batch Rb
        Where Rb.Purchase_Return_Id = $1
    `, returnID)

	if batchQueryErr != nil {
		return nil, fmt.Errorf("gagal mengambil detail batch pengembalian: %w", batchQueryErr)
	}
	defer rows.Close()

	var batchDetails []BatchReturnResponse
	for rows.Next() {
		var batchReturn BatchReturnResponse
		scanErr := rows.Scan(
			&batchReturn.BatchID,
			&batchReturn.ReturnQuantity,
		)
		if scanErr != nil {
			return nil, fmt.Errorf("gagal membaca data batch pengembalian: %w", scanErr)
		}
		batchDetails = append(batchDetails, batchReturn)
	}

	if rowErr := rows.Err(); rowErr != nil {
		return nil, fmt.Errorf("gagal memproses data batch pengembalian: %w", rowErr)
	}

	returnDetail.BatchDetails = batchDetails

	return &returnDetail, nil
}

// CancelReturn cancels a purchase order return
func (r *PurchaseOrderRepositoryImpl) CancelReturn(returnID string, userID string) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Check if return exists and is pending
		var returnStatus string
		var purchaseOrderID, productDetailID string
		var returnQuantity float64

		checkErr := tx.QueryRow(`
            Select Return_Status, Purchase_Order_Id, Product_Detail_Id, Return_Quantity
            From Purchase_Order_Return 
            Where Id = $1
        `, returnID).Scan(&returnStatus, &purchaseOrderID, &productDetailID, &returnQuantity)

		if checkErr != nil {
			if errors.Is(checkErr, sql.ErrNoRows) {
				return fmt.Errorf("pengembalian dengan ID %s tidak ditemukan", returnID)
			}
			return fmt.Errorf("gagal memeriksa status pengembalian: %w", checkErr)
		}

		if returnStatus != "pending" {
			return fmt.Errorf("hanya pengembalian dengan status 'pending' yang dapat dibatalkan")
		}

		// Update return status to cancelled
		_, updateErr := tx.Exec(`
            Update Purchase_Order_Return
            Set Return_Status = 'cancelled', Cancelled_By = $1, Cancelled_At = Current_Timestamp
            Where Id = $2
        `, userID, returnID)

		if updateErr != nil {
			return fmt.Errorf("gagal memperbarui status pengembalian: %w", updateErr)
		}

		// Update purchase order detail to decrease returned quantity
		_, detailUpdateErr := tx.Exec(`
            Update Purchase_Order_Detail
            Set Total_Returned_Quantity = Total_Returned_Quantity - $1
            Where Id = $2
        `, returnQuantity, productDetailID)

		if detailUpdateErr != nil {
			return fmt.Errorf("gagal memperbarui jumlah pengembalian pada detail pesanan: %w", detailUpdateErr)
		}

		// Check if order status needs to be updated from returned/partially_returned to checked
		var orderStatus string
		var totalReceived, totalReturned float64

		statusErr := tx.QueryRow(`
            Select Po.Status, 
                Sum(Pod.Received_Quantity) As Totalreceived,
                Sum(Pod.Total_Returned_Quantity) As Totalreturned
            From Purchase_Order Po 
            Join Purchase_Order_Detail Pod On Po.Id = Pod.Purchase_Order_Id
            Where Po.Id = $1
            Group By Po.Status
        `, purchaseOrderID).Scan(&orderStatus, &totalReceived, &totalReturned)

		if statusErr != nil {
			return fmt.Errorf("gagal memeriksa status pesanan: %w", statusErr)
		}

		// If order was in a returned status and now there's no returns, update to checked
		if (orderStatus == "returned" || orderStatus == "partially_returned") && totalReturned == 0 {
			_, orderUpdateErr := tx.Exec(`
                Update Purchase_Order
                Set Status = 'checked', Return_Cancelled_At = Current_Timestamp, Return_Cancelled_By = $1
                Where Id = $2
            `, userID, purchaseOrderID)

			if orderUpdateErr != nil {
				return fmt.Errorf("gagal memperbarui status pesanan: %w", orderUpdateErr)
			}
		} else if orderStatus == "returned" && totalReturned > 0 && totalReturned < totalReceived {
			// If partial returns remain, update to partially_returned
			_, orderUpdateErr := tx.Exec(`
                Update Purchase_Order
                Set Status = 'partially_returned'
                Where Id = $1
            `, purchaseOrderID)

			if orderUpdateErr != nil {
				return fmt.Errorf("gagal memperbarui status pesanan: %w", orderUpdateErr)
			}
		}

		return nil
	})
}

// GetByOrderID fetches all items for a specific purchase order
func (r *PurchaseOrderRepositoryImpl) GetByOrderID(orderID string) ([]GetPurchaseOrderItemResponse, error) {
	rows, err := r.db.Query(`
        Select Pod.Id, Pod.Product_Id, P.Name, Pod.Requested_Quantity, Pod.Unit_Price, 
               (Pod.Requested_Quantity * Pod.Unit_Price) As Subtotal
        From Purchase_Order_Detail Pod
        Join Product P On Pod.Product_Id = P.Id
        Where Pod.Purchase_Order_Id = $1
    `, orderID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil detail pesanan: %w", err)
	}
	defer rows.Close()

	var items []GetPurchaseOrderItemResponse
	for rows.Next() {
		var item GetPurchaseOrderItemResponse
		errScan := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
		)
		if errScan != nil {
			return nil, fmt.Errorf("gagal memproses data detail pesanan: %w", errScan)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal memproses hasil query: %w", err)
	}

	return items, nil
}

// RemoveItem removes an item from a purchase order
func (r *PurchaseOrderRepositoryImpl) RemoveItem(itemID string) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Check if item exists and get order ID
		var orderID string
		var orderStatus string

		err := tx.QueryRow(`
            Select Pod.Purchase_Order_Id, Po.Status
            From Purchase_Order_Detail Pod
            Join Purchase_Order Po On Pod.Purchase_Order_Id = Po.Id
            Where Pod.Id = $1
        `, itemID).Scan(&orderID, &orderStatus)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("item pesanan tidak ditemukan")
			}
			return fmt.Errorf("gagal memeriksa detail item: %w", err)
		}

		// Only allow removals if order is still in 'ordered' status
		if orderStatus != "ordered" {
			return fmt.Errorf("item hanya dapat dihapus pada pesanan dengan status 'ordered', status saat ini: %s", orderStatus)
		}

		// Delete item
		result, err := tx.Exec(`
            Delete From Purchase_Order_Detail
            Where Id = $1
        `, itemID)

		if err != nil {
			return fmt.Errorf("gagal menghapus item pesanan: %w", err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("gagal memeriksa hasil penghapusan: %w", err)
		}

		if affected == 0 {
			return fmt.Errorf("item pesanan tidak ditemukan")
		}

		// Update total amount in purchase order
		_, err = tx.Exec(`
            Update Purchase_Order
            Set Total_Amount = Coalesce((
                Select Sum(Requested_Quantity * Unit_Price)
                From Purchase_Order_Detail
                Where Purchase_Order_Id = $1
            ), 0),
            Updated_At = Current_Timestamp
            Where Id = $1
        `, orderID)

		if err != nil {
			return fmt.Errorf("gagal memperbarui total pesanan: %w", err)
		}

		return nil
	})
}

// AddItem adds a new item to an existing purchase order
func (r *PurchaseOrderRepositoryImpl) AddItem(orderID string, req CreatePurchaseOrderItemRequest) error {
	txErr := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Verify purchase order exists and is in "ordered" status
		var status string
		orderErr := tx.QueryRow(`
            Select Status 
            From Purchase_Order 
            Where Id = $1
        `, orderID).Scan(&status)

		if orderErr != nil {
			if errors.Is(orderErr, sql.ErrNoRows) {
				return fmt.Errorf("pesanan dengan ID %s tidak ditemukan", orderID)
			}
			return fmt.Errorf("gagal memeriksa status pesanan: %w", orderErr)
		}

		if status != "ordered" {
			return fmt.Errorf("item hanya dapat ditambahkan ke pesanan dengan status 'ordered', status saat ini: %s", status)
		}

		// Verify product exists
		var productName string
		productErr := tx.QueryRow(`
            Select Name 
            From Product 
            Where Id = $1 And Deleted_At Is Null
        `, req.ProductID).Scan(&productName)

		if productErr != nil {
			if errors.Is(productErr, sql.ErrNoRows) {
				return fmt.Errorf("produk dengan ID %s tidak ditemukan", req.ProductID)
			}
			return fmt.Errorf("gagal memeriksa produk: %w", productErr)
		}

		// Check if product already exists in this order
		var existingDetailID string
		dupErr := tx.QueryRow(`
            Select Id 
            From Purchase_Order_Detail 
            Where Purchase_Order_Id = $1 And Product_Id = $2
        `, orderID, req.ProductID).Scan(&existingDetailID)

		if dupErr == nil {
			// Product already exists in this order
			return fmt.Errorf("produk %s sudah ada dalam pesanan ini", productName)
		} else if !errors.Is(dupErr, sql.ErrNoRows) {
			return fmt.Errorf("gagal memeriksa duplikasi produk: %w", dupErr)
		}

		// Insert new detail
		var detailID string
		insertErr := tx.QueryRow(`
            Insert Into Purchase_Order_Detail 
                (Purchase_Order_Id, Product_Id, Requested_Quantity, Unit_Price)
            Values ($1, $2, $3, $4)
            Returning Id
        `, orderID, req.ProductID, req.Quantity, req.Price).Scan(&detailID)

		if insertErr != nil {
			return fmt.Errorf("gagal menambahkan item pesanan: %w", insertErr)
		}

		// Update order total amount
		_, updateOrderErr := tx.Exec(`
            Update Purchase_Order 
            Set Total_Amount = (
                Select Sum(Requested_Quantity * Unit_Price) 
                From Purchase_Order_Detail 
                Where Purchase_Order_Id = $1
            ), Updated_At = Current_Timestamp
            Where Id = $1
        `, orderID)

		if updateOrderErr != nil {
			return fmt.Errorf("gagal memperbarui total pesanan: %w", updateOrderErr)
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

// UpdateItem updates an existing purchase order item
func (r *PurchaseOrderRepositoryImpl) UpdateItem(req UpdatePurchaseOrderItemRequest) error {
	txErr := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Verify item exists
		var orderID, productID string
		var currentQty, currentPrice float64

		itemErr := tx.QueryRow(`
            Select Pod.Purchase_Order_Id, Pod.Product_Id, Pod.Requested_Quantity, Pod.Unit_Price
            From Purchase_Order_Detail Pod
            Where Pod.Id = $1
        `, req.ID).Scan(&orderID, &productID, &currentQty, &currentPrice)

		if itemErr != nil {
			if errors.Is(itemErr, sql.ErrNoRows) {
				return fmt.Errorf("item pesanan dengan ID %s tidak ditemukan", req.ID)
			}
			return fmt.Errorf("gagal memeriksa item pesanan: %w", itemErr)
		}

		// Verify purchase order is in "ordered" status
		var status string
		orderErr := tx.QueryRow(`
            Select Status 
            From Purchase_Order 
            Where Id = $1
        `, orderID).Scan(&status)

		if orderErr != nil {
			return fmt.Errorf("gagal memeriksa status pesanan: %w", orderErr)
		}

		if status != "ordered" {
			return fmt.Errorf("item hanya dapat diperbarui pada pesanan dengan status 'ordered', status saat ini: %s", status)
		}

		// Prepare update values
		updatedQty := currentQty
		if req.Quantity > 0 {
			updatedQty = req.Quantity
		}

		updatedPrice := currentPrice
		if req.Price > 0 {
			updatedPrice = req.Price
		}

		// Update the item
		updateErr := tx.QueryRow(`
            Update Purchase_Order_Detail 
            Set Requested_Quantity = $1, Unit_Price = $2, Updated_At = Current_Timestamp
            Where Id = $3
            Returning Product_Id
        `, updatedQty, updatedPrice, req.ID).Scan(&productID)

		if updateErr != nil {
			return fmt.Errorf("gagal memperbarui item pesanan: %w", updateErr)
		}

		// Update order total amount
		_, updateOrderErr := tx.Exec(`
            Update Purchase_Order 
            Set Total_Amount = (
                Select Sum(Requested_Quantity * Unit_Price) 
                From Purchase_Order_Detail 
                Where Purchase_Order_Id = $1
            ), Updated_At = Current_Timestamp
            Where Id = $1
        `, orderID)

		if updateOrderErr != nil {
			return fmt.Errorf("gagal memperbarui total pesanan: %w", updateOrderErr)
		}

		// Get product name
		var productName string
		nameErr := tx.QueryRow(`
            Select Name 
            From Product 
            Where Id = $1
        `, productID).Scan(&productName)

		if nameErr != nil {
			return fmt.Errorf("gagal mendapatkan nama produk: %w", nameErr)
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

// GetInventoryLogs fetches all inventory logs with filtering and pagination
