package sales

import (
	"sinartimur-go/utils"
	"time"
)

// SalesOrder represents a sales order entity from the database
type SalesOrder struct {
	ID             string     `json:"id"`
	SerialID       string     `json:"serial_id"`
	CustomerID     string     `json:"customer_id"`
	CustomerName   string     `json:"customer_name"`
	OrderDate      time.Time  `json:"order_date"`
	Status         string     `json:"status"`
	PaymentMethod  string     `json:"payment_method"`
	PaymentDueDate *time.Time `json:"payment_due_date,omitempty"`
	TotalAmount    float64    `json:"total_amount"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CancelledAt    *time.Time `json:"cancelled_at,omitempty"`
}

// SalesOrderDetail represents a sales order detail entity
type SalesOrderDetail struct {
	ID           string    `json:"id"`
	SalesOrderID string    `json:"sales_order_id"`
	BatchID      string    `json:"batch_id"`
	ProductID    string    `json:"product_id"`
	Quantity     float64   `json:"quantity"`
	UnitPrice    float64   `json:"unit_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SalesOrderStorage represents storage allocation for a sales order detail
type SalesOrderStorage struct {
	ID                 string    `json:"id"`
	SalesOrderDetailID string    `json:"sales_order_detail_id"`
	StorageID          string    `json:"storage_id"`
	BatchID            string    `json:"batch_id"`
	Quantity           float64   `json:"quantity"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// GetSalesOrdersRequest defines the parameters for fetching sales orders
type GetSalesOrdersRequest struct {
	CustomerID    string `json:"customer_id,omitempty" validate:"omitempty,uuid"`
	Status        string `json:"status,omitempty" validate:"omitempty,oneof=order invoice delivery partial_return return cancel"`
	PaymentMethod string `json:"payment_method,omitempty" validate:"omitempty,oneof=cash paylater"`
	StartDate     string `json:"start_date,omitempty" validate:"omitempty,rfc3339"`
	EndDate       string `json:"end_date,omitempty" validate:"omitempty,rfc3339"`
	SerialID      string `json:"serial_id,omitempty"`
	utils.PaginationParameter
}

// SalesOrderPaginatedResponse defines a paginated response for sales orders
type SalesOrderPaginatedResponse struct {
	utils.PaginationParameter
	Items []GetSalesOrdersResponse `json:"items"`
}

// GetSalesOrdersResponse defines the response for fetching sales orders
type GetSalesOrdersResponse struct {
	ID             string  `json:"id"`
	SerialID       string  `json:"serial_id"`
	CustomerID     string  `json:"customer_id"`
	CustomerName   string  `json:"customer_name"`
	OrderDate      string  `json:"order_date"`
	Status         string  `json:"status"`
	PaymentMethod  string  `json:"payment_method"`
	PaymentDueDate string  `json:"payment_due_date,omitempty"`
	TotalAmount    float64 `json:"total_amount"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	CancelledAt    string  `json:"cancelled_at,omitempty"`
}

// GetSalesOrderDetailResponse defines the response for fetching a sales order's details
type GetSalesOrderDetailResponse struct {
	ID                 string                      `json:"id"`
	SalesOrderID       string                      `json:"sales_order_id"`
	ProductID          string                      `json:"product_id"`
	ProductName        string                      `json:"product_name"`
	BatchID            string                      `json:"batch_id"`
	BatchSKU           string                      `json:"batch_sku"`
	Quantity           float64                     `json:"quantity"`
	UnitPrice          float64                     `json:"unit_price"`
	TotalPrice         float64                     `json:"total_price"`
	StorageAllocations []SalesOrderStorageResponse `json:"storage_allocations"`
}

// SalesOrderStorageResponse defines the storage allocation response
type SalesOrderStorageResponse struct {
	ID          string  `json:"id"`
	StorageID   string  `json:"storage_id"`
	StorageName string  `json:"storage_name"`
	Quantity    float64 `json:"quantity"`
}

// CreateSalesOrderRequest defines the request for creating a sales order
type CreateSalesOrderRequest struct {
	CustomerID     string                  `json:"customer_id" validate:"required,uuid"`
	PaymentMethod  string                  `json:"payment_method" validate:"required,oneof=cash paylater"`
	PaymentDueDate string                  `json:"payment_due_date,omitempty" validate:"omitempty,rfc3339"`
	Items          []SalesOrderItemRequest `json:"items" validate:"required,min=1,dive"`
	CreateInvoice  bool                    `json:"create_invoice" validate:"omitempty"`
}

// SalesOrderItemRequest defines an item in a create sales order request
type SalesOrderItemRequest struct {
	ProductID          string                     `json:"product_id" validate:"required,uuid"`
	BatchID            string                     `json:"batch_id" validate:"required,uuid"`
	Quantity           float64                    `json:"quantity" validate:"required,gt=0"`
	UnitPrice          float64                    `json:"unit_price" validate:"required,gt=0"`
	StorageAllocations []StorageAllocationRequest `json:"storage_allocations" validate:"required,min=1,dive"`
}

// StorageAllocationRequest defines a storage allocation for an item
type StorageAllocationRequest struct {
	StorageID string  `json:"storage_id" validate:"required,uuid"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
}

// CreateSalesOrderResponse defines the response for creating a sales order
type CreateSalesOrderResponse struct {
	ID              string  `json:"id"`
	SerialID        string  `json:"serial_id"`
	CustomerID      string  `json:"customer_id"`
	CustomerName    string  `json:"customer_name"`
	Status          string  `json:"status"`
	PaymentMethod   string  `json:"payment_method"`
	PaymentDueDate  string  `json:"payment_due_date,omitempty"`
	TotalAmount     float64 `json:"total_amount"`
	CreatedAt       string  `json:"created_at"`
	InvoiceID       string  `json:"invoice_id,omitempty"`
	InvoiceSerialID string  `json:"invoice_serial_id,omitempty"`
}

// UpdateSalesOrderRequest defines the request for updating a sales order
type UpdateSalesOrderRequest struct {
	ID             string `json:"id" validate:"required,uuid"`
	CustomerID     string `json:"customer_id,omitempty" validate:"omitempty,uuid"`
	PaymentMethod  string `json:"payment_method,omitempty" validate:"omitempty,oneof=cash paylater"`
	PaymentDueDate string `json:"payment_due_date,omitempty" validate:"omitempty,rfc3339"`
}

// UpdateSalesOrderResponse defines the response for updating a sales order
type UpdateSalesOrderResponse struct {
	ID             string `json:"id"`
	SerialID       string `json:"serial_id"`
	CustomerID     string `json:"customer_id"`
	Status         string `json:"status"`
	PaymentMethod  string `json:"payment_method"`
	PaymentDueDate string `json:"payment_due_date,omitempty"`
	UpdatedAt      string `json:"updated_at"`
}

// AddItemToSalesOrderRequest defines the request for adding an item to an existing sales order
type AddItemToSalesOrderRequest struct {
	SalesOrderID       string                     `json:"sales_order_id" validate:"required,uuid"`
	ProductID          string                     `json:"product_id" validate:"required,uuid"`
	BatchID            string                     `json:"batch_id" validate:"required,uuid"`
	Quantity           float64                    `json:"quantity" validate:"required,gt=0"`
	UnitPrice          float64                    `json:"unit_price" validate:"required,gt=0"`
	StorageAllocations []StorageAllocationRequest `json:"storage_allocations" validate:"required,min=1,dive"`
}

// UpdateItemRequest defines the request for updating an item in a sales order
type UpdateItemRequest struct {
	SalesOrderID       string                     `json:"sales_order_id" validate:"required,uuid"`
	DetailID           string                     `json:"detail_id" validate:"required,uuid"`
	Quantity           float64                    `json:"quantity" validate:"omitempty,gt=0"`
	UnitPrice          float64                    `json:"unit_price" validate:"omitempty,gt=0"`
	StorageAllocations []StorageAllocationRequest `json:"storage_allocations" validate:"omitempty,dive"`
}

// DeleteItemRequest defines the request for deleting an item from a sales order
type DeleteItemRequest struct {
	SalesOrderID string `json:"sales_order_id" validate:"required,uuid"`
	DetailID     string `json:"detail_id" validate:"required,uuid"`
}

// CancelSalesOrderRequest defines the request for cancelling a sales order
type CancelSalesOrderRequest struct {
	SalesOrderID string `json:"sales_order_id" validate:"required,uuid"`
}

// UpdateItemResponse defines the response for updating or adding an item
type UpdateItemResponse struct {
	DetailID    string  `json:"detail_id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	BatchID     string  `json:"batch_id"`
	BatchSKU    string  `json:"batch_sku"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

// SalesInvoice represents a sales invoice entity from the database
type SalesInvoice struct {
	ID           string     `json:"id"`
	SalesOrderID string     `json:"sales_order_id"`
	SerialID     string     `json:"serial_id"`
	InvoiceDate  time.Time  `json:"invoice_date"`
	TotalAmount  float64    `json:"total_amount"`
	CreatedBy    string     `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CancelledAt  *time.Time `json:"cancelled_at,omitempty"`
	CancelledBy  *string    `json:"cancelled_by,omitempty"`
}

// GetSalesInvoicesRequest defines parameters for fetching sales invoices
type GetSalesInvoicesRequest struct {
	CustomerID string `json:"customer_id,omitempty" validate:"omitempty,uuid"`
	StartDate  string `json:"start_date,omitempty" validate:"omitempty,rfc3339"`
	EndDate    string `json:"end_date,omitempty" validate:"omitempty,rfc3339"`
	SerialID   string `json:"serial_id,omitempty"`
	Status     string `json:"status,omitempty" validate:"omitempty,oneof=active cancelled partially_returned returned"`
	utils.PaginationParameter
}

// GetSalesInvoicesResponse defines the response for fetching sales invoices
type GetSalesInvoicesResponse struct {
	ID               string  `json:"id"`
	SerialID         string  `json:"serial_id"`
	SalesOrderID     string  `json:"sales_order_id"`
	SalesOrderSerial string  `json:"sales_order_serial"`
	CustomerID       string  `json:"customer_id"`
	CustomerName     string  `json:"customer_name"`
	InvoiceDate      string  `json:"invoice_date"`
	TotalAmount      float64 `json:"total_amount"`
	Status           string  `json:"status"`
	HasDeliveryNote  bool    `json:"has_delivery_note"`
	CreatedBy        string  `json:"created_by"`
	CreatedAt        string  `json:"created_at"`
	CancelledAt      string  `json:"cancelled_at,omitempty"`
}

// SalesInvoicePaginatedResponse defines a paginated response for sales invoices
type SalesInvoicePaginatedResponse struct {
	utils.PaginationParameter
	Items []GetSalesInvoicesResponse `json:"items"`
}

// SalesInvoiceItemResponse defines the detail items in a sales invoice
type SalesInvoiceItemResponse struct {
	ID             string                      `json:"id"`
	SalesOrderID   string                      `json:"sales_order_id"`
	ProductID      string                      `json:"product_id"`
	ProductName    string                      `json:"product_name"`
	BatchID        string                      `json:"batch_id"`
	BatchSKU       string                      `json:"batch_sku"`
	Quantity       float64                     `json:"quantity"`
	ReturnedQty    float64                     `json:"returned_qty"`
	UnitPrice      float64                     `json:"unit_price"`
	TotalPrice     float64                     `json:"total_price"`
	StorageDetails []SalesOrderStorageResponse `json:"storage_details"`
}

// CreateSalesInvoiceRequest defines the request for creating a sales invoice
type CreateSalesInvoiceRequest struct {
	SalesOrderID string `json:"sales_order_id" validate:"required,uuid"`
}

// CreateSalesInvoiceResponse defines the response for creating a sales invoice
type CreateSalesInvoiceResponse struct {
	ID               string  `json:"id"`
	SerialID         string  `json:"serial_id"`
	SalesOrderID     string  `json:"sales_order_id"`
	SalesOrderSerial string  `json:"sales_order_serial"`
	CustomerID       string  `json:"customer_id"`
	CustomerName     string  `json:"customer_name"`
	InvoiceDate      string  `json:"invoice_date"`
	TotalAmount      float64 `json:"total_amount"`
	Status           string  `json:"status"`
	CreatedBy        string  `json:"created_by"`
	CreatedAt        string  `json:"created_at"`
}

// CancelSalesInvoiceRequest defines the request for cancelling a sales invoice
type CancelSalesInvoiceRequest struct {
	InvoiceID string `json:"invoice_id" validate:"required,uuid"`
}

// ReturnInvoiceItemsRequest defines the request for returning items from a sales invoice
type ReturnInvoiceItemsRequest struct {
	InvoiceID    string                     `json:"invoice_id" validate:"required,uuid"`
	ReturnItems  []InvoiceReturnItemRequest `json:"return_items" validate:"required,min=1,dive"`
	ReturnReason string                     `json:"return_reason,omitempty"`
}

// InvoiceReturnItemRequest defines a single item to be returned
type InvoiceReturnItemRequest struct {
	DetailID       string                 `json:"detail_id" validate:"required,uuid"`
	Quantity       float64                `json:"quantity" validate:"required,gt=0"`
	StorageReturns []StorageReturnRequest `json:"storage_returns" validate:"required,min=1,dive"`
}

// StorageReturnRequest defines storage details for a returned item
type StorageReturnRequest struct {
	StorageID string  `json:"storage_id" validate:"required,uuid"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
}

// ReturnInvoiceItemsResponse defines the response for returning items
type ReturnInvoiceItemsResponse struct {
	ReturnID      string  `json:"return_id"`
	InvoiceID     string  `json:"invoice_id"`
	ReturnedItems int     `json:"returned_items"`
	TotalQuantity float64 `json:"total_quantity"`
	ReturnDate    string  `json:"return_date"`
	ReturnStatus  string  `json:"return_status"`
	IsFullReturn  bool    `json:"is_full_return"`
}

// CancelInvoiceReturnRequest defines the request for cancelling a return
type CancelInvoiceReturnRequest struct {
	ReturnID string `json:"return_id" validate:"required,uuid"`
}

// SalesOrderReturn represents a sales order return entity
type SalesOrderReturn struct {
	ID                string     `json:"id"`
	ReturnSource      string     `json:"return_source"`
	DeliveryNoteID    *string    `json:"delivery_note_id,omitempty"`
	SalesOrderID      string     `json:"sales_order_id"`
	SalesDetailID     string     `json:"sales_detail_id"`
	ReturnQuantity    float64    `json:"return_quantity"`
	RemainingQuantity float64    `json:"remaining_quantity"`
	ReturnReason      *string    `json:"return_reason,omitempty"`
	ReturnStatus      string     `json:"return_status"`
	ReturnedBy        *string    `json:"returned_by,omitempty"`
	CancelledBy       *string    `json:"cancelled_by,omitempty"`
	ReturnedAt        time.Time  `json:"returned_at"`
	CancelledAt       *time.Time `json:"cancelled_at,omitempty"`
}

// SalesOrderReturnBatch represents a batch record for a sales return
type SalesOrderReturnBatch struct {
	ID             string    `json:"id"`
	SalesReturnID  string    `json:"sales_return_id"`
	BatchID        *string   `json:"batch_id,omitempty"`
	ReturnQuantity float64   `json:"return_quantity"`
	CreatedAt      time.Time `json:"created_at"`
}

// GetSalesReturnsResponse defines the response for fetching sales returns
type GetSalesReturnsResponse struct {
	ID                 string  `json:"id"`
	ReturnSource       string  `json:"return_source"`
	SalesOrderID       string  `json:"sales_order_id"`
	SalesOrderSerial   string  `json:"sales_order_serial"`
	InvoiceID          string  `json:"invoice_id,omitempty"`
	InvoiceSerial      string  `json:"invoice_serial,omitempty"`
	DeliveryNoteID     string  `json:"delivery_note_id,omitempty"`
	DeliveryNoteSerial string  `json:"delivery_note_serial,omitempty"`
	ReturnedItems      int     `json:"returned_items"`
	TotalQuantity      float64 `json:"total_quantity"`
	ReturnReason       string  `json:"return_reason,omitempty"`
	ReturnStatus       string  `json:"return_status"`
	ReturnedBy         string  `json:"returned_by"`
	ReturnedAt         string  `json:"returned_at"`
	CancelledAt        string  `json:"cancelled_at,omitempty"`
	CancelledBy        string  `json:"cancelled_by,omitempty"`
}

// SalesReturnItemResponse defines the details of returned items
type SalesReturnItemResponse struct {
	ReturnID    string  `json:"return_id"`
	DetailID    string  `json:"detail_id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	BatchID     string  `json:"batch_id"`
	BatchSKU    string  `json:"batch_sku"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

// DeliveryNote represents a delivery note entity from the database
type DeliveryNote struct {
	ID             string     `json:"id"`
	SerialID       string     `json:"serial_id"`
	SalesOrderID   string     `json:"sales_order_id"`
	SalesInvoiceID string     `json:"sales_invoice_id"`
	DeliveryDate   time.Time  `json:"delivery_date"`
	DriverName     string     `json:"driver_name"`
	RecipientName  string     `json:"recipient_name"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CancelledAt    *time.Time `json:"cancelled_at,omitempty"`
	CancelledBy    *string    `json:"cancelled_by,omitempty"`
}

// CreateDeliveryNoteRequest defines the request for creating a delivery note
type CreateDeliveryNoteRequest struct {
	SalesInvoiceID string `json:"sales_invoice_id" validate:"required,uuid"`
	DriverName     string `json:"driver_name" validate:"required"`
	RecipientName  string `json:"recipient_name" validate:"required"`
	DeliveryDate   string `json:"delivery_date,omitempty" validate:"omitempty,rfc3339"`
}

// CreateDeliveryNoteResponse defines the response for creating a delivery note
type CreateDeliveryNoteResponse struct {
	ID                 string `json:"id"`
	SerialID           string `json:"serial_id"`
	SalesOrderID       string `json:"sales_order_id"`
	SalesOrderSerial   string `json:"sales_order_serial"`
	SalesInvoiceID     string `json:"sales_invoice_id"`
	SalesInvoiceSerial string `json:"sales_invoice_serial"`
	DeliveryDate       string `json:"delivery_date"`
	DriverName         string `json:"driver_name"`
	RecipientName      string `json:"recipient_name"`
	CreatedBy          string `json:"created_by"`
	CreatedAt          string `json:"created_at"`
}

// CancelDeliveryNoteRequest defines the request for cancelling a delivery note
type CancelDeliveryNoteRequest struct {
	DeliveryNoteID string `json:"delivery_note_id" validate:"required,uuid"`
}

// GetDeliveryNotesRequest defines parameters for fetching delivery notes
type GetDeliveryNotesRequest struct {
	SalesOrderID   string `json:"sales_order_id,omitempty" validate:"omitempty,uuid"`
	SalesInvoiceID string `json:"sales_invoice_id,omitempty" validate:"omitempty,uuid"`
	StartDate      string `json:"start_date,omitempty" validate:"omitempty,rfc3339"`
	EndDate        string `json:"end_date,omitempty" validate:"omitempty,rfc3339"`
	SerialID       string `json:"serial_id,omitempty"`
	Status         string `json:"status,omitempty" validate:"omitempty,oneof=active cancelled partially_returned returned"`
	utils.PaginationParameter
}

// GetDeliveryNotesResponse defines the response for fetching delivery notes
type GetDeliveryNotesResponse struct {
	ID                 string `json:"id"`
	SerialID           string `json:"serial_id"`
	SalesOrderID       string `json:"sales_order_id"`
	SalesOrderSerial   string `json:"sales_order_serial"`
	SalesInvoiceID     string `json:"sales_invoice_id"`
	SalesInvoiceSerial string `json:"sales_invoice_serial"`
	DeliveryDate       string `json:"delivery_date"`
	DriverName         string `json:"driver_name"`
	RecipientName      string `json:"recipient_name"`
	Status             string `json:"status"`
	CreatedBy          string `json:"created_by"`
	CreatedAt          string `json:"created_at"`
	CancelledAt        string `json:"cancelled_at,omitempty"`
}

// DeliveryNotePaginatedResponse defines a paginated response for delivery notes
type DeliveryNotePaginatedResponse struct {
	TotalItems  int                        `json:"total_items"`
	TotalPages  int                        `json:"total_pages"`
	CurrentPage int                        `json:"current_page"`
	PageSize    int                        `json:"page_size"`
	HasNext     bool                       `json:"has_next"`
	HasPrevious bool                       `json:"has_previous"`
	Items       []GetDeliveryNotesResponse `json:"items"`
}

// ReturnDeliveryItemsRequest defines the request for returning items from a delivery note
type ReturnDeliveryItemsRequest struct {
	DeliveryNoteID string                      `json:"delivery_note_id" validate:"required,uuid"`
	ReturnItems    []DeliveryReturnItemRequest `json:"return_items" validate:"required,min=1,dive"`
	ReturnReason   string                      `json:"return_reason,omitempty"`
}

// DeliveryReturnItemRequest defines a single item to be returned from delivery
type DeliveryReturnItemRequest struct {
	DetailID       string                 `json:"detail_id" validate:"required,uuid"`
	Quantity       float64                `json:"quantity" validate:"required,gt=0"`
	StorageReturns []StorageReturnRequest `json:"storage_returns" validate:"required,min=1,dive"`
}

// ReturnDeliveryItemsResponse defines the response for returning delivery items
type ReturnDeliveryItemsResponse struct {
	ReturnID       string  `json:"return_id"`
	DeliveryNoteID string  `json:"delivery_note_id"`
	ReturnedItems  int     `json:"returned_items"`
	TotalQuantity  float64 `json:"total_quantity"`
	ReturnDate     string  `json:"return_date"`
	ReturnStatus   string  `json:"return_status"`
	IsFullReturn   bool    `json:"is_full_return"`
}

// CancelDeliveryReturnRequest defines the request for cancelling a delivery return
type CancelDeliveryReturnRequest struct {
	ReturnID string `json:"return_id" validate:"required,uuid"`
}

// GetDeliveryReturnsResponse defines the response for fetching delivery returns
type GetDeliveryReturnsResponse struct {
	ID                 string  `json:"id"`
	DeliveryNoteID     string  `json:"delivery_note_id"`
	DeliveryNoteSerial string  `json:"delivery_note_serial"`
	SalesOrderID       string  `json:"sales_order_id"`
	SalesOrderSerial   string  `json:"sales_order_serial"`
	ReturnedItems      int     `json:"returned_items"`
	TotalQuantity      float64 `json:"total_quantity"`
	ReturnReason       string  `json:"return_reason,omitempty"`
	ReturnStatus       string  `json:"return_status"`
	ReturnedBy         string  `json:"returned_by"`
	ReturnedAt         string  `json:"returned_at"`
	CancelledAt        string  `json:"cancelled_at,omitempty"`
	CancelledBy        string  `json:"cancelled_by,omitempty"`
}

// DeliveryReturnItemResponse defines the details of returned delivery items
type DeliveryReturnItemResponse struct {
	ReturnID    string  `json:"return_id"`
	DetailID    string  `json:"detail_id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	BatchID     string  `json:"batch_id"`
	BatchSKU    string  `json:"batch_sku"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}
