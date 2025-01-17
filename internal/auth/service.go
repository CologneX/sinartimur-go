package auth

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"sinartimur-go/config"
	"sinartimur-go/utils"
	"time"
)

// AuthService is a service that handles user authentication
type AuthService struct {
	repo        UserRepository
	redisClient *config.RedisClient
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(repo UserRepository, redisClient *config.RedisClient) *AuthService {
	return &AuthService{repo: repo, redisClient: redisClient}
}

// CreateUserService registers a new user
func (s *AuthService) CreateUserService(request RegisterUserRequest) (int, error) {
	// Check if password and confirm password match
	if request.Password != request.ConfirmPassword {
		return http.StatusBadRequest, errors.New("Password dan konfirmasi password tidak sama")
	}

	// Insert user to database
	err := s.repo.CreateUser(request.Username, utils.HashPassword(request.Password))
	if err != nil {
		// Check returned error if Unique Constraint Violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return http.StatusConflict, errors.New("User sudah terdaftar")
		}
		return http.StatusInternalServerError, fmt.Errorf("Gagal membuat user: %w", err)
	}
	return http.StatusOK, nil
}

// LoginUserService logs in a user
func (s *AuthService) LoginUserService(username, password string) (int, string, string, error) {
	// Fetch user from database
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return http.StatusUnauthorized, "", "", errors.New("username atau password salah")
	}

	// Verify password
	if !utils.ComparePasswords(user.PasswordHash, password) {
		return http.StatusUnauthorized, "", "", errors.New("username atau password salah")
	}
	// Generate tokens
	accessToken, err := GenerateAccessToken(user.ID.String())
	if err != nil {
		return http.StatusInternalServerError, "", "", fmt.Errorf("Failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateRefreshToken(user.ID.String())
	if err != nil {
		return http.StatusInternalServerError, "", "", fmt.Errorf("Failed to generate refresh token: %w", err)
	}

	// Store refresh token in Redis
	err = s.redisClient.Set(user.ID.String(), refreshToken, time.Hour*24*7)
	if err != nil {
		return http.StatusInternalServerError, "", "", fmt.Errorf("Failed to store refresh token: %w", err)
	}

	return http.StatusOK, accessToken, refreshToken, nil
}

// RefreshAuthService refreshes the access token
//func (s *AuthService) RefreshAuthService(refreshToken string) (int, error) {
//	// Validate refresh token
//	token, err := ValidateToken(refreshToken)
//	if err != nil {
//		return http.StatusUnauthorized, errors.New("Invalid refresh token")
//	}
//
//	claims, ok := token.Claims.(jwt.MapClaims)
//	if !ok || !token.Valid {
//		return http.StatusUnauthorized, errors.New("Invalid refresh token")
//	}
//
//	userID := claims["user_id"].(string)
//
//	// Check refresh token in Redis
//	storedToken, err := s.redisClient.Get(userID)
//	if err != nil || storedToken != refreshToken {
//		return http.StatusUnauthorized, errors.New("Invalid refresh token")
//	}
//
//	// Generate new access token
//	accessToken, err := GenerateAccessToken(userID)
//	if err != nil {
//		return http.StatusInternalServerError, fmt.Errorf("Failed to generate access token: %w", err)
//	}
//
//	// Set new access token cookie
//	http.SetCookie(w, &http.Cookie{
//		Name:     "access_token",
//		Value:    accessToken,
//		Expires:  time.Now().Add(time.Minute * 15),
//		HttpOnly: true,
//		Secure:   true,
//		SameSite: http.SameSiteStrictMode,
//	})
//
//	return http.StatusOK, nil
//}
