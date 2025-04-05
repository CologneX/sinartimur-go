package sales

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sinartimur-go/utils"
	"strconv"
	"strings"
	"time"
)

type SalesRepository interface {
	// Sales Order operations
	GetSalesOrders(req GetSalesOrdersRequest) ([]GetSalesOrdersResponse, int, error)
	GetSalesOrderByID(id string) (*SalesOrder, error)
	GetSalesOrderDetails(salesOrderID string) ([]GetSalesOrderDetail, error)
	CreateSalesOrder(req CreateSalesOrderRequest, userID string) (*CreateSalesOrderResponse, error)
	UpdateSalesOrder(req UpdateSalesOrderRequest) (*UpdateSalesOrderResponse, error)
	CancelSalesOrder(req CancelSalesOrderRequest, userID string) error

	// Sales Order Item operations
	AddItemToSalesOrder(req AddItemToSalesOrderRequest) (*UpdateItemResponse, error)
	UpdateSalesOrderItem(req UpdateItemRequest) (*UpdateItemResponse, error)
	DeleteSalesOrderItem(req DeleteItemRequest) error

	// Invoice operations
	GetSalesInvoices(req GetSalesInvoicesRequest) ([]GetSalesInvoicesResponse, int, error)
	CreateSalesInvoice(req CreateSalesInvoiceRequest, userID string) (*CreateSalesInvoiceResponse, error)
	CancelSalesInvoice(req CancelSalesInvoiceRequest, userID string) error

	// Return operations
	ReturnInvoiceItems(req ReturnInvoiceItemsRequest, userID string) (*ReturnInvoiceItemsResponse, error)
	CancelInvoiceReturn(req CancelInvoiceReturnRequest, userID string) error

	// Delivery Note operations
	CreateDeliveryNote(req CreateDeliveryNoteRequest, userID string) (*CreateDeliveryNoteResponse, error)
	CancelDeliveryNote(req CancelDeliveryNoteRequest, userID string) error
}

type SalesRepositoryImpl struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) SalesRepository {
	return &SalesRepositoryImpl{db: db}
}

