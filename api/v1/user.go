package v1

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/user"
	"sinartimur-go/utils"
)

func RegisterUserRoutes(router *mux.Router, userService *user.UserService) {
	router.HandleFunc("/admin/user", CreateUserHandler(userService)).Methods("POST")
	router.HandleFunc("/admin/users", GetAllUsersHandler(userService)).Methods("GET")
	router.HandleFunc("/admin/user/{id}", UpdateUserHandler(userService)).Methods("PUT")
}

func CreateUserHandler(userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req user.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		// Validate request
		if req.Username == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Username tidak boleh kosong"})
			return
		}

		if req.Password == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Password tidak boleh kosong"})
			return
		}

		if req.Password != req.ConfirmPassword {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Password harus sama"})
			return
		}

		httpCode, err := userService.CreateUser(req)
		if err != nil {
			utils.WriteJSON(w, httpCode, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User berhasil didaftarkan"})
	}
}

func GetAllUsersHandler(userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		users, err := userService.GetAllUsers(search)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, users)
	}
}

func UpdateUserHandler(userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req user.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		// Get ID from URL
		var err error
		req.ID, err = uuid.Parse(mux.Vars(r)["id"])
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
			return
		}

		// Validate request
		if req.ID.String() == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ID tidak boleh kosong"})
			return
		}

		if req.Username == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Username tidak boleh kosong"})
			return
		}

		httpCode, err := userService.Update(req)
		if err != nil {
			utils.WriteJSON(w, httpCode, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User berhasil diupdate"})
	}
}
