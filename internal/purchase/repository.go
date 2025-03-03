package purchase

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
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

// NewSupplierRepository creates a new instance of SupplierRepositoryImpl
func NewSupplierRepository(db *sql.DB) SupplierRepository {
	return &SupplierRepositoryImpl{db: db}
}

// GetAll fetches all suppliers with filtering and pagination
func (r *SupplierRepositoryImpl) GetAll(req GetSupplierRequest) ([]GetSupplierResponse, int, error) {
	var suppliers []GetSupplierResponse
	var totalItems int

	// Base queries
	query := `Select Id, Name, Address, Telephone, Created_At, Updated_At 
              From Supplier 
              Where Deleted_At Is Null`

	countQuery := `Select Count(Id) 
                  From Supplier 
                  Where Deleted_At Is Null`

	// Apply filters
	var filterParams []interface{}
	paramCount := 1
	if req.Name != "" {
		query += ` AND name ILIKE $` + strconv.Itoa(paramCount)
		countQuery += ` AND name ILIKE $` + strconv.Itoa(paramCount)
		filterParams = append(filterParams, "%"+req.Name+"%")
		paramCount++
	}

	if req.Telephone != "" {
		query += ` AND telephone ILIKE $` + strconv.Itoa(paramCount)
		countQuery += ` AND telephone ILIKE $` + strconv.Itoa(paramCount)
		filterParams = append(filterParams, "%"+req.Telephone+"%")
		paramCount++
	}

	// Add sorting
	//query += ` ORDER BY created_at DESC`

	// Execute count query
	err := r.db.QueryRow(countQuery, filterParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination to main query
	query += ` LIMIT $` + strconv.Itoa(paramCount) + ` OFFSET $` + strconv.Itoa(paramCount+1)
	queryParams := append(filterParams, req.PageSize, (req.Page-1)*req.PageSize)

	// Execute main query with all parameters including pagination
	rows, err := r.db.Query(query, queryParams...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Process results
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
			return nil, 0, errScan
		}

		supplier.CreatedAt = createdAt.Format(time.RFC3339)
		supplier.UpdatedAt = updatedAt.Format(time.RFC3339)
		suppliers = append(suppliers, supplier)
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
		if err == sql.ErrNoRows {
			return nil, errors.New("Supplier tidak ditemukan")
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
		if err == sql.ErrNoRows {
			return nil, errors.New("Supplier tidak ditemukan")
		}
		return nil, err
	}

	supplier.CreatedAt = createdAt.Format(time.RFC3339)
	supplier.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &supplier, nil
}

// Create inserts a new supplier
func (r *SupplierRepositoryImpl) Create(req CreateSupplierRequest) (*GetSupplierResponse, error) {
	var supplier GetSupplierResponse

	err := r.db.QueryRow(`
		Insert Into Supplier (Name, Address, Telephone) 
		Values ($1, $2, $3) 
		Returning Id, Name, Address, Telephone, Created_At, Updated_At
	`, req.Name, req.Address, req.Telephone).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.Address,
		&supplier.Telephone,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &supplier, nil
}

// Update updates an existing supplier
func (r *SupplierRepositoryImpl) Update(req UpdateSupplierRequest) (*GetSupplierResponse, error) {
	// First check if supplier exists
	_, err := r.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	query := `Update Supplier Set Updated_At = Now()`
	var params []interface{}
	paramCount := 1

	if req.Name != "" {
		query += `, name = $` + string(rune(paramCount))
		params = append(params, req.Name)
		paramCount++
	}

	if req.Address != "" {
		query += `, address = $` + string(rune(paramCount))
		params = append(params, req.Address)
		paramCount++
	}

	if req.Telephone != "" {
		query += `, telephone = $` + string(rune(paramCount))
		params = append(params, req.Telephone)
		paramCount++
	}

	// Add WHERE clause and RETURNING
	query += ` WHERE id = $` + string(rune(paramCount)) + ` AND deleted_at IS NULL 
		RETURNING id, name, address, telephone, created_at, updated_at`
	params = append(params, req.ID)

	// Execute update
	var supplier GetSupplierResponse
	var createdAt, updatedAt time.Time

	err = r.db.QueryRow(query, params...).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.Address,
		&supplier.Telephone,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	supplier.CreatedAt = createdAt.Format(time.RFC3339)
	supplier.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &supplier, nil
}

// Delete soft-deletes a supplier
func (r *SupplierRepositoryImpl) Delete(id string) error {
	// First check if supplier exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Soft delete the supplier
	_, err = r.db.Exec(`
		Update Supplier 
		Set Deleted_At = Now() 
		Where Id = $1 And Deleted_At Is Null
	`, id)

	return err
}

// PurchaseOrderRepository interface defines methods for purchase order operations
type PurchaseOrderRepository interface {
	GetAll(req GetPurchaseOrderRequest) ([]GetPurchaseOrderResponse, int, error)
	GetByID(id string) (*GetPurchaseOrderResponse, error)
	GetDetailByID(id string) (*PurchaseOrderDetailResponse, error)
	Create(req CreatePurchaseOrderRequest, userID string) (*GetPurchaseOrderResponse, error)
	Update(req UpdatePurchaseOrderRequest) (*GetPurchaseOrderResponse, error)
	Cancel(id string) error
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
	var orders []GetPurchaseOrderResponse
	var totalItems int

	var params []interface{}
	paramCount := 1

	// Initialize builders for both the main query and count query
	var query strings.Builder
	var countQuery strings.Builder

	// Base queries
	query.WriteString(`
        Select Po.Id, U.Username, Po.Supplier_Id, S.Name, Po.Order_Date, Po.Status,
               Po.Total_Amount, Po.Created_At, Po.Updated_At
        From Purchase_Order Po
        Left Join Appuser U On Po.Created_By = U.Id
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Where Po.Cancelled_At Is Null`)

	countQuery.WriteString(`
        Select Count(Po.Id)
        From Purchase_Order Po
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Where Po.Cancelled_At Is Null`)

	// Apply filters
	if req.SupplierName != "" {
		filter := " AND S.Name ILIKE $" + strconv.Itoa(paramCount)
		query.WriteString(filter)
		countQuery.WriteString(filter)
		params = append(params, "%"+req.SupplierName+"%")
		paramCount++
	}

	// Rest of the filtering logic remains the same...

	// Execute count query
	err := r.db.QueryRow(countQuery.String(), params...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination to main query
	query.WriteString(" LIMIT $" + strconv.Itoa(paramCount))
	params = append(params, req.PageSize)
	paramCount++

	query.WriteString(" OFFSET $" + strconv.Itoa(paramCount))
	params = append(params, (req.Page-1)*req.PageSize)

	// Execute main query
	rows, err := r.db.Query(query.String(), params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var order GetPurchaseOrderResponse

		errScan := rows.Scan(
			&order.ID,
			&order.CreatedBy,
			&order.SupplierID,
			&order.SupplierName,
			&order.OrderDate,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if errScan != nil {
			return nil, 0, errScan
		}

		orders = append(orders, order)
	}

	return orders, totalItems, nil
}

// GetByID fetches a purchase order by ID
func (r *PurchaseOrderRepositoryImpl) GetByID(id string) (*GetPurchaseOrderResponse, error) {
	var order GetPurchaseOrderResponse

	err := r.db.QueryRow(`
        Select Po.Id, U.Username, Po.Supplier_Id, S.Name, Po.Order_Date, Po.Status,
               Po.Total_Amount, Po.Created_At, Po.Updated_At
        From Purchase_Order Po
        Left Join Appuser U On Po.Created_By = U.Id
        Left Join Supplier S On Po.Supplier_Id = S.Id
        Where Po.Id = $1 And Po.Cancelled_At Is Null
    `, id).Scan(
		&order.ID,
		&order.CreatedBy,
		&order.SupplierID,
		&order.SupplierName,
		&order.OrderDate,
		&order.Status,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Purchase Order tidak ditemukan!")
		}
		return nil, err
	}

	return &order, nil
}

// GetDetailByID fetches a purchase order with its items by ID
func (r *PurchaseOrderRepositoryImpl) GetDetailByID(id string) (*PurchaseOrderDetailResponse, error) {
	// First get the order
	order, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Then get the items
	var query strings.Builder
	query.WriteString(`
        Select Pod.Id, Pod.Product_Id, Pod.Quantity, Pod.Price,
               (Pod.Quantity * Pod.Price) As Subtotal,
               Coalesce(P.Name, '') As Description
        From Purchase_Order_Detail Pod
        Left Join Product P On Pod.Product_Id = P.Id
        Where Pod.Purchase_Order_Id = $1
        Order By Pod.Created_At Desc
    `)

	rows, err := r.db.Query(query.String(), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []GetPurchaseOrderItemResponse
	for rows.Next() {
		var item GetPurchaseOrderItemResponse
		errScan := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
			&item.ProductName,
		)
		if errScan != nil {
			return nil, errScan
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Create the combined response with both order details and items
	return &PurchaseOrderDetailResponse{
		GetPurchaseOrderResponse: *order,
		Items:                    items,
	}, nil
}

// Create inserts a new purchase order with its items
func (r *PurchaseOrderRepositoryImpl) Create(req CreatePurchaseOrderRequest, userID string) (*GetPurchaseOrderResponse, error) {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Calculate total amount from items
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += float64(item.Quantity) * item.Price
	}

	// Get supplier name for record
	var supplierName string
	err = tx.QueryRow(`Select Name From Supplier Where Id = $1`, req.SupplierID).Scan(&supplierName)
	if err != nil {
		return nil, err
	}

	var orderID string
	var orderDate, createdAt, updatedAt time.Time

	// Insert the purchase order
	err = tx.QueryRow(`
        Insert Into Purchase_Order (Supplier_Id, Order_Date, Status, Total_Amount, Created_By)
        Values ($1, $2, $3, $4, $5)
        Returning Id, Order_Date, Created_At, Updated_At
    `, req.SupplierID, req.OrderDate, req.Status, totalAmount, userID).Scan(
		&orderID,
		&orderDate,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Insert the purchase order items
	for _, item := range req.Items {
		_, err = tx.Exec(`
            Insert Into Purchase_Order_Detail (Purchase_Order_Id, Product_Id, Quantity, Price)
            Values ($1, $2, $3, $4)
        `, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Return the created purchase order
	return &GetPurchaseOrderResponse{
		ID:           orderID,
		SupplierID:   req.SupplierID,
		SupplierName: supplierName,
		OrderDate:    orderDate.Format(time.RFC3339),
		Status:       req.Status,
		TotalAmount:  totalAmount,
		CreatedAt:    createdAt.Format(time.RFC3339),
		UpdatedAt:    updatedAt.Format(time.RFC3339),
	}, nil
}

// Update updates an existing purchase order (only date as specified)
func (r *PurchaseOrderRepositoryImpl) Update(req UpdatePurchaseOrderRequest) (*GetPurchaseOrderResponse, error) {
	// Build dynamic update query
	var query strings.Builder
	var params []interface{}
	paramCount := 1

	query.WriteString(`Update Purchase_Order Set Updated_At = Now()`)

	if req.OrderDate != "" {
		query.WriteString(`, Order_Date = $` + strconv.Itoa(paramCount))
		params = append(params, req.OrderDate)
		paramCount++
	}

	if req.Status != "" {
		query.WriteString(`, Status = $` + strconv.Itoa(paramCount))
		params = append(params, req.Status)
		paramCount++
	}

	if req.SupplierID != "" {
		query.WriteString(`, Supplier_Id = $` + strconv.Itoa(paramCount))
		params = append(params, req.SupplierID)
		paramCount++
	}

	// Add WHERE clause and RETURNING
	query.WriteString(` WHERE Id = $` + strconv.Itoa(paramCount) + ` AND Cancelled_At IS NULL
        RETURNING Id`)
	params = append(params, req.ID)

	// Execute update
	var orderID string
	err := r.db.QueryRow(query.String(), params...).Scan(&orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Purchase Order tidak dapat diubah")
		}
		return nil, err
	}

	// Fetch the updated purchase order to return
	return r.GetByID(orderID)
}

// Cancel soft-deletes a purchase order by setting cancelled_at timestamp
func (r *PurchaseOrderRepositoryImpl) Cancel(id string) error {
	// First check if purchase order exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Soft cancel the purchase order
	result, err := r.db.Exec(`
		Update Purchase_Order
		Set Cancelled_At = Now(), Status = 'cancelled'
		Where Id = $1 And Cancelled_At Is Null
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Purchase order tidak ditemukan atau sudah dibatalkan")
	}

	return nil
}

// PurchaseOrderDetailRepository interface defines methods for purchase order detail operations
type PurchaseOrderDetailRepository interface {
	CheckProductByID(productID string) (bool, error)
	GetAllByOrderID(orderID string, page, pageSize int) ([]GetPurchaseOrderItemResponse, int, error)
	GetByID(id string) (*GetPurchaseOrderItemResponse, error)
	Create(orderID string, req CreatePurchaseOrderItemRequest) (*GetPurchaseOrderItemResponse, error)
	Update(req UpdatePurchaseOrderItemRequest) (*GetPurchaseOrderItemResponse, error)
	Delete(id string) error
}

// PurchaseOrderDetailRepositoryImpl implements PurchaseOrderDetailRepository
type PurchaseOrderDetailRepositoryImpl struct {
	db *sql.DB
}

// NewPurchaseOrderDetailRepository creates a new instance of PurchaseOrderDetailRepositoryImpl
func NewPurchaseOrderDetailRepository(db *sql.DB) PurchaseOrderDetailRepository {
	return &PurchaseOrderDetailRepositoryImpl{db: db}
}

// GetAllByOrderID fetches all purchase order details for a specific order with pagination
func (r *PurchaseOrderDetailRepositoryImpl) GetAllByOrderID(orderID string, page, pageSize int) ([]GetPurchaseOrderItemResponse, int, error) {
	var items []GetPurchaseOrderItemResponse
	var totalItems int

	// Initialize query builders
	var countQuery strings.Builder
	var query strings.Builder

	// Build count query
	countQuery.WriteString(`
        Select Count(Id)
        From Purchase_Order_Detail
        Where Purchase_Order_Id = $1
    `)

	// Get total count
	err := r.db.QueryRow(countQuery.String(), orderID).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// Build main query
	query.WriteString(`
        Select Pod.Id, Pod.Product_Id, Pod.Quantity, Pod.Price,
               (Pod.Quantity * Pod.Price) As Subtotal,
               Coalesce(P.Name, '') As Description
        From Purchase_Order_Detail Pod
        Left Join Product P On Pod.Product_Id = P.Id
        Where Pod.Purchase_Order_Id = $1
        Order By Pod.Created_At Desc
    `)

	// Add pagination
	query.WriteString(" LIMIT $2 OFFSET $3")

	// Calculate offset
	offset := (page - 1) * pageSize

	// Execute query with pagination
	rows, err := r.db.Query(query.String(), orderID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var item GetPurchaseOrderItemResponse
		errScan := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
			&item.ProductName,
		)
		if errScan != nil {
			return nil, 0, errScan
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, totalItems, nil
}

// GetByID fetches a purchase order detail by ID
func (r *PurchaseOrderDetailRepositoryImpl) GetByID(id string) (*GetPurchaseOrderItemResponse, error) {
	var item GetPurchaseOrderItemResponse

	err := r.db.QueryRow(`
		Select Pod.Id, Pod.Product_Id, Pod.Quantity, Pod.Price,
		       (Pod.Quantity * Pod.Price) As Subtotal,
		       P.Name As Description
		From Purchase_Order_Detail Pod
		Left Join Product P On Pod.Product_Id = P.Id
		Where Pod.Id = $1
	`, id).Scan(
		&item.ID,
		&item.ProductID,
		&item.Quantity,
		&item.Price,
		&item.Subtotal,
		&item.ProductName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Item Purchase Order tidak ditemukan!")
		}
		return nil, err
	}

	return &item, nil
}

// Create inserts a new purchase order detail
func (r *PurchaseOrderDetailRepositoryImpl) Create(orderID string, req CreatePurchaseOrderItemRequest) (*GetPurchaseOrderItemResponse, error) {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Verify order exists and is not cancelled
	var orderExists bool
	err = tx.QueryRow(`
		Select Exists(Select 1 From Purchase_Order Where Id = $1 And Cancelled_At Is Null)
	`, orderID).Scan(&orderExists)

	if err != nil {
		return nil, err
	}

	if !orderExists {
		return nil, errors.New("Purchase order tidak ditemukan atau sudah dibatalkan")
	}

	// Insert the detail
	var detailID string
	var subtotal float64
	err = tx.QueryRow(`
		Insert Into Purchase_Order_Detail (Purchase_Order_Id, Product_Id, Quantity, Price)
		Values ($1, $2, $3, $4)
		Returning Id, (Quantity * Price) As Subtotal
	`, orderID, req.ProductID, req.Quantity, req.Price).Scan(&detailID, &subtotal)

	if err != nil {
		return nil, err
	}

	// Update the total amount in the purchase order
	_, err = tx.Exec(`
		Update Purchase_Order
		Set Total_Amount = (Select Sum(Quantity * Price) From Purchase_Order_Detail Where Purchase_Order_Id = $1),
		    Updated_At = Now()
		Where Id = $1
	`, orderID)

	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Get product name for productName
	var productName string
	err = r.db.QueryRow(`
		Select Name From Product Where Id = $1
	`, req.ProductID).Scan(&productName)

	if err != nil {
		// Not critical, proceed without productName
		productName = ""
	}

	// Return the created item
	return &GetPurchaseOrderItemResponse{
		ID:          detailID,
		ProductID:   req.ProductID,
		Quantity:    float64(req.Quantity),
		Price:       req.Price,
		Subtotal:    subtotal,
		ProductName: productName,
	}, nil
}

// Update updates an existing purchase order detail
func (r *PurchaseOrderDetailRepositoryImpl) Update(req UpdatePurchaseOrderItemRequest) (*GetPurchaseOrderItemResponse, error) {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get purchase order ID and product ID for this detail
	var orderID, productID string
	err = tx.QueryRow(`Select Purchase_Order_Id, Product_Id From Purchase_Order_Detail Where Id = $1`, req.ID).Scan(&orderID, &productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Item tidak ditemukan")
		}
		return nil, err
	}

	// Verify order is not cancelled
	var cancelledAt sql.NullTime
	err = tx.QueryRow(`Select Cancelled_At From Purchase_Order Where Id = $1`, orderID).Scan(&cancelledAt)
	if err != nil {
		return nil, err
	}

	if cancelledAt.Valid {
		return nil, errors.New("Tidak dapat mengubah item pada Purchase Order yang telah dibatalkan")
	}

	// Build dynamic update query
	var query strings.Builder
	var params []interface{}
	paramCount := 1

	query.WriteString(`Update Purchase_Order_Detail Set Updated_At = Now()`)

	if req.Quantity > 0 {
		query.WriteString(`, Quantity = $` + strconv.Itoa(paramCount))
		params = append(params, req.Quantity)
		paramCount++
	}

	if req.Price > 0 {
		query.WriteString(`, Price = $` + strconv.Itoa(paramCount))
		params = append(params, req.Price)
		paramCount++
	}

	// Add WHERE clause and RETURNING
	query.WriteString(` WHERE Id = $` + strconv.Itoa(paramCount) + `
        RETURNING Quantity, Price`)
	params = append(params, req.ID)

	// Execute update
	var quantity float64
	var price float64
	err = tx.QueryRow(query.String(), params...).Scan(&quantity, &price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Item tidak dapat diubah")
		}
		return nil, err
	}

	// Calculate subtotal
	subtotal := quantity * price

	// Update the total amount in the purchase order
	_, err = tx.Exec(`
        Update Purchase_Order
        Set Total_Amount = (
            Select Sum(Quantity * Price)
            From Purchase_Order_Detail
            Where Purchase_Order_Id = $1
        ),
        Updated_At = Now()
        Where Id = $1
    `, orderID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Get product name for productName
	var productName string
	err = r.db.QueryRow(`Select Coalesce(Name, '') From Product Where Id = $1`, productID).Scan(&productName)
	if err != nil {
		// Not critical, proceed without productName
		productName = ""
	}

	// Return the updated item
	return &GetPurchaseOrderItemResponse{
		ID:          req.ID,
		ProductID:   productID,
		Quantity:    quantity,
		Price:       price,
		Subtotal:    subtotal,
		ProductName: productName,
	}, nil
}

// Delete removes a purchase order detail
func (r *PurchaseOrderDetailRepositoryImpl) Delete(id string) error {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get purchase order ID for this detail
	var orderID string
	err = tx.QueryRow(`Select Purchase_Order_Id From Purchase_Order_Detail Where Id = $1`, id).Scan(&orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("Item tidak ditemukan")
		}
		return err
	}

	// Verify order is not cancelled
	var cancelledAt sql.NullTime
	err = tx.QueryRow(`Select Cancelled_At From Purchase_Order Where Id = $1`, orderID).Scan(&cancelledAt)
	if err != nil {
		return err
	}

	if cancelledAt.Valid {
		return errors.New("Tidak dapat menghapus item pada Purchase Order yang telah dibatalkan")
	}

	// Delete the detail
	result, err := tx.Exec(`Delete From Purchase_Order_Detail Where Id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Item tidak ditemukan")
	}

	// Check if there are any details remaining and update the order's total amount
	var detailCount int
	err = tx.QueryRow(`Select Count(*) From Purchase_Order_Detail Where Purchase_Order_Id = $1`, orderID).Scan(&detailCount)
	if err != nil {
		return err
	}

	if detailCount > 0 {
		// Update the total amount in the purchase order
		_, err = tx.Exec(`
            Update Purchase_Order 
            Set Total_Amount = (
                Select Sum(Quantity * Price) 
                From Purchase_Order_Detail 
                Where Purchase_Order_Id = $1
            ),
            Updated_At = Now()
            Where Id = $1
        `, orderID)
		if err != nil {
			return err
		}
	} else {
		// If no details left, set total to 0
		_, err = tx.Exec(`Update Purchase_Order Set Total_Amount = 0, Updated_At = Now() Where Id = $1`, orderID)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// CheckProductByID checks if a product with the given ID exists
func (r *PurchaseOrderDetailRepositoryImpl) CheckProductByID(productID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`Select Exists(Select 1 From Product Where Id = $1)`, productID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
