package v1

import (
	"net/http"
	"sinartimur-go/internal/employee"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateEmployeeHandler creates a new employee
func CreateEmployeeHandler(employeeService *employee.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req employee.CreateEmployeeRequest

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		// Call the service
		err := employeeService.CreateEmployee(req)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Employee berhasil didaftarkan"))
	}
}

// UpdateEmployeeHandler updates an employee
func UpdateEmployeeHandler(employeeService *employee.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req employee.UpdateEmployeeRequest
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

		errService := employeeService.UpdateEmployee(req)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Employee berhasil diperbaharui"))
	}
}

// DeleteEmployeeHandler soft deletes an employee
func DeleteEmployeeHandler(employeeService *employee.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"general": "ID tidak valid",
			}))
			return
		}

		errService := employeeService.DeleteEmployee(employee.DeleteEmployeeRequest{ID: id})
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Employee berhasil dihapus"))
	}
}

// GetAllEmployeesHandler fetches all employees
func GetAllEmployeesHandler(employeeService *employee.EmployeeService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		var req employee.GetAllEmployeeRequest
		req.Name = r.URL.Query().Get("name")
		req.Page = page
		req.PageSize = pageSize
		req.SortBy = sortBy
		req.SortOrder = sortOrder

		// validate struct
		err := utils.ValidateStruct(req)
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, err))
			return
		}

		employees, totalItems, errService := employeeService.GetAllEmployees(req)
		if err != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WritePaginationJSON(w, http.StatusOK, page, totalItems, pageSize, employees)
	})
}
