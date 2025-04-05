package sales

import (
	"errors"
	"fmt"
	"sinartimur-go/utils"
	"time"
)

// SalesService is the service for the Product domain.
type SalesService struct {
	repo SalesRepository
}

// NewSalesService creates a new instance of SalesService
func NewSalesService(repo SalesRepository) *SalesService {
	return &SalesService{repo: repo}
}

// GetSalesOrders retrieves a paginated list of sales orders with optional filtering
func (s *SalesService) GetSalesOrders(req GetSalesOrdersRequest) ([]GetSalesOrdersResponse, int, error) {
	return s.repo.GetSalesOrders(req)
}

// GetSalesOrderDetail retrieves detailed information about a sales order including its items
func (s *SalesService) GetSalesOrderDetail(orderID string) (*GetSalesOrderDetailResponse, error) {
	// Get the sales order header information
	orderHeader, err := s.repo.GetSalesOrderByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("sales order tidak ditemukan: %w", err)
	}

	// Get the sales order details/items
	orderItems, err := s.repo.GetSalesOrderDetails(orderID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil detail pesanan: %w", err)
	}

	// Combine the results
	result := &GetSalesOrderDetailResponse{
		ID:             orderHeader.ID,
		SerialID:       orderHeader.SerialID,
		CustomerID:     orderHeader.CustomerID,
		CustomerName:   orderHeader.CustomerName,
		OrderDate:      orderHeader.OrderDate,
		Status:         orderHeader.Status,
		TotalAmount:    orderHeader.TotalAmount,
		PaymentMethod:  orderHeader.PaymentMethod,
		PaymentDueDate: orderHeader.PaymentDueDate,
		CreatedBy:      orderHeader.CreatedBy,
		CreatedAt:      orderHeader.CreatedAt,
		UpdatedAt:      orderHeader.UpdatedAt,
		Items:          orderItems,
	}

	return result, nil
}

// CreateSalesOrder creates a new sales order with items and optional invoice creation
func (s *SalesService) CreateSalesOrder(req CreateSalesOrderRequest, userID string) (*CreateSalesOrderResponse, error) {
	// Validate payment information
	if req.PaymentMethod == "paylater" && req.PaymentDueDate == "" {
		return nil, fmt.Errorf("tanggal jatuh tempo pembayaran diperlukan untuk metode pembayaran paylater")
	}

	// Validate due date is in the future for paylater
	if req.PaymentDueDate != "" {
		dueDate, err := time.Parse(time.RFC3339, req.PaymentDueDate)
		if err != nil {
			return nil, fmt.Errorf("format tanggal jatuh tempo tidak valid: %w", err)
		}

		if dueDate.Before(time.Now()) {
			return nil, fmt.Errorf("tanggal jatuh tempo harus di masa depan")
		}
	}

	// Validate items
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("pesanan harus memiliki minimal satu item")
	}

	for i, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("kuantitas item #%d harus lebih dari 0", i+1)
		}
		if item.UnitPrice < 0 {
			return nil, fmt.Errorf("harga satuan item #%d tidak boleh negatif", i+1)
		}

		// Check storage allocations
		var totalAllocated float64
		for _, alloc := range item.StorageAllocations {
			if alloc.Quantity <= 0 {
				return nil, fmt.Errorf("kuantitas alokasi penyimpanan harus lebih dari 0")
			}
			totalAllocated += alloc.Quantity
		}

		// Ensure total allocation matches requested quantity
		if totalAllocated != item.Quantity {
			return nil, fmt.Errorf(
				"total alokasi penyimpanan (%.2f) tidak sama dengan kuantitas yang diminta (%.2f) untuk item #%d",
				totalAllocated, item.Quantity, i+1)
		}
	}

	return s.repo.CreateSalesOrder(req, userID)
}

// UpdateSalesOrder updates basic information of a sales order
func (s *SalesService) UpdateSalesOrder(req UpdateSalesOrderRequest) (*UpdateSalesOrderResponse, error) {
	// Validate order ID
	if req.ID == "" {
		return nil, fmt.Errorf("ID pesanan tidak boleh kosong")
	}

	// Validate payment due date format if provided
	if req.PaymentMethod == "paylater" && req.PaymentDueDate != "" {
		dueDate, err := time.Parse(time.RFC3339, req.PaymentDueDate)
		if err != nil {
			return nil, fmt.Errorf("format tanggal jatuh tempo tidak valid: %w", err)
		}

		if dueDate.Before(time.Now()) {
			return nil, fmt.Errorf("tanggal jatuh tempo harus di masa depan")
		}
	}

	return s.repo.UpdateSalesOrder(req)
}

