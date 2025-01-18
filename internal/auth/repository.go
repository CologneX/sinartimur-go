package auth

import "database/sql"

type UserRepository interface {
	Create(username, hashedPassword string) error
	GetByUsername(username string) (*User, error)
	GetRolesByID(userID string) ([]string, error)
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

// GetRolesByID fetches roles by user ID
func (r *userRepositoryImpl) GetRolesByID(userID string) ([]string, error) {
	var roles []string
	rows, err := r.db.Query("SELECT r.name FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var role string
		err := rows.Scan(&role)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
