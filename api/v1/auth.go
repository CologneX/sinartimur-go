package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/auth"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
	"time"
)

var validate *validator.Validate

func RegisterAuthRoutes(router *mux.Router, userService *auth.AuthService) {
	router.HandleFunc("/auth/login", LoginHandler(userService)).Methods("GET")
	router.HandleFunc("/auth/refresh", RefreshTokenHandler(userService)).Methods("GET")
}

// LoginHandler logs in a user
func LoginHandler(userService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.LoginUserRequest
		validationErrors := utils.DecodeAndValidate(r, &req)
		if validationErrors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, validationErrors))
			return
		}

		username := req.Username
		password := req.Password

		accessToken, refreshToken, err, roles := userService.LoginUser(username, password)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		// Set cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().Add(time.Minute * 15),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(time.Hour * 24 * 7),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		})

		// Return JSON response with username and roles
		response := map[string]interface{}{
			"username": &username,
			"roles":    &roles,
		}
		utils.WriteJSON(w, http.StatusOK, response)
	}
}

func RefreshTokenHandler(userService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{"error": "Refresh token not found"}))
			return
		}

		refreshToken := refreshTokenCookie.Value
		accessToken, serviceErr := userService.RefreshAuth(refreshToken)
		if serviceErr != nil {
			utils.ErrorJSON(w, serviceErr)
			return
		}

		// Set new access token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().Add(time.Minute * 15),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		})

		w.WriteHeader(http.StatusOK)
	}
}
