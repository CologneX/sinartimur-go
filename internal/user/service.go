package user

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"sinartimur-go/utils"
)

// UserService is a service that handles user authentication
type UserService struct {
	repo UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUserService registers a new user
func (s *UserService) CreateUser(request CreateUserRequest) (int, error) {
	// Check if password and confirm password match
	if request.Password != request.ConfirmPassword {
		return http.StatusBadRequest, errors.New("Password dan konfirmasi password tidak sama")
	}

	// Insert user to database
	err := s.repo.Create(request.Username, utils.HashPassword(request.Password))
	if err != nil {
		// Check returned error if Unique Constraint Violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return http.StatusConflict, errors.New("User sudah terdaftar")
		}
		return http.StatusInternalServerError, fmt.Errorf("Gagal membuat user: %w", err)
	}
	return http.StatusOK, nil
}
