package user

import (
	"database/sql"
	"fmt"
)

type UserRepository interface {
	Create(req CreateUserRequest) error
	GetByUsername(username string) (*GetUserResponse, error)
	GetByID(id string) (*GetUserResponse, error)
	Update(req UpdateUserRequest) error
	GetAll(search string) ([]*GetUserResponse, error)
	UpdateCredential(req UpdateUserCredentialRequest) error
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create creates a new user
func (r *userRepositoryImpl) Create(req CreateUserRequest) error {
	_, err := r.db.Exec("Insert Into Appuser (Username, Password_Hash, Is_Admin, Is_Hr, Is_Finance, Is_Inventory, Is_Sales, Is_Purchase) Values ($1, $2, $3, $4, $5, $6, $7, $8)",
		req.Username, req.Password, req.IsAdmin, req.IsHr, req.IsFinance, req.IsInventory, req.IsSales, req.IsPurchase)
	if err != nil {
		return err
	}
	return nil
}

// GetByUsername fetches a user by username
func (r *userRepositoryImpl) GetByUsername(username string) (*GetUserResponse, error) {
	user := &GetUserResponse{}
	err := r.db.QueryRow("Select Id, Username, Is_Active, Created_At, Updated_At From Appuser Where Username = $1 And Is_Active = True", username).Scan(
		&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID fetches a user by ID
func (r *userRepositoryImpl) GetByID(id string) (*GetUserResponse, error) {
	user := &GetUserResponse{}
	err := r.db.QueryRow("Select Id, Username, Is_Active, Created_At, Updated_At From Appuser Where Id = $1 And Is_Active = True", id).Scan(
		&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates a user
func (r *userRepositoryImpl) Update(req UpdateUserRequest) error {
	fmt.Println(req)
	_, err := r.db.Exec("Update Appuser Set Username = $1, Is_Admin = $2, Is_Hr = $3, Is_Finance = $4, Is_Inventory = $5, Is_Sales = $6, Is_Purchase = $7, Is_Active = $8, Updated_At = Now() Where Id = $9",
		req.Username, req.IsAdmin, req.IsHr, req.IsFinance, req.IsInventory, req.IsSales, req.IsPurchase, req.IsActive, req.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAll fetches all users with its roles
func (r *userRepositoryImpl) GetAll(search string) ([]*GetUserResponse, error) {
	query := `
		Select Id, Username, Is_Active, Created_At, Updated_At,
		       Is_Admin, Is_Hr, Is_Finance, Is_Inventory, Is_Sales, Is_Purchase
		From Appuser
		Where Username Ilike '%' || $1 || '%'
	`
	rows, err := r.db.Query(query, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*GetUserResponse
	for rows.Next() {
		user := &GetUserResponse{}
		var isAdmin, isHr, isFinance, isInventory, isSales, isPurchase bool
		err = rows.Scan(&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
			&isAdmin, &isHr, &isFinance, &isInventory, &isSales, &isPurchase)
		if err != nil {
			return nil, err
		}

		var roles []string
		if isAdmin {
			roles = append(roles, "admin")
		}
		if isHr {
			roles = append(roles, "hr")
		}
		if isFinance {
			roles = append(roles, "finance")
		}
		if isInventory {
			roles = append(roles, "inventory")
		}
		if isSales {
			roles = append(roles, "sales")
		}
		if isPurchase {
			roles = append(roles, "purchase")
		}
		user.Role = &roles

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateCredential updates user's password
func (r *userRepositoryImpl) UpdateCredential(req UpdateUserCredentialRequest) error {
	_, err := r.db.Exec("Update Appuser Set Password_Hash = $1, Updated_At = Now() Where Id = $2", req.Password, req.ID)
	if err != nil {
		return err
	}
	return nil
}

//func (r *userRepositoryImpl) GetAll(search string) ([]*GetUserResponse, error) {
//	query := `
//		SELECT u.id, u.username, u.is_active, u.created_at, u.updated_at,
//		       COALESCE(json_agg(json_build_object('id', ur.id, 'name', r.name)) FILTER (WHERE r.name IS NOT NULL), '[]') AS roles
//		FROM users u
//		LEFT JOIN user_roles ur ON u.id = ur.user_id
//		LEFT JOIN roles r ON ur.role_id = r.id
//		WHERE u.username ILIKE '%' || $1 || '%'
//		GROUP BY u.id
//	`
//	rows, err := r.db.Query(query, search)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var users []*GetUserResponse
//	for rows.Next() {
//		user := &GetUserResponse{}
//		var rolesJSON []byte
//		err = rows.Scan(&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &rolesJSON)
//		if err != nil {
//			return nil, err
//		}
//		var roles []UserRole
//		err = json.Unmarshal(rolesJSON, &roles)
//		if err != nil {
//			return nil, err
//		}
//		user.Role = &roles
//		users = append(users, user)
//	}
//
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//
//	return users, nil
//}
