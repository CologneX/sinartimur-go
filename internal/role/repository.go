package role

import (
	"database/sql"
	"sinartimur-go/internal/user"
)

type RoleRepository interface {
	//Create(request CreateRoleRequest) error
	//Delete(request DeleteEmployeeRequest) error
	//Update(request UpdateRoleRequest) error
	//GetAll(name string) ([]GetAllRoleRequest, error)
	//GetByID(id string) (*GetRoleRequest, error)
	//GetByName(name string) (*GetRoleRequest, error)
	//AddRoleToUser(req AssignRoleRequest) error
	//RemoveRoleFromUser(req UnassignRoleRequest) error
	//GetRoleByUserIDAndRoleID(userID, roleID string) (*GetRoleRequest, error)
	//GetUserRoleByID(id string) (*GetRoleRequest, error)
	GetUserByID(id string) (*user.GetUserResponse, error)
}

type roleRepositoryImpl struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepositoryImpl{db: db}
}

//// Create creates a new role
//func (r *roleRepositoryImpl) Create(request CreateRoleRequest) error {
//	_, err := r.db.Exec("INSERT INTO roles (name, description) VALUES ($1, $2)", request.Name, request.Description)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// Update updates a role
//func (r *roleRepositoryImpl) Update(request UpdateRoleRequest) error {
//	_, err := r.db.Exec("UPDATE roles SET name = $1, description = $2, updated_at = NOW() WHERE id = $3", request.Name, request.Description, request.ID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// GetAll fetches all roles
//func (r *roleRepositoryImpl) GetAll(name string) ([]GetAllRoleRequest, error) {
//	var roles []GetAllRoleRequest
//	var rows *sql.Rows
//	var err error
//
//	if name != "" {
//		rows, err = r.db.Query("SELECT id, name, description, created_at, updated_at FROM roles WHERE name ILIKE $1", "%"+name+"%")
//	} else {
//		rows, err = r.db.Query("SELECT id, name, description, created_at, updated_at FROM roles")
//	}
//
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		var role GetAllRoleRequest
//		err = rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
//		if err != nil {
//			return nil, err
//		}
//		roles = append(roles, role)
//	}
//
//	return roles, nil
//}

//// GetByName fetches a role by name
//func (r *roleRepositoryImpl) GetByName(name string) (*GetRoleRequest, error) {
//	role := &GetRoleRequest{}
//	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1", name).Scan(
//		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
//	if err != nil {
//		return nil, err
//	}
//	return role, nil
//}

//// GetByID fetches a role by ID
//func (r *roleRepositoryImpl) GetByID(id string) (*GetRoleRequest, error) {
//	role := &GetRoleRequest{}
//	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1", id).Scan(
//		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
//	if err != nil {
//		return nil, err
//	}
//	return role, nil
//}

//// AddRoleToUser assigns a role to a user
//func (r *roleRepositoryImpl) AddRoleToUser(req AssignRoleRequest) error {
//	_, err := r.db.Exec("INSERT INTO user_roles (user_id, role_id, assigned_at) VALUES ($1, $2, NOW())", req.UserID, req.RoleID)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//// RemoveRoleFromUser unassigns a role from a user
//func (r *roleRepositoryImpl) RemoveRoleFromUser(req UnassignRoleRequest) error {
//	_, err := r.db.Exec("DELETE FROM user_roles WHERE id = $1", req.ID)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//// GetRoleByUserIDAndRoleID fetches a role by user ID and role ID
//func (r *roleRepositoryImpl) GetRoleByUserIDAndRoleID(userID, roleID string) (*GetRoleRequest, error) {
//	role := &GetRoleRequest{}
//	err := r.db.QueryRow("SELECT r.id, r.name, r.description, r.created_at, r.updated_at FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1 AND ur.role_id = $2", userID, roleID).Scan(
//		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
//	if err != nil {
//		return nil, err
//	}
//	return role, nil
//}

//// GetUserRoleByID fetches a user's role by ID
//func (r *roleRepositoryImpl) GetUserRoleByID(id string) (*GetRoleRequest, error) {
//	role := &GetRoleRequest{}
//	err := r.db.QueryRow("SELECT r.id, r.name, r.description, r.created_at, r.updated_at FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.id = $1", id).Scan(
//		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
//	if err != nil {
//		return nil, err
//	}
//	return role, nil
//}

// GetUserByID fetches a user by ID
func (r *roleRepositoryImpl) GetUserByID(id string) (*user.GetUserResponse, error) {
	u := &user.GetUserResponse{}
	err := r.db.QueryRow("Select Id, Username, Is_Active, Created_At, Updated_At From Appuser Where Id = $1", id).Scan(
		&u.ID, &u.Username, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
