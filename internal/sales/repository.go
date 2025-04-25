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
	ReturnItemFromSalesOrder(req ReturnItemRequest, userID string) (*ReturnInvoiceItemsResponse, error)
	CancelSalesOrderReturn(req CancelReturnRequest, userID string) error

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
        Select Bs.Id, Pb.Id As Batch_Id, Pb.Sku, Pb.Product_Id, P.Name As Product_Name, 
               Pb.Current_Quantity, Pb.Unit_Price, Pb.Created_At, 
               Bs.Storage_Id, S.Name As Storage_Name, S.Location As Storage_Location, Bs.Quantity
        From Product_Batch Pb
        Join Product P On Pb.Product_Id = P.Id
        Join Batch_Storage Bs On Pb.Id = Bs.Batch_Id
        Join Storage S On Bs.Storage_Id = S.Id
        Where Pb.Current_Quantity > 0 And Bs.Quantity > 0
    `)

	// Add search filter if provided
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		qb.Query.WriteString(" AND (pb.sku ILIKE $" + strconv.Itoa(qb.Count) + " OR p.name ILIKE $" + strconv.Itoa(qb.Count) + ")")
		qb.Params = append(qb.Params, searchTerm)
		qb.Count++
	}

	// Get count first (count distinct storage_ids to get number of storage groups)
	countQuery := "Select Count(Distinct Bs.Storage_Id) From Product_Batch Pb Join Batch_Storage Bs On Pb.Id = Bs.Batch_Id Join Product P On Pb.Product_Id = P.Id Where Pb.Current_Quantity > 0 And Bs.Quantity > 0"

	// Add search condition to count query if needed
	if req.Search != "" {
		countQuery += " And (pb.sku Ilike $1 Or p.name Ilike $1)"
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
        Select So.Id, So.Serial_Id, So.Customer_Id, C.Name As Customer_Name, 
               So.Order_Date, So.Status, So.Payment_Method, So.Payment_Due_Date, 
               So.Total_Amount, So.Created_At, So.Updated_At, So.Cancelled_At,
               (Select Si.Id From Sales_Invoice Si Where Si.Sales_Order_Id = So.Id And Si.Cancelled_At Is Null Limit 1) As Sales_Invoice_Id,
               (Select Dn.Id From Delivery_Note Dn Where Dn.Sales_Order_Id = So.Id And Dn.Cancelled_At Is Null Limit 1) As Delivery_Note_Id
        From Sales_Order So
        Join Customer C On So.Customer_Id = C.Id
        Where 1=1`

	// Create count query
	countQuery := `Select Count(*) From Sales_Order So Join Customer C On So.Customer_Id = C.Id Where 1=1`

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
            So.Id, So.Serial_Id, So.Customer_Id, C.Name, C.Telephone, C.Address,
            So.Order_Date, So.Status, So.Payment_Method, So.Payment_Due_Date, 
            So.Total_Amount, So.Created_By, Au.Username, So.Created_At, So.Updated_At, So.Cancelled_At,
            Si.Id, Si.Serial_Id AS Sales_Invoice_Serial_Id,
            Dn.Id, Dn.Serial_Id AS Delivery_Note_Serial_Id,
            So.Cancelled_At,
            Au.Username AS Created_By_Name, 
            Au2.Username AS Cancelled_By_Name
        FROM Sales_Order So
        LEFT JOIN Customer C ON So.Customer_Id = C.Id
        LEFT JOIN Appuser Au ON So.Created_By = Au.Id
        LEFT JOIN Appuser Au2 ON So.Cancelled_By = Au2.Id
        LEFT JOIN LATERAL (
            SELECT Id, Serial_Id 
            FROM Sales_Invoice 
            WHERE Sales_Order_Id = $1 AND Cancelled_At IS NULL
            ORDER BY Created_At DESC
            LIMIT 1
        ) Si ON TRUE
        LEFT JOIN LATERAL (
            SELECT Id, Serial_Id
            FROM Delivery_Note
            WHERE Sales_Order_Id = $1 AND Cancelled_At IS NULL
            ORDER BY Created_At DESC
            LIMIT 1
        ) Dn ON TRUE
        WHERE So.Id = $1
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
		&response.CreatedByName,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.CancelledBy,
		&response.SalesInvoiceID,
		&response.SalesInvoiceSerialID,
		&response.DeliveryNoteID,
		&response.DeliveryNoteSerialID,
		&response.CancelledAt,
		&response.CreatedByName,
		&response.CancelledByName,
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

	// If the order has a returned status, fetch return information for each item
	if response.Status == "returned" || response.Status == "partially_returned" {
		// Create a map to store return information by detail ID
		type returnInfo struct {
			ReturnID       string
			ReturnQuantity float64
			ReturnReason   string
			ReturnedAt     string
			ReturnedBy     string
		}

		returnMap := make(map[string]returnInfo)

		// Query to get return information
		returnsRows, errRet := r.db.Query(`
            SELECT 
                Sor.Id, Sor.Sales_Detail_Id, Sor.Return_Quantity, 
                Sor.Return_Reason, Sor.Returned_At, Au.Username AS Returned_By
            FROM Sales_Order_Return Sor
            LEFT JOIN Appuser Au ON Sor.Returned_By = Au.Id
            WHERE Sor.Sales_Order_Id = $1 
            AND Sor.Return_Status = 'completed'
            AND Sor.Cancelled_At IS NULL
        `, salesOrderID)

		if errRet != nil {
			return nil, fmt.Errorf("error fetching return information: %w", errRet)
		}
		defer returnsRows.Close()

		// Process return information
		for returnsRows.Next() {
			var ri returnInfo
			var detailID string
			if errScan := returnsRows.Scan(&ri.ReturnID, &detailID, &ri.ReturnQuantity, &ri.ReturnReason, &ri.ReturnedAt, &ri.ReturnedBy); errScan != nil {
				return nil, fmt.Errorf("error scanning return data: %w", errScan)
			}
			returnMap[detailID] = ri
		}

		if err = returnsRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating return rows: %w", err)
		}

		// Add return information to each item
		for i := range items {
			if ri, exists := returnMap[items[i].ID]; exists {
				items[i].ReturnID = ri.ReturnID
				items[i].ReturnQuantity = ri.ReturnQuantity
				items[i].ReturnReason = ri.ReturnReason
				items[i].ReturnedAt = ri.ReturnedAt
				items[i].ReturnedBy = ri.ReturnedBy
			}
		}
	}
	response.Items = items
	return &response, nil
}

