package v1

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/sales"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
	"strconv"
)

// GetSalesOrdersHandler retrieves a list of sales orders with pagination and filtering
func GetSalesOrdersHandler(salesService *sales.SalesService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		// Extract filter parameters
		req := sales.GetSalesOrdersRequest{
			CustomerID:    r.URL.Query().Get("customer_id"),
			Status:        r.URL.Query().Get("status"),
			PaymentMethod: r.URL.Query().Get("payment_method"),
			StartDate:     r.URL.Query().Get("start_date"),
			EndDate:       r.URL.Query().Get("end_date"),
			SerialID:      r.URL.Query().Get("serial_id"),
			PaginationParameter: utils.PaginationParameter{
				Page:      page,
				PageSize:  pageSize,
				SortBy:    sortBy,
				SortOrder: sortOrder,
			},
		}

		// Validate filter parameters if provided
		if errors := utils.ValidateStruct(req); errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to get data
		orders, totalCount, err := salesService.GetSalesOrders(req)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusInternalServerError,
				Details: map[string]string{
					"general": "Gagal mengambil daftar pesanan: " + err.Error(),
				},
			})
			return
		}

		// Return paginated response
		utils.WritePaginationJSON(w, http.StatusOK, page, totalCount, pageSize, orders)
	})
}

// CreateSalesOrderHandler handles creation of a new sales order
func CreateSalesOrderHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Parse and validate request body
		var req sales.CreateSalesOrderRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to create order
		response, err := salesService.CreateSalesOrder(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal membuat pesanan: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusCreated, response)
	}
}

// UpdateSalesOrderHandler handles updating basic sales order information
func UpdateSalesOrderHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req sales.UpdateSalesOrderRequest

		// Get order ID from URL path parameters
		vars := mux.Vars(r)
		orderID := vars["id"]
		req.ID = orderID

		// Parse and validate request body
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to update order
		response, err := salesService.UpdateSalesOrder(req)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal memperbarui pesanan: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusOK, response)
	}
}

// AddItemToSalesOrderHandler handles adding a new item to an existing sales order
func AddItemToSalesOrderHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Parse and validate request body
		var req sales.AddItemToSalesOrderRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to add item
		response, err := salesService.AddItemToSalesOrder(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal menambahkan item: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusCreated, response)
	}
}

// UpdateSalesOrderItemHandler handles updating an existing item in a sales order
func UpdateSalesOrderItemHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Parse and validate request body
		var req sales.UpdateItemRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to update item
		response, err := salesService.UpdateSalesOrderItem(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal memperbarui item: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusOK, response)
	}
}

// DeleteSalesOrderItemHandler handles removing an item from a sales order
func DeleteSalesOrderItemHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Extract order ID and detail ID from URL path parameters
		vars := mux.Vars(r)
		salesOrderID := vars["order_id"]
		detailID := vars["detail_id"]

		if salesOrderID == "" || detailID == "" {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "ID pesanan dan ID detail harus diisi",
				},
			})
			return
		}

		req := sales.DeleteItemRequest{
			SalesOrderID: salesOrderID,
			DetailID:     detailID,
		}

		// Call service to delete item
		err := salesService.DeleteSalesOrderItem(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal menghapus item: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "Item berhasil dihapus dari pesanan",
		})
	}
}

// GetSalesOrderDetailsHandler retrieves detailed information about a specific sales order
func GetSalesOrderDetailsHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract order ID from URL path parameters
		vars := mux.Vars(r)
		orderID := vars["id"]

		if orderID == "" {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "ID pesanan harus diisi",
				},
			})
			return
		}

		// Call service to get order details
		details, err := salesService.GetSalesOrderDetail(orderID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal mengambil detail pesanan: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusOK, details)
	}
}

// CancelSalesOrderHandler handles cancelling a sales order
func CancelSalesOrderHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Parse and validate request body
		var req sales.CancelSalesOrderRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to cancel order
		err := salesService.CancelSalesOrder(req, userID)
		if err != nil {

			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal membatalkan pesanan: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "Pesanan berhasil dibatalkan",
		})
	}
}

