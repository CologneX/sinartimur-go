package user

import (
	"database/sql"
	"fmt"
	"sinartimur-go/utils"
)

type UserRepository interface {
	Create(req CreateUserRequest) error
	GetByUsername(username string) (*GetUserResponse, error)
	GetByID(id string) (*GetUserResponse, error)
	Update(req UpdateUserRequest) error
	GetAll(req GetAllUserRequest) ([]*GetUserResponse, int, error)
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
	_, err := r.db.Exec("Update Appuser Set Username = $1, Is_Admin = $2, Is_Hr = $3, Is_Finance = $4, Is_Inventory = $5, Is_Sales = $6, Is_Purchase = $7, Is_Active = $8, Updated_At = Now() Where Id = $9",
		req.Username, req.IsAdmin, req.IsHr, req.IsFinance, req.IsInventory, req.IsSales, req.IsPurchase, req.IsActive, req.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetAll fetches all users with its roles
func (r *userRepositoryImpl) GetAll(req GetAllUserRequest) ([]*GetUserResponse, int, error) {
	// Build the base query
	queryBuilder := utils.NewQueryBuilder(`
		Select Id, Username, Is_Active, Created_At, Updated_At,
		       Is_Admin, Is_Hr, Is_Finance, Is_Inventory, Is_Sales, Is_Purchase
		From Appuser
		Where 1=1
	`)

	// Add search filter if provided
	if req.Search != "" {
		queryBuilder.AddFilter("Username ILIKE ", "%"+req.Search+"%")
	}

	// Build count query to get total items
	countQuery, countParams := queryBuilder.Build()
	countQuery = fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", countQuery)

	// Execute count query
	var totalItems int
	err := r.db.QueryRow(countQuery, countParams...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total pengguna: %w", err)
	}

	// Add sorting if provided
	if req.SortBy != "" {
		direction := "ASC"
		if req.SortOrder == "desc" {
			direction = "DESC"
		}
		queryBuilder.Query.WriteString(fmt.Sprintf(" ORDER BY %s %s", req.SortBy, direction))
	}

	// Add pagination
	queryBuilder.AddPagination(req.PageSize, req.Page)

	// Execute final query
	query, params := queryBuilder.Build()
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data pengguna: %w", err)
	}
	defer rows.Close()

	var users []*GetUserResponse
	for rows.Next() {
		user := &GetUserResponse{}
		var isAdmin, isHr, isFinance, isInventory, isSales, isPurchase bool
		err = rows.Scan(
			&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
			&isAdmin, &isHr, &isFinance, &isInventory, &isSales, &isPurchase,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal membaca data pengguna: %w", err)
		}

		// Construct roles array
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
			roles = append(roles, "customer")
		}
		if isPurchase {
			roles = append(roles, "purchase")
		}
		user.Role = &roles

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("terjadi kesalahan saat membaca data pengguna: %w", err)
	}

	return users, totalItems, nil
}

//func (r *userRepositoryImpl) GetAll(search string) ([]*GetUserResponse, error) {
//	query := `
//		Select Id, Username, Is_Active, Created_At, Updated_At,
//		       Is_Admin, Is_Hr, Is_Finance, Is_Inventory, Is_Sales, Is_Purchase
//		From Appuser
//		Where Username Ilike '%' || $1 || '%'
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
//		var isAdmin, isHr, isFinance, isInventory, isSales, isPurchase bool
//		err = rows.Scan(&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
//			&isAdmin, &isHr, &isFinance, &isInventory, &isSales, &isPurchase)
//		if err != nil {
//			return nil, err
//		}
//
//		var roles []string
//		if isAdmin {
//			roles = append(roles, "admin")
//		}
//		if isHr {
//			roles = append(roles, "hr")
//		}
//		if isFinance {
//			roles = append(roles, "finance")
//		}
//		if isInventory {
//			roles = append(roles, "inventory")
//		}
//		if isSales {
//			roles = append(roles, "customer")
//		}
//		if isPurchase {
//			roles = append(roles, "purchase")
//		}
//		user.Role = &roles
//
//		users = append(users, user)
//	}
//
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//
//	return users, nil
//}

// UpdateCredential updates user's password
func (r *userRepositoryImpl) UpdateCredential(req UpdateUserCredentialRequest) error {
	_, err := r.db.Exec("Update Appuser Set Password_Hash = $1, Updated_At = Now() Where Id = $2", req.Password, req.ID)
	if err != nil {
		return err
	}
	return nil
}
