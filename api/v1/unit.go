package v1

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/unit"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
)

// GetAllUnitHandler fetches all units
func GetAllUnitHandler(unitService *unit.UnitService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req unit.GetUnitRequest
		req.Name = r.URL.Query().Get("name")
		// Validate req
		validationErrors := utils.ValidateStruct(&req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}
		units, err := unitService.GetAllUnit(req)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, units)
	}
}

// CreateUnitHandler creates a new unit
func CreateUnitHandler(unitService *unit.UnitService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req unit.CreateUnitRequest

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		createUnit, errSer := unitService.CreateUnit(req)
		if errSer != nil {
			utils.ErrorJSON(w, errSer)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": fmt.Sprintf("Unit %s berhasil didaftarkan", createUnit.Name),
		})
	}
}

// UpdateUnitHandler updates a unit
func UpdateUnitHandler(unitService *unit.UnitService) http.HandlerFunc {
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
		var req unit.UpdateUnitRequest
		req.ID = id

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		updateUnit, errService := unitService.UpdateUnit(req)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Unit berhasil diupdate ke %s", updateUnit.Name)})
	}
}

// DeleteUnitHandler soft deletes a unit
func DeleteUnitHandler(unitService *unit.UnitService) http.HandlerFunc {
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
		var req unit.DeleteUnitRequest
		req.ID = id

		validationErrors := utils.ValidateStruct(&req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		errService := unitService.DeleteUnit(req)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Unit berhasil dihapus"})
	}
}