// GetSalesInvoicesHandler retrieves a paginated list of sales invoices with filtering
func GetSalesInvoicesHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract pagination parameters
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = utils.DefaultPage
		}

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		if pageSize == 0 {
			pageSize = utils.DefaultPageSize
		}

		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := r.URL.Query().Get("sort_order")

		// Extract filter parameters
		req := sales.GetSalesInvoicesRequest{
			CustomerID: r.URL.Query().Get("customer_id"),
			Status:     r.URL.Query().Get("status"),
			StartDate:  r.URL.Query().Get("start_date"),
			EndDate:    r.URL.Query().Get("end_date"),
			SerialID:   r.URL.Query().Get("serial_id"),
			PaginationParameter: utils.PaginationParameter{
				Page:      page,
				PageSize:  pageSize,
				SortBy:    sortBy,
				SortOrder: sortOrder,
			},
		}

		// Validate filter parameters if provided
		if errors := utils.ValidateStruct(req); errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to get data
		invoices, totalCount, err := salesService.GetSalesInvoices(req)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal mengambil daftar faktur: " + err.Error(),
				},
			})
			return
		}

		// Return paginated response
		utils.WritePaginationJSON(w, http.StatusOK, page, totalCount, pageSize, invoices)
	}
}

// CreateSalesInvoiceHandler handles creation of a new sales invoice from an order
func CreateSalesInvoiceHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Parse and validate request body
		var req sales.CreateSalesInvoiceRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to create invoice
		response, err := salesService.CreateSalesInvoice(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal membuat faktur: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusCreated, response)
	}
}

// CancelSalesInvoiceHandler handles cancellation of a sales invoice
func CancelSalesInvoiceHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Parse and validate request body
		var req sales.CancelSalesInvoiceRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to cancel invoice
		err := salesService.CancelSalesInvoice(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal membatalkan faktur: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "Faktur berhasil dibatalkan",
		})
	}
}

// ReturnInvoiceItemsHandler handles processing returns for invoice items
func ReturnInvoiceItemsHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Parse and validate request body
		var req sales.ReturnInvoiceItemsRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to process returns
		response, err := salesService.ReturnInvoiceItems(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal memproses pengembalian: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusCreated, response)
	}
}

// CancelInvoiceReturnHandler handles cancellation of a previously processed return
func CancelInvoiceReturnHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Extract return ID from URL path parameters
		vars := mux.Vars(r)
		returnID := vars["return_id"]

		if returnID == "" {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "ID pengembalian harus diisi",
				},
			})
			return
		}

		req := sales.CancelInvoiceReturnRequest{
			ReturnID: returnID,
		}

		// Call service to cancel return
		err := salesService.CancelInvoiceReturn(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal membatalkan pengembalian: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "Pengembalian berhasil dibatalkan",
		})
	}
}

// CreateDeliveryNoteHandler handles creation of a new delivery note for a sales invoice
func CreateDeliveryNoteHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Parse and validate request body
		var req sales.CreateDeliveryNoteRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details:    errors,
			})
			return
		}

		// Call service to create delivery note
		response, err := salesService.CreateDeliveryNote(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal membuat surat jalan: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusCreated, response)
	}
}

// CancelDeliveryNoteHandler handles cancellation of a delivery note
func CancelDeliveryNoteHandler(salesService *sales.SalesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from token
		claims := r.Context().Value("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		// Extract delivery note ID from URL path parameters
		vars := mux.Vars(r)
		deliveryNoteID := vars["delivery_note_id"]

		if deliveryNoteID == "" {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "ID surat jalan harus diisi",
				},
			})
			return
		}

		req := sales.CancelDeliveryNoteRequest{
			DeliveryNoteID: deliveryNoteID,
		}

		// Call service to cancel delivery note
		err := salesService.CancelDeliveryNote(req, userID)
		if err != nil {
			utils.ErrorJSON(w, &dto.APIError{
				StatusCode: http.StatusBadRequest,
				Details: map[string]string{
					"general": "Gagal membatalkan surat jalan: " + err.Error(),
				},
			})
			return
		}

		// Return success response
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "Surat jalan berhasil dibatalkan",
		})
	}
}
