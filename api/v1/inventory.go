package v1

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/inventory"
	"sinartimur-go/utils"
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
