package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"sinartimur-go/config"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
	"time"
)

// AuthService is a service that handles user authentication
type AuthService struct {
	repo        AuthRepository
	redisClient *config.RedisClient
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(repo AuthRepository, redisClient *config.RedisClient) *AuthService {
	return &AuthService{repo: repo, redisClient: redisClient}
}

// LoginUser logs in a user
func (s *AuthService) LoginUser(username, password string) (string, string, *dto.APIError, []*string) {
	// Fetch user from database
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", "", &dto.APIError{
			StatusCode: http.StatusUnauthorized,
			Details: map[string]string{
				"general": "Username atau password salah",
			},
		}, nil
	}

	// Verify password
	if !utils.ComparePasswords(user.PasswordHash, password) {
		return "", "", &dto.APIError{
			StatusCode: http.StatusUnauthorized,
			Details: map[string]string{
				"general": "Username atau password salah",
			},
		}, nil
	}

	// Write use role if Is_{Role} is true
	var roles []*string
	if user.IsAdmin {
		role := "admin"
		roles = append(roles, &role)
	}
	if user.IsHr {
		role := "hr"
		roles = append(roles, &role)
	}

	if user.IsFinance {
		role := "finance"
		roles = append(roles, &role)
	}

	if user.IsInventory {
		role := "inventory"
		roles = append(roles, &role)
	}

	if user.IsSales {
		role := "sales"
		roles = append(roles, &role)
	}

	if user.IsPurchase {
		role := "purchase"
		roles = append(roles, &role)
	}

	if len(roles) == 0 {
		roles = nil
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID.String(), roles)
	if err != nil {
		return "", "", &dto.APIError{
			StatusCode: http.StatusInternalServerError,
			Details: map[string]string{
				"general": "Gagal login. Silahkan coba lagi",
			},
		}, nil
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), roles)
	if err != nil {
		return "", "", &dto.APIError{
			StatusCode: http.StatusInternalServerError,
			Details: map[string]string{
				"general": "Gagal login. Silahkan coba lagi",
			},
		}, nil
	}

	// Store refresh token in Redis
	err = s.redisClient.Set(user.ID.String(), refreshToken, time.Hour*24*7)
	if err != nil {
		return "", "", &dto.APIError{
			StatusCode: http.StatusInternalServerError,
			Details: map[string]string{
				"general": "Gagal login. Silahkan coba lagi",
			},
		}, nil
	}
	return accessToken, refreshToken, nil, roles
}

// RefreshAuth refreshes the access token
func (s *AuthService) RefreshAuth(refreshToken string) (string, *dto.APIError) {
	// Validate refresh token
	token, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return "", &dto.APIError{
			StatusCode: http.StatusUnauthorized,
			Details: map[string]string{
				"general": "Token tidak valid",
			},
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", &dto.APIError{
			StatusCode: http.StatusUnauthorized,
			Details: map[string]string{
				"general": "Token tidak valid",
			},
		}
	}

	userID := claims["user_id"].(string)
	rolesClaim, ok := claims["roles"].([]interface{})
	var roles []*string
	if ok {
		for _, role := range rolesClaim {
			roles = append(roles, role.(*string))
		}
	}

	// Check refresh token in Redis
	storedToken, err := s.redisClient.Get(userID)
	if err != nil || storedToken != refreshToken {
		return "", &dto.APIError{
			StatusCode: http.StatusUnauthorized,
			Details: map[string]string{
				"general": "Token tidak valid",
			},
		}
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(userID, roles)
	if err != nil {
		return "", &dto.APIError{
			StatusCode: http.StatusInternalServerError,
			Details: map[string]string{
				"general": "Gagal refresh token",
			},
		}
	}

	return accessToken, nil
}
