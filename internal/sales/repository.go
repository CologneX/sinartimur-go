package sales

import (
	"database/sql"
	"errors"
	"fmt"
	"sinartimur-go/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SalesRepository interface {
	// Sales Order operations
	GetSalesOrders(req GetSalesOrdersRequest) ([]GetSalesOrdersResponse, int, error)
	GetSalesOrderByID(id string) (*SalesOrder, error)
	GetSalesOrderItems(salesOrderID string) ([]SalesOrderItem, error)
	GetSalesOrderWithDetails(salesOrderID string) (*GetSalesOrderDetailResponse, error)
	CreateSalesOrder(req CreateSalesOrderRequest, userID string) (*CreateSalesOrderResponse, error)
	UpdateSalesOrder(req UpdateSalesOrderRequest) (*UpdateSalesOrderResponse, error)
	CancelSalesOrder(req CancelSalesOrderRequest, userID string) error

	// Sales Order Item operations
	AddItemToSalesOrder(req AddSalesOrderItemRequest) (*UpdateAndCreateItemResponse, error)
	UpdateSalesOrderItem(req UpdateSalesOrderItemRequest) (*UpdateAndCreateItemResponse, error)
	DeleteSalesOrderItem(req DeleteSalesOrderItemRequest) error

	// Invoice operations
	GetSalesInvoices(req GetSalesInvoicesRequest) ([]GetSalesInvoicesResponse, int, error)
	CreateSalesInvoice(req CreateSalesInvoiceRequest, userID string,
		tx *sql.Tx) (*CreateSalesInvoiceResponse, error)
	CancelSalesInvoice(req CancelSalesInvoiceRequest, userID string) error

	// Return operations
	ReturnInvoiceItems(req ReturnInvoiceItemsRequest, userID string) (*ReturnInvoiceItemsResponse, error)
	CancelInvoiceReturn(req CancelInvoiceReturnRequest, userID string) error

	// Delivery Note operations
	CreateDeliveryNote(req CreateDeliveryNoteRequest, userID string) (*CreateDeliveryNoteResponse, error)
	CancelDeliveryNote(req CancelDeliveryNoteRequest, userID string) error

	// Batch operations
	GetAllBatches(req GetAllBatchesRequest) ([]GetAllBatchesResponse, int, error)
}

type SalesRepositoryImpl struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) SalesRepository {
	return &SalesRepositoryImpl{db: db}
}

// GetAllBatches retrieves all product batches with pagination and filtering
// grouped by storage location for sales order creation
func (r *SalesRepositoryImpl) GetAllBatches(req GetAllBatchesRequest) ([]GetAllBatchesResponse, int, error) {
	// Map to hold grouped results by storage
	storageMap := make(map[string]*GetAllBatchesResponse)
	var totalItems int

	// Build base query with joins to get product name, storage information and include storage location
	qb := utils.NewQueryBuilder(`
        SELECT bs.id, pb.id as batch_id, pb.sku, pb.product_id, p.name as product_name, 
               pb.current_quantity, pb.unit_price, pb.created_at, 
               bs.storage_id, s.name as storage_name, s.location as storage_location, bs.quantity
        FROM product_batch pb
        JOIN product p ON pb.product_id = p.id
        JOIN batch_storage bs ON pb.id = bs.batch_id
        JOIN storage s ON bs.storage_id = s.id
        WHERE pb.current_quantity > 0 AND bs.quantity > 0
    `)

	// Add search filter if provided
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		qb.Query.WriteString(" AND (pb.sku ILIKE $" + strconv.Itoa(qb.Count) + " OR p.name ILIKE $" + strconv.Itoa(qb.Count) + ")")
		qb.Params = append(qb.Params, searchTerm)
		qb.Count++
	}

	// Get count first (count distinct storage_ids to get number of storage groups)
	countQuery := "SELECT COUNT(DISTINCT bs.storage_id) FROM product_batch pb " +
		"JOIN batch_storage bs ON pb.id = bs.batch_id " +
		"JOIN product p ON pb.product_id = p.id " +
		"WHERE pb.current_quantity > 0 AND bs.quantity > 0"

	// Add search condition to count query if needed
	if req.Search != "" {
		countQuery += " AND (pb.sku ILIKE $1 OR p.name ILIKE $1)"
		err := r.db.QueryRow(countQuery, "%"+req.Search+"%").Scan(&totalItems)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.db.QueryRow(countQuery).Scan(&totalItems)
		if err != nil {
			return nil, 0, err
		}
	}

	// Add sorting
	sortBy := req.SortBy
	sortOrder := req.SortOrder

	if sortBy == "" {
		sortBy = "s.name" // Default sort by storage name
	}

	if sortOrder == "" {
		sortOrder = "asc"
	}

	// Sanitize the sort fields to prevent SQL injection
	validSortFields := map[string]string{
		"sku":              "pb.sku",
		"product_name":     "p.name",
		"current_quantity": "bs.quantity",
		"unit_price":       "pb.unit_price",
		"created_at":       "pb.created_at",
		"storage_name":     "s.name",
	}

	// Use the mapped field if valid, otherwise default to storage name
	if sortField, ok := validSortFields[sortBy]; ok {
		sortBy = sortField
	} else {
		sortBy = "s.name"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	qb.Query.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder))

	// Don't add pagination yet - we'll fetch all results first and then paginate the grouped response

	// Execute query
	query, params := qb.Build()
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Map results to storage groups
	for rows.Next() {
		var storageID, storageName, storageLocation, id, batchID, sku, productID, productName string
		var quantity, unitPrice float64
		var createdAt time.Time

		err := rows.Scan(
			&id,
			&batchID,
			&sku,
			&productID,
			&productName,
			&quantity,
			&unitPrice,
			&createdAt,
			&storageID,
			&storageName,
			&storageLocation,
			&quantity,
		)
		if err != nil {
			return nil, 0, err
		}

		// Get or create storage group
		storageGroup, exists := storageMap[storageID]
		if !exists {
			storageGroup = &GetAllBatchesResponse{
				StorageID:                 storageID,
				StorageName:               storageName,
				StorageLocation:           storageLocation,
				GetAllBatchesStorageItems: []GetAllBatchesStorageItem{},
			}
			storageMap[storageID] = storageGroup
		}

		// Add batch item to storage group
		storageGroup.GetAllBatchesStorageItems = append(storageGroup.GetAllBatchesStorageItems, GetAllBatchesStorageItem{
			BatchStorageID: id,
			BatchID:        batchID,
			BatchSKU:       sku,
			Quantity:       quantity,
			Price:          unitPrice,
			ProductID:      productID,
			ProductName:    productName,
			CreatedAt:      createdAt.Format(time.RFC3339),
		})
	}

	// Convert map to slice
	result := make([]GetAllBatchesResponse, 0, len(storageMap))
	for _, storageGroup := range storageMap {
		result = append(result, *storageGroup)
	}

	// Sort the result slice by storage name to ensure consistent ordering
	sort.Slice(result, func(i, j int) bool {
		if sortOrder == "asc" {
			return result[i].StorageName < result[j].StorageName
		}
		return result[i].StorageName > result[j].StorageName
	})

	// Apply pagination to the final result
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
	if start >= len(result) {
		return []GetAllBatchesResponse{}, totalItems, nil
	}
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], totalItems, nil
}

// GetSalesOrders retrieves a paginated list of sales orders with filtering options
func (r *SalesRepositoryImpl) GetSalesOrders(req GetSalesOrdersRequest) ([]GetSalesOrdersResponse, int, error) {
	// Build base query for fetching sales orders
	baseQuery := `
        SELECT so.id, so.serial_id, so.customer_id, c.name AS customer_name, 
               so.order_date, so.status, so.payment_method, so.payment_due_date, 
               so.total_amount, so.created_at, so.updated_at, so.cancelled_at,
               (SELECT si.id FROM sales_invoice si WHERE si.sales_order_id = so.id AND si.cancelled_at IS NULL LIMIT 1) AS sales_invoice_id,
               (SELECT dn.id FROM delivery_note dn WHERE dn.sales_order_id = so.id AND dn.cancelled_at IS NULL LIMIT 1) AS delivery_note_id
        FROM sales_order so
        JOIN customer c ON so.customer_id = c.id
        WHERE 1=1`

	// Create count query
	countQuery := `SELECT COUNT(*) FROM sales_order so JOIN customer c ON so.customer_id = c.id WHERE 1=1`

	// Initialize query builders
	qb := utils.NewQueryBuilder(baseQuery)
	countQb := utils.NewQueryBuilder(countQuery)

	// Add filters to both queries
	if req.CustomerID != "" {
		condition := "so.customer_id ="
		qb.AddFilter(condition, req.CustomerID)
		countQb.AddFilter(condition, req.CustomerID)
	}

	if req.Status != "" {
		condition := "so.status ="
		qb.AddFilter(condition, req.Status)
		countQb.AddFilter(condition, req.Status)
	}

	if req.PaymentMethod != "" {
		condition := "so.payment_method ="
		qb.AddFilter(condition, req.PaymentMethod)
		countQb.AddFilter(condition, req.PaymentMethod)
	}

	if req.SerialID != "" {
		condition := "so.serial_id ILIKE"
		likeValue := "%" + req.SerialID + "%"
		qb.AddFilter(condition, likeValue)
		countQb.AddFilter(condition, likeValue)
	}

	if req.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, req.StartDate)
		if err == nil {
			condition := "so.order_date >="
			qb.AddFilter(condition, startDate)
			countQb.AddFilter(condition, startDate)
		}
	}

	if req.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err == nil {
			condition := "so.order_date <="
			qb.AddFilter(condition, endDate)
			countQb.AddFilter(condition, endDate)
		}
	}

	// Add sorting
	sortBy := "so.created_at"
	if req.SortBy != "" {
		sortBy = fmt.Sprintf("so.%s", req.SortBy)
	}

	sortOrder := "DESC"
	if req.SortOrder != "" {
		sortOrder = req.SortOrder
	}

	// Append sorting to main query only (not count query)
	qb.Query.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder))

	// Add pagination to main query only
	qb.AddPagination(req.PageSize, req.Page)

	// Build final queries
	query, params := qb.Build()
	countQuery, countParams := countQb.Build()

	// Execute count query
	var totalItems int
	err := r.db.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting sales orders: %w", err)
	}

	// Execute main query
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching sales orders: %w", err)
	}
	defer rows.Close()

	// Parse results
	var orders []GetSalesOrdersResponse
	for rows.Next() {
		var order GetSalesOrdersResponse

		errScan := rows.Scan(
			&order.ID,
			&order.SerialID,
			&order.CustomerID,
			&order.CustomerName,
			&order.OrderDate,
			&order.Status,
			&order.PaymentMethod,
			&order.PaymentDueDate,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.CancelledAt,
			&order.SalesInvoiceID,
			&order.DeliveryNoteID,
		)
		if errScan != nil {
			return nil, 0, fmt.Errorf("error scanning sales order row: %w", errScan)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating sales order rows: %w", err)
	}

	return orders, totalItems, nil
}