// CancelSalesOrder cancels a sales order and restores inventory
func (s *SalesService) CancelSalesOrder(req CancelSalesOrderRequest, userID string) error {
	// Validate order ID
	if req.SalesOrderID == "" {
		return fmt.Errorf("ID pesanan tidak boleh kosong")
	}

	return s.repo.CancelSalesOrder(req, userID)
}

// AddItemToSalesOrder adds a new item to an existing sales order
func (s *SalesService) AddItemToSalesOrder(req AddItemToSalesOrderRequest) (*UpdateItemResponse, error) {
	// Validate basic parameters
	if req.SalesOrderID == "" {
		return nil, fmt.Errorf("ID pesanan tidak boleh kosong")
	}

	if req.ProductID == "" || req.BatchID == "" {
		return nil, fmt.Errorf("ID produk dan batch harus diisi")
	}

	// Validate quantity
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("kuantitas harus lebih dari 0")
	}

	// Validate unit price
	if req.UnitPrice < 0 {
		return nil, fmt.Errorf("harga satuan tidak boleh negatif")
	}

	// Check storage allocations
	if len(req.StorageAllocations) == 0 {
		return nil, fmt.Errorf("minimal satu lokasi penyimpanan harus dipilih")
	}

	var totalAllocated float64
	for _, alloc := range req.StorageAllocations {
		if alloc.Quantity <= 0 {
			return nil, fmt.Errorf("kuantitas alokasi penyimpanan harus lebih dari 0")
		}
		totalAllocated += alloc.Quantity
	}

	// Ensure total allocation matches requested quantity
	if totalAllocated != req.Quantity {
		return nil, fmt.Errorf(
			"total alokasi penyimpanan (%.2f) tidak sama dengan kuantitas yang diminta (%.2f)",
			totalAllocated, req.Quantity)
	}

	return s.repo.AddItemToSalesOrder(req)
}

// UpdateSalesOrderItem updates an existing item in a sales order
func (s *SalesService) UpdateSalesOrderItem(req UpdateItemRequest) (*UpdateItemResponse, error) {
	// Validate basic parameters
	if req.SalesOrderID == "" || req.DetailID == "" {
		return nil, fmt.Errorf("ID pesanan dan ID detail harus diisi")
	}

	// If storage allocations are provided, validate them
	if len(req.StorageAllocations) > 0 {
		var totalAllocated float64
		for _, alloc := range req.StorageAllocations {
			if alloc.Quantity <= 0 {
				return nil, fmt.Errorf("kuantitas alokasi penyimpanan harus lebih dari 0")
			}
			totalAllocated += alloc.Quantity
		}

		// If quantity is provided, ensure it matches storage allocations
		if req.Quantity > 0 && totalAllocated != req.Quantity {
			return nil, fmt.Errorf(
				"total alokasi penyimpanan (%.2f) tidak sama dengan kuantitas yang diminta (%.2f)",
				totalAllocated, req.Quantity)
		}
	}

	// Validate quantity if provided
	if req.Quantity < 0 {
		return nil, fmt.Errorf("kuantitas tidak boleh negatif")
	}

	// Validate unit price if provided
	if req.UnitPrice < 0 {
		return nil, fmt.Errorf("harga satuan tidak boleh negatif")
	}

	return s.repo.UpdateSalesOrderItem(req)
}

// DeleteSalesOrderItem removes an item from a sales order and restores inventory
func (s *SalesService) DeleteSalesOrderItem(req DeleteItemRequest) error {
	// Validate parameters
	if req.SalesOrderID == "" || req.DetailID == "" {
		return fmt.Errorf("ID pesanan dan ID detail harus diisi")
	}

	return s.repo.DeleteSalesOrderItem(req)
}

// GetSalesInvoices returns a paginated list of sales invoices
func (s *SalesService) GetSalesInvoices(req GetSalesInvoicesRequest) ([]GetSalesInvoicesResponse, int, error) {
	// Call repository to fetch invoices
	invoices, totalItems, err := s.repo.GetSalesInvoices(req)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar faktur: %w", err)
	}

	return invoices, totalItems, nil
}