// GetSalesOrders retrieves a paginated list of sales orders with filtering options
func (r *SalesRepositoryImpl) GetSalesOrders(req GetSalesOrdersRequest) ([]GetSalesOrdersResponse, int, error) {
	// Build base query for fetching sales orders
	baseQuery := `
        SELECT so.id, so.serial_id, so.customer_id, c.name AS customer_name, 
               so.order_date, so.status, so.payment_method, so.payment_due_date, 
               so.total_amount, so.created_at, so.updated_at, so.cancelled_at
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
		var orderDate time.Time
		var updatedAt time.Time
		var paymentDueDate, cancelledAt sql.NullTime

		errScan := rows.Scan(
			&order.ID,
			&order.SerialID,
			&order.CustomerID,
			&order.CustomerName,
			&orderDate,
			&order.Status,
			&order.PaymentMethod,
			&paymentDueDate,
			&order.TotalAmount,
			&order.CreatedAt,
			&updatedAt,
			&cancelledAt,
		)
		if errScan != nil {
			return nil, 0, fmt.Errorf("error scanning sales order row: %w", errScan)
		}

		// Format dates for response
		order.OrderDate = orderDate.Format(time.RFC3339)
		order.UpdatedAt = updatedAt.Format(time.RFC3339)

		if paymentDueDate.Valid {
			order.PaymentDueDate = paymentDueDate.Time.Format(time.RFC3339)
		}

		if cancelledAt.Valid {
			order.CancelledAt = cancelledAt.Time.Format(time.RFC3339)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating sales order rows: %w", err)
	}

	return orders, totalItems, nil
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
	var paymentDueDate, cancelledAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.SerialID,
		&order.CustomerID,
		&order.CustomerName,
		&order.OrderDate,
		&order.Status,
		&order.PaymentMethod,
		&paymentDueDate,
		&order.TotalAmount,
		&order.CreatedBy,
		&order.CreatedAt,
		&order.UpdatedAt,
		&cancelledAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sales order not found: %s", id)
		}
		return nil, fmt.Errorf("error fetching sales order: %w", err)
	}

	return &order, nil
}

// GetSalesOrderDetails retrieves the details (items) for a specific sales order
func (r *SalesRepositoryImpl) GetSalesOrderDetails(salesOrderID string) ([]GetSalesOrderDetail, error) {
	query := `
        SELECT sod.id, sod.sales_order_id, sod.product_id, p.name AS product_name, 
               sod.batch_id, pb.sku AS batch_sku, sod.quantity, sod.unit_price,
               (sod.quantity * sod.unit_price) AS total_price
        FROM sales_order_detail sod
        JOIN product p ON sod.product_id = p.id
        JOIN product_batch pb ON sod.batch_id = pb.id
        WHERE sod.sales_order_id = $1`

	rows, err := r.db.Query(query, salesOrderID)
	if err != nil {
		return nil, fmt.Errorf("error fetching sales order details: %w", err)
	}
	defer rows.Close()

	var details []GetSalesOrderDetail
	for rows.Next() {
		var detail GetSalesOrderDetail
		errScan := rows.Scan(
			&detail.ID,
			&detail.SalesOrderID,
			&detail.ProductID,
			&detail.ProductName,
			&detail.BatchID,
			&detail.BatchSKU,
			&detail.Quantity,
			&detail.UnitPrice,
			&detail.TotalPrice,
		)
		if errScan != nil {
			return nil, fmt.Errorf("error scanning sales order detail row: %w", errScan)
		}

		// Fetch storage allocations for this detail
		storageQuery := `
            SELECT sos.id, sos.storage_id, s.name AS storage_name, sos.quantity
            FROM sales_order_storage sos
            JOIN storage s ON sos.storage_id = s.id
            WHERE sos.sales_order_detail_id = $1`

		storageRows, errStorage := r.db.Query(storageQuery, detail.ID)
		if errStorage != nil {
			return nil, fmt.Errorf("error fetching storage allocations: %w", errStorage)
		}

		for storageRows.Next() {
			var storage SalesOrderStorageResponse
			errScan = storageRows.Scan(
				&storage.ID,
				&storage.StorageID,
				&storage.StorageName,
				&storage.Quantity,
			)
			if errScan != nil {
				storageRows.Close()
				return nil, fmt.Errorf("error scanning storage allocation row: %w", errScan)
			}
			detail.StorageAllocations = append(detail.StorageAllocations, storage)
		}
		storageRows.Close()

		details = append(details, detail)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sales order detail rows: %w", err)
	}

	return details, nil
}

// CreateSalesOrder creates a new sales order with its details
func (r *SalesRepositoryImpl) CreateSalesOrder(req CreateSalesOrderRequest, userID string) (*CreateSalesOrderResponse, error) {
	var response CreateSalesOrderResponse

	err := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Insert sales order
		var orderID, serialID string
		var orderDate time.Time
		var paymentDueDate sql.NullTime

		// Calculate total amount for the order
		var totalAmount float64
		for _, item := range req.Items {
			totalAmount += item.Quantity * item.UnitPrice
		}

		// Convert payment due date if provided
		if req.PaymentDueDate != "" {
			parsedDate, errParse := time.Parse(time.RFC3339, req.PaymentDueDate)
			if errParse != nil {
				return fmt.Errorf("format tanggal pembayaran tidak valid: %w", errParse)
			}
			paymentDueDate.Valid = true
			paymentDueDate.Time = parsedDate
		}

		// Insert sales order
		orderQuery := `
			INSERT INTO sales_order (customer_id, payment_method, payment_due_date, created_by, status, total_amount)
			VALUES ($1, $2, $3, $4, 'order', $5)
			RETURNING id, serial_id, order_date, created_at`

		errOrder := tx.QueryRow(
			orderQuery,
			req.CustomerID,
			req.PaymentMethod,
			paymentDueDate,
			userID,
			totalAmount,
		).Scan(&orderID, &serialID, &orderDate, &totalAmount, &response.CreatedAt)

		if errOrder != nil {
			return fmt.Errorf("gagal membuat pesanan: %w", errOrder)
		}

		// Get customer name
		var customerName string
		errCustomer := tx.QueryRow("SELECT name FROM customer WHERE id = $1", req.CustomerID).Scan(&customerName)
		if errCustomer != nil {
			return fmt.Errorf("pelanggan tidak ditemukan: %w", errCustomer)
		}

		// Process each item
		for _, item := range req.Items {
			// Verify product exists
			var productName string
			errProduct := tx.QueryRow("SELECT name FROM product WHERE id = $1", item.ProductID).Scan(&productName)
			if errProduct != nil {
				return fmt.Errorf("produk tidak ditemukan: %w", errProduct)
			}

			// Verify batch exists and has enough quantity
			var currentQuantity float64
			errBatch := tx.QueryRow(
				"SELECT current_quantity FROM product_batch WHERE id = $1 AND product_id = $2",
				item.BatchID, item.ProductID,
			).Scan(&currentQuantity)

			if errBatch != nil {
				return fmt.Errorf("batch produk tidak ditemukan: %w", errBatch)
			}

			if currentQuantity < item.Quantity {
				return fmt.Errorf("stok produk %s tidak mencukupi (tersedia: %.2f, diminta: %.2f)",
					productName, currentQuantity, item.Quantity)
			}

			// Insert order detail
			var detailID string
			detailQuery := `
				INSERT INTO sales_order_detail (sales_order_id, product_id, batch_id, quantity, unit_price)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id`

			errDetail := tx.QueryRow(
				detailQuery,
				orderID,
				item.ProductID,
				item.BatchID,
				item.Quantity,
				item.UnitPrice,
			).Scan(&detailID)

			if errDetail != nil {
				return fmt.Errorf("gagal menambahkan item pesanan: %w", errDetail)
			}

			// Process storage allocations
			var totalAllocated float64
			for _, allocation := range item.StorageAllocations {
				// Verify storage exists
				var storageExists bool
				errStorage := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM storage WHERE id = $1)", allocation.StorageID).Scan(&storageExists)
				if errStorage != nil {
					return fmt.Errorf("gagal memeriksa penyimpanan: %w", errStorage)
				}

				if !storageExists {
					return fmt.Errorf("penyimpanan tidak ditemukan")
				}

				// Verify storage has enough stock of this batch
				var storageQuantity float64
				errBatchInStorage := tx.QueryRow(
					"SELECT quantity FROM batch_storage WHERE batch_id = $1 AND storage_id = $2",
					item.BatchID, allocation.StorageID,
				).Scan(&storageQuantity)

				if errBatchInStorage != nil {
					return fmt.Errorf("batch tidak tersedia di penyimpanan ini: %w", errBatchInStorage)
				}

				if storageQuantity < allocation.Quantity {
					return fmt.Errorf("stok di penyimpanan tidak mencukupi (tersedia: %.2f, diminta: %.2f)",
						storageQuantity, allocation.Quantity)
				}

				// Insert allocation
				_, errInsertAlloc := tx.Exec(
					"INSERT INTO sales_order_storage (sales_order_detail_id, storage_id, batch_id, quantity) VALUES ($1, $2, $3, $4)",
					detailID,
					allocation.StorageID,
					item.BatchID,
					allocation.Quantity,
				)

				if errInsertAlloc != nil {
					return fmt.Errorf("gagal menyimpan alokasi penyimpanan: %w", errInsertAlloc)
				}

				// Update batch_storage quantity
				_, errUpdateStorage := tx.Exec(
					"UPDATE batch_storage SET quantity = quantity - $1 WHERE batch_id = $2 AND storage_id = $3",
					allocation.Quantity,
					item.BatchID,
					allocation.StorageID,
				)

				if errUpdateStorage != nil {
					return fmt.Errorf("gagal memperbarui stok di penyimpanan: %w", errUpdateStorage)
				}

				totalAllocated += allocation.Quantity
			}

			// Verify total allocation matches requested quantity
			if totalAllocated != item.Quantity {
				return fmt.Errorf("total alokasi (%.2f) tidak sama dengan kuantitas yang diminta (%.2f)",
					totalAllocated, item.Quantity)
			}

			// Update product_batch current_quantity
			_, errUpdateBatch := tx.Exec(
				"UPDATE product_batch SET current_quantity = current_quantity - $1 WHERE id = $2",
				item.Quantity,
				item.BatchID,
			)

			if errUpdateBatch != nil {
				return fmt.Errorf("gagal memperbarui stok batch: %w", errUpdateBatch)
			}
		}

		// Create invoice if requested
		if req.CreateInvoice {
			var invoiceID, invoiceSerialID string
			errInvoice := tx.QueryRow(
				"INSERT INTO sales_invoice (sales_order_id, created_by, total_amount) VALUES ($1, $2, $3) RETURNING id, serial_id",
				orderID,
				userID,
				totalAmount,
			).Scan(&invoiceID, &invoiceSerialID)

			if errInvoice != nil {
				return fmt.Errorf("gagal membuat faktur: %w", errInvoice)
			}

			response.InvoiceID = invoiceID
			response.InvoiceSerialID = invoiceSerialID
		}

		// Set response data
		response.ID = orderID
		response.SerialID = serialID
		response.CustomerID = req.CustomerID
		response.CustomerName = customerName
		response.Status = "order" // Default status for new orders
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

	query := "UPDATE sales_order SET "
	params := []interface{}{}
	paramCount := 1
	setValues := []string{}

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
		paymentDueDate, errParse := time.Parse(time.RFC3339, req.PaymentDueDate)
		if errParse != nil {
			return nil, fmt.Errorf("format tanggal pembayaran tidak valid: %w", errParse)
		}
		setValues = append(setValues, fmt.Sprintf("payment_due_date = $%d", paramCount))
		params = append(params, paymentDueDate)
		paramCount++
	}

	// Add updated_at timestamp
	setValues = append(setValues, "updated_at = NOW()")

	// If no fields to update
	if len(setValues) <= 1 {
		return nil, fmt.Errorf("tidak ada data yang diperbarui")
	}

	// Construct final query
	query += strings.Join(setValues, ", ") + " WHERE id = $" + strconv.Itoa(paramCount) + " RETURNING updated_at"
	params = append(params, req.ID)

	var updatedAt time.Time
	errUpdate := r.db.QueryRow(query, params...).Scan(&updatedAt)
	if errUpdate != nil {
		return nil, fmt.Errorf("gagal memperbarui pesanan: %w", errUpdate)
	}

	// Get updated order data for response
	queryOrder := `
		SELECT so.id, so.serial_id, so.customer_id, c.name AS customer_name, 
		       so.status, so.payment_method, so.payment_due_date
		FROM sales_order so
		JOIN customer c ON so.customer_id = c.id
		WHERE so.id = $1`

	var paymentDueDate sql.NullTime
	errFetch := r.db.QueryRow(queryOrder, req.ID).Scan(
		&response.ID,
		&response.SerialID,
		&response.CustomerID,
		&response.Status,
		&response.PaymentMethod,
		&paymentDueDate,
	)

	if errFetch != nil {
		return nil, fmt.Errorf("gagal mengambil data pesanan yang diperbarui: %w", errFetch)
	}

	if paymentDueDate.Valid {
		response.PaymentDueDate = paymentDueDate.Time.Format(time.RFC3339)
	}
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

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
            SELECT sod.id, sod.batch_id, sod.quantity 
            FROM sales_order_detail sod
            WHERE sod.sales_order_id = $1
        `, req.SalesOrderID)

		if errDetails != nil {
			return fmt.Errorf("gagal mengambil detail pesanan: %w", errDetails)
		}
		defer rows.Close()

		for rows.Next() {
			var detailID, batchID string
			var quantity float64

			errScan := rows.Scan(&detailID, &batchID, &quantity)
			if errScan != nil {
				return fmt.Errorf("gagal membaca detail pesanan: %w", errScan)
			}

			// Restore inventory in product_batch
			_, errRestore := tx.Exec(
				"UPDATE product_batch SET current_quantity = current_quantity + $1 WHERE id = $2",
				quantity, batchID,
			)
			if errRestore != nil {
				return fmt.Errorf("gagal mengembalikan stok produk: %w", errRestore)
			}

			// Get storage allocations for this detail
			storageRows, errStorage := tx.Query(`
                SELECT storage_id, quantity 
                FROM sales_order_storage 
                WHERE sales_order_detail_id = $1
            `, detailID)

			if errStorage != nil {
				return fmt.Errorf("gagal mengambil alokasi penyimpanan: %w", errStorage)
			}

			for storageRows.Next() {
				var storageID string
				var storageQty float64

				errStorageScan := storageRows.Scan(&storageID, &storageQty)
				if errStorageScan != nil {
					storageRows.Close()
					return fmt.Errorf("gagal membaca alokasi penyimpanan: %w", errStorageScan)
				}

				// Restore inventory in batch_storage
				_, errBatchStorage := tx.Exec(`
                    UPDATE batch_storage 
                    SET quantity = quantity + $1 
                    WHERE batch_id = $2 AND storage_id = $3
                `, storageQty, batchID, storageID)

				if errBatchStorage != nil {
					storageRows.Close()
					return fmt.Errorf("gagal mengembalikan stok di lokasi penyimpanan: %w", errBatchStorage)
				}
			}
			storageRows.Close()
		}

		if err := rows.Err(); err != nil {
			return fmt.Errorf("terjadi kesalahan saat memproses detail pesanan: %w", err)
		}

		return nil
	})
}