// GetSalesOrderByID retrieves a single sales order by its ID
func (r *SalesRepositoryImpl) GetSalesOrderByID(id string) (*SalesOrder, error) {
	query := `
        Select So.Id, So.Serial_Id, So.Customer_Id, C.Name As Customer_Name, 
               So.Order_Date, So.Status, So.Payment_Method, So.Payment_Due_Date, 
               So.Total_Amount, So.Created_By, So.Created_At, So.Updated_At, So.Cancelled_At
        From Sales_Order So
        Join Customer C On So.Customer_Id = C.Id
        Where So.Id = $1`

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
        Select Sod.Id, Sod.Sales_Order_Id, 
               P.Id As Product_Id, P.Name As Product_Name, 
               U.Name As Product_Unit,
               Pb.Id As Batch_Id, Pb.Sku As Batch_Sku, 
               Sod.Batch_Storage_Id, 
               S.Id As Storage_Id, S.Name As Storage_Name,
               Sod.Quantity, Sod.Unit_Price,
               (Sod.Quantity * Sod.Unit_Price) As Total_Price,
               Bs.Quantity + Sod.Quantity As Max_Quantity
        From Sales_Order_Detail Sod
        Join Batch_Storage Bs On Sod.Batch_Storage_Id = Bs.Id
        Join Product_Batch Pb On Bs.Batch_Id = Pb.Id
        Join Product P On Pb.Product_Id = P.Id
        Join Unit U On P.Unit_Id = U.Id
        Join Storage S On Bs.Storage_Id = S.Id
        Where Sod.Sales_Order_Id = $1`

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
		Insert Into Sales_Order (Customer_Id, Serial_Id, Payment_Method, Payment_Due_Date, Created_By, Status, Total_Amount)
		Values ($1, $2, $3, $4, $5, 'order', $6)
		Returning Id, Serial_Id, Order_Date, Created_At, Status`

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
		errCustomer := tx.QueryRow("Select Name From Customer Where Id = $1", req.CustomerID).Scan(&customerName)
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
				Select Bs.Quantity, Bs.Batch_Id, Pb.Product_Id, P.Name, Pb.Sku, Pb.Unit_Price
				From Batch_Storage Bs
				Join Product_Batch Pb On Bs.Batch_Id = Pb.Id
				Join Product P On Pb.Product_Id = P.Id
				Where Bs.Id = $1
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
				Insert Into Sales_Order_Detail 
				(Sales_Order_Id, Batch_Storage_Id, Quantity, Unit_Price) 
				Values ($1, $2, $3, $4) 
				Returning Id`,
				orderID, item.BatchStorageID, item.Quantity, item.UnitPrice).Scan(&detailID)

			if errDetail != nil {
				return fmt.Errorf("gagal menambahkan detail pesanan: %w", errDetail)
			}

			// Update batch_storage quantity
			_, errBatchStorage := tx.Exec(`
					Update Batch_Storage 
					Set Quantity = Quantity - $1 
					Where Id = $2
				`, item.Quantity, item.BatchStorageID)

			if errBatchStorage != nil {
				return fmt.Errorf("gagal memperbarui kuantitas batch storage: %w", errBatchStorage)
			}

			// Update product_batch current_quantity
			_, errBatchQuantity := tx.Exec(`
					Update Product_Batch 
					Set Current_Quantity = Current_Quantity - $1 
					Where Id = $2
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
	errCheck := r.db.QueryRow("Select Status From Sales_Order Where Id = $1", req.ID).Scan(&status)
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
	query := "Update Sales_Order Set " + strings.Join(setValues, ", ") + " WHERE id = $" + strconv.Itoa(paramCount) + " RETURNING id, serial_id, customer_id, status, payment_method, payment_due_date"
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
		"Select Status From Sales_Order Where Id = $1",
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
			"Update Sales_Order Set Status = 'cancel', Cancelled_At = Now(), Cancelled_By = $1 Where Id = $2",
			userID, req.SalesOrderID,
		)
		if errCancel != nil {
			return fmt.Errorf("gagal membatalkan pesanan: %w", errCancel)
		}

		// Get all order details to restore inventory
		rows, errDetails := tx.Query(`
            Select Sod.Id, Sod.Batch_Id, Sod.Batch_Storage_Id, Sod.Quantity 
            From Sales_Order_Detail Sod
            Where Sod.Sales_Order_Id = $1
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
				"Update Product_Batch Set Current_Quantity = Current_Quantity + $1 Where Id = $2",
				detail.quantity, detail.batchID,
			)
			if errRestore != nil {
				return fmt.Errorf("gagal mengembalikan stok produk: %w", errRestore)
			}

			// Restore inventory in batch_storage directly using batch_storage_id
			_, errBatchStorage := tx.Exec(`
                Update Batch_Storage 
                Set Quantity = Quantity + $1 
                Where Id = $2
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
	errCheckItem := r.db.QueryRow("Select Id From Sales_Order_Detail Where Sales_Order_Id = $1 And Batch_Storage_Id = $2", req.SalesOrderID, req.BatchStorageID).Scan(&existingItemID)
	if errCheckItem != nil && !errors.Is(errCheckItem, sql.ErrNoRows) {
		return nil, fmt.Errorf("gagal memeriksa item pesanan: %w", errCheckItem)
	}
	if existingItemID != "" {
		return nil, fmt.Errorf("item ini sudah ada dalam pesanan")
	}

	// Check if sales order exists and if it can be modified
	errCheck := r.db.QueryRow("Select Status From Sales_Order Where Id = $1", req.SalesOrderID).Scan(&status)
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
        Select Bs.Id, Bs.Batch_Id, Pb.Product_Id, P.Name, Pb.Sku, Bs.Storage_Id, Bs.Quantity, Pb.Unit_Price
        From Batch_Storage Bs
        Join Product_Batch Pb On Bs.Batch_Id = Pb.Id
        Join Product P On Pb.Product_Id = P.Id
        Where Bs.Id = $1
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
            Insert Into Sales_Order_Detail 
            (Sales_Order_Id, Batch_Storage_Id, Quantity, Unit_Price) 
            Values ($1, $2, $3, $4) 
            Returning Id`,
			req.SalesOrderID, batchStorageID, req.Quantity, req.UnitPrice,
		).Scan(&detailID)
		if errDetail != nil {
			return fmt.Errorf("gagal menambahkan item ke pesanan: %w", errDetail)
		}

		// Update product_batch quantity
		_, errBatchUpdate := tx.Exec(
			"Update Product_Batch Set Current_Quantity = Current_Quantity - $1 Where Id = $2",
			req.Quantity, batchID,
		)
		if errBatchUpdate != nil {
			return fmt.Errorf("gagal memperbarui stok batch: %w", errBatchUpdate)
		}

		// Update batch_storage quantity
		_, errBatchStorageUpdate := tx.Exec(
			"Update Batch_Storage Set Quantity = Quantity - $1 Where Id = $2",
			req.Quantity, batchStorageID,
		)
		if errBatchStorageUpdate != nil {
			return fmt.Errorf("gagal memperbarui stok di lokasi penyimpanan: %w", errBatchStorageUpdate)
		}

		// Update the order's total_amount
		_, errTotal := tx.Exec(`
            Update Sales_Order 
            Set Total_Amount = (
                Select Coalesce(Sum(Quantity * Unit_Price), 0) 
                From Sales_Order_Detail 
                Where Sales_Order_Id = $1
            ),
            Updated_At = Now()
            Where Id = $1`,
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
	err := r.db.QueryRow("Select Status From Sales_Order Where Id = $1", req.SalesOrderID).Scan(&status)
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
            Select 
                Bs.Batch_Id, Sod.Batch_Storage_Id, Sod.Quantity, Sod.Unit_Price
            From Sales_Order_Detail Sod
            Join Batch_Storage Bs On Sod.Batch_Storage_Id = Bs.Id
            Where Sod.Id = $1 And Sod.Sales_Order_Id = $2
        `, req.DetailID, req.SalesOrderID).Scan(&batchID, &batchStorageID, &quantity, &unitPrice)

		if queryErr != nil {
			if errors.Is(queryErr, sql.ErrNoRows) {
				return fmt.Errorf("item tidak ditemukan dalam pesanan")
			}
			return fmt.Errorf("gagal mengambil detail item: %w", queryErr)
		}

		// Update batch_storage to restore quantity directly using batch_storage_id
		_, updateErr := tx.Exec(`
            Update Batch_Storage
            Set Quantity = Quantity + $1
            Where Id = $2
        `, quantity, batchStorageID)
		if updateErr != nil {
			return fmt.Errorf("gagal memulihkan kuantitas di penyimpanan: %w", updateErr)
		}

		// Update product_batch to restore total quantity
		_, updateErr = tx.Exec(`
            Update Product_Batch
            Set Current_Quantity = Current_Quantity + $1
            Where Id = $2
        `, quantity, batchID)
		if updateErr != nil {
			return fmt.Errorf("gagal memulihkan kuantitas batch: %w", updateErr)
		}

		// Delete the order detail
		_, deleteDetailErr := tx.Exec(`
            Delete From Sales_Order_Detail
            Where Id = $1
        `, req.DetailID)
		if deleteDetailErr != nil {
			return fmt.Errorf("gagal menghapus item pesanan: %w", deleteDetailErr)
		}

		// Update total order amount
		_, updateOrderErr := tx.Exec(`
            Update Sales_Order
            Set Total_Amount = Total_Amount - $1,
                Updated_At = Current_Timestamp
            Where Id = $2
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
	errCheck := r.db.QueryRow("Select Status From Sales_Order Where Id = $1", req.SalesOrderID).Scan(&status)
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
		Select Sod.Quantity, Sod.Unit_Price, Sod.Batch_Storage_Id 
		From Sales_Order_Detail Sod
		Where Sod.Id = $1 And Sod.Sales_Order_Id = $2
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
		Select Bs.Batch_Id, Pb.Product_Id, P.Name, Pb.Sku, Bs.Storage_Id
		From Batch_Storage Bs
		Join Product_Batch Pb On Bs.Batch_Id = Pb.Id
		Join Product P On Pb.Product_Id = P.Id
		Where Bs.Id = $1
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
			Select Bs.Batch_Id, Pb.Product_Id, Bs.Storage_Id
			From Batch_Storage Bs
			Join Product_Batch Pb On Bs.Batch_Id = Pb.Id
			Where Bs.Id = $1
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
				Update Sales_Order_Detail 
				Set Unit_Price = $1, Updated_At = Now()
				Where Id = $2
			`, newPrice, req.DetailID)
			if errUpdate != nil {
				return fmt.Errorf("gagal memperbarui harga item: %w", errUpdate)
			}

			// Update the order's total_amount
			_, errTotal := tx.Exec(`
				Update Sales_Order Set Total_Amount = (
					Select Coalesce(Sum(Quantity * Unit_Price), 0)
					From Sales_Order_Detail
					Where Sales_Order_Id = $1
				), Updated_At = Now()
				Where Id = $1
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
					Update Product_Batch 
					Set Current_Quantity = Current_Quantity + $1 
					Where Id = $2
				`, currentQty, batchID)
				if errRestore != nil {
					return fmt.Errorf("gagal mengembalikan stok batch lama: %w", errRestore)
				}

				_, errRestoreStorage := tx.Exec(`
					Update Batch_Storage 
					Set Quantity = Quantity + $1 
					Where Id = $2
				`, currentQty, batchStorageID)
				if errRestoreStorage != nil {
					return fmt.Errorf("gagal mengembalikan stok di lokasi penyimpanan lama: %w", errRestoreStorage)
				}

				// Validate new batch storage has enough quantity
				var availableQty float64
				errAvail := tx.QueryRow(`
					Select Quantity From Batch_Storage Where Id = $1
				`, req.BatchStorageID).Scan(&availableQty)
				if errAvail != nil {
					return fmt.Errorf("gagal memeriksa ketersediaan stok di lokasi baru: %w", errAvail)
				}

				if availableQty < newQty {
					return fmt.Errorf("stok tidak cukup di lokasi baru: tersedia %g, diminta %g", availableQty, newQty)
				}

				// Take quantity from new batch and storage
				_, errDeduct := tx.Exec(`
					Update Product_Batch 
					Set Current_Quantity = Current_Quantity - $1 
					Where Id = $2
				`, newQty, newBatchID)
				if errDeduct != nil {
					return fmt.Errorf("gagal mengurangi stok batch baru: %w", errDeduct)
				}

				_, errDeductStorage := tx.Exec(`
					Update Batch_Storage 
					Set Quantity = Quantity - $1 
					Where Id = $2
				`, newQty, req.BatchStorageID)
				if errDeductStorage != nil {
					return fmt.Errorf("gagal mengurangi stok di lokasi penyimpanan baru: %w", errDeductStorage)
				}

				// Update sales order detail with new batch storage
				_, errUpdateDetail := tx.Exec(`
					Update Sales_Order_Detail 
					Set Batch_Storage_Id = $1, Quantity = $2, Unit_Price = $3, Updated_At = Now()
					Where Id = $4
				`, req.BatchStorageID, newQty, newPrice, req.DetailID)
				if errUpdateDetail != nil {
					return fmt.Errorf("gagal memperbarui detail pesanan: %w", errUpdateDetail)
				}

				// Update response with new batch info
				batchID = newBatchID
				// Get updated batch SKU
				if err := tx.QueryRow("Select Sku From Product_Batch Where Id = $1", newBatchID).Scan(&batchSKU); err != nil {
					return fmt.Errorf("gagal mendapatkan informasi SKU batch baru: %w", err)
				}
			} else {
				// Just updating quantity of the same batch_storage
				if qtyDifference != 0 {
					// Check if we have enough quantity if increasing
					if qtyDifference > 0 {
						var availableQty float64
						errAvail := tx.QueryRow(`
							Select Quantity From Batch_Storage Where Id = $1
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
						Update Product_Batch 
						Set Current_Quantity = Current_Quantity - $1 
						Where Id = $2
					`, qtyDifference, batchID)
					if err != nil {
						return fmt.Errorf("gagal memperbarui stok batch: %w", err)
					}

					// Update batch_storage quantity
					_, err = tx.Exec(`
						Update Batch_Storage 
						Set Quantity = Quantity - $1 
						Where Id = $2
					`, qtyDifference, batchStorageID)
					if err != nil {
						return fmt.Errorf("gagal memperbarui stok di lokasi penyimpanan: %w", err)
					}
				}

				// Update sales order detail with new quantity/price
				_, err := tx.Exec(`
					Update Sales_Order_Detail 
					Set Quantity = $1, Unit_Price = $2, Updated_At = Now()
					Where Id = $3
				`, newQty, newPrice, req.DetailID)
				if err != nil {
					return fmt.Errorf("gagal memperbarui detail pesanan: %w", err)
				}
			}

			// Update the order's total_amount
			_, err := tx.Exec(`
				Update Sales_Order Set Total_Amount = (
					Select Coalesce(Sum(Quantity * Unit_Price), 0)
					From Sales_Order_Detail
					Where Sales_Order_Id = $1
				), Updated_At = Now()
				Where Id = $1
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
        Select Si.Id, Si.Serial_Id, Si.Sales_Order_Id, So.Serial_Id As Sales_Order_Serial,
               So.Customer_Id, C.Name As Customer_Name, Si.Invoice_Date, Si.Total_Amount,
               Case 
                 When Si.Cancelled_At Is Not Null Then 'cancelled'
                 When Exists(Select 1 From Sales_Order_Return Sor 
                            Join Sales_Order_Detail Sod On Sor.Sales_Detail_Id = Sod.Id
                            Where Sod.Sales_Order_Id = So.Id And Sor.Return_Status = 'returned') 
                    And Not Exists(Select 1 From Sales_Order_Detail Sod 
                                  Where Sod.Sales_Order_Id = So.Id 
                                  And Not Exists(Select 1 From Sales_Order_Return Sor 
                                               Where Sor.Sales_Detail_Id = Sod.Id And Sor.Return_Status = 'returned')) 
                    Then 'returned'
                 When Exists(Select 1 From Sales_Order_Return Sor 
                            Join Sales_Order_Detail Sod On Sor.Sales_Detail_Id = Sod.Id
                            Where Sod.Sales_Order_Id = So.Id And Sor.Return_Status = 'returned') Then 'partially_returned'
                 Else 'active'
               End As Status,
               Exists(Select 1 From Delivery_Note Dn Where Dn.Sales_Invoice_Id = Si.Id And Dn.Cancelled_At Is Null) As Has_Delivery_Note,
               Si.Created_By, Si.Created_At, Si.Cancelled_At
        From Sales_Invoice Si
        Join Sales_Order So On Si.Sales_Order_Id = So.Id
        Join Customer C On So.Customer_Id = C.Id
        Where 1=1`

	// Create count query
	countQuery := `
        Select Count(*) 
        From Sales_Invoice Si
        Join Sales_Order So On Si.Sales_Order_Id = So.Id
        Join Customer C On So.Customer_Id = C.Id
        Where 1=1`

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
            Select So.Status, So.Serial_Id, So.Customer_Id, C.Name, So.Total_Amount
            From Sales_Order So
            Join Customer C On So.Customer_Id = C.Id
            Where So.Id = $1
        `, req.SalesOrderID).Scan(&orderStatus, &orderSerial, &customerId, &customerName, &totalAmount)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("sales order tidak ditemukan")
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
            Select Id From Sales_Invoice 
            Where Sales_Order_Id = $1 And Cancelled_At Is Null
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
            Insert Into Sales_Invoice (
                Sales_Order_Id, Serial_Id, Total_Amount, Created_By
            ) Values ($1, $2, $3, $4)
            Returning Id, Serial_Id, Invoice_Date
        `, req.SalesOrderID, serialID, totalAmount, userID).Scan(&invoiceID, &serialID, &invoiceDate)

		if err != nil {
			return fmt.Errorf("gagal membuat faktur: %w", err)
		}

		// Update order status to 'invoice'
		_, err = tx.Exec(`
            Update Sales_Order 
            Set Status = 'invoice', Updated_At = Now() 
            Where Id = $1
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
            Select 
                Sod.Id, 
                Sod.Batch_Storage_Id, 
                Bs.Batch_Id,
                Bs.Storage_Id,
                Sod.Quantity, 
                Sod.Unit_Price
            From Sales_Order_Detail Sod
            Join Batch_Storage Bs On Sod.Batch_Storage_Id = Bs.Id
            Where Sod.Sales_Order_Id = $1
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
                Insert Into Inventory_Log (
                    Batch_Id, 
                    Storage_Id, 
                    User_Id, 
                    Sales_Order_Id, 
                    Action, 
                    Quantity, 
                    Description
                ) Values ($1, $2, $3, $4, $5, $6, $7)
            `,
				item.batchID,
				item.storageID,
				userID,
				req.SalesOrderID,
				"remove",
				item.quantity,
				fmt.Sprintf("Penjualan %s", serialID))

			if err != nil {
				return fmt.Errorf("gagal mencatat log inventaris: %w", err)
			}
		}

		// Create financial transaction log
		_, err = tx.Exec(`
            Insert Into Financial_Transaction_Log (
                User_Id,
                Amount,
                Type,
                Sales_Order_Id,
                Description,
                Transaction_Date,
                Is_System
            ) Values ($1, $2, $3, $4, $5, $6, $7)
        `,
			userID,
			totalAmount,
			"debit",
			req.SalesOrderID,
			fmt.Sprintf("Penjualan  %s", serialID),
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
		_, err := tx.Exec("Savepoint Create_Invoice")
		if err != nil {
			return nil, fmt.Errorf("failed to create savepoint: %w", err)
		}

		// Execute the invoice creation
		err = createInvoiceFunc(tx)

		if err != nil {
			// If there's an error, roll back to the savepoint
			_, rbErr := tx.Exec("Rollback To Savepoint Create_Invoice")
			if rbErr != nil {
				return nil, fmt.Errorf("error creating invoice: %v, and failed to rollback: %w", err, rbErr)
			}
			return nil, err
		}

		// Release the savepoint on success
		_, err = tx.Exec("Release Savepoint Create_Invoice")
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
        Select Sales_Order_Id, Cancelled_At Is Not Null 
        From Sales_Invoice 
        Where Id = $1
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
        Select Exists(
            Select 1 From Delivery_Note 
            Where Sales_Invoice_Id = $1 And Cancelled_At Is Null
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
            Update Sales_Invoice 
            Set Cancelled_At = Now(), Cancelled_By = $1 
            Where Id = $2
        `, userID, req.InvoiceID)

		if err != nil {
			return fmt.Errorf("gagal membatalkan faktur: %w", err)
		}

		// Revert sales order to 'order' status
		_, err = tx.Exec(`
            Update Sales_Order 
            Set Status = 'order', Updated_At = Now() 
            Where Id = $1
        `, salesOrderID)

		if err != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
		}

		return nil
	})
}

// ReturnItemFromSalesOrder handles returning items from a sales order (invoice or delivery)
func (r *SalesRepositoryImpl) ReturnItemFromSalesOrder(req ReturnItemRequest, userID string) (*ReturnInvoiceItemsResponse, error) {
	var response ReturnInvoiceItemsResponse

	// First, verify the sales order exists and get its status
	var salesOrderStatus string
	var isInvoiceCancelled bool
	var hasDeliveryNote bool
	var deliveryNoteID sql.NullString

	err := r.db.QueryRow(`
        Select 
            So.Status, 
            Exists(Select 1 From Sales_Invoice Si Where Si.Sales_Order_Id = So.Id And Si.Cancelled_At Is Null) As Has_Invoice,
            Exists(Select 1 From Delivery_Note Dn Where Dn.Sales_Order_Id = So.Id And Dn.Cancelled_At Is Null) As Has_Delivery,
            (Select Dn.Id From Delivery_Note Dn Where Dn.Sales_Order_Id = So.Id And Dn.Cancelled_At Is Null Limit 1) As Delivery_Note_Id
        From Sales_Order So
        Where So.Id = $1
    `, req.SalesOrderID).Scan(&salesOrderStatus, &isInvoiceCancelled, &hasDeliveryNote, &deliveryNoteID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sales order with ID %s not found", req.SalesOrderID)
		}
		return nil, fmt.Errorf("failed to check sales order: %w", err)
	}

	// Determine return source based on order status
	var returnSource string
	if salesOrderStatus == "delivery" {
		returnSource = "delivery"
		if !deliveryNoteID.Valid {
			return nil, fmt.Errorf("order is in delivery status but no active delivery note found")
		}
	} else if salesOrderStatus == "invoice" {
		returnSource = "invoice"
	} else {
		return nil, fmt.Errorf("cannot process return for order in '%s' status", salesOrderStatus)
	}

	// Verify the detail ID belongs to this order and get necessary info
	var detailExists bool
	var currentQuantity, previousReturnedQty float64
	var batchStorageID, batchID, productID string

	err = r.db.QueryRow(`
        Select 
            Exists(Select 1 From Sales_Order_Detail Where Id = $1 And Sales_Order_Id = $2),
            (Select Quantity From Sales_Order_Detail Where Id = $1),
            Coalesce((Select Sum(Return_Quantity) From Sales_Order_Return 
                Where Sales_Detail_Id = $1 And Return_Status = 'completed'), 0),
            (Select Batch_Storage_Id From Sales_Order_Detail Where Id = $1),
            (Select Bs.Batch_Id From Sales_Order_Detail Sod 
                Join Batch_Storage Bs On Sod.Batch_Storage_Id = Bs.Id
                Where Sod.Id = $1),
            (Select Pb.Product_Id From Sales_Order_Detail Sod 
                Join Batch_Storage Bs On Sod.Batch_Storage_Id = Bs.Id
                Join Product_Batch Pb On Bs.Batch_Id = Pb.Id
                Where Sod.Id = $1)
    `, req.SalesOrderDetailID, req.SalesOrderID).Scan(
		&detailExists, &currentQuantity, &previousReturnedQty,
		&batchStorageID, &batchID, &productID)

	if err != nil {
		return nil, fmt.Errorf("failed to verify item with ID %s: %w", req.SalesOrderDetailID, err)
	}

	if !detailExists {
		return nil, fmt.Errorf("item with ID %s not found in this order", req.SalesOrderDetailID)
	}

	// Check if return quantity is valid
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("return quantity must be greater than 0")
	}

	if req.Quantity > (currentQuantity - previousReturnedQty) {
		return nil, fmt.Errorf("return quantity %f exceeds available quantity %f",
			req.Quantity, (currentQuantity - previousReturnedQty))
	}

	// Execute transaction
	err = utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Generate return ID
		returnID := uuid.New().String()
		response.ReturnID = returnID

		// Calculate remaining quantity
		remainingQty := currentQuantity - previousReturnedQty - req.Quantity

		// Insert return record
		if _, err := tx.Exec(`
            Insert Into Sales_Order_Return (
                Id, Return_Source, Sales_Order_Id, Sales_Detail_Id,
                Return_Quantity, Remaining_Quantity, Return_Reason, 
                Return_Status, Returned_By, Returned_At, Delivery_Note_Id
            ) Values (
                $1, $2, $3, $4, $5, $6, $7, 'completed', $8, Now(), $9
            )
        `, returnID, returnSource, req.SalesOrderID, req.SalesOrderDetailID,
			req.Quantity, remainingQty, req.ReturnReason, userID,
			sql.NullString{String: deliveryNoteID.String, Valid: deliveryNoteID.Valid}); err != nil {
			return fmt.Errorf("failed to create return record: %w", err)
		}

		// Insert return batch record
		batchReturnID := uuid.New().String()
		if _, err := tx.Exec(`
            Insert Into Sales_Order_Return_Batch (
                Id, Sales_Return_Id, Batch_Id, Return_Quantity
            ) Values ($1, $2, $3, $4)
        `, batchReturnID, returnID, batchID, req.Quantity); err != nil {
			return fmt.Errorf("failed to record batch return: %w", err)
		}

		// Get product information for later use in inventory log
		var productName string
		if err := tx.QueryRow(`Select Name From Product Where Id = $1`, productID).Scan(&productName); err != nil {
			return fmt.Errorf("failed to get product name: %w", err)
		}

		// Determine where to return inventory (default to original batch_storage)
		var storageID string
		if err := tx.QueryRow(`Select Storage_Id From Batch_Storage Where Id = $1`, batchStorageID).Scan(&storageID); err != nil {
			return fmt.Errorf("failed to get storage ID: %w", err)
		}

		// Update batch storage quantity
		if _, err := tx.Exec(`
            Insert Into Batch_Storage (Batch_Id, Storage_Id, Quantity)
            Values ($1, $2, $3)
            On Conflict (Batch_Id, Storage_Id) 
            Do Update Set Quantity = Batch_Storage.Quantity + $3
        `, batchID, storageID, req.Quantity); err != nil {
			return fmt.Errorf("failed to update storage quantity: %w", err)
		}

		// Update batch current quantity
		if _, err := tx.Exec(`
            Update Product_Batch
            Set Current_Quantity = Current_Quantity + $1, Updated_At = Now()
            Where Id = $2
        `, req.Quantity, batchID); err != nil {
			return fmt.Errorf("failed to update batch quantity: %w", err)
		}

		// Get sales order serial ID
		var serialID string
		if err := tx.QueryRow(`SELECT Serial_Id FROM Sales_Order WHERE Id = $1`, req.SalesOrderID).Scan(&serialID); err != nil {
			return fmt.Errorf("failed to get sales order serial ID: %w", err)
		}

		// Insert with the modified description
		if _, err := tx.Exec(`
    Insert Into Inventory_Log (
        Batch_Id, Storage_Id, User_Id, Sales_Order_Id, Action,
        Quantity, Description
    ) Values (
        $1, $2, $3, $4, 'add',
        $5, $6
    )
`, batchID, storageID, userID, req.SalesOrderID,
			req.Quantity, fmt.Sprintf("Retur Barang Penjualan %s", serialID)); err != nil {
			return fmt.Errorf("failed to create inventory log: %w", err)
		}

		// 2) financial log for the return (refund)
		var unitPrice float64
		if err := tx.QueryRow(`
        SELECT Sod.Unit_Price
        FROM Sales_Order_Detail Sod
        WHERE Sod.Id = $1 AND Sod.Sales_Order_Id = $2
    `, req.SalesOrderDetailID, req.SalesOrderID).Scan(&unitPrice); err != nil {
			return fmt.Errorf("gagal mengambil harga satuan untuk retur: %w", err)
		}

		if _, err := tx.Exec(`
        INSERT INTO Financial_Transaction_Log (
            User_Id, Amount, Type, Sales_Order_Id,
            Description, Transaction_Date, Is_System
        ) VALUES ($1,$2,$3,$4,$5,$6,$7)
    `, userID,
			req.Quantity*unitPrice,
			"credit",
			req.SalesOrderID,
			fmt.Sprintf("Retur Barang Penjualan %s", serialID),
			time.Now(),
			true,
		); err != nil {
			return fmt.Errorf("gagal mencatat transaksi keuangan retur: %w", err)
		}

		// Check if this is a full return (all items returned)
		var isFullReturn bool
		if err := tx.QueryRow(`
            With Total_Items As (
                Select Count(*) As Count
                From Sales_Order_Detail
                Where Sales_Order_Id = $1
            ),
            Fully_Returned_Items As (
                Select Count(*) As Count
                From Sales_Order_Detail Sod
                Where Sales_Order_Id = $1
                And (
                    Select Coalesce(Sum(Return_Quantity), 0) 
                    From Sales_Order_Return 
                    Where Sales_Detail_Id = Sod.Id 
                    And Return_Status = 'completed'
                    And Cancelled_At Is Null
                ) >= Sod.Quantity
            )
            Select Ti.Count = Fri.Count
            From Total_Items Ti, Fully_Returned_Items Fri
        `, req.SalesOrderID).Scan(&isFullReturn); err != nil {
			return fmt.Errorf("failed to check full return status: %w", err)
		}

		response.IsFullReturn = isFullReturn

		// Update sales order status based on return status
		var newStatus string
		if isFullReturn {
			newStatus = "returned"
		} else {
			newStatus = "partially_returned"
		}

		if _, err := tx.Exec(`
            Update Sales_Order
            Set Status = $1, Updated_At = Now()
            Where Id = $2
        `, newStatus, req.SalesOrderID); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Set response fields
	response.SalesOrderID = req.SalesOrderID
	response.ReturnedItems = 2
	response.TotalQuantity = 2
	response.ReturnDate = time.Now().Format(time.RFC3339)
	response.ReturnStatus = "completed"

	return &response, nil
}

// CancelSalesOrderReturn cancels a previously processed return
func (r *SalesRepositoryImpl) CancelSalesOrderReturn(req CancelReturnRequest, userID string) error {
	// Verify the return exists and is not already cancelled
	var salesOrderID string
	var isCancelled bool
	var returnSource string
	var salesDetailID string
	var returnQuantity float64

	err := r.db.QueryRow(`
        Select R.Sales_Order_Id, R.Cancelled_At Is Not Null, R.Return_Source, R.Sales_Detail_Id, R.Return_Quantity
        From Sales_Order_Return R
        Where R.Id = $1
    `, req.ReturnID).Scan(&salesOrderID, &isCancelled, &returnSource, &salesDetailID, &returnQuantity)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("return not found")
		}
		return fmt.Errorf("failed to check return: %w", err)
	}

	// Check if already cancelled
	if isCancelled {
		return errors.New("this return has already been cancelled")
	}

	// Execute transaction
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// First, collect all return batch records
		type batchReturn struct {
			id       string
			batchId  string
			quantity float64
		}

		var batchReturns []batchReturn

		// Get all batch returns
		rows, err := tx.Query(`
            Select 
                Rb.Id,
                Rb.Batch_Id,
                Rb.Return_Quantity
            From Sales_Order_Return_Batch Rb
            Where Rb.Sales_Return_Id = $1
        `, req.ReturnID)

		if err != nil {
			return fmt.Errorf("failed to get return batch records: %w", err)
		}

		// Important: collect all data before closing rows
		for rows.Next() {
			var br batchReturn
			if err := rows.Scan(&br.id, &br.batchId, &br.quantity); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan batch return data: %w", err)
			}
			batchReturns = append(batchReturns, br)
		}

		// Close rows before further DB operations
		rows.Close()

		if err = rows.Err(); err != nil {
			return fmt.Errorf("error while reading batch returns: %w", err)
		}

		// Process each batch return after closing the first result set
		for _, br := range batchReturns {
			// Update batch quantity - remove the returned quantity
			if _, err := tx.Exec(`
                Update Product_Batch
                Set Current_Quantity = Current_Quantity - $1,
                    Updated_At = Now()
                Where Id = $2
            `, br.quantity, br.batchId); err != nil {
				return fmt.Errorf("failed to update batch quantity: %w", err)
			}

			// Next, collect all storage records for this batch in a separate slice
			type storageRecord struct {
				id        string
				storageId string
				quantity  float64
			}

			var storageRecords []storageRecord

			// Query storage records
			storageRows, err := tx.Query(`
                Select Bs.Id, Bs.Storage_Id, Bs.Quantity
                From Batch_Storage Bs
                Where Bs.Batch_Id = $1
                Order By Bs.Quantity Desc
            `, br.batchId)

			if err != nil {
				return fmt.Errorf("failed to get storage records: %w", err)
			}

			// Collect all storage records before closing
			for storageRows.Next() {
				var sr storageRecord
				if err := storageRows.Scan(&sr.id, &sr.storageId, &sr.quantity); err != nil {
					storageRows.Close()
					return fmt.Errorf("failed to scan storage data: %w", err)
				}
				storageRecords = append(storageRecords, sr)
			}

			// Close the storage rows cursor
			storageRows.Close()

			if err = storageRows.Err(); err != nil {
				return fmt.Errorf("error reading storage records: %w", err)
			}

			// Now process the storage records outside the cursor loop
			var remainingQty = br.quantity
			for _, sr := range storageRecords {
				if remainingQty <= 0 {
					break
				}

				// Determine how much to take from this storage
				var qtyToRemove float64
				if sr.quantity >= remainingQty {
					qtyToRemove = remainingQty
					remainingQty = 0
				} else {
					qtyToRemove = sr.quantity
					remainingQty -= sr.quantity
				}

				// Update storage quantity
				if _, err := tx.Exec(`
                    Update Batch_Storage
                    Set Quantity = Quantity - $1
                    Where Id = $2
                `, qtyToRemove, sr.id); err != nil {
					return fmt.Errorf("failed to update storage quantity: %w", err)
				}

				// Get sales order serial ID
				var serialID string
				if err := tx.QueryRow(`SELECT Serial_Id FROM Sales_Order WHERE Id = $1`, salesOrderID).Scan(&serialID); err != nil {
					return fmt.Errorf("failed to get sales order serial ID: %w", err)
				}

				// Log inventory change
				if _, err := tx.Exec(`
                    Insert Into Inventory_Log (
                        Batch_Id, Storage_Id, User_Id, Sales_Order_Id, Action, 
                        Quantity, Description
                    ) Values (
                        $1, $2, $3, $4, 'remove', $5, 
                        $6
                    )
                `, br.batchId, sr.storageId, userID, salesOrderID, qtyToRemove,
					fmt.Sprintf("Batal Retur Barang Penjualan %s", serialID)); err != nil {
					return fmt.Errorf("failed to log inventory change: %w", err)
				}

				var unitPrice float64
				if err := tx.QueryRow(`
        SELECT Sod.Unit_Price
        FROM Sales_Order_Detail Sod
        WHERE Sod.Id = $1
    `, salesDetailID).Scan(&unitPrice); err != nil {
					return fmt.Errorf("gagal mengambil harga satuan untuk batal retur: %w", err)
				}

				if _, err := tx.Exec(`
        INSERT INTO Financial_Transaction_Log (
            User_Id, Amount, Type, Sales_Order_Id,
            Description, Transaction_Date, Is_System
        ) VALUES ($1,$2,$3,$4,$5,$6,$7)
    `, userID,
					returnQuantity*unitPrice,
					"debit",
					salesOrderID,
					fmt.Sprintf("Batal Retur %s", req.ReturnID),
					time.Now(),
					true,
				); err != nil {
					return fmt.Errorf("gagal mencatat transaksi keuangan batal retur: %w", err)
				}
			}

			// Ensure all quantity was processed
			if remainingQty > 0 {
				return fmt.Errorf("insufficient storage quantities found for batch %s", br.batchId)
			}
		}

		// Mark return as cancelled
		if _, err := tx.Exec(`
            Update Sales_Order_Return
            Set Return_Status = 'cancelled',
                Cancelled_At = Now(),
                Cancelled_By = $1
            Where Id = $2
        `, userID, req.ReturnID); err != nil {
			return fmt.Errorf("failed to mark return as cancelled: %w", err)
		}

		// Continue with the rest of the function to update order status...

		// The remaining code is unchanged
		var hasRemainingReturns bool
		if err := tx.QueryRow(`
            Select Exists (
                Select 1 From Sales_Order_Return
                Where Sales_Order_Id = $1
                And Return_Status = 'completed'
                And Cancelled_At Is Null
            )
        `, salesOrderID).Scan(&hasRemainingReturns); err != nil {
			return fmt.Errorf("failed to check remaining returns: %w", err)
		}

		var newStatus string
		if hasRemainingReturns {
			// Check if it's a full return or partial
			var isFullReturn bool
			if err := tx.QueryRow(`
                With Total_Items As (
                    Select Count(*) As Count
                    From Sales_Order_Detail
                    Where Sales_Order_Id = $1
                ),
                Fully_Returned_Items As (
                    Select Count(*) As Count
                    From Sales_Order_Detail Sod
                    Where Sales_Order_Id = $1
                    And (
                        Select Coalesce(Sum(Return_Quantity), 0) 
                        From Sales_Order_Return 
                        Where Sales_Detail_Id = Sod.Id 
                        And Return_Status = 'completed'
                        And Cancelled_At Is Null
                    ) >= Sod.Quantity
                )
                Select Ti.Count = Fri.Count
                From Total_Items Ti, Fully_Returned_Items Fri
            `, salesOrderID).Scan(&isFullReturn); err != nil {
				return fmt.Errorf("failed to check full return status: %w", err)
			}

			if isFullReturn {
				newStatus = "returned"
			} else {
				newStatus = "partially_returned"
			}
		} else {
			// Determine appropriate status based on delivery note/invoice existence
			var hasDelivery bool
			var hasInvoice bool

			if err := tx.QueryRow(`
                Select 
                    Exists(Select 1 From Delivery_Note Where Sales_Order_Id = $1 And Cancelled_At Is Null),
                    Exists(Select 1 From Sales_Invoice Where Sales_Order_Id = $1 And Cancelled_At Is Null)
            `, salesOrderID).Scan(&hasDelivery, &hasInvoice); err != nil {
				return fmt.Errorf("failed to check order documents: %w", err)
			}

			if hasDelivery {
				newStatus = "delivery"
			} else if hasInvoice {
				newStatus = "invoice"
			} else {
				newStatus = "order"
			}
		}

		// Update order status
		if _, err := tx.Exec(`
            Update Sales_Order
            Set Status = $1,
                Updated_At = Now()
            Where Id = $2
        `, newStatus, salesOrderID); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
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
        Select 
            S.Id, 
            I.Serial_Id, 
            S.Serial_Id,
            C.Name As Customer_Name,
            I.Cancelled_At Is Not Null,
            Exists (
                Select 1 From Delivery_Note Dn 
                Where Dn.Sales_Invoice_Id = I.Id And Dn.Cancelled_At Is Null
            )
        From Sales_Invoice I
        Join Sales_Order S On I.Sales_Order_Id = S.Id
        Join Customer C On S.Customer_Id = C.Id
        Where I.Id = $1
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
            Insert Into Delivery_Note (
                Serial_Id, Sales_Order_Id, Sales_Invoice_Id, 
                Delivery_Date, Driver_Name, Recipient_Name, Created_By
            ) Values ($1, $2, $3, $4, $5, $6, $7)
            Returning Id
        `, serialID, salesOrderID, req.SalesInvoiceID, req.DeliveryDate,
			req.DriverName, req.RecipientName, userID).Scan(&deliveryNoteID)

		if err != nil {
			return fmt.Errorf("gagal membuat surat jalan: %w", err)
		}

		// Update sales order status to 'delivery'
		if _, err = tx.Exec(`
            Update Sales_Order 
            Set Status = 'delivery', Updated_At = Now() 
            Where Id = $1
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
        Select 
            Dn.Sales_Order_Id, 
            Dn.Cancelled_At Is Not Null,
            Exists (
                Select 1 From Sales_Order_Return R
                Where R.Delivery_Note_Id = Dn.Id 
                And R.Cancelled_At Is Null
                And R.Return_Status = 'completed'
            )
        From Delivery_Note Dn
        Where Dn.Id = $1
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
            Update Delivery_Note
            Set Cancelled_At = Now(),
                Cancelled_By = $1
            Where Id = $2
        `, userID, req.DeliveryNoteID); err != nil {
			return fmt.Errorf("gagal membatalkan surat jalan: %w", err)
		}

		// Update sales order status back to 'invoice'
		if _, err := tx.Exec(`
            Update Sales_Order
            Set Status = 'invoice',
                Updated_At = Now()
            Where Id = $1
        `, salesOrderID); err != nil {
			return fmt.Errorf("gagal memperbarui status pesanan: %w", err)
		}

		return nil
	})
}