// GetSalesOrderWithDetails gets both order header and details for a specific order
func (r *SalesRepositoryImpl) GetSalesOrderWithDetails(salesOrderID string) (*GetSalesOrderDetailResponse, error) {
	var response GetSalesOrderDetailResponse

	// Get order header information with optimized query using a single join for related documents
	err := r.db.QueryRow(`
        SELECT 
            so.id, so.serial_id, so.customer_id, c.name, c.telephone, c.address,
            so.order_date, so.status, so.payment_method, so.payment_due_date, 
            so.total_amount, so.created_by, so.created_at, so.updated_at, so.cancelled_at,
            si.id, si.serial_id as sales_invoice_serial_id,
            dn.id, dn.serial_id as delivery_note_serial_id
        FROM sales_order so
        JOIN customer c ON so.customer_id = c.id
        LEFT JOIN LATERAL (
            SELECT id, serial_id 
            FROM sales_invoice 
            WHERE sales_order_id = $1 AND cancelled_at IS NULL
            ORDER BY created_at DESC
            LIMIT 1
        ) si ON true
        LEFT JOIN LATERAL (
            SELECT id, serial_id
            FROM delivery_note
            WHERE sales_order_id = $1 AND cancelled_at IS NULL
            ORDER BY created_at DESC
            LIMIT 1
        ) dn ON true
        WHERE so.id = $1
    `, salesOrderID).Scan(
		&response.ID,
		&response.SerialID,
		&response.CustomerID,
		&response.CustomerName,
		&response.CustomerPhone,
		&response.CustomerAddress,
		&response.OrderDate,
		&response.Status,
		&response.PaymentMethod,
		&response.PaymentDueDate,
		&response.TotalAmount,
		&response.CreatedBy,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.CancelledAt,
		&response.SalesInvoiceID,
		&response.SalesInvoiceSerialID,
		&response.DeliveryNoteID,
		&response.DeliveryNoteSerialID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sales order tidak ditemukan: %s", salesOrderID)
		}
		return nil, fmt.Errorf("error fetching sales order: %w", err)
	}

	// Get order items/details
	items, err := r.GetSalesOrderItems(salesOrderID)
	if err != nil {
		return nil, fmt.Errorf("error fetching order items: %w", err)
	}
	fmt.Println(response.SalesInvoiceSerialID, response.DeliveryNoteSerialID)

	response.Items = items
	return &response, nil
}

// GetSalesOrderByID retrieves a single sales order by its ID
func (r *SalesRepositoryImpl) GetSalesOrderByID(id string) (*SalesOrder, error) {
	query := `
        SELECT so.id, so.serial_id, so.customer_id, c.name AS customer_name, 
               so.order_date, so.status, so.payment_method, so.payment_due_date, 
               so.total_amount, so.created_by, so.created_at, so.updated_at, so.cancelled_at
        FROM sales_order so
        JOIN customer c ON so.customer_id = c.id
        WHERE so.id = $1`

	var order SalesOrder

	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.SerialID,
		&order.CustomerID,
		&order.CustomerName,
		&order.OrderDate,
		&order.Status,
		&order.PaymentMethod,
		&order.PaymentDueDate,
		&order.TotalAmount,
		&order.CreatedBy,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.CancelledAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sales order not found: %s", id)
		}
		return nil, fmt.Errorf("error fetching sales order: %w", err)
	}

	return &order, nil
}

