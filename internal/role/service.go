package role

import (
	"sinartimur-go/pkg/dto"
)

// RoleService is a service that provides role operations
type RoleService struct {
	repo RoleRepository
}

// NewRoleService creates a new instance of RoleService
func NewRoleService(repo RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(request CreateRoleRequest) *dto.APIError {
	// Check if role exists
	_, err := s.repo.GetByName(request.Name)
	if err == nil {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"name": "Role sudah ada",
			},
		}
	}

	err = s.repo.Create(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return nil
}

// UpdateRole updates a role
func (s *RoleService) UpdateRole(request UpdateRoleRequest) *dto.APIError {
	// Check if role exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Role tidak ditemukan",
			},
		}
	}

	// Check if role with the same name already exists
	_, err = s.repo.GetByName(request.Name)
	if err == nil {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"name": "Role sudah ada",
			},
		}
	}
	err = s.repo.Update(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return nil
}

// GetAllRoles fetches all roles
func (s *RoleService) GetAllRoles(name string) ([]GetAllRoleRequest, *dto.APIError) {
	roles, err := s.repo.GetAll(name)
	if err != nil {
		return nil,
			&dto.APIError{
				StatusCode: 500,
				Details: map[string]string{
					"general": "Kesalahan Server",
				},
			}
	}
	return roles, nil
}

// GetRoleByID fetches a role by ID
func (s *RoleService) GetRoleByID(id string) (*GetRoleRequest, *dto.APIError) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil,
			&dto.APIError{
				StatusCode: 404,
				Details: map[string]string{
					"general": "Role tidak ditemukan",
				},
			}
	}
	return role, nil
}

// AssignRoleToUser assigns a role to a user
func (s *RoleService) AssignRoleToUser(request AssignRoleRequest) *dto.APIError {
	// Check if user exists
	_, err := s.repo.GetUserByID(request.UserID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"user": "User tidak ditemukan",
			},
		}
	}
	// Check if role exists
	_, err = s.repo.GetByID(request.RoleID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"role": "Role tidak ditemukan",
			},
		}
	}
	// Check if user already has the role
	role, err := s.repo.GetRoleByUserIDAndRoleID(request.UserID.String(), request.RoleID.String())
	if role != nil {
		return &dto.APIError{
			StatusCode: 409,
			Details: map[string]string{
				"role": "Role sudah diassign ke user",
			},
		}
	}
	err = s.repo.AddRoleToUser(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return nil
}

// UnassignRoleFromUser unassigns a role from a user
func (s *RoleService) UnassignRoleFromUser(request UnassignRoleRequest) *dto.APIError {
	// Check if user-role exists
	_, err := s.repo.GetUserRoleByID(request.ID.String())
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "User Role tidak ditemukan",
			},
		}
	}
	err = s.repo.RemoveRoleFromUser(request)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}
	return nil
}
