package v1

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/purchase"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
)

// CreatePurchaseOrderHandler creates a new purchase order
func CreatePurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req purchase.CreatePurchaseOrderRequest

		validationErrors := utils.DecodeAndValidate(r, &req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		// Get user ID from context
		userID := r.Context().Value("user_id").(string)
		order, apiError := purchaseOrderService.CreatePurchaseOrder(req, userID)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, order)
	}
}

// GetAllPurchaseOrdersHandler fetches all purchase orders
func GetAllPurchaseOrdersHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		req := purchase.GetPurchaseOrderRequest{
			SupplierName: r.URL.Query().Get("supplier_name"),
			OrderDate:    r.URL.Query().Get("order_date"),
			Status:       r.URL.Query().Get("status"),
			FromDate:     r.URL.Query().Get("from_date"),
			ToDate:       r.URL.Query().Get("to_date"),
			PaginationParameter: utils.PaginationParameter{
				Page:      page,
				PageSize:  pageSize,
				SortBy:    sortBy,
				SortOrder: sortOrder,
			},
		}

		validationErrors := utils.ValidateStruct(req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		orders, totalItems, apiError := purchaseOrderService.GetAllPurchaseOrders(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WritePaginationJSON(w, http.StatusOK, req.Page, totalItems, req.PageSize, orders)
	})
}

// GetPurchaseOrderByIDHandler GetPurchaseOrderDetailByIDHandler fetches a purchase order with its items by ID
func GetPurchaseOrderByIDHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}

		order, apiError := purchaseOrderService.GetPurchaseOrderByID(id.String())
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, order)
	}
}

// GetPurchaseOrderDetailHandler CancelPurchaseOrderHandler cancels a purchase order
func GetPurchaseOrderDetailHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}

		orderDetail, apiError := purchaseOrderService.GetPurchaseOrderDetailByID(id.String())
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, orderDetail)
	}
}

// UpdatePurchaseOrderHandler CancelPurchaseOrderHandler cancels a purchase order
func UpdatePurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}

		var req purchase.UpdatePurchaseOrderRequest

		req.ID = id.String()
		validationErrors := utils.DecodeAndValidate(r, &req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		order, apiError := purchaseOrderService.UpdatePurchaseOrder(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, order)
	}
}

// CancelPurchaseOrderHandler cancels a purchase order
func CancelPurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}

		apiError := purchaseOrderService.CancelPurchaseOrder(id.String())
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Purchase order berhasil dibatalkan!"})
	}
}

// ReceivePurchaseOrderHandler receives a purchase order
func ReceivePurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}

		apiError := purchaseOrderService.ReceivePurchaseOrder(id.String())
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Purchase order berhasil diterima!"})
	}
}

// AddPurchaseOrderDetailHandler Purchase Order Detail Handlers
func AddPurchaseOrderDetailHandler(purchaseOrderDetailService *purchase.PurchaseOrderDetailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		orderID, err := uuid.Parse(params["order_id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID order tidak valid",
			}))
			return
		}

		var req purchase.CreatePurchaseOrderItemRequest

		validationErrors := utils.DecodeAndValidate(r, &req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		detail, apiError := purchaseOrderDetailService.CreateDetail(orderID.String(), req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, detail)
	}
}

// UpdatePurchaseOrderDetailHandler Purchase Order Detail Handlers
func UpdatePurchaseOrderDetailHandler(purchaseOrderDetailService *purchase.PurchaseOrderDetailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}

		var req purchase.UpdatePurchaseOrderItemRequest

		req.ID = id.String()
		validationErrors := utils.DecodeAndValidate(r, &req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		detail, apiError := purchaseOrderDetailService.UpdateDetail(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, detail)
	}
}

// DeletePurchaseOrderDetailHandler Purchase Order Detail Handlers
func DeletePurchaseOrderDetailHandler(purchaseOrderDetailService *purchase.PurchaseOrderDetailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}

		apiError := purchaseOrderDetailService.DeleteDetail(id.String())
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Item purchase order berhasil dihapus!"})
	}
}

// CreateSupplierHandler creates a new supplier
func CreateSupplierHandler(supplierService *purchase.SupplierService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req purchase.CreateSupplierRequest
		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		supplier, apiError := supplierService.CreateSupplier(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, supplier)
	}
}

// UpdateSupplierHandler updates a supplier
func UpdateSupplierHandler(supplierService *purchase.SupplierService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		var req purchase.UpdateSupplierRequest

		req.ID = id
		validationErrors := utils.DecodeAndValidate(r, &req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		supplier, apiError := supplierService.UpdateSupplier(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, supplier)
	}
}

// GetSupplierByIDHandler fetches a supplier by ID
func GetSupplierByIDHandler(supplierService *purchase.SupplierService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}

		supplier, apiError := supplierService.GetSupplierByID(id.String())
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, supplier)
	}
}

// GetAllSuppliersHandler fetches all suppliers
func GetAllSuppliersHandler(supplierService *purchase.SupplierService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		var req purchase.GetSupplierRequest
		req.Name, req.Telephone = r.URL.Query().Get("name"), r.URL.Query().Get("telephone")
		req.Page = page
		req.PageSize = pageSize
		req.SortBy = sortBy
		req.SortOrder = sortOrder
		// Validate req
		validationErrors := utils.ValidateStruct(req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		suppliers, totalItems, apiError := supplierService.GetAllSuppliers(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WritePaginationJSON(w, http.StatusOK, req.Page, totalItems, req.PageSize, suppliers)
	})
}

// DeleteSupplierHandler deletes a supplier
func DeleteSupplierHandler(supplierService *purchase.SupplierService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from parameter
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}
		var req purchase.DeleteSupplierRequest
		req.ID = id.String()

		apiError := supplierService.DeleteSupplier(req.ID)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Supplier berhasil dihapus"})
	}
}