// AddItemToSalesOrder adds a new item to an existing sales order
func (r *SalesRepositoryImpl) AddItemToSalesOrder(req AddItemToSalesOrderRequest) (*UpdateItemResponse, error) {
	var response UpdateItemResponse
	var status string

	// Check if sales order exists and if it can be modified
	errCheck := r.db.QueryRow("SELECT status FROM sales_order WHERE id = $1", req.SalesOrderID).Scan(&status)
	if errCheck != nil {
		if errors.Is(errCheck, sql.ErrNoRows) {
			return nil, fmt.Errorf("pesanan penjualan dengan ID tersebut tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa status pesanan: %w", errCheck)
	}

	// Only allow items to be added to orders with status 'order'
	if status != "order" {
		return nil, fmt.Errorf("hanya pesanan dengan status 'order' yang dapat diubah")
	}

	// Verify product exists
	var productName string
	errProduct := r.db.QueryRow("SELECT name FROM product WHERE id = $1 AND deleted_at IS NULL", req.ProductID).Scan(&productName)
	if errProduct != nil {
		if errors.Is(errProduct, sql.ErrNoRows) {
			return nil, fmt.Errorf("produk dengan ID tersebut tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa produk: %w", errProduct)
	}

	// Verify batch exists and belongs to the product
	var batchSKU string
	errBatch := r.db.QueryRow(
		"SELECT sku FROM product_batch WHERE id = $1 AND product_id = $2",
		req.BatchID, req.ProductID,
	).Scan(&batchSKU)
	if errBatch != nil {
		if errors.Is(errBatch, sql.ErrNoRows) {
			return nil, fmt.Errorf("batch tidak ditemukan atau tidak terkait dengan produk ini")
		}
		return nil, fmt.Errorf("gagal memeriksa batch produk: %w", errBatch)
	}

	// Check batch has enough quantity
	var currentQuantity float64
	errQuantity := r.db.QueryRow(
		"SELECT current_quantity FROM product_batch WHERE id = $1",
		req.BatchID,
	).Scan(&currentQuantity)
	if errQuantity != nil {
		return nil, fmt.Errorf("gagal memeriksa ketersediaan stok: %w", errQuantity)
	}

	if currentQuantity < req.Quantity {
		return nil, fmt.Errorf("jumlah yang diminta (%g) melebihi stok yang tersedia (%g)", req.Quantity, currentQuantity)
	}

	// Validate storage allocations total matches the requested quantity
	var totalAllocated float64
	for _, alloc := range req.StorageAllocations {
		totalAllocated += alloc.Quantity
	}

	if totalAllocated != req.Quantity {
		return nil, fmt.Errorf("total alokasi penyimpanan (%g) tidak sama dengan jumlah yang diminta (%g)", totalAllocated, req.Quantity)
	}

	// Execute transaction
	err := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Create a new sales order detail entry
		var detailID string
		errDetail := tx.QueryRow(`
            INSERT INTO sales_order_detail 
            (sales_order_id, product_id, batch_id, quantity, unit_price) 
            VALUES ($1, $2, $3, $4, $5) 
            RETURNING id`,
			req.SalesOrderID, req.ProductID, req.BatchID, req.Quantity, req.UnitPrice,
		).Scan(&detailID)
		if errDetail != nil {
			return fmt.Errorf("gagal menambahkan item ke pesanan: %w", errDetail)
		}

		// Update product_batch quantity
		_, errBatchUpdate := tx.Exec(
			"UPDATE product_batch SET current_quantity = current_quantity - $1 WHERE id = $2",
			req.Quantity, req.BatchID,
		)
		if errBatchUpdate != nil {
			return fmt.Errorf("gagal memperbarui stok batch: %w", errBatchUpdate)
		}

		// Process storage allocations
		for _, alloc := range req.StorageAllocations {
			// Verify storage exists
			var storageName string
			errStorageCheck := tx.QueryRow(
				"SELECT name FROM storage WHERE id = $1 AND deleted_at IS NULL",
				alloc.StorageID,
			).Scan(&storageName)
			if errStorageCheck != nil {
				if errors.Is(errStorageCheck, sql.ErrNoRows) {
					return fmt.Errorf("lokasi penyimpanan dengan ID %s tidak ditemukan", alloc.StorageID)
				}
				return fmt.Errorf("gagal memeriksa lokasi penyimpanan: %w", errStorageCheck)
			}

			// Verify batch storage allocation has enough quantity
			var storageQuantity float64
			errStorageQty := tx.QueryRow(
				"SELECT quantity FROM batch_storage WHERE batch_id = $1 AND storage_id = $2",
				req.BatchID, alloc.StorageID,
			).Scan(&storageQuantity)
			if errStorageQty != nil {
				if errors.Is(errStorageQty, sql.ErrNoRows) {
					return fmt.Errorf("batch tidak tersedia di lokasi penyimpanan yang dipilih")
				}
				return fmt.Errorf("gagal memeriksa ketersediaan di lokasi penyimpanan: %w", errStorageQty)
			}

			if storageQuantity < alloc.Quantity {
				return fmt.Errorf("jumlah alokasi di %s (%g) melebihi stok yang tersedia (%g)", storageName, alloc.Quantity, storageQuantity)
			}

			// Create sales_order_storage entry
			_, errStorage := tx.Exec(`
                INSERT INTO sales_order_storage 
                (sales_order_detail_id, storage_id, batch_id, quantity) 
                VALUES ($1, $2, $3, $4)`,
				detailID, alloc.StorageID, req.BatchID, alloc.Quantity,
			)
			if errStorage != nil {
				return fmt.Errorf("gagal mencatat alokasi penyimpanan: %w", errStorage)
			}

			// Update batch_storage quantity
			_, errBatchStorage := tx.Exec(
				"UPDATE batch_storage SET quantity = quantity - $1 WHERE batch_id = $2 AND storage_id = $3",
				alloc.Quantity, req.BatchID, alloc.StorageID,
			)
			if errBatchStorage != nil {
				return fmt.Errorf("gagal memperbarui stok di lokasi penyimpanan: %w", errBatchStorage)
			}
		}

		// Update the order's total_amount
		_, errTotal := tx.Exec(`
            UPDATE sales_order 
            SET total_amount = (
                SELECT COALESCE(SUM(quantity * unit_price), 0) 
                FROM sales_order_detail 
                WHERE sales_order_id = $1 AND cancelled_at IS NULL
            ),
            updated_at = NOW()
            WHERE id = $1`,
			req.SalesOrderID,
		)
		if errTotal != nil {
			return fmt.Errorf("gagal memperbarui total pesanan: %w", errTotal)
		}

		// Set response values
		response.DetailID = detailID
		response.ProductID = req.ProductID
		response.ProductName = productName
		response.BatchID = req.BatchID
		response.BatchSKU = batchSKU
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
func (r *SalesRepositoryImpl) DeleteSalesOrderItem(req DeleteItemRequest) error {
	var status string
	var productName, batchSKU string
	var quantity float64

	// Check if sales order exists and if it can be modified
	errCheck := r.db.QueryRow("SELECT status FROM sales_order WHERE id = $1", req.SalesOrderID).Scan(&status)
	if errCheck != nil {
		if errors.Is(errCheck, sql.ErrNoRows) {
			return fmt.Errorf("pesanan dengan ID %s tidak ditemukan", req.SalesOrderID)
		}
		return fmt.Errorf("gagal memeriksa pesanan: %w", errCheck)
	}

	// Only allow items to be deleted from orders with status 'order'
	if status != "order" {
		return fmt.Errorf("tidak dapat menghapus item: pesanan sudah dalam status %s", status)
	}

	// Check if detail exists and get its data for inventory restoration
	errDetail := r.db.QueryRow(`
        SELECT sod.quantity, pb.sku, p.name 
        FROM sales_order_detail sod
        JOIN product_batch pb ON sod.batch_id = pb.id
        JOIN product p ON pb.product_id = p.id
        WHERE sod.id = $1 AND sod.sales_order_id = $2`,
		req.DetailID, req.SalesOrderID).Scan(&quantity, &batchSKU, &productName)

	if errDetail != nil {
		if errors.Is(errDetail, sql.ErrNoRows) {
			return fmt.Errorf("item dengan ID %s tidak ditemukan dalam pesanan ini", req.DetailID)
		}
		return fmt.Errorf("gagal mendapatkan informasi item: %w", errDetail)
	}

	// Execute transaction
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Get all storage allocations for this detail
		rows, err := tx.Query(`
            SELECT storage_id, batch_id, quantity 
            FROM sales_order_storage 
            WHERE sales_order_detail_id = $1`, req.DetailID)
		if err != nil {
			return fmt.Errorf("gagal mendapatkan alokasi penyimpanan: %w", err)
		}
		defer rows.Close()

		// Restore quantities to batch and storage
		for rows.Next() {
			var storageID, batchID string
			var allocQty float64

			if errScan := rows.Scan(&storageID, &batchID, &allocQty); errScan != nil {
				return fmt.Errorf("gagal membaca data alokasi: %w", errScan)
			}

			// Restore batch quantity
			_, errBatch := tx.Exec(
				"UPDATE product_batch SET current_quantity = current_quantity + $1, updated_at = NOW() WHERE id = $2",
				allocQty, batchID)
			if errBatch != nil {
				return fmt.Errorf("gagal memulihkan kuantitas batch: %w", errBatch)
			}

			// Restore storage quantity
			_, errStorage := tx.Exec(
				"UPDATE batch_storage SET quantity = quantity + $1, updated_at = NOW() WHERE batch_id = $2 AND storage_id = $3",
				allocQty, batchID, storageID)
			if errStorage != nil {
				return fmt.Errorf("gagal memulihkan kuantitas di penyimpanan: %w", errStorage)
			}
		}

		if errScan := rows.Err(); errScan != nil {
			return fmt.Errorf("error pada iterasi alokasi penyimpanan: %w", errScan)
		}

		// Delete storage allocations (will be deleted by CASCADE, but let's be explicit)
		_, errDeleteStorage := tx.Exec("DELETE FROM sales_order_storage WHERE sales_order_detail_id = $1",
			req.DetailID)
		if errDeleteStorage != nil {
			return fmt.Errorf("gagal menghapus alokasi penyimpanan: %w", errDeleteStorage)
		}

		// Delete the sales order detail
		_, errDeleteDetail := tx.Exec("DELETE FROM sales_order_detail WHERE id = $1", req.DetailID)
		if errDeleteDetail != nil {
			return fmt.Errorf("gagal menghapus item pesanan: %w", errDeleteDetail)
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
			req.SalesOrderID)
		if errTotal != nil {
			return fmt.Errorf("gagal memperbarui total pesanan: %w", errTotal)
		}

		return nil
	})
}

// UpdateSalesOrderItem updates an item in a sales order with new quantity or price
func (r *SalesRepositoryImpl) UpdateSalesOrderItem(req UpdateItemRequest) (*UpdateItemResponse, error) {
	var response UpdateItemResponse
	var status string
	var currentQty, currentPrice float64
	var productID, productName, batchID, batchSKU string

	// Check if sales order exists and if it's in a modifiable state
	errCheck := r.db.QueryRow("SELECT status FROM sales_order WHERE id = $1", req.SalesOrderID).Scan(&status)
	if errCheck != nil {
		if errors.Is(errCheck, sql.ErrNoRows) {
			return nil, fmt.Errorf("pesanan penjualan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa status pesanan: %w", errCheck)
	}

	// Only allow items to be updated if order's status is 'order'
	if status != "order" {
		return nil, fmt.Errorf("hanya pesanan dengan status 'order' yang dapat diubah")
	}

	// Get current detail information
	errDetail := r.db.QueryRow(`
        SELECT d.quantity, d.unit_price, d.product_id, p.name, d.batch_id, pb.sku 
        FROM sales_order_detail d
        JOIN product p ON d.product_id = p.id
        JOIN product_batch pb ON d.batch_id = pb.id
        WHERE d.id = $1 AND d.sales_order_id = $2
    `, req.DetailID, req.SalesOrderID).Scan(&currentQty, &currentPrice, &productID, &productName, &batchID, &batchSKU)

	if errDetail != nil {
		if errors.Is(errDetail, sql.ErrNoRows) {
			return nil, fmt.Errorf("item pesanan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mendapatkan informasi item: %w", errDetail)
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

	// Validate batch availability if quantity is being increased
	if newQty > currentQty {
		var batchCurrentQty float64
		errBatch := r.db.QueryRow("SELECT current_quantity FROM product_batch WHERE id = $1", batchID).Scan(&batchCurrentQty)
		if errBatch != nil {
			return nil, fmt.Errorf("gagal memeriksa ketersediaan batch: %w", errBatch)
		}

		additionalQty := newQty - currentQty
		if batchCurrentQty < additionalQty {
			return nil, fmt.Errorf("stok batch tidak mencukupi (tersedia: %.2f, dibutuhkan tambahan: %.2f)", batchCurrentQty, additionalQty)
		}
	}

	// Execute transaction
	err := utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		qtyDifference := newQty - currentQty

		// Handle storage allocation changes if quantity changed and storage allocations provided
		if qtyDifference != 0 && len(req.StorageAllocations) > 0 {
			// Delete existing storage allocations if we're completely replacing them
			_, err := tx.Exec("DELETE FROM sales_order_storage WHERE sales_order_detail_id = $1", req.DetailID)
			if err != nil {
				return fmt.Errorf("gagal menghapus alokasi penyimpanan lama: %w", err)
			}

			// Restore quantities to batch and storages
			_, errRestoreBatch := tx.Exec("UPDATE product_batch SET current_quantity = current_quantity + $1 WHERE id = $2",
				currentQty, batchID)
			if errRestoreBatch != nil {
				return fmt.Errorf("gagal mengembalikan kuantitas batch: %w", errRestoreBatch)
			}

			// Get all existing storage allocations
			rows, errGet := tx.Query(`
                SELECT storage_id, quantity FROM sales_order_storage
                WHERE sales_order_detail_id = $1
            `, req.DetailID)
			if errGet == nil {
				defer rows.Close()
				for rows.Next() {
					var storageID string
					var storageQty float64
					if errScan := rows.Scan(&storageID, &storageQty); errScan != nil {
						return fmt.Errorf("gagal membaca alokasi penyimpanan: %w", errScan)
					}

					// Restore storage quantities
					_, errRestoreStorage := tx.Exec(`
                        UPDATE batch_storage SET quantity = quantity + $1
                        WHERE batch_id = $2 AND storage_id = $3
                    `, storageQty, batchID, storageID)
					if errRestoreStorage != nil {
						return fmt.Errorf("gagal mengembalikan kuantitas ke penyimpanan: %w", errRestoreStorage)
					}
				}
			}

			// Process new storage allocations
			var totalAllocated float64
			for _, alloc := range req.StorageAllocations {
				totalAllocated += alloc.Quantity

				// Verify storage has sufficient quantity
				var availableQty float64
				errCheck = tx.QueryRow(`
                    SELECT quantity FROM batch_storage
                    WHERE batch_id = $1 AND storage_id = $2
                `, batchID, alloc.StorageID).Scan(&availableQty)
				if errCheck != nil {
					if errors.Is(sql.ErrNoRows, errCheck) {
						return fmt.Errorf("batch tidak tersedia di lokasi penyimpanan yang dipilih")
					}
					return fmt.Errorf("gagal memeriksa ketersediaan: %w", errCheck)
				}

				if availableQty < alloc.Quantity {
					return fmt.Errorf("stok tidak mencukupi di lokasi penyimpanan (tersedia: %.2f, dibutuhkan: %.2f)",
						availableQty, alloc.Quantity)
				}

				// Create storage allocation
				_, errStorage := tx.Exec(`
                    INSERT INTO sales_order_storage 
                    (sales_order_detail_id, storage_id, batch_id, quantity) 
                    VALUES ($1, $2, $3, $4)
                `, req.DetailID, alloc.StorageID, batchID, alloc.Quantity)
				if errStorage != nil {
					return fmt.Errorf("gagal mencatat alokasi penyimpanan: %w", errStorage)
				}

				// Reduce storage quantity
				_, errUpdate := tx.Exec(`
                    UPDATE batch_storage SET quantity = quantity - $1
                    WHERE batch_id = $2 AND storage_id = $3
                `, alloc.Quantity, batchID, alloc.StorageID)
				if errUpdate != nil {
					return fmt.Errorf("gagal memperbarui kuantitas penyimpanan: %w", errUpdate)
				}
			}

			// Verify the total allocated quantity matches the new quantity
			if totalAllocated != newQty {
				return fmt.Errorf("total alokasi penyimpanan (%.2f) tidak sesuai dengan kuantitas yang diminta (%.2f)",
					totalAllocated, newQty)
			}

			// Update batch quantity
			_, errBatch := tx.Exec("UPDATE product_batch SET current_quantity = current_quantity - $1 WHERE id = $2",
				newQty, batchID)
			if errBatch != nil {
				return fmt.Errorf("gagal memperbarui kuantitas batch: %w", errBatch)
			}
		} else if qtyDifference != 0 {
			// If quantity changed but no new storage allocations provided, just adjust the batch quantity
			_, errBatch := tx.Exec("UPDATE product_batch SET current_quantity = current_quantity - $1 WHERE id = $2",
				qtyDifference, batchID)
			if errBatch != nil {
				return fmt.Errorf("gagal memperbarui kuantitas batch: %w", errBatch)
			}
		}

		// Update the sales order detail
		_, errUpdate := tx.Exec(`
            UPDATE sales_order_detail 
            SET quantity = $1, unit_price = $2, updated_at = NOW()
            WHERE id = $3
        `, newQty, newPrice, req.DetailID)
		if errUpdate != nil {
			return fmt.Errorf("gagal memperbarui item pesanan: %w", errUpdate)
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
			return fmt.Errorf("gagal memperbarui total pesanan: %w", errTotal)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Set response values
	response.DetailID = req.DetailID
	response.ProductID = productID
	response.ProductName = productName
	response.BatchID = batchID
	response.BatchSKU = batchSKU
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

// CreateSalesInvoice creates a new invoice from a sales order
func (r *SalesRepositoryImpl) CreateSalesInvoice(req CreateSalesInvoiceRequest, userID string) (*CreateSalesInvoiceResponse, error) {
	var response CreateSalesInvoiceResponse

	// Check if sales order exists and is in 'order' status
	var orderStatus, orderSerial string
	var customerId string
	var customerName string
	var totalAmount float64

	err := r.db.QueryRow(`
        SELECT so.status, so.serial_id, so.customer_id, c.name, so.total_amount
        FROM sales_order so
        JOIN customer c ON so.customer_id = c.id
        WHERE so.id = $1
    `, req.SalesOrderID).Scan(&orderStatus, &orderSerial, &customerId, &customerName, &totalAmount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("pesanan penjualan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal memeriksa pesanan: %w", err)
	}

	// Validate order status
	if orderStatus != "order" {
		return nil, fmt.Errorf("tidak dapat membuat faktur: pesanan dalam status %s", orderStatus)
	}

	// Check if invoice already exists for this order
	var existingInvoice string
	err = r.db.QueryRow(`
        SELECT id FROM sales_invoice 
        WHERE sales_order_id = $1 AND cancelled_at IS NULL
    `, req.SalesOrderID).Scan(&existingInvoice)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("gagal memeriksa faktur yang ada: %w", err)
	}

	if existingInvoice != "" {
		return nil, fmt.Errorf("faktur sudah ada untuk pesanan ini")
	}

	// Execute transaction
	err = utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		var invoiceID, serialID string
		var invoiceDate time.Time

		// Create invoice
		err := tx.QueryRow(`
            INSERT INTO sales_invoice (
                sales_order_id, total_amount, created_by
            ) VALUES ($1, $2, $3)
            RETURNING id, serial_id, invoice_date
        `, req.SalesOrderID, totalAmount, userID).Scan(&invoiceID, &serialID, &invoiceDate)

		if err != nil {
			return fmt.Errorf("gagal membuat faktur: %w", err)
		}

		// Update sales order status to 'invoice'
		_, err = tx.Exec(`
            UPDATE sales_order 
            SET status = 'invoice', updated_at = NOW() 
            WHERE id = $1
        `, req.SalesOrderID)

		if err != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
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
            i.sales_order_id, 
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

	// Process delivery date
	deliveryDate := time.Now()
	if req.DeliveryDate != "" {
		parsedDate, err := time.Parse(time.RFC3339, req.DeliveryDate)
		if err != nil {
			return nil, fmt.Errorf("format tanggal pengiriman tidak valid: %w", err)
		}
		deliveryDate = parsedDate
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
        `, serialID, salesOrderID, req.SalesInvoiceID, deliveryDate,
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
			DeliveryDate:       deliveryDate.Format(time.RFC3339),
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
