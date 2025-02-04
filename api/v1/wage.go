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
	return func(w http.ResponseWriter, r *http.Request) {
		employeeID := r.URL.Query().Get("employee_id")
		// Transform year into int
		year, err := strconv.Atoi(r.URL.Query().Get("year"))
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"year": "Tahun tidak valid",
			}))
			return
		}

		// Transform month into int
		month, err := strconv.Atoi(r.URL.Query().Get("month"))
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"month": "Bulan tidak valid",
			}))
			return
		}

		wages, errService := wageService.GetAllWages(employeeID, year, month)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, wages)
	}
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
