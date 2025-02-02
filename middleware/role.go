package middleware

import (
	"net/http"
	"sinartimur-go/utils"
)

//func HRRoleMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		accessTokenCookie, err := r.Cookie("access_token")
//		if err != nil {
//			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Access token not found"})
//			return
//		}
//
//		accessToken := accessTokenCookie.Value
//		claims, err := utils.GetClaims(accessToken)
//		if err != nil {
//			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid access token"})
//			return
//		}
//
//		rolesClaim, ok := claims["roles"].([]interface{})
//		for _, role := range rolesClaim {
//			roleStr := role.(string)
//			if roleStr == "admin" {
//				next.ServeHTTP(w, r)
//				return
//			}
//			if roleStr == "hr" {
//				next.ServeHTTP(w, r)
//				return
//			}
//		}
//
//		if !ok {
//			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid access token"})
//			return
//		}
//
//		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Akses Tidak Diizinkan"})
//	})
//}

// RoleMiddleware is a middleware that checks if the user has the required role from params
func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessTokenCookie, err := r.Cookie("access_token")
			if err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Access token not found"})
				return
			}

			accessToken := accessTokenCookie.Value
			claims, err := utils.GetClaims(accessToken)
			if err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid access token"})
				return
			}

			rolesClaim, ok := claims["roles"].([]interface{})
			for _, role := range rolesClaim {
				roleStr := role.(string)
				if roleStr == "admin" {
					next.ServeHTTP(w, r)
					return
				}
				for _, requiredRole := range roles {
					if roleStr == requiredRole {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			if !ok {
				utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Akses Tidak Diizinkan"})
				return
			}

			utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Akses Tidak Diizinkan"})
		})
	}
}