// GetSalesOrderItems retrieves the details (items) for a specific sales order
func (r *SalesRepositoryImpl) GetSalesOrderItems(salesOrderID string) ([]SalesOrderItem, error) {
	query := `
        SELECT sod.id, sod.sales_order_id, 
               p.id AS product_id, p.name AS product_name, 
               u.name AS product_unit,
               pb.id AS batch_id, pb.sku AS batch_sku, 
               sod.batch_storage_id, 
               s.id AS storage_id, s.name AS storage_name,
               sod.quantity, sod.unit_price,
               (sod.quantity * sod.unit_price) AS total_price,
               bs.quantity + sod.quantity AS max_quantity
        FROM sales_order_detail sod
        JOIN batch_storage bs ON sod.batch_storage_id = bs.id
        JOIN product_batch pb ON bs.batch_id = pb.id
        JOIN product p ON pb.product_id = p.id
        JOIN unit u ON p.unit_id = u.id
        JOIN storage s ON bs.storage_id = s.id
        WHERE sod.sales_order_id = $1`

	rows, err := r.db.Query(query, salesOrderID)
	if err != nil {
		return nil, fmt.Errorf("error fetching sales order details: %w", err)
	}
	defer rows.Close()

	var items []SalesOrderItem
	for rows.Next() {
		var item SalesOrderItem

		errScan := rows.Scan(
			&item.ID,
			&item.SalesOrderID,
			&item.ProductID,
			&item.ProductName,
			&item.ProductUnit,
			&item.BatchID,
			&item.BatchSKU,
			&item.BatchStorageID,
			&item.StorageID,
			&item.StorageName,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
			&item.MaxQuantity,
		)
		if errScan != nil {
			return nil, fmt.Errorf("error scanning sales order detail row: %w", errScan)
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sales order detail rows: %w", err)
	}

	return items, nil
}

// CreateSalesOrder creates a new sales order with its details
func (r *SalesRepositoryImpl) CreateSalesOrder(req CreateSalesOrderRequest, userID string) (*CreateSalesOrderResponse, error) {
	var response CreateSalesOrderResponse

	err := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Insert sales order
		var orderID string
		var orderDate time.Time
		var paymentDueDate sql.NullTime
		var status string

		// Calculate total amount for the order
		var totalAmount float64
		for _, item := range req.Items {
			totalAmount += item.Quantity * item.UnitPrice
		}

		// Convert payment due date if provided
		if req.PaymentDueDate != "" {
			dueDate, err := time.Parse(time.RFC3339, req.PaymentDueDate)
			if err != nil {
				return fmt.Errorf("format tanggal jatuh tempo tidak valid: %w", err)
			}
			paymentDueDate = sql.NullTime{Time: dueDate, Valid: true}
		}

		serialID, errSerial := utils.GenerateNextSerialID(tx, "SO")
		if errSerial != nil {
			return fmt.Errorf("gagal membuat serial ID: %w", errSerial)
		}

		// Insert sales order
		orderQuery := `
		INSERT INTO sales_order (customer_id, serial_id, payment_method, payment_due_date, created_by, status, total_amount)
		VALUES ($1, $2, $3, $4, $5, 'order', $6)
		RETURNING id, serial_id, order_date, created_at, status`

		errOrder := tx.QueryRow(
			orderQuery,
			req.CustomerID,
			serialID,
			req.PaymentMethod,
			paymentDueDate,
			userID,
			totalAmount,
		).Scan(&orderID, &serialID, &orderDate, &response.CreatedAt, &status)

		if errOrder != nil {
			return fmt.Errorf("gagal membuat pesanan: %w", errOrder)
		}

		// Get customer name
		var customerName string
		errCustomer := tx.QueryRow("SELECT name FROM customer WHERE id = $1", req.CustomerID).Scan(&customerName)
		if errCustomer != nil {
			return fmt.Errorf("gagal mendapatkan data pelanggan: %w", errCustomer)
		}

		// Process each item
		for _, item := range req.Items {
			// We'll use batch_storage_id directly
			var batchID, productID string
			var availableQty, unitPrice float64
			var productName, batchSKU string

			// Get batch and product information from batch_storage ID
			errBatch := tx.QueryRow(`
				SELECT bs.quantity, bs.batch_id, pb.product_id, p.name, pb.sku, pb.unit_price
				FROM batch_storage bs
				JOIN product_batch pb ON bs.batch_id = pb.id
				JOIN product p ON pb.product_id = p.id
				WHERE bs.id = $1
			`, item.BatchStorageID).Scan(&availableQty, &batchID, &productID, &productName, &batchSKU, &unitPrice)

			if errBatch != nil {
				if errors.Is(errBatch, sql.ErrNoRows) {
					return fmt.Errorf("batch storage dengan ID %s tidak ditemukan", item.BatchStorageID)
				}
				return fmt.Errorf("gagal mengambil informasi batch: %w", errBatch)
			}

			// Verify batch has enough quantity
			if availableQty < item.Quantity {
				return fmt.Errorf("stok tidak cukup untuk %s: tersedia %g, diminta %g", productName, availableQty, item.Quantity)
			}

			// Insert order detail with batch_storage_id
			var detailID string
			errDetail := tx.QueryRow(`
				INSERT INTO sales_order_detail 
				(sales_order_id, batch_storage_id, quantity, unit_price) 
				VALUES ($1, $2, $3, $4) 
				RETURNING id`,
				orderID, item.BatchStorageID, item.Quantity, item.UnitPrice).Scan(&detailID)

			if errDetail != nil {
				return fmt.Errorf("gagal menambahkan detail pesanan: %w", errDetail)
			}

			// Update batch_storage quantity
			_, errBatchStorage := tx.Exec(`
					UPDATE batch_storage 
					SET quantity = quantity - $1 
					WHERE id = $2
				`, item.Quantity, item.BatchStorageID)

			if errBatchStorage != nil {
				return fmt.Errorf("gagal memperbarui kuantitas batch storage: %w", errBatchStorage)
			}

			// Update product_batch current_quantity
			_, errBatchQuantity := tx.Exec(`
					UPDATE product_batch 
					SET current_quantity = current_quantity - $1 
					WHERE id = $2
				`, item.Quantity, batchID)

			if errBatchQuantity != nil {
				return fmt.Errorf("gagal memperbarui kuantitas batch: %w", errBatchQuantity)
			}

			// If req.CreateInvoice is true, create a invoice
			if req.CreateInvoice {
				// use CreateSalesInvoice function to create invoice
				invoiceReq := CreateSalesInvoiceRequest{
					SalesOrderID: orderID,
				}

				_, errCreateInvoice := r.CreateSalesInvoice(invoiceReq, userID, tx)
				if errCreateInvoice != nil {
					return fmt.Errorf("gagal membuat faktur penjualan: %w", errCreateInvoice)
				}
			}
		}

		// Set response data
		response.ID = orderID
		response.SerialID = serialID
		response.CustomerID = req.CustomerID
		response.CustomerName = customerName
		response.Status = status
		response.PaymentMethod = req.PaymentMethod
		if paymentDueDate.Valid {
			response.PaymentDueDate = paymentDueDate.Time.Format(time.RFC3339)
		}
		response.TotalAmount = totalAmount
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateSalesOrder updates an existing sales order
func (r *SalesRepositoryImpl) UpdateSalesOrder(req UpdateSalesOrderRequest) (*UpdateSalesOrderResponse, error) {
	var response UpdateSalesOrderResponse

	// Check if order exists
	var status string
	errCheck := r.db.QueryRow("SELECT status FROM sales_order WHERE id = $1", req.ID).Scan(&status)
	if errCheck != nil {
		if errors.Is(errCheck, sql.ErrNoRows) {
			return nil, fmt.Errorf("pesanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa pesanan: %w", errCheck)
	}

	// Validate order can be updated based on status
	if status != "order" {
		return nil, fmt.Errorf("hanya pesanan dengan status 'order' yang dapat diperbarui")
	}

	// Build dynamic SQL query
	setValues := []string{}
	params := []interface{}{}
	paramCount := 1

	// Conditionally add fields to update
	if req.CustomerID != "" {
		setValues = append(setValues, fmt.Sprintf("customer_id = $%d", paramCount))
		params = append(params, req.CustomerID)
		paramCount++
	}

	if req.PaymentMethod != "" {
		setValues = append(setValues, fmt.Sprintf("payment_method = $%d", paramCount))
		params = append(params, req.PaymentMethod)
		paramCount++
	}

	if req.PaymentDueDate != "" {
		setValues = append(setValues, fmt.Sprintf("payment_due_date = $%d", paramCount))
		params = append(params, req.PaymentDueDate)
		paramCount++
	}

	// Construct final query
	query := "UPDATE sales_order SET " + strings.Join(setValues, ", ") + " WHERE id = $" + strconv.Itoa(paramCount) + " RETURNING id, serial_id, customer_id, status, payment_method, payment_due_date"
	params = append(params, req.ID)

	errUpdate := r.db.QueryRow(query, params...).Scan(&response.ID, &response.SerialID, &response.CustomerID,
		&response.Status, &response.PaymentMethod, &response.PaymentDueDate)
	if errUpdate != nil {
		return nil, fmt.Errorf("gagal memperbarui pesanan: %w", errUpdate)
	}

	return &response, nil
}

// CancelSalesOrder cancels a sales order if it's in a cancellable state
func (r *SalesRepositoryImpl) CancelSalesOrder(req CancelSalesOrderRequest, userID string) error {
	// First check if the order exists and its status
	var status string
	errCheck := r.db.QueryRow(
		"SELECT status FROM sales_order WHERE id = $1",
		req.SalesOrderID,
	).Scan(&status)

	if errCheck != nil {
		if errors.Is(errCheck, sql.ErrNoRows) {
			return fmt.Errorf("pesanan dengan ID tersebut tidak ditemukan")
		}
		return fmt.Errorf("gagal memeriksa status pesanan: %w", errCheck)
	}

	// Only allow cancellation for orders in 'order' status
	if status != "order" {
		return fmt.Errorf("hanya pesanan dengan status 'order' yang dapat dibatalkan")
	}

	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Mark the sales order as cancelled
		_, errCancel := tx.Exec(
			"UPDATE sales_order SET status = 'cancel', cancelled_at = NOW(), cancelled_by = $1 WHERE id = $2",
			userID, req.SalesOrderID,
		)
		if errCancel != nil {
			return fmt.Errorf("gagal membatalkan pesanan: %w", errCancel)
		}

		// Get all order details to restore inventory
		rows, errDetails := tx.Query(`
            SELECT sod.id, sod.batch_id, sod.batch_storage_id, sod.quantity 
            FROM sales_order_detail sod
            WHERE sod.sales_order_id = $1
        `, req.SalesOrderID)

		if errDetails != nil {
			return fmt.Errorf("gagal mengambil detail pesanan: %w", errDetails)
		}
		defer rows.Close()

		// Collect all details first before processing them
		type orderDetail struct {
			detailID       string
			batchID        string
			batchStorageID string
			quantity       float64
		}

		var details []orderDetail
		for rows.Next() {
			var detail orderDetail
			if err := rows.Scan(&detail.detailID, &detail.batchID, &detail.batchStorageID, &detail.quantity); err != nil {
				return fmt.Errorf("gagal membaca detail pesanan: %w", err)
			}
			details = append(details, detail)
		}

		if err := rows.Err(); err != nil {
			return fmt.Errorf("terjadi kesalahan saat memproses detail pesanan: %w", err)
		}

		// Now process each detail
		for _, detail := range details {
			// Restore inventory in product_batch
			_, errRestore := tx.Exec(
				"UPDATE product_batch SET current_quantity = current_quantity + $1 WHERE id = $2",
				detail.quantity, detail.batchID,
			)
			if errRestore != nil {
				return fmt.Errorf("gagal mengembalikan stok produk: %w", errRestore)
			}

			// Restore inventory in batch_storage directly using batch_storage_id
			_, errBatchStorage := tx.Exec(`
                UPDATE batch_storage 
                SET quantity = quantity + $1 
                WHERE id = $2
            `, detail.quantity, detail.batchStorageID)
			if errBatchStorage != nil {
				return fmt.Errorf("gagal mengembalikan stok di lokasi penyimpanan: %w", errBatchStorage)
			}
		}

		return nil
	})
}

// AddItemToSalesOrder adds a new item to an existing sales order
func (r *SalesRepositoryImpl) AddItemToSalesOrder(req AddSalesOrderItemRequest) (*UpdateAndCreateItemResponse, error) {
	var response UpdateAndCreateItemResponse
	var status string

	// Check if item already exists in the order
	var existingItemID string
	errCheckItem := r.db.QueryRow("SELECT id FROM sales_order_detail WHERE sales_order_id = $1 AND batch_storage_id = $2", req.SalesOrderID, req.BatchStorageID).Scan(&existingItemID)
	if errCheckItem != nil && !errors.Is(errCheckItem, sql.ErrNoRows) {
		return nil, fmt.Errorf("gagal memeriksa item pesanan: %w", errCheckItem)
	}
	if existingItemID != "" {
		return nil, fmt.Errorf("item ini sudah ada dalam pesanan")
	}

	// Check if sales order exists and if it can be modified
	errCheck := r.db.QueryRow("SELECT status FROM sales_order WHERE id = $1", req.SalesOrderID).Scan(&status)
	if errCheck != nil {
		if errors.Is(errCheck, sql.ErrNoRows) {
			return nil, fmt.Errorf("pesanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa status pesanan: %w", errCheck)
	}

	// Only allow items to be added to orders with status 'order'
	if status != "order" {
		return nil, fmt.Errorf("hanya pesanan dengan status 'order' yang dapat diubah")
	}

	// Get batch_storage information including product and batch details
	var batchStorageID, batchID, productID, productName, batchSKU string
	var storageID string
	var storageQty, unitPrice float64

	errBatchStorage := r.db.QueryRow(`
        SELECT bs.id, bs.batch_id, pb.product_id, p.name, pb.sku, bs.storage_id, bs.quantity, pb.unit_price
        FROM batch_storage bs
        JOIN product_batch pb ON bs.batch_id = pb.id
        JOIN product p ON pb.product_id = p.id
        WHERE bs.id = $1
    `, req.BatchStorageID).Scan(
		&batchStorageID,
		&batchID,
		&productID,
		&productName,
		&batchSKU,
		&storageID,
		&storageQty,
		&unitPrice,
	)

	if errBatchStorage != nil {
		if errors.Is(errBatchStorage, sql.ErrNoRows) {
			return nil, fmt.Errorf("batch storage tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil informasi batch storage: %w", errBatchStorage)
	}

	// Check if quantity requested is available
	if storageQty < req.Quantity {
		return nil, fmt.Errorf("jumlah yang diminta (%g) melebihi stok yang tersedia (%g)", req.Quantity, storageQty)
	}

	// Use provided unit price instead of batch price if specified
	if req.UnitPrice <= 0 {
		req.UnitPrice = unitPrice
	}

	// Execute transaction
	err := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Create a new sales order detail entry with batch_storage_id
		var detailID string
		errDetail := tx.QueryRow(`
            INSERT INTO sales_order_detail 
            (sales_order_id, batch_storage_id, quantity, unit_price) 
            VALUES ($1, $2, $3, $4) 
            RETURNING id`,
			req.SalesOrderID, batchStorageID, req.Quantity, req.UnitPrice,
		).Scan(&detailID)
		if errDetail != nil {
			return fmt.Errorf("gagal menambahkan item ke pesanan: %w", errDetail)
		}

		// Update product_batch quantity
		_, errBatchUpdate := tx.Exec(
			"UPDATE product_batch SET current_quantity = current_quantity - $1 WHERE id = $2",
			req.Quantity, batchID,
		)
		if errBatchUpdate != nil {
			return fmt.Errorf("gagal memperbarui stok batch: %w", errBatchUpdate)
		}

		// Update batch_storage quantity
		_, errBatchStorageUpdate := tx.Exec(
			"UPDATE batch_storage SET quantity = quantity - $1 WHERE id = $2",
			req.Quantity, batchStorageID,
		)
		if errBatchStorageUpdate != nil {
			return fmt.Errorf("gagal memperbarui stok di lokasi penyimpanan: %w", errBatchStorageUpdate)
		}

		// Update the order's total_amount
		_, errTotal := tx.Exec(`
            UPDATE sales_order 
            SET total_amount = (
                SELECT COALESCE(SUM(quantity * unit_price), 0) 
                FROM sales_order_detail 
                WHERE sales_order_id = $1
            ),
            updated_at = NOW()
            WHERE id = $1`,
			req.SalesOrderID,
		)
		if errTotal != nil {
			return fmt.Errorf("gagal memperbarui total harga pesanan: %w", errTotal)
		}

		// Set response values
		response.DetailID = detailID
		response.ProductID = productID
		response.ProductName = productName
		response.BatchID = batchID
		response.BatchSKU = batchSKU
		response.BatchStorageID = batchStorageID
		response.Quantity = req.Quantity
		response.UnitPrice = req.UnitPrice
		response.TotalPrice = req.Quantity * req.UnitPrice

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteSalesOrderItem deletes an item from a sales order and restores inventory
func (r *SalesRepositoryImpl) DeleteSalesOrderItem(req DeleteSalesOrderItemRequest) error {
	// Check if sales order exists and is in editable state
	var status string
	err := r.db.QueryRow("SELECT status FROM sales_order WHERE id = $1", req.SalesOrderID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("pesanan tidak ditemukan")
		}
		return fmt.Errorf("gagal memeriksa pesanan: %w", err)
	}

	// Only allow deletion if order is in initial state
	if status != "order" {
		return fmt.Errorf("item hanya dapat dihapus pada pesanan dengan status 'order'")
	}

	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Get details for the item to be deleted including batch_storage_id
		var batchID string
		var batchStorageID string
		var quantity, unitPrice float64
		queryErr := tx.QueryRow(`
            SELECT 
                bs.batch_id, sod.batch_storage_id, sod.quantity, sod.unit_price
            FROM sales_order_detail sod
            JOIN batch_storage bs ON sod.batch_storage_id = bs.id
            WHERE sod.id = $1 AND sod.sales_order_id = $2
        `, req.DetailID, req.SalesOrderID).Scan(&batchID, &batchStorageID, &quantity, &unitPrice)

		if queryErr != nil {
			if errors.Is(queryErr, sql.ErrNoRows) {
				return fmt.Errorf("item tidak ditemukan dalam pesanan")
			}
			return fmt.Errorf("gagal mengambil detail item: %w", queryErr)
		}

		// Update batch_storage to restore quantity directly using batch_storage_id
		_, updateErr := tx.Exec(`
            UPDATE batch_storage
            SET quantity = quantity + $1
            WHERE id = $2
        `, quantity, batchStorageID)
		if updateErr != nil {
			return fmt.Errorf("gagal memulihkan kuantitas di penyimpanan: %w", updateErr)
		}

		// Update product_batch to restore total quantity
		_, updateErr = tx.Exec(`
            UPDATE product_batch
            SET current_quantity = current_quantity + $1
            WHERE id = $2
        `, quantity, batchID)
		if updateErr != nil {
			return fmt.Errorf("gagal memulihkan kuantitas batch: %w", updateErr)
		}

		// Delete the order detail
		_, deleteDetailErr := tx.Exec(`
            DELETE FROM sales_order_detail
            WHERE id = $1
        `, req.DetailID)
		if deleteDetailErr != nil {
			return fmt.Errorf("gagal menghapus item pesanan: %w", deleteDetailErr)
		}

		// Update total order amount
		_, updateOrderErr := tx.Exec(`
            UPDATE sales_order
            SET total_amount = total_amount - $1,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = $2
        `, quantity*unitPrice, req.SalesOrderID)
		if updateOrderErr != nil {
			return fmt.Errorf("gagal memperbarui total pesanan: %w", updateOrderErr)
		}

		return nil
	})
}

// UpdateSalesOrderItem updates an item in a sales order with new quantity or price
func (r *SalesRepositoryImpl) UpdateSalesOrderItem(req UpdateSalesOrderItemRequest) (*UpdateAndCreateItemResponse, error) {
	var response UpdateAndCreateItemResponse
	var status string
	var currentQty, currentPrice float64
	var batchStorageID string

	// Check if sales order exists and if it's in a modifiable state
	errCheck := r.db.QueryRow("SELECT status FROM sales_order WHERE id = $1", req.SalesOrderID).Scan(&status)
	if errCheck != nil {
		if errors.Is(errCheck, sql.ErrNoRows) {
			return nil, fmt.Errorf("pesanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa status pesanan: %w", errCheck)
	}

	// Only allow items to be updated if order's status is 'order'
	if status != "order" {
		return nil, fmt.Errorf("hanya pesanan dengan status 'order' yang dapat diubah")
	}

	// Get current detail information including batch_storage_id
	errDetail := r.db.QueryRow(`
		SELECT sod.quantity, sod.unit_price, sod.batch_storage_id 
		FROM sales_order_detail sod
		WHERE sod.id = $1 AND sod.sales_order_id = $2
	`, req.DetailID, req.SalesOrderID).Scan(
		&currentQty,
		&currentPrice,
		&batchStorageID,
	)

	if errDetail != nil {
		if errors.Is(errDetail, sql.ErrNoRows) {
			return nil, fmt.Errorf("item pesanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mendapatkan informasi item: %w", errDetail)
	}

	// Check if batch_storage_id is changing
	isChangingStorage := req.BatchStorageID != "" && req.BatchStorageID != batchStorageID

	// Get current batch info
	var batchID, productID, storageID string
	var productName, batchSKU string

	errBatchInfo := r.db.QueryRow(`
		SELECT bs.batch_id, pb.product_id, p.name, pb.sku, bs.storage_id
		FROM batch_storage bs
		JOIN product_batch pb ON bs.batch_id = pb.id
		JOIN product p ON pb.product_id = p.id
		WHERE bs.id = $1
	`, batchStorageID).Scan(
		&batchID,
		&productID,
		&productName,
		&batchSKU,
		&storageID,
	)

	if errBatchInfo != nil {
		return nil, fmt.Errorf("gagal mendapatkan informasi batch: %w", errBatchInfo)
	}

	// If changing storage, get new batch storage info
	var newBatchID, newProductID, newStorageID string
	if isChangingStorage {
		errBatchStorage := r.db.QueryRow(`
			SELECT bs.batch_id, pb.product_id, bs.storage_id
			FROM batch_storage bs
			JOIN product_batch pb ON bs.batch_id = pb.id
			WHERE bs.id = $1
		`, req.BatchStorageID).Scan(&newBatchID, &newProductID, &newStorageID)

		if errBatchStorage != nil {
			if errors.Is(errBatchStorage, sql.ErrNoRows) {
				return nil, fmt.Errorf("lokasi batch baru tidak ditemukan")
			}
			return nil, fmt.Errorf("gagal mendapatkan informasi lokasi batch baru: %w", errBatchStorage)
		}

		// Verify product is the same when changing batch storage
		if newProductID != productID {
			return nil, fmt.Errorf("tidak dapat mengubah lokasi penyimpanan ke produk yang berbeda")
		}
	}

	// Set new values or keep current ones if not provided
	newQty := currentQty
	newPrice := currentPrice

	if req.Quantity > 0 {
		newQty = req.Quantity
	}

	if req.UnitPrice > 0 {
		newPrice = req.UnitPrice
	}

	// If quantity is unchanged and only price is updated, and no storage change, simple update
	if newQty == currentQty && !isChangingStorage {
		err := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
			// Update the sales order detail with new price
			_, errUpdate := tx.Exec(`
				UPDATE sales_order_detail 
				SET unit_price = $1, updated_at = NOW()
				WHERE id = $2
			`, newPrice, req.DetailID)
			if errUpdate != nil {
				return fmt.Errorf("gagal memperbarui harga item: %w", errUpdate)
			}

			// Update the order's total_amount
			_, errTotal := tx.Exec(`
				UPDATE sales_order SET total_amount = (
					SELECT COALESCE(SUM(quantity * unit_price), 0)
					FROM sales_order_detail
					WHERE sales_order_id = $1
				), updated_at = NOW()
				WHERE id = $1
			`, req.SalesOrderID)
			if errTotal != nil {
				return fmt.Errorf("gagal memperbarui total harga pesanan: %w", errTotal)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	} else {
		// Calculate quantity difference
		qtyDifference := newQty - currentQty

		// Execute complex update transaction
		err := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
			// If we're changing storage location
			if isChangingStorage {
				// Return quantity to original batch and storage
				_, errRestore := tx.Exec(`
					UPDATE product_batch 
					SET current_quantity = current_quantity + $1 
					WHERE id = $2
				`, currentQty, batchID)
				if errRestore != nil {
					return fmt.Errorf("gagal mengembalikan stok batch lama: %w", errRestore)
				}

				_, errRestoreStorage := tx.Exec(`
					UPDATE batch_storage 
					SET quantity = quantity + $1 
					WHERE id = $2
				`, currentQty, batchStorageID)
				if errRestoreStorage != nil {
					return fmt.Errorf("gagal mengembalikan stok di lokasi penyimpanan lama: %w", errRestoreStorage)
				}

				// Validate new batch storage has enough quantity
				var availableQty float64
				errAvail := tx.QueryRow(`
					SELECT quantity FROM batch_storage WHERE id = $1
				`, req.BatchStorageID).Scan(&availableQty)
				if errAvail != nil {
					return fmt.Errorf("gagal memeriksa ketersediaan stok di lokasi baru: %w", errAvail)
				}

				if availableQty < newQty {
					return fmt.Errorf("stok tidak cukup di lokasi baru: tersedia %g, diminta %g", availableQty, newQty)
				}

				// Take quantity from new batch and storage
				_, errDeduct := tx.Exec(`
					UPDATE product_batch 
					SET current_quantity = current_quantity - $1 
					WHERE id = $2
				`, newQty, newBatchID)
				if errDeduct != nil {
					return fmt.Errorf("gagal mengurangi stok batch baru: %w", errDeduct)
				}

				_, errDeductStorage := tx.Exec(`
					UPDATE batch_storage 
					SET quantity = quantity - $1 
					WHERE id = $2
				`, newQty, req.BatchStorageID)
				if errDeductStorage != nil {
					return fmt.Errorf("gagal mengurangi stok di lokasi penyimpanan baru: %w", errDeductStorage)
				}

				// Update sales order detail with new batch storage
				_, errUpdateDetail := tx.Exec(`
					UPDATE sales_order_detail 
					SET batch_storage_id = $1, quantity = $2, unit_price = $3, updated_at = NOW()
					WHERE id = $4
				`, req.BatchStorageID, newQty, newPrice, req.DetailID)
				if errUpdateDetail != nil {
					return fmt.Errorf("gagal memperbarui detail pesanan: %w", errUpdateDetail)
				}

				// Update response with new batch info
				batchID = newBatchID
				// Get updated batch SKU
				if err := tx.QueryRow("SELECT sku FROM product_batch WHERE id = $1", newBatchID).Scan(&batchSKU); err != nil {
					return fmt.Errorf("gagal mendapatkan informasi SKU batch baru: %w", err)
				}
			} else {
				// Just updating quantity of the same batch_storage
				if qtyDifference != 0 {
					// Check if we have enough quantity if increasing
					if qtyDifference > 0 {
						var availableQty float64
						errAvail := tx.QueryRow(`
							SELECT quantity FROM batch_storage WHERE id = $1
						`, batchStorageID).Scan(&availableQty)
						if errAvail != nil {
							return fmt.Errorf("gagal memeriksa ketersediaan stok: %w", errAvail)
						}

						if availableQty < qtyDifference {
							return fmt.Errorf("stok tidak mencukupi, tersedia: %g, diminta tambahan: %g", availableQty, qtyDifference)
						}
					}

					// Update product_batch total quantity
					_, err := tx.Exec(`
						UPDATE product_batch 
						SET current_quantity = current_quantity - $1 
						WHERE id = $2
					`, qtyDifference, batchID)
					if err != nil {
						return fmt.Errorf("gagal memperbarui stok batch: %w", err)
					}

					// Update batch_storage quantity
					_, err = tx.Exec(`
						UPDATE batch_storage 
						SET quantity = quantity - $1 
						WHERE id = $2
					`, qtyDifference, batchStorageID)
					if err != nil {
						return fmt.Errorf("gagal memperbarui stok di lokasi penyimpanan: %w", err)
					}
				}

				// Update sales order detail with new quantity/price
				_, err := tx.Exec(`
					UPDATE sales_order_detail 
					SET quantity = $1, unit_price = $2, updated_at = NOW()
					WHERE id = $3
				`, newQty, newPrice, req.DetailID)
				if err != nil {
					return fmt.Errorf("gagal memperbarui detail pesanan: %w", err)
				}
			}

			// Update the order's total_amount
			_, err := tx.Exec(`
				UPDATE sales_order SET total_amount = (
					SELECT COALESCE(SUM(quantity * unit_price), 0)
					FROM sales_order_detail
					WHERE sales_order_id = $1
				), updated_at = NOW()
				WHERE id = $1
			`, req.SalesOrderID)
			if err != nil {
				return fmt.Errorf("gagal memperbarui total harga pesanan: %w", err)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	// Set response values
	response.DetailID = req.DetailID
	response.ProductID = productID
	response.ProductName = productName
	response.BatchID = batchID
	response.BatchSKU = batchSKU
	if isChangingStorage {
		response.BatchStorageID = req.BatchStorageID
	} else {
		response.BatchStorageID = batchStorageID
	}
	response.Quantity = newQty
	response.UnitPrice = newPrice
	response.TotalPrice = newQty * newPrice

	return &response, nil
}

// GetSalesInvoices retrieves a paginated list of sales invoices with filtering options
func (r *SalesRepositoryImpl) GetSalesInvoices(req GetSalesInvoicesRequest) ([]GetSalesInvoicesResponse, int, error) {
	// Build base query for fetching sales invoices
	baseQuery := `
        SELECT si.id, si.serial_id, si.sales_order_id, so.serial_id AS sales_order_serial,
               so.customer_id, c.name AS customer_name, si.invoice_date, si.total_amount,
               CASE 
                 WHEN si.cancelled_at IS NOT NULL THEN 'cancelled'
                 WHEN EXISTS(SELECT 1 FROM sales_order_return sor 
                            JOIN sales_order_detail sod ON sor.sales_detail_id = sod.id
                            WHERE sod.sales_order_id = so.id AND sor.return_status = 'returned') 
                    AND NOT EXISTS(SELECT 1 FROM sales_order_detail sod 
                                  WHERE sod.sales_order_id = so.id 
                                  AND NOT EXISTS(SELECT 1 FROM sales_order_return sor 
                                               WHERE sor.sales_detail_id = sod.id AND sor.return_status = 'returned')) 
                    THEN 'returned'
                 WHEN EXISTS(SELECT 1 FROM sales_order_return sor 
                            JOIN sales_order_detail sod ON sor.sales_detail_id = sod.id
                            WHERE sod.sales_order_id = so.id AND sor.return_status = 'returned') THEN 'partially_returned'
                 ELSE 'active'
               END AS status,
               EXISTS(SELECT 1 FROM delivery_note dn WHERE dn.sales_invoice_id = si.id AND dn.cancelled_at IS NULL) AS has_delivery_note,
               si.created_by, si.created_at, si.cancelled_at
        FROM sales_invoice si
        JOIN sales_order so ON si.sales_order_id = so.id
        JOIN customer c ON so.customer_id = c.id
        WHERE 1=1`

	// Create count query
	countQuery := `
        SELECT COUNT(*) 
        FROM sales_invoice si
        JOIN sales_order so ON si.sales_order_id = so.id
        JOIN customer c ON so.customer_id = c.id
        WHERE 1=1`

	// Initialize query builders
	qb := utils.NewQueryBuilder(baseQuery)
	countQb := utils.NewQueryBuilder(countQuery)

	// Add filters to both queries
	if req.CustomerID != "" {
		qb.AddFilter("so.customer_id =", req.CustomerID)
		countQb.AddFilter("so.customer_id =", req.CustomerID)
	}

	if req.Status != "" {
		if req.Status == "cancelled" {
			qb.AddFilter("si.cancelled_at IS NOT", nil)
			countQb.AddFilter("si.cancelled_at IS NOT", nil)
		} else if req.Status == "active" {
			qb.AddFilter("si.cancelled_at IS", nil)
			countQb.AddFilter("si.cancelled_at IS", nil)
			// Additional conditions for active status can be added here if needed
		}
		// Handle partially_returned and returned statuses with additional subqueries if needed
	}

	if req.SerialID != "" {
		qb.AddFilter("si.serial_id ILIKE", "%"+req.SerialID+"%")
		countQb.AddFilter("si.serial_id ILIKE", "%"+req.SerialID+"%")
	}

	if req.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, req.StartDate)
		if err == nil {
			qb.AddFilter("si.invoice_date >=", startDate)
			countQb.AddFilter("si.invoice_date >=", startDate)
		}
	}

	if req.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err == nil {
			// Add one day to include the entire end date
			endDate = endDate.Add(24 * time.Hour)
			qb.AddFilter("si.invoice_date <", endDate)
			countQb.AddFilter("si.invoice_date <", endDate)
		}
	}

	// Add sorting
	sortBy := "si.created_at"
	if req.SortBy != "" {
		switch req.SortBy {
		case "serial_id", "invoice_date", "total_amount", "created_at":
			sortBy = "si." + req.SortBy
		case "customer_name":
			sortBy = "c.name"
		case "status":
			sortBy = "status"
		}
	}

	sortOrder := "DESC"
	if req.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	// Append sorting to main query only (not count query)
	qb.Query.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder))

	// Add pagination to main query
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = utils.DefaultPageSize
	}
	page := req.Page
	if page <= 0 {
		page = utils.DefaultPage
	}
	qb.AddPagination(pageSize, page)

	// Build final queries
	query, params := qb.Build()
	countQuery, countParams := countQb.Build()

	// Execute count query
	var totalItems int
	err := r.db.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting sales invoices: %w", err)
	}

	// Execute main query
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying sales invoices: %w", err)
	}
	defer rows.Close()

	// Parse results
	var invoices []GetSalesInvoicesResponse
	for rows.Next() {
		var invoice GetSalesInvoicesResponse
		var invoiceDate, createdAt time.Time
		var cancelledAt sql.NullTime
		var hasDeliveryNote bool

		err := rows.Scan(
			&invoice.ID,
			&invoice.SerialID,
			&invoice.SalesOrderID,
			&invoice.SalesOrderSerial,
			&invoice.CustomerID,
			&invoice.CustomerName,
			&invoiceDate,
			&invoice.TotalAmount,
			&invoice.Status,
			&hasDeliveryNote,
			&invoice.CreatedBy,
			&createdAt,
			&cancelledAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning sales invoice row: %w", err)
		}

		// Format dates for response
		invoice.InvoiceDate = invoiceDate.Format(time.RFC3339)
		invoice.CreatedAt = createdAt.Format(time.RFC3339)
		invoice.HasDeliveryNote = hasDeliveryNote

		if cancelledAt.Valid {
			invoice.CancelledAt = cancelledAt.Time.Format(time.RFC3339)
		}

		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating sales invoice rows: %w", err)
	}

	return invoices, totalItems, nil
}

