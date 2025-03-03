package middleware

import (
	"context"
	"net/http"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("access_token")
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{
				"general": "Access token not found",
			}))
			return
		}

		accessToken := accessTokenCookie.Value
		claims, err := utils.GetClaims(accessToken)
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{
				"general": "Invalid access token",
			}))
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
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
