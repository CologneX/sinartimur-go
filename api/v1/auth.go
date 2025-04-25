package v1

import (
	"net/http"
	"net/url"
	"sinartimur-go/internal/auth"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
	"time"
)

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

		accessToken, refreshToken, userID, err, roles := userService.LoginUser(username, password)
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
		response := auth.LoginUserResponse{
			Id:       userID,
			Username: username,
			Roles:    roles,
		}

		// convert response to json string
		responseJSON, err := utils.ToJSON(response)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
		responseJSON = url.QueryEscape(responseJSON)

		http.SetCookie(w, &http.Cookie{
			Name:     "user",
			Value:    responseJSON,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		})

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

// LogoutHandler revokes the refresh token in Redis and clears all cookies.
func LogoutHandler(userService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fetch the refresh token cookie
		refreshCookie, err := r.Cookie("refresh_token")
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{"error": "Refresh token missing"}))
			return
		}

		// revoke it in Redis
		if err := userService.Logout(refreshCookie.Value); err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		// clear all cookies
		names := []string{"access_token", "refresh_token", "user"}
		for _, name := range names {
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: name != "user",
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
		}

		w.WriteHeader(http.StatusOK)
	}
}
