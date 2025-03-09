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
		apiError := purchaseOrderService.Create(req, userID)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, utils.WriteMessage("Purchase Order berhasil dibuat"))
	}
}

// UpdatePurchaseOrderHandler updates a purchase order
func UpdatePurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req purchase.UpdatePurchaseOrderRequest

		validationErrors := utils.DecodeAndValidate(r, &req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		apiError := purchaseOrderService.Update(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Purchase Order berhasil diupdate"))
	}
}

// DeletePurchaseOrderHandler deletes a purchase order
func DeletePurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		apiError := purchaseOrderService.Delete(id)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Purchase Order berhasil dihapus"))
	}
}

// ReceivePurchaseOrderHandler receives a purchase order
func ReceivePurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []purchase.ReceivedItemRequest
		validationErrors := utils.DecodeAndValidate(r, &req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		params := mux.Vars(r)
		id := params["id"]

		// Get user ID from context
		userID := r.Context().Value("user_id").(string)
		apiError := purchaseOrderService.ReceivePurchaseOrder(id, userID, req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Purchase Order berhasil diterima"))
	}
}

// CheckPurchaseOrderHandler checks a purchase order
func CheckPurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		// Get user ID from context
		userID := r.Context().Value("user_id").(string)
		apiError := purchaseOrderService.CheckPurchaseOrder(id, userID)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Purchase Order berhasil dicek"))
	}
}

// CancelPurchaseOrderHandler cancels a purchase order
func CancelPurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		// Get user ID from context
		userID := r.Context().Value("user_id").(string)
		apiError := purchaseOrderService.CancelPurchaseOrder(id, userID)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Purchase Order berhasil dibatalkan"))
	}
}

// CreatePurchaseOrderReturnHandler creates a new purchase order return
func CreatePurchaseOrderReturnHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req purchase.CreatePurchaseOrderReturnRequest

		validationErrors := utils.DecodeAndValidate(r, &req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		// Get user ID from context
		userID := r.Context().Value("user_id").(string)
		apiError := purchaseOrderService.CreateReturn(req, userID)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, utils.WriteMessage("Purchase Order Return berhasil dibuat"))
	}
}

// GetAllPurchaseOrderReturnHandler fetches all purchase order returns
func GetAllPurchaseOrderReturnHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		req := purchase.GetPurchaseOrderReturnRequest{
			FromDate: r.URL.Query().Get("from_date"),
			ToDate:   r.URL.Query().Get("to_date"),
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

		returns, totalItems, apiError := purchaseOrderService.GetAllReturns(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WritePaginationJSON(w, http.StatusOK, page, totalItems, pageSize, returns)
	})
}

// CancelPurchaseOrderReturnHandler cancels a purchase order return
func CancelPurchaseOrderReturnHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		// Get user ID from context
		userID := r.Context().Value("user_id").(string)
		apiError := purchaseOrderService.CancelReturn(id, userID)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Purchase Order Return berhasil dibatalkan"))
	}
}

// DeletePurchaseOrderItemHandler deletes a purchase order item
func DeletePurchaseOrderItemHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		apiError := purchaseOrderService.RemovePurchaseOrderItem(id)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Item Purchase Order berhasil dihapus"))
	}
}

// CreatePurchaseOrderItemHandler creates a new purchase order item
func CreatePurchaseOrderItemHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req purchase.CreatePurchaseOrderItemRequest

		validationErrors := utils.DecodeAndValidate(r, &req)
		params := mux.Vars(r)
		orderID := params["order_id"]

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		apiError := purchaseOrderService.AddPurchaseOrderItem(orderID, req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, utils.WriteMessage("Item Purchase Order berhasil dibuat"))
	}
}

// UpdatePurchaseOrderItemHandler updates a purchase order item
func UpdatePurchaseOrderItemHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req purchase.UpdatePurchaseOrderItemRequest

		validationErrors := utils.DecodeAndValidate(r, &req)

		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		apiError := purchaseOrderService.UpdatePurchaseOrderItem(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Item Purchase Order berhasil diupdate"))
	}
}

// GetAllPurchaseOrderHandler fetch all purchase orders
func GetAllPurchaseOrderHandler(purchaseOrderService *purchase.PurchaseOrderService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		req := purchase.GetPurchaseOrderRequest{
			SupplierID: r.URL.Query().Get("supplier_id"),
			Status:     r.URL.Query().Get("status"),
			FromDate:   r.URL.Query().Get("from_date"),
			ToDate:     r.URL.Query().Get("to_date"),
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
		orders, totalItems, apiError := purchaseOrderService.GetAllPurchaseOrder(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WritePaginationJSON(w, http.StatusOK, page, totalItems, pageSize, orders)
	})
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

		apiError := supplierService.CreateSupplier(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, utils.WriteMessage("Supplier berhasil dibuat"))
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

		apiError := supplierService.UpdateSupplier(req)
		if apiError != nil {
			utils.ErrorJSON(w, apiError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Supplier berhasil diupdate"))
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

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Supplier berhasil dihapus"))
	}
}
