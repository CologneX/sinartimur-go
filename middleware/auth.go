package middleware

import (
	"net/http"
	"sinartimur-go/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("access_token")
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Access token not found"})
			return
		}

		accessToken := accessTokenCookie.Value
		_, err = utils.ValidateToken(accessToken)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid access token"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
