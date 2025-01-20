package user

import (
	"database/sql"
	"github.com/lib/pq"
)

type UserRepository interface {
	Create(username, hashedPassword string) *pq.Error
	GetByUsername(username string) (*GetUserResponse, *pq.Error)
	Update(req UpdateUserRequest) *pq.Error
	GetAll(search string) ([]*GetUserResponse, *pq.Error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create creates a new user
func (r *userRepositoryImpl) Create(
	username, hashedPassword string,
) *pq.Error {
	_, err := r.db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, hashedPassword)
	if err != nil {
		return err.(*pq.Error)
	}
	return nil
}

// GetByUsername fetches a user by username
func (r *userRepositoryImpl) GetByUsername(username string) (*GetUserResponse, *pq.Error) {
	user := &GetUserResponse{}
	err := r.db.QueryRow("SELECT id, username, is_active, created_at, updated_at FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err.(*pq.Error)
	}
	return user, nil
}

// Update updates a user
func (r *userRepositoryImpl) Update(req UpdateUserRequest) *pq.Error {
	_, err := r.db.Exec("UPDATE users SET username = $1, is_active = $2 WHERE id = $3", req.Username, req.IsActive, req.ID)
	if err != nil {
		return err.(*pq.Error)
	}
	return nil
}

// GetAll fetches all users
func (r *userRepositoryImpl) GetAll(search string) ([]*GetUserResponse, *pq.Error) {
	// Get all users and its roles

	return users, nil
}