// CreateSalesInvoice creates a new invoice from a order
func (r *SalesRepositoryImpl) CreateSalesInvoice(req CreateSalesInvoiceRequest, userID string, tx *sql.Tx) (*CreateSalesInvoiceResponse, error) {
	var response CreateSalesInvoiceResponse

	// Function to execute the invoice creation logic
	createInvoiceFunc := func(tx *sql.Tx) error {
		// Check if order exists and is in 'order' status
		var orderStatus, orderSerial string
		var customerId string
		var customerName string
		var totalAmount float64

		err := tx.QueryRow(`
            SELECT so.status, so.serial_id, so.customer_id, c.name, so.total_amount
            FROM sales_order so
            JOIN customer c ON so.customer_id = c.id
            WHERE so.id = $1
        `, req.SalesOrderID).Scan(&orderStatus, &orderSerial, &customerId, &customerName, &totalAmount)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("pesanan penjualan tidak ditemukan")
			}
			return fmt.Errorf("gagal memeriksa pesanan: %w", err)
		}

		// Validate order status
		if orderStatus != "order" {
			return fmt.Errorf("pesanan dalam status %s", orderStatus)
		}

		// Check if invoice already exists for this order
		var existingInvoice string
		err = tx.QueryRow(`
            SELECT id FROM sales_invoice 
            WHERE sales_order_id = $1 AND cancelled_at IS NULL
        `, req.SalesOrderID).Scan(&existingInvoice)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("gagal memeriksa faktur yang ada: %w", err)
		}

		if existingInvoice != "" {
			return fmt.Errorf("faktur sudah ada untuk pesanan ini")
		}

		var invoiceID string
		var invoiceDate time.Time

		// Generate a new serial ID for the invoice
		serialID, err := utils.GenerateNextSerialID(tx, "SI")
		if err != nil {
			return fmt.Errorf("gagal membuat ID faktur: %w", err)
		}

		// Create invoice
		err = tx.QueryRow(`
            INSERT INTO sales_invoice (
                sales_order_id, serial_id, total_amount, created_by
            ) VALUES ($1, $2, $3, $4)
            RETURNING id, serial_id, invoice_date
        `, req.SalesOrderID, serialID, totalAmount, userID).Scan(&invoiceID, &serialID, &invoiceDate)

		if err != nil {
			return fmt.Errorf("gagal membuat faktur: %w", err)
		}

		// Update order status to 'invoice'
		_, err = tx.Exec(`
            UPDATE sales_order 
            SET status = 'invoice', updated_at = NOW() 
            WHERE id = $1
        `, req.SalesOrderID)

		if err != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
		}

		// Get all order details for logging - collect items first to avoid connection issues
		type detailItem struct {
			detailID       string
			batchStorageID string
			batchID        string
			storageID      string
			quantity       float64
			unitPrice      float64
		}

		var items []detailItem

		rows, err := tx.Query(`
            SELECT 
                sod.id, 
                sod.batch_storage_id, 
                bs.batch_id,
                bs.storage_id,
                sod.quantity, 
                sod.unit_price
            FROM sales_order_detail sod
            JOIN batch_storage bs ON sod.batch_storage_id = bs.id
            WHERE sod.sales_order_id = $1
        `, req.SalesOrderID)

		if err != nil {
			return fmt.Errorf("gagal mengambil detail pesanan untuk pencatatan: %w", err)
		}

		// Important: properly close rows when done
		defer rows.Close()

		// Read all items at once to avoid connection issues
		for rows.Next() {
			var item detailItem
			if err := rows.Scan(
				&item.detailID,
				&item.batchStorageID,
				&item.batchID,
				&item.storageID,
				&item.quantity,
				&item.unitPrice,
			); err != nil {
				return fmt.Errorf("gagal memindai detail item: %w", err)
			}
			items = append(items, item)
		}

		// Check for any errors during row iteration
		if err = rows.Err(); err != nil {
			return fmt.Errorf("terjadi kesalahan saat membaca detail item: %w", err)
		}

		// Now that rows are closed, process the items for logging
		for _, item := range items {
			// Log inventory movement
			_, err = tx.Exec(`
                INSERT INTO inventory_log (
                    batch_id, 
                    storage_id, 
                    user_id, 
                    sales_order_id, 
                    action, 
                    quantity, 
                    description
                ) VALUES ($1, $2, $3, $4, $5, $6, $7)
            `,
				item.batchID,
				item.storageID,
				userID,
				req.SalesOrderID,
				"invoice",
				item.quantity,
				fmt.Sprintf("Pembuatan faktur %s", serialID))

			if err != nil {
				return fmt.Errorf("gagal mencatat log inventaris: %w", err)
			}
		}

		// Create financial transaction log
		_, err = tx.Exec(`
            INSERT INTO financial_transaction_log (
                user_id,
                amount,
                type,
                sales_order_id,
                description,
                transaction_date,
                is_system
            ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        `,
			userID,
			totalAmount,
			"sales_invoice",
			req.SalesOrderID,
			fmt.Sprintf("Pembuatan faktur penjualan %s", serialID),
			invoiceDate,
			true)

		if err != nil {
			return fmt.Errorf("gagal mencatat transaksi keuangan: %w", err)
		}

		// Set response
		response.ID = invoiceID
		response.SerialID = serialID
		response.SalesOrderID = req.SalesOrderID
		response.SalesOrderSerial = orderSerial
		response.CustomerID = customerId
		response.CustomerName = customerName
		response.InvoiceDate = invoiceDate.Format(time.RFC3339)
		response.TotalAmount = totalAmount
		response.Status = "active"
		response.CreatedBy = userID
		response.CreatedAt = invoiceDate.Format(time.RFC3339)

		return nil
	}

	// If a transaction was provided, use it directly
	if tx != nil {
		// Create a savepoint to allow rolling back just this part of the transaction
		_, err := tx.Exec("SAVEPOINT create_invoice")
		if err != nil {
			return nil, fmt.Errorf("failed to create savepoint: %w", err)
		}

		// Execute the invoice creation
		err = createInvoiceFunc(tx)

		if err != nil {
			// If there's an error, roll back to the savepoint
			_, rbErr := tx.Exec("ROLLBACK TO SAVEPOINT create_invoice")
			if rbErr != nil {
				return nil, fmt.Errorf("error creating invoice: %v, and failed to rollback: %w", err, rbErr)
			}
			return nil, err
		}

		// Release the savepoint on success
		_, err = tx.Exec("RELEASE SAVEPOINT create_invoice")
		if err != nil {
			return nil, fmt.Errorf("failed to release savepoint: %w", err)
		}

		return &response, nil
	}

	// Otherwise, create a new transaction
	err := utils.WithTransaction(r.db, func(newTx *sql.Tx) error {
		return createInvoiceFunc(newTx)
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CancelSalesInvoice cancels an existing sales invoice
func (r *SalesRepositoryImpl) CancelSalesInvoice(req CancelSalesInvoiceRequest, userID string) error {
	// Check if invoice exists
	var salesOrderID string
	var cancelled bool

	err := r.db.QueryRow(`
        SELECT sales_order_id, cancelled_at IS NOT NULL 
        FROM sales_invoice 
        WHERE id = $1
    `, req.InvoiceID).Scan(&salesOrderID, &cancelled)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("faktur tidak ditemukan")
		}
		return fmt.Errorf("gagal memeriksa faktur: %w", err)
	}

	// Check if already cancelled
	if cancelled {
		return fmt.Errorf("faktur ini sudah dibatalkan")
	}

	// Check if delivery note exists
	var deliveryExists bool
	err = r.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM delivery_note 
            WHERE sales_invoice_id = $1 AND cancelled_at IS NULL
        )
    `, req.InvoiceID).Scan(&deliveryExists)

	if err != nil {
		return fmt.Errorf("gagal memeriksa surat jalan: %w", err)
	}

	if deliveryExists {
		return fmt.Errorf("tidak dapat membatalkan faktur karena sudah memiliki surat jalan aktif")
	}

	// Execute transaction
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Mark invoice as cancelled
		_, err := tx.Exec(`
            UPDATE sales_invoice 
            SET cancelled_at = NOW(), cancelled_by = $1 
            WHERE id = $2
        `, userID, req.InvoiceID)

		if err != nil {
			return fmt.Errorf("gagal membatalkan faktur: %w", err)
		}

		// Revert sales order to 'order' status
		_, err = tx.Exec(`
            UPDATE sales_order 
            SET status = 'order', updated_at = NOW() 
            WHERE id = $1
        `, salesOrderID)

		if err != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
		}

		return nil
	})
}

// ReturnInvoiceItems handles returning items from a sales invoice
func (r *SalesRepositoryImpl) ReturnInvoiceItems(req ReturnInvoiceItemsRequest, userID string) (*ReturnInvoiceItemsResponse, error) {
	var response ReturnInvoiceItemsResponse
	var salesOrderID string
	var totalReturnQuantity float64
	var isInvoiceCancelled bool

	// Check if invoice exists and get sales order ID
	err := r.db.QueryRow(`
        SELECT i.sales_order_id, i.cancelled_at IS NOT NULL
        FROM sales_invoice i
        WHERE i.id = $1
    `, req.InvoiceID).Scan(&salesOrderID, &isInvoiceCancelled)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("faktur dengan ID %s tidak ditemukan", req.InvoiceID)
		}
		return nil, fmt.Errorf("gagal memeriksa faktur: %w", err)
	}

	// Check if invoice is already cancelled
	if isInvoiceCancelled {
		return nil, errors.New("faktur sudah dibatalkan, tidak dapat diproses pengembalian")
	}

	// Check if a delivery note has been created
	var hasDeliveryNote bool
	err = r.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM delivery_note 
            WHERE sales_invoice_id = $1 AND cancelled_at IS NULL
        )
    `, req.InvoiceID).Scan(&hasDeliveryNote)

	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa surat jalan: %w", err)
	}

	if hasDeliveryNote {
		return nil, errors.New("faktur telah memiliki surat jalan, pengembalian harus dilakukan melalui surat jalan")
	}

	// Verify all detail IDs exist and belong to this sales order
	for _, item := range req.ReturnItems {
		var detailExists bool
		var currentQuantity float64
		var previousReturnedQty float64

		err = r.db.QueryRow(`
            SELECT EXISTS(
                SELECT 1 FROM sales_order_detail 
                WHERE id = $1 AND sales_order_id = $2
            ), 
            (SELECT quantity FROM sales_order_detail WHERE id = $1),
            COALESCE(
                (SELECT SUM(return_quantity) FROM sales_order_return 
                WHERE sales_detail_id = $1 AND return_status = 'completed'),
                0
            )
        `, item.DetailID, salesOrderID).Scan(&detailExists, &currentQuantity, &previousReturnedQty)

		if err != nil {
			return nil, fmt.Errorf("gagal memverifikasi item dengan ID %s: %w", item.DetailID, err)
		}

		if !detailExists {
			return nil, fmt.Errorf("item dengan ID %s tidak ditemukan dalam pesanan ini", item.DetailID)
		}

		// Check if return quantity is valid
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("kuantitas pengembalian harus lebih dari 0")
		}

		if item.Quantity > (currentQuantity - previousReturnedQty) {
			return nil, fmt.Errorf("kuantitas pengembalian %f melebihi kuantitas yang tersedia %f",
				item.Quantity, (currentQuantity - previousReturnedQty))
		}

		totalReturnQuantity += item.Quantity

		// Validate storage returns match the total quantity
		var totalStorageQty float64
		for _, storage := range item.StorageReturns {
			totalStorageQty += storage.Quantity
		}

		if !utils.FloatEquals(totalStorageQty, item.Quantity) {
			return nil, fmt.Errorf("total kuantitas dari semua penyimpanan (%f) tidak sama dengan kuantitas pengembalian (%f)",
				totalStorageQty, item.Quantity)
		}
	}

	// Execute transaction
	err = utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Generate return ID
		returnID := uuid.New().String()
		response.ReturnID = returnID

		// Process each return item
		for _, item := range req.ReturnItems {
			// Get product and batch information
			var productID, batchID string
			if err := tx.QueryRow(`
                SELECT product_id, batch_id 
                FROM sales_order_detail 
                WHERE id = $1
            `, item.DetailID).Scan(&productID, &batchID); err != nil {
				return fmt.Errorf("gagal mendapatkan informasi produk: %w", err)
			}

			// Create return record for this item
			remainingQty := 0.0
			if err := tx.QueryRow(`
                SELECT quantity - $1 - COALESCE(
                    (SELECT SUM(return_quantity) FROM sales_order_return 
                    WHERE sales_detail_id = $2 AND return_status = 'completed'),
                    0
                )
                FROM sales_order_detail
                WHERE id = $2
            `, item.Quantity, item.DetailID).Scan(&remainingQty); err != nil {
				return fmt.Errorf("gagal menghitung sisa kuantitas: %w", err)
			}

			// Insert return record
			if _, err := tx.Exec(`
                INSERT INTO sales_order_return (
                    id, return_source, sales_order_id, sales_detail_id,
                    return_quantity, remaining_quantity, return_reason, 
                    return_status, returned_by, returned_at
                ) VALUES (
                    $1, 'invoice', $2, $3, $4, $5, $6, 'completed', $7, NOW()
                )
            `, returnID, salesOrderID, item.DetailID, item.Quantity, remainingQty,
				req.ReturnReason, userID); err != nil {
				return fmt.Errorf("gagal membuat catatan pengembalian: %w", err)
			}

			// Process storage returns and restore inventory
			for _, storage := range item.StorageReturns {
				// Insert return batch record
				batchReturnID := uuid.New().String()
				if _, err := tx.Exec(`
                    INSERT INTO sales_order_return_batch (
                        id, sales_return_id, batch_id, return_quantity
                    ) VALUES ($1, $2, $3, $4)
                `, batchReturnID, returnID, batchID, storage.Quantity); err != nil {
					return fmt.Errorf("gagal mencatat batch pengembalian: %w", err)
				}

				// Update batch storage quantity
				if _, err := tx.Exec(`
                    INSERT INTO batch_storage (batch_id, storage_id, quantity)
                    VALUES ($1, $2, $3)
                    ON CONFLICT (batch_id, storage_id) 
                    DO UPDATE SET quantity = batch_storage.quantity + $3
                `, batchID, storage.StorageID, storage.Quantity); err != nil {
					return fmt.Errorf("gagal memperbarui jumlah di penyimpanan: %w", err)
				}

				// Update batch current quantity
				if _, err := tx.Exec(`
                    UPDATE product_batch
                    SET current_quantity = current_quantity + $1, updated_at = NOW()
                    WHERE id = $2
                `, storage.Quantity, batchID); err != nil {
					return fmt.Errorf("gagal memperbarui jumlah batch: %w", err)
				}

				// Create inventory log
				if _, err := tx.Exec(`
                    INSERT INTO inventory_log (
                        batch_id, storage_id, user_id, sales_order_id, action,
                        quantity, description
                    ) VALUES (
                        $1, $2, $3, $4, 'return',
                        $5, $6
                    )
                `, batchID, storage.StorageID, userID, salesOrderID,
					storage.Quantity, fmt.Sprintf("Pengembalian dari faktur")); err != nil {
					return fmt.Errorf("gagal membuat log inventaris: %w", err)
				}
			}
		}

		// Check if this is a full return (all items returned)
		var isFullReturn bool
		var returnedItemsCount int
		if err := tx.QueryRow(`
            WITH total_items AS (
                SELECT COUNT(*) as count, SUM(quantity) as total_qty
                FROM sales_order_detail
                WHERE sales_order_id = $1
            ),
            returned_items AS (
                SELECT 
                    COUNT(DISTINCT d.id) as count,
                    SUM(CASE WHEN r.remaining_quantity <= 0 THEN 1 ELSE 0 END) as fully_returned
                FROM sales_order_detail d
                JOIN sales_order_return r ON d.id = r.sales_detail_id
                WHERE d.sales_order_id = $1 AND r.return_status = 'completed'
            )
            SELECT 
                CASE WHEN ti.count = ri.fully_returned THEN TRUE ELSE FALSE END,
                ri.count
            FROM total_items ti, returned_items ri
        `, salesOrderID).Scan(&isFullReturn, &returnedItemsCount); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("gagal memeriksa status pengembalian penuh: %w", err)
			}
			isFullReturn = false
			returnedItemsCount = len(req.ReturnItems)
		}

		response.IsFullReturn = isFullReturn

		// Update sales order status if all items are fully returned
		if isFullReturn {
			if _, err := tx.Exec(`
                UPDATE sales_order
                SET status = 'return', updated_at = NOW()
                WHERE id = $1
            `, salesOrderID); err != nil {
				return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
			}
		} else {
			// Otherwise mark as partially returned
			if _, err := tx.Exec(`
                UPDATE sales_order
                SET status = 'partially_return', updated_at = NOW()
                WHERE id = $1 AND status != 'return'
            `, salesOrderID); err != nil {
				return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Set response fields
	response.InvoiceID = req.InvoiceID
	response.ReturnedItems = len(req.ReturnItems)
	response.TotalQuantity = totalReturnQuantity
	response.ReturnDate = time.Now().Format(time.RFC3339)
	response.ReturnStatus = "completed"

	return &response, nil
}

// CancelInvoiceReturn cancels a previously processed return on a sales invoice
func (r *SalesRepositoryImpl) CancelInvoiceReturn(req CancelInvoiceReturnRequest, userID string) error {
	// Verify the return exists and is not already cancelled
	var salesOrderID string
	var isCancelled bool
	var returnSource string

	err := r.db.QueryRow(`
        SELECT r.sales_order_id, r.cancelled_at IS NOT NULL, r.return_source
        FROM sales_order_return r
        WHERE r.id = $1
    `, req.ReturnID).Scan(&salesOrderID, &isCancelled, &returnSource)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("pengembalian tidak ditemukan")
		}
		return fmt.Errorf("gagal memeriksa pengembalian: %w", err)
	}

	// Check if return source is from invoice (not delivery note)
	if returnSource != "invoice" {
		return errors.New("pengembalian ini berasal dari surat jalan, bukan faktur")
	}

	// Check if already cancelled
	if isCancelled {
		return errors.New("pengembalian sudah dibatalkan sebelumnya")
	}

	// Execute transaction
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Get all return items to revert the inventory
		rows, err := tx.Query(`
            SELECT 
                r.id,
                r.sales_detail_id,
                r.return_quantity,
                rb.id as batch_return_id,
                rb.batch_id,
                rb.return_quantity as batch_return_quantity,
                sod.product_id
            FROM sales_order_return r
            JOIN sales_order_return_batch rb ON r.id = rb.sales_return_id
            JOIN sales_order_detail sod ON r.sales_detail_id = sod.id
            WHERE r.id = $1
        `, req.ReturnID)

		if err != nil {
			return fmt.Errorf("gagal mengambil detail pengembalian: %w", err)
		}
		defer rows.Close()

		type returnDetail struct {
			returnID       string
			detailID       string
			returnQuantity float64
			batchReturnID  string
			batchID        string
			batchReturnQty float64
			productID      string
		}

		var returnDetails []returnDetail
		for rows.Next() {
			var detail returnDetail
			if err := rows.Scan(
				&detail.returnID,
				&detail.detailID,
				&detail.returnQuantity,
				&detail.batchReturnID,
				&detail.batchID,
				&detail.batchReturnQty,
				&detail.productID,
			); err != nil {
				return fmt.Errorf("gagal memindai data pengembalian: %w", err)
			}
			returnDetails = append(returnDetails, detail)
		}

		if err = rows.Err(); err != nil {
			return fmt.Errorf("error saat membaca data pengembalian: %w", err)
		}

		// Process each return detail to revert inventory
		for _, detail := range returnDetails {
			// Update batch quantity
			if _, err := tx.Exec(`
                UPDATE product_batch
                SET current_quantity = current_quantity - $1,
                    updated_at = NOW()
                WHERE id = $2
            `, detail.batchReturnQty, detail.batchID); err != nil {
				return fmt.Errorf("gagal memperbarui jumlah batch: %w", err)
			}

			// Update sales order detail to remove returned quantity
			if _, err := tx.Exec(`
                UPDATE sales_order_detail
                SET quantity = quantity + $1,
                    updated_at = NOW()
                WHERE id = $2
            `, detail.returnQuantity, detail.detailID); err != nil {
				return fmt.Errorf("gagal memperbarui detail pesanan: %w", err)
			}

			// Log inventory change
			if _, err := tx.Exec(`
                INSERT INTO inventory_log (
                    batch_id, user_id, sales_order_id, action, 
                    quantity, description, log_date
                ) VALUES (
                    $1, $2, $3, 'remove', $4, 
                    'Pembatalan pengembalian barang', NOW()
                )
            `, detail.batchID, userID, salesOrderID, detail.batchReturnQty); err != nil {
				return fmt.Errorf("gagal mencatat perubahan inventaris: %w", err)
			}
		}

		// Mark return as cancelled
		if _, err := tx.Exec(`
            UPDATE sales_order_return
            SET return_status = 'cancelled',
                cancelled_at = NOW(),
                cancelled_by = $1
            WHERE id = $2
        `, userID, req.ReturnID); err != nil {
			return fmt.Errorf("gagal membatalkan pengembalian: %w", err)
		}

		// Determine current status based on remaining returns
		var hasActiveReturns bool
		if err := tx.QueryRow(`
            SELECT EXISTS (
                SELECT 1 FROM sales_order_return
                WHERE sales_order_id = $1
                AND return_status = 'completed'
                AND cancelled_at IS NULL
            )
        `, salesOrderID).Scan(&hasActiveReturns); err != nil {
			return fmt.Errorf("gagal memeriksa status pengembalian: %w", err)
		}

		// Update sales order status based on remaining returns
		var newStatus string
		if hasActiveReturns {
			newStatus = "partially_return"
		} else {
			newStatus = "invoice" // Return to invoice status if no active returns
		}

		if _, err := tx.Exec(`
            UPDATE sales_order
            SET status = $1,
                updated_at = NOW()
            WHERE id = $2
        `, newStatus, salesOrderID); err != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
		}

		return nil
	})
}

