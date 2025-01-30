package v1

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/employee"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
)

func RegisterEmployeeRoutes(router *mux.Router, employeeService *employee.EmployeeService) {
	router.HandleFunc("/employee", CreateEmployeeHandler(employeeService)).Methods("POST")
	router.HandleFunc("/employee/{id}", UpdateEmployeeHandler(employeeService)).Methods("PUT")
	router.HandleFunc("/employee/{id}", DeleteEmployeeHandler(employeeService)).Methods("DELETE")
	router.HandleFunc("/employees", GetAllEmployeesHandler(employeeService)).Methods("GET")
}

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
			//utils.WriteJSON(w, err.StatusCode, map[string]interface{}{
			//	"errors": err.Details,
			//})
			utils.ErrorJSON(w, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Employee berhasil didaftarkan",
		})
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

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Employee berhasil diupdate"})
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

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Employee berhasil dihapus"})
	}
}

// GetAllEmployeesHandler fetches all employees
func GetAllEmployeesHandler(employeeService *employee.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		employees, err := employeeService.GetAllEmployees(name)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, employees)
	}
}
