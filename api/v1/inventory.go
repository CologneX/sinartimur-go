package v1

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/inventory"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
	"time"
)

// GetAllStoragesHandler returns all storage locations with pagination
func GetAllStoragesHandler(storageService *inventory.StorageService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		var req inventory.GetStorageRequest

		// Extract query parameters for filtering
		req.Name = r.URL.Query().Get("name")
		req.Location = r.URL.Query().Get("location")

		// Set pagination parameters
		req.Page = page
		req.PageSize = pageSize
		req.SortBy = sortBy
		req.SortOrder = sortOrder

		// Validate
		if err := utils.ValidateStruct(req); err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, err))
			return
		}

		// Get storages with pagination
		storages, totalItems, apiErr := storageService.GetAllStorages(req)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WritePaginationJSON(w, http.StatusOK, page, totalItems, pageSize, storages)
	})
}

// GetStorageByIDHandler returns a single storage location by ID
func GetStorageByIDHandler(storageService *inventory.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		if err := utils.ValidateStruct(id); err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, err))
			return
		}

		storage, apiErr := storageService.GetStorageByID(id)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, storage)
	}
}

// CreateStorageHandler creates a new storage location
func CreateStorageHandler(storageService *inventory.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req inventory.CreateStorageRequest

		if err := utils.DecodeAndValidate(r, &req); err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, err))
			return
		}

		_, apiErr := storageService.CreateStorage(req)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, utils.WriteMessage("Gudang berhasil dibuat"))
	}
}

// UpdateStorageHandler updates an existing storage location
func UpdateStorageHandler(storageService *inventory.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var req inventory.UpdateStorageRequest
		if err := utils.DecodeAndValidate(r, &req); err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, err))
			return
		}

		// Ensure ID from URL matches body or set it if not provided
		idUUID, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}
		req.ID = idUUID

		_, apiErr := storageService.UpdateStorage(req)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Gudang berhasil diperbarui"))
	}
}

// DeleteStorageHandler deletes a storage location
func DeleteStorageHandler(storageService *inventory.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		if err := utils.ValidateStruct(id); err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, err))
			return
		}

		apiErr := storageService.DeleteStorage(id)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Gudang berhasil dihapus"))
	}
}

// MoveBatchHandler moves a product batch between storage locations
func MoveBatchHandler(storageService *inventory.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req inventory.MoveBatchRequest
		fmt.Println(&req)

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		// Get userID from context (provided by auth middleware)
		userID := r.Context().Value("user_id").(string)

		apiErr := storageService.MoveBatch(req, userID)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Batch produk berhasil dipindahkan"))
	}
}

// GetAllInventoryLogHandler handles requests to get inventory logs
func GetAllInventoryLogHandler(storageService *inventory.StorageService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		var req inventory.GetInventoryLogsRequest

		// Extract query parameters for filtering
		req.BatchID = r.URL.Query().Get("batch_id")
		req.ProductID = r.URL.Query().Get("product_id")
		req.StorageID = r.URL.Query().Get("storage_id")
		req.TargetStorageID = r.URL.Query().Get("target_storage_id")
		req.UserID = r.URL.Query().Get("user_id")
		req.Action = r.URL.Query().Get("action")
		req.FromDate = r.URL.Query().Get("from_date")
		req.ToDate = r.URL.Query().Get("to_date")

		// Set pagination parameters
		req.Page = page
		req.PageSize = pageSize
		req.SortBy = sortBy
		req.SortOrder = sortOrder

		// Validate
		if err := utils.ValidateStruct(req); err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, err))
			return
		}

		// Get last refresh timestamp
		lastRefreshed, apiErr := storageService.GetInventoryLogLastRefreshed()
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		// Get inventory logs with pagination
		logs, totalItems, apiErr := storageService.GetInventoryLogs(req)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		// Create response with logs and last refresh info
		response := struct {
			Data          []inventory.GetInventoryLogResponse `json:"data"`
			LastRefreshed *time.Time                          `json:"last_refreshed"`
			Page          int                                 `json:"page"`
			TotalItems    int                                 `json:"total_items"`
			PageSize      int                                 `json:"page_size"`
		}{
			Data:          logs,
			LastRefreshed: lastRefreshed,
			Page:          page,
			TotalItems:    totalItems,
			PageSize:      pageSize,
		}

		utils.WriteJSON(w, http.StatusOK, response)
	})
}

// RefreshInventoryLogViewHandler handles requests to refresh the inventory log materialized view
func RefreshInventoryLogViewHandler(storageService *inventory.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiError := storageService.RefreshInventoryLogView()
		if apiError != nil {
			return
		}
	}
}
