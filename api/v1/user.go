package v1

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/user"
	"sinartimur-go/utils"
)

var validate *validator.Validate

func RegisterUserRoutes(router *mux.Router, userService *user.UserService) {
	router.HandleFunc("/users", RegisterUserHandler(userService)).Methods("POST")
}

func RegisterUserHandler(userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req user.RegisterUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		validate = validator.New()
		if err := validate.Struct(req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		err := userService.RegisterUser(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User berhasil didaftarkan"})
	}
}
