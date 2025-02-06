package auth

import "database/sql"

type AuthRepository interface {
	GetByUsername(username string) (*User, error)
	//GetRolesByID(userID string) ([]string, error)
}

type authRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}

// GetByUsername fetches a user by username
func (r *authRepositoryImpl) GetByUsername(username string) (*User, error) {
	// Query and scan user
	user := &User{}
	err := r.db.QueryRow("Select * From Users Where Username = $1 And Is_Active = True", username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.IsActive, &user.IsAdmin, &user.IsHr, &user.IsFinance, &user.IsInventory, &user.IsSales, &user.IsPurchase, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//// GetRolesByID fetches role by user ID
//func (r *authRepositoryImpl) GetRolesByID(userID string) ([]string, error) {
//	var roles []string
//	rows, err := r.db.Query("SELECT r.name FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1", userID)
//	if err != nil {
//		return nil, err
//	}
//	for rows.Next() {
//		var role string
//		err = rows.Scan(&role)
//		if err != nil {
//			return nil, err
//		}
//		roles = append(roles, role)
//	}
//	return roles, nil
//}
