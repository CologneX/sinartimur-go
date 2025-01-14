package user

import "fmt"

type UserService struct {
	repo UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(request RegisterUserRequest) error {
	// Business logic for registering a user
	fmt.Println(request)
	return nil
}