// CreateDeliveryNote creates a new delivery note from a sales invoice
func (r *SalesRepositoryImpl) CreateDeliveryNote(req CreateDeliveryNoteRequest, userID string) (*CreateDeliveryNoteResponse, error) {
	// Verify the invoice exists and is not cancelled
	var salesOrderID string
	var invoiceSerialID string
	var salesOrderSerialID string
	var hasDeliveryNote bool
	var invoiceCancelled bool
	var customerName string

	err := r.db.QueryRow(`
        SELECT 
            s.id, 
            i.serial_id, 
            s.serial_id,
            c.name as customer_name,
            i.cancelled_at IS NOT NULL,
            EXISTS (
                SELECT 1 FROM delivery_note dn 
                WHERE dn.sales_invoice_id = i.id AND dn.cancelled_at IS NULL
            )
        FROM sales_invoice i
        JOIN sales_order s ON i.sales_order_id = s.id
        JOIN customer c ON s.customer_id = c.id
        WHERE i.id = $1
    `, req.SalesInvoiceID).Scan(&salesOrderID, &invoiceSerialID, &salesOrderSerialID, &customerName, &invoiceCancelled, &hasDeliveryNote)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("faktur penjualan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa faktur penjualan: %w", err)
	}

	if invoiceCancelled {
		return nil, errors.New("tidak dapat membuat surat jalan dari faktur yang sudah dibatalkan")
	}

	if hasDeliveryNote {
		return nil, errors.New("faktur ini sudah memiliki surat jalan aktif")
	}

	// Create transaction to handle serial number generation and delivery note creation
	var response CreateDeliveryNoteResponse
	err = utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Get next serial number for delivery note
		serialID, err := utils.GenerateNextSerialID(tx, "DN")
		if err != nil {
			return fmt.Errorf("gagal membuat nomor surat jalan: %w", err)
		}

		// Create delivery note
		var deliveryNoteID string
		err = tx.QueryRow(`
            INSERT INTO delivery_note (
                serial_id, sales_order_id, sales_invoice_id, 
                delivery_date, driver_name, recipient_name, created_by
            ) VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id
        `, serialID, salesOrderID, req.SalesInvoiceID, req.DeliveryDate,
			req.DriverName, req.RecipientName, userID).Scan(&deliveryNoteID)

		if err != nil {
			return fmt.Errorf("gagal membuat surat jalan: %w", err)
		}

		// Update sales order status to 'delivery'
		if _, err = tx.Exec(`
            UPDATE sales_order 
            SET status = 'delivery', updated_at = NOW() 
            WHERE id = $1
        `, salesOrderID); err != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
		}

		// Set response fields
		response = CreateDeliveryNoteResponse{
			ID:                 deliveryNoteID,
			SerialID:           serialID,
			SalesOrderID:       salesOrderID,
			SalesOrderSerial:   salesOrderSerialID,
			SalesInvoiceID:     req.SalesInvoiceID,
			SalesInvoiceSerial: invoiceSerialID,
			DeliveryDate:       req.DeliveryDate,
			DriverName:         req.DriverName,
			RecipientName:      req.RecipientName,
			CreatedBy:          userID,
			CreatedAt:          time.Now().Format(time.RFC3339),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CancelDeliveryNote cancels an existing delivery note
func (r *SalesRepositoryImpl) CancelDeliveryNote(req CancelDeliveryNoteRequest, userID string) error {
	// Verify the delivery note exists and is not already cancelled
	var salesOrderID string
	var hasReturn bool
	var isCancelled bool

	err := r.db.QueryRow(`
        SELECT 
            dn.sales_order_id, 
            dn.cancelled_at IS NOT NULL,
            EXISTS (
                SELECT 1 FROM sales_order_return r
                WHERE r.delivery_note_id = dn.id 
                AND r.cancelled_at IS NULL
                AND r.return_status = 'completed'
            )
        FROM delivery_note dn
        WHERE dn.id = $1
    `, req.DeliveryNoteID).Scan(&salesOrderID, &isCancelled, &hasReturn)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("surat jalan tidak ditemukan")
		}
		return fmt.Errorf("gagal memeriksa surat jalan: %w", err)
	}

	if isCancelled {
		return errors.New("surat jalan sudah dibatalkan sebelumnya")
	}

	if hasReturn {
		return errors.New("tidak dapat membatalkan surat jalan yang memiliki pengembalian aktif")
	}

	// Execute transaction
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Mark delivery note as cancelled
		if _, err := tx.Exec(`
            UPDATE delivery_note
            SET cancelled_at = NOW(),
                cancelled_by = $1
            WHERE id = $2
        `, userID, req.DeliveryNoteID); err != nil {
			return fmt.Errorf("gagal membatalkan surat jalan: %w", err)
		}

		// Update sales order status back to 'invoice'
		if _, err := tx.Exec(`
            UPDATE sales_order
            SET status = 'invoice',
                updated_at = NOW()
            WHERE id = $1
        `, salesOrderID); err != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
		}

		return nil
	})
}
