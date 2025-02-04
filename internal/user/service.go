package user

import (
	"fmt"
	"net/http"
	"sinartimur-go/pkg/dto"
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
func (s *UserService) CreateUser(request CreateUserRequest) *dto.APIError {
	// Check if user exists
	_, err := s.repo.GetByUsername(request.Username)
	if err == nil {
		return &dto.APIError{
			StatusCode: http.StatusConflict,
			Details: map[string]string{
				"username": "Username sudah terdaftar",
			},
		}
	}
	// Insert user to database
	err = s.repo.Create(request.Username, utils.HashPassword(request.Password))
	if err != nil {
		return &dto.APIError{
			StatusCode: http.StatusInternalServerError,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return nil
}

// Update updates a user
func (s *UserService) Update(request UpdateUserRequest) *dto.APIError {
	// Check if user exists with the same username
	user, err := s.repo.GetByUsername(request.Username)
	if err == nil && user.ID != request.ID {
		return &dto.APIError{
			StatusCode: http.StatusConflict,
			Details: map[string]string{
				"username": "Username sudah terdaftar",
			},
		}
	}
	// UpdateDetail user in database
	err = s.repo.Update(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: http.StatusInternalServerError,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return nil
}

// GetAllUsers fetches all users
func (s *UserService) GetAllUsers(search string) ([]*GetUserResponse, error) {
	users, err := s.repo.GetAll(search)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengambil data user: %w", err)
	}
	return users, nil
}
