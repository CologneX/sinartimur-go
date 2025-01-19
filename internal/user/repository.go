package user

import "database/sql"

type UserRepository interface {
	Create(username, hashedPassword string) error
	GetByUsername(username string) (*User, error)
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
) error {
	_, err := r.db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

// GetByUsername fetches a user by username
func (r *userRepositoryImpl) GetByUsername(username string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow("SELECT id, username, password_hash, is_active, created_at, updated_at FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
