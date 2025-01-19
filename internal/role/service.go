package role

import (
	"errors"
	"fmt"
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
func (s *RoleService) CreateRole(request CreateRoleRequest) error {
	err := s.repo.Create(request)
	if err != nil {
		return err
	}
	return nil
}

// UpdateRole updates a role
func (s *RoleService) UpdateRole(request UpdateRoleRequest) error {
	// Check if role exists
	fmt.Println(request.ID.String())
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return errors.New("Role tidak ditemukan")
	}
	err = s.repo.Update(request)
	if err != nil {
		return err
	}
	return nil
}

// GetAllRoles fetches all roles
func (s *RoleService) GetAllRoles(name string) ([]GetAllRoleRequest, error) {
	roles, err := s.repo.GetAll(name)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByID fetches a role by ID
func (s *RoleService) GetRoleByID(id string) (*GetRoleRequest, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// AssignRoleToUser assigns a role to a user
func (s *RoleService) AssignRoleToUser(request AssignRoleRequest) error {
	// Check if role exists
	_, err := s.repo.GetByID(request.RoleID.String())
	if err != nil {
		return errors.New("Role tidak ditemukan")
	}
	err = s.repo.AddRoleToUser(request)
	if err != nil {
		return err
	}
	return nil
}

// UnassignRoleFromUser unassigns a role from a user
func (s *RoleService) UnassignRoleFromUser(request UnassignRoleRequest) error {
	// Check if role exists
	_, err := s.repo.GetByID(request.ID.String())
	if err != nil {
		return errors.New("Role tidak ditemukan")
	}
	err = s.repo.RemoveRoleFromUser(request)
	if err != nil {
		return err
	}
	return nil
}
