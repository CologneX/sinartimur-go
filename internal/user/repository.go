package user

import (
	"database/sql"
	"encoding/json"
)

type UserRepository interface {
	Create(username, hashedPassword string) error
	GetByUsername(username string) (*GetUserResponse, error)
	Update(req UpdateUserRequest) error
	GetAll(search string) ([]*GetUserResponse, error)
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
func (r *userRepositoryImpl) GetByUsername(username string) (*GetUserResponse, error) {
	user := &GetUserResponse{}
	err := r.db.QueryRow("SELECT id, username, is_active, created_at, updated_at FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates a user
func (r *userRepositoryImpl) Update(req UpdateUserRequest) error {
	_, err := r.db.Exec("UPDATE users SET username = $1, is_active = $2, updated_at = now() WHERE id = $3", req.Username, req.IsActive, req.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAll fetches all users
func (r *userRepositoryImpl) GetAll(search string) ([]*GetUserResponse, error) {
	query := `
		SELECT u.id, u.username, u.is_active, u.created_at, u.updated_at,
		       COALESCE(json_agg(json_build_object('id', ur.id, 'name', r.name)) FILTER (WHERE r.name IS NOT NULL), '[]') AS roles
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.username ILIKE '%' || $1 || '%'
		GROUP BY u.id
	`
	rows, err := r.db.Query(query, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*GetUserResponse
	for rows.Next() {
		user := &GetUserResponse{}
		var rolesJSON []byte
		err = rows.Scan(&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &rolesJSON)
		if err != nil {
			return nil, err
		}
		var roles []UserRole
		err = json.Unmarshal(rolesJSON, &roles)
		if err != nil {
			return nil, err
		}
		user.Role = &roles
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