// CreateSalesInvoice creates a new invoice for a sales order
func (s *SalesService) CreateSalesInvoice(req CreateSalesInvoiceRequest, userID string) (*CreateSalesInvoiceResponse, error) {
	// Validate request
	if req.SalesOrderID == "" {
		return nil, errors.New("ID pesanan penjualan wajib diisi")
	}

	// Call repository to create invoice
	response, err := s.repo.CreateSalesInvoice(req, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat faktur penjualan: %w", err)
	}

	return response, nil
}

// CancelSalesInvoice cancels an existing sales invoice
func (s *SalesService) CancelSalesInvoice(req CancelSalesInvoiceRequest, userID string) error {
	// Validate request
	if req.InvoiceID == "" {
		return errors.New("ID faktur wajib diisi")
	}

	// Call repository to cancel invoice
	err := s.repo.CancelSalesInvoice(req, userID)
	if err != nil {
		return fmt.Errorf("gagal membatalkan faktur: %w", err)
	}

	return nil
}

// ReturnInvoiceItems processes returns for invoice items
func (s *SalesService) ReturnInvoiceItems(req ReturnInvoiceItemsRequest, userID string) (*ReturnInvoiceItemsResponse, error) {
	// Validate request
	if req.InvoiceID == "" {
		return nil, errors.New("ID faktur wajib diisi")
	}

	if len(req.ReturnItems) == 0 {
		return nil, errors.New("tidak ada item yang ditentukan untuk pengembalian")
	}

	for _, item := range req.ReturnItems {
		if item.DetailID == "" {
			return nil, errors.New("ID detail item wajib diisi")
		}
		if item.Quantity <= 0 {
			return nil, errors.New("jumlah pengembalian harus lebih dari 0")
		}

		totalStorageQty := 0.0
		for _, storage := range item.StorageReturns {
			if storage.StorageID == "" {
				return nil, errors.New("ID lokasi penyimpanan wajib diisi")
			}
			if storage.Quantity <= 0 {
				return nil, errors.New("jumlah pengembalian di penyimpanan harus lebih dari 0")
			}
			totalStorageQty += storage.Quantity
		}

		if !utils.FloatEquals(totalStorageQty, item.Quantity) {
			return nil, fmt.Errorf("total jumlah pengembalian di penyimpanan (%f) tidak sama dengan jumlah yang diminta (%f)",
				totalStorageQty, item.Quantity)
		}
	}

	// Call repository to process returns
	response, err := s.repo.ReturnInvoiceItems(req, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses pengembalian: %w", err)
	}

	return response, nil
}

// CancelInvoiceReturn cancels a previously processed return
func (s *SalesService) CancelInvoiceReturn(req CancelInvoiceReturnRequest, userID string) error {
	// Validate request
	if req.ReturnID == "" {
		return errors.New("ID pengembalian wajib diisi")
	}

	// Call repository to cancel return
	err := s.repo.CancelInvoiceReturn(req, userID)
	if err != nil {
		return fmt.Errorf("gagal membatalkan pengembalian: %w", err)
	}

	return nil
}

// CreateDeliveryNote creates a new delivery note for a sales invoice
func (s *SalesService) CreateDeliveryNote(req CreateDeliveryNoteRequest, userID string) (*CreateDeliveryNoteResponse, error) {
	// Validate request
	if req.SalesInvoiceID == "" {
		return nil, errors.New("ID faktur penjualan wajib diisi")
	}

	// Validate driver and recipient names
	if req.DriverName == "" {
		return nil, errors.New("nama pengemudi wajib diisi")
	}

	if req.RecipientName == "" {
		return nil, errors.New("nama penerima wajib diisi")
	}

	// If delivery date is provided, validate format
	if req.DeliveryDate != "" {
		_, err := time.Parse(time.RFC3339, req.DeliveryDate)
		if err != nil {
			return nil, errors.New("format tanggal pengiriman tidak valid")
		}
	}

	// Call repository to create delivery note
	response, err := s.repo.CreateDeliveryNote(req, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat surat jalan: %w", err)
	}

	return response, nil
}

// CancelDeliveryNote cancels an existing delivery note
func (s *SalesService) CancelDeliveryNote(req CancelDeliveryNoteRequest, userID string) error {
	// Validate request
	if req.DeliveryNoteID == "" {
		return errors.New("ID surat jalan wajib diisi")
	}

	// Call repository to cancel delivery note
	err := s.repo.CancelDeliveryNote(req, userID)
	if err != nil {
		return fmt.Errorf("gagal membatalkan surat jalan: %w", err)
	}

	return nil
}
