package middleware

import (
	"context"
	"net/http"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
	"strings"
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
			// Extract user ID from claims and add to context
			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{
					"general": "Invalid token claims",
				}))
				return
			}

			// Create new context with user_id and pass to next handler
			ctx := context.WithValue(r.Context(), "user_id", userID)

			rolesClaim, ok := claims["roles"].([]interface{})
			for _, role := range rolesClaim {
				roleStr := role.(string)
				if roleStr == "admin" {
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				for _, requiredRole := range roles {
					if roleStr == requiredRole {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}

			if !ok {
				// utils.WriteJSON(w, http.StatusUnauthorized, utils.ErrorJSON(w, dto.))
				utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{
					"general": "Akses Tidak Diizinkan",
				}))
				return
			}
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{
				"general": "Akses Tidak Diizinkan",
			}))
			return
		})
	}
}

// Create a new middleware function that allows specific roles for specific paths
func RoleOrPathMiddleware(pathRoleMap map[string][]string, defaultRoles ...string) func(http.Handler) http.Handler {
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

			// Extract user ID from claims and add to context
			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{
					"general": "Invalid token claims",
				}))
				return
			}

			// Create new context with user_id and pass to next handler
			ctx := context.WithValue(r.Context(), "user_id", userID)

			// Get user roles from claims
			rolesClaim, ok := claims["roles"].([]interface{})
			if !ok {
				utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{
					"general": "Akses Tidak Diizinkan",
				}))
				return
			}

			path := r.URL.Path

			// Admin always has access
			for _, role := range rolesClaim {
				roleStr := role.(string)
				if roleStr == "admin" {
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			// Check if this path has specific role requirements
			allowedRoles := defaultRoles
			for pathPrefix, roles := range pathRoleMap {
				if strings.HasPrefix(path, pathPrefix) {
					allowedRoles = append(allowedRoles, roles...)
					break
				}
			}

			// Check if user has any of the required roles
			for _, role := range rolesClaim {
				roleStr := role.(string)
				for _, requiredRole := range allowedRoles {
					if roleStr == requiredRole {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}

			utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{
				"general": "Akses Tidak Diizinkan",
			}))
			return
		})
	}
}
