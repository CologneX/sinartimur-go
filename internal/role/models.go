package role

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type UserRole struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	RoleID     uuid.UUID `json:"role_id"`
	AssignedAt string    `json:"assigned_at"`
}

type GetAllRoleRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
type GetRoleRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateRoleRequest struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
}

type DeleteRoleRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

type AssignRoleRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
	RoleID uuid.UUID `json:"role_id" validate:"required,uuid"`
}

type UnassignRoleRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}
