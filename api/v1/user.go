package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/user"
	"sinartimur-go/utils"
)

func RegisterUserRoutes(router *mux.Router, userService *user.UserService) {
	router.HandleFunc("/admin/user/create", CreateUserHandler(userService)).Methods("POST")
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
