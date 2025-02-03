package v1

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/user"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
)

func CreateUserHandler(userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req user.CreateUserRequest
		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		err := userService.CreateUser(req)
		if err != nil {
			utils.ErrorJSON(w, err)
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

		errService := userService.Update(req)
		if errService != nil {
			utils.ErrorJSON(w, errService)
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "User berhasil diupdate"})
	}
}
