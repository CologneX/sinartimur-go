package v1

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/wage"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
	"strconv"
)

// CreateWageHandler creates a new employee
func CreateWageHandler(wageService *wage.WageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req wage.CreateWageRequest

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		// Call the service
		err := wageService.CreateWage(req)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Gaji berhasil ditambahkan",
		})
	}
}

// UpdateWageHandler updates an employee
func UpdateWageHandler(wageService *wage.WageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req wage.UpdateWageDetailRequest
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}
		req.ID = id

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		errService := wageService.UpdateWage(req)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Gaji berhasil diupdate"})
	}
}

// GetAllWagesHandler fetches all wages
func GetAllWagesHandler(wageService *wage.WageService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		var req wage.GetWageRequest
		req.EmployeeId = r.URL.Query().Get("employee_id")
		req.Year, _ = strconv.Atoi(r.URL.Query().Get("year"))
		req.Month, _ = strconv.Atoi(r.URL.Query().Get("month"))
		req.Page = page
		req.PageSize = pageSize
		req.SortBy = sortBy
		req.SortOrder = sortOrder

		validationErrors := utils.ValidateStruct(req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		wages, totalItems, errService := wageService.GetAllWages(req)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WritePaginationJSON(w, http.StatusOK, req.Page, totalItems, req.PageSize, wages)
	})
}

// GetWageDetailHandler fetches wage details
func GetWageDetailHandler(wageService *wage.WageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		wageDetails, errService := wageService.GetWageDetail(id)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, wageDetails)
	}
}

// DeleteWageHandler deletes a wage
func DeleteWageHandler(wageService *wage.WageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req wage.DeleteWageRequest
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}
		req.ID = id

		errService := wageService.DeleteWage(req)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Gaji berhasil dihapus"})
	}
}
