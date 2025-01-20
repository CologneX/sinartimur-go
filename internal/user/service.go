package user

import (
	"errors"
	"fmt"
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

// CreateUser registers a new user
func (s *UserService) CreateUser(request CreateUserRequest) (int, error) {
	// Check if password and confirm password match
	if request.Password != request.ConfirmPassword {
		return http.StatusBadRequest, errors.New("Password dan konfirmasi password tidak sama")
	}

	// Insert user to database
	err := s.repo.Create(request.Username, utils.HashPassword(request.Password))
	if err != nil {
		// Check returned error if Unique Constraint Violation
		if err.Code == "23505" {
			return http.StatusConflict, errors.New("Username sudah terdaftar")
		}
		return http.StatusInternalServerError, fmt.Errorf("Gagal membuat user: %w", err)
	}
	return http.StatusOK, nil
}

// Update updates a user
func (s *UserService) Update(request UpdateUserRequest) (int, error) {
	// Update user in database
	err := s.repo.Update(request)
	if err != nil {
		if err.Code == "23505" {
			return http.StatusConflict, errors.New("Username sudah terdaftar")
		}
		return http.StatusInternalServerError, fmt.Errorf("Gagal mengupdate user: %w", err)
	}
	return http.StatusOK, nil
}

// GetAllUsers fetches all users
func (s *UserService) GetAllUsers(search string) ([]*GetUserResponse, error) {
	users, err := s.repo.GetAll(search)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengambil data user: %w", err)
	}
	return users, nil
}
