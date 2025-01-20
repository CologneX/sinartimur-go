package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"sinartimur-go/config"
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
func (s *AuthService) LoginUser(username, password string) (int, string, string, error, []string) {
	// Fetch user from database
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return http.StatusUnauthorized, "", "", errors.New("username atau password salah"), nil
	}

	// Verify password
	if !utils.ComparePasswords(user.PasswordHash, password) {
		return http.StatusUnauthorized, "", "", errors.New("username atau password salah"), nil
	}

	// Get user roles
	roles, err := s.repo.GetRolesByID(user.ID.String())
	if err != nil {
		return http.StatusInternalServerError, "", "", errors.New("Gagal mengambil role user"), nil
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID.String())
	if err != nil {
		return http.StatusInternalServerError, "", "", errors.New("Gagal login. Silahkan coba lagi"), nil
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return http.StatusInternalServerError, "", "", errors.New("Gagal login. Silahkan coba lagi"), nil
	}

	// Store refresh token in Redis
	err = s.redisClient.Set(user.ID.String(), refreshToken, time.Minute*1)
	if err != nil {
		return http.StatusInternalServerError, "", "", fmt.Errorf("Failed to store refresh token: %w", err), nil
	}

	return http.StatusOK, accessToken, refreshToken, nil, roles
}

// RefreshAuth refreshes the access token
func (s *AuthService) RefreshAuth(refreshToken string) (int, string, error) {
	// Validate refresh token
	token, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return http.StatusUnauthorized, "", errors.New("Invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return http.StatusUnauthorized, "", errors.New("Invalid refresh token")
	}

	userID := claims["user_id"].(string)

	// Check refresh token in Redis
	storedToken, err := s.redisClient.Get(userID)
	if err != nil || storedToken != refreshToken {
		return http.StatusUnauthorized, "", errors.New("Invalid refresh token")
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("Failed to generate access token: %w", err)
	}

	return http.StatusOK, accessToken, nil
}
