package v1

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/role"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
)

func RegisterRoleRoutes(router *mux.Router, roleService *role.RoleService) {
	router.HandleFunc("/admin/role", CreateRoleHandler(roleService)).Methods("POST")
	router.HandleFunc("/admin/role/{id}", UpdateRoleHandler(roleService)).Methods("PUT")
	router.HandleFunc("/admin/roles", GetAllRolesHandler(roleService)).Methods("GET")
	router.HandleFunc("/admin/role/assign", AssignRoleToUserHandler(roleService)).Methods("POST")
	router.HandleFunc("/admin/role/unassign", UnassignRoleFromUserHandler(roleService)).Methods("POST")
}

// CreateRoleHandler creates a new role
func CreateRoleHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req role.CreateRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}
		// Validate request
		if req.Name == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Nama role tidak boleh kosong"})
			return
		}

		err := roleService.CreateRole(req)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Role berhasil didaftarkan"})
	}
}

// UpdateRoleHandler updates a role
func UpdateRoleHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req role.UpdateRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		// Validate request
		if req.Name == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Nama role tidak boleh kosong"})
			return
		}

		// Get role ID from search query
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
			return
		}
		req.ID = id

		err = roleService.UpdateRole(req)
		var apiErr *dto.APIError
		if errors.As(err, &apiErr) {
			utils.WriteJSON(w, apiErr.StatusCode, map[string]string{"error": apiErr.Message})
			return
		} else if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Role berhasil diupdate"})
	}
}

// GetAllRolesHandler fetches all roles
func GetAllRolesHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get search query
		name := r.URL.Query().Get("name")
		roles, err := roleService.GetAllRoles(name)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusOK, roles)
	}
}

// AssignRoleToUserHandler assigns a role to a user
func AssignRoleToUserHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req role.AssignRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		// Validate request
		if req.RoleID.String() == "" || req.UserID.String() == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak boleh kosong"})
			return
		}

		err := roleService.AssignRoleToUser(req)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Role berhasil ditambahkan ke user"})
	}
}

// UnassignRoleFromUserHandler unassigns a role from a user
func UnassignRoleFromUserHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req role.UnassignRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		// Validate request
		if req.ID.String() == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak boleh kosong"})
			return
		}

		err := roleService.UnassignRoleFromUser(req)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Role berhasil dihapus dari user"})
	}
}
