package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/role"
	"sinartimur-go/utils"
)

func RegisterRoleRoutes(router *mux.Router, roleService *role.RoleService) {
	router.HandleFunc("/role", CreateRoleHandler(roleService)).Methods("POST")
	router.HandleFunc("/role/{id}", UpdateRoleHandler(roleService)).Methods("PUT")
	router.HandleFunc("/roles", GetAllRolesHandler(roleService)).Methods("GET")
}

// CreateRoleHandler creates a new role
func CreateRoleHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req role.CreateRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
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

		err := roleService.UpdateRole(req)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
