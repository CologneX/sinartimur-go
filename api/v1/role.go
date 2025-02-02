package v1

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/role"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
)

func RegisterRoleRoutes(router *mux.Router, roleService *role.RoleService) {
	router.HandleFunc("/role", CreateRoleHandler(roleService)).Methods("POST")
	router.HandleFunc("/role/{id}", UpdateRoleHandler(roleService)).Methods("PUT")
	router.HandleFunc("/roles", GetAllRolesHandler(roleService)).Methods("GET")
	router.HandleFunc("/role/assign", AssignRoleToUserHandler(roleService)).Methods("POST")
	router.HandleFunc("/role/unassign", UnassignRoleFromUserHandler(roleService)).Methods("POST")
}

// CreateRoleHandler creates a new role
func CreateRoleHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req role.CreateRoleRequest
		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}
		err := roleService.CreateRole(req)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Role berhasil didaftarkan"})
	}
}

// UpdateRoleHandler updates a role
func UpdateRoleHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req role.UpdateRoleRequest
		// Get role ID from search query
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

		serviceErr := roleService.UpdateRole(req)
		if serviceErr != nil {
			utils.ErrorJSON(w, serviceErr)
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
			utils.ErrorJSON(w, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, roles)
	}
}

// AssignRoleToUserHandler assigns a role to a user
func AssignRoleToUserHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req role.AssignRoleRequest
		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		err := roleService.AssignRoleToUser(req)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Role berhasil ditambahkan ke user"})
	}
}

// UnassignRoleFromUserHandler unassigns a role from a user
func UnassignRoleFromUserHandler(roleService *role.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req role.UnassignRoleRequest
		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		err := roleService.UnassignRoleFromUser(req)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Role berhasil dihapus dari user"})
	}
}
