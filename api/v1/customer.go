package v1

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/customer"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
)

// CreateCustomerHandler handles the creation of a new customer
func CreateCustomerHandler(customerService *customer.CustomerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req customer.CreateCustomerRequest

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		apiErr := customerService.CreateCustomer(req)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, utils.WriteMessage("Pelanggan berhasil dibuat"))
	}
}

// GetAllCustomersHandler fetches all customers with pagination
func GetAllCustomersHandler(customerService *customer.CustomerService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		var req customer.GetCustomerRequest
		req.Name = r.URL.Query().Get("name")
		req.Address = r.URL.Query().Get("address")
		req.Telephone = r.URL.Query().Get("telephone")
		req.Page = page
		req.PageSize = pageSize
		req.SortBy = sortBy
		req.SortOrder = sortOrder

		// Validate request
		validationErrors := utils.ValidateStruct(req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		customers, totalItems, apiErr := customerService.GetAllCustomers(req)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WritePaginationJSON(w, http.StatusOK, req.Page, totalItems, req.PageSize, customers)
	})
}

// GetCustomerByIDHandler fetches a customer by ID
func GetCustomerByIDHandler(customerService *customer.CustomerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		customer, apiErr := customerService.GetCustomerByID(id)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, customer)
	}
}

// GetCustomerByNameHandler fetches a customer by name
func GetCustomerByNameHandler(customerService *customer.CustomerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name := params["name"]

		customer, apiErr := customerService.GetCustomerByName(name)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, customer)
	}
}

// UpdateCustomerHandler updates a customer
func UpdateCustomerHandler(customerService *customer.CustomerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from parameter
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"id": "ID pelanggan tidak valid",
			}))
			return
		}

		var req customer.UpdateCustomerRequest
		req.ID = id

		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		apiErr := customerService.UpdateCustomer(req)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Pelanggan berhasil diperbarui"))
	}
}

// DeleteCustomerHandler soft deletes a customer
func DeleteCustomerHandler(customerService *customer.CustomerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from parameter
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"id": "ID pelanggan tidak valid",
			}))
			return
		}

		req := customer.DeleteCustomerRequest{
			ID: id,
		}

		validationErrors := utils.ValidateStruct(&req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		apiErr := customerService.DeleteCustomer(req)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Pelanggan berhasil dihapus"))
	}
}
