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

	// Hash password
	request.Password = utils.HashPassword(request.Password)

	// Insert user to database
	err = s.repo.Create(request)
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
	user, err := s.repo.GetByID(request.ID.String())
	if err != nil || user.ID != request.ID {
		return &dto.APIError{
			StatusCode: http.StatusNotFound,
			Details: map[string]string{
				"username": "User tidak ditemukan",
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

// UpdateCredential updates user's password
func (s *UserService) UpdateCredential(request UpdateUserCredentialRequest) *dto.APIError {
	// Check if user exists with the same username
	user, err := s.repo.GetByID(request.ID.String())
	if err != nil || user.ID != request.ID {
		return &dto.APIError{
			StatusCode: http.StatusNotFound,
			Details: map[string]string{
				"username": "User tidak ditemukan",
			},
		}
	}

	// Hash password
	request.Password = utils.HashPassword(request.Password)

	// UpdateDetail user's password in database
	errSer := s.repo.UpdateCredential(request)
	if errSer != nil {
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
