package v1

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/employee"
	"sinartimur-go/utils"
	"time"
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}
		// Check every field in the request manually
		if req.Name == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Nama tidak boleh kosong"})
			return
		}
		if req.Position == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Posisi Karyawan tidak boleh kosong"})
			return
		}

		if req.HiredDate == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Tanggal harus diisi"})
			return
		}

		if _, err := time.Parse(time.RFC3339, req.HiredDate); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format tanggal salah"})
			return
		}

		err := employeeService.CreateEmployee(req)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Employee berhasil didaftarkan"})
	}
}

// UpdateEmployeeHandler updates an employee
func UpdateEmployeeHandler(employeeService *employee.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req employee.UpdateEmployeeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		// Check every field in the request manually
		if req.Name == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Nama tidak boleh kosong"})
			return
		}
		if req.Position == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Posisi Karyawan tidak boleh kosong"})
			return
		}
		if req.HiredDate == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Tanggal harus diisi"})
			return
		}
		if _, err := time.Parse(time.RFC3339, req.HiredDate); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format tanggal salah"})
			return
		}

		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
			return
		}
		req.ID = id

		err = employeeService.UpdateEmployee(req)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
			return
		}

		err = employeeService.DeleteEmployee(employee.DeleteEmployeeRequest{ID: id})
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, employees)
	}
}
