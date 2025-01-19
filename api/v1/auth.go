package v1

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"sinartimur-go/internal/auth"
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		validate = validator.New()
		if err := validate.Struct(req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Data tidak valid"})
			return
		}

		username := req.Username
		password := req.Password

		status, accessToken, refreshToken, err, roles := userService.LoginUser(username, password)
		if err != nil {
			utils.WriteJSON(w, status, map[string]string{"error": err.Error()})
			return
		}

		// Set cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().Add(time.Minute * 1),
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
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Refresh token not found"})
			return
		}

		refreshToken := refreshTokenCookie.Value
		status, accessToken, err := userService.RefreshAuth(refreshToken)
		if err != nil {
			utils.WriteJSON(w, status, map[string]string{"error": err.Error()})
			return
		}

		// Set new access token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Expires:  time.Now().Add(time.Minute * 1),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		})

		w.WriteHeader(http.StatusOK)
	}
}
