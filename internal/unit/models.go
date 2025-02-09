package unit

import "github.com/google/uuid"

// Unit is the model for the Unit table.
type Unit struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=50"`
	Description string    `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at" validate:"rfc3339"`
	UpdatedAt   string    `json:"updated_at" validate:"rfc3339"`
	DeletedAt   *string   `json:"deleted_at,omitempty" validate:"omitempty,rfc3339"`
}

// GetUnitResponse is the response payload for the GetUnit endpoint.
type GetUnitResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at" validate:"rfc3339"`
	UpdatedAt   string    `json:"updated_at" validate:"rfc3339"`
}

// GetUnitRequest is the payload for fetching units.
type GetUnitRequest struct {
	Name string `json:"name,omitempty"`
}

// CreateUnitRequest is the payload for creating a new unit.
type CreateUnitRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description,omitempty"`
}

// UpdateUnitRequest is the payload for updating an existing unit.
type UpdateUnitRequest struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required,max=50"`
	Description string    `json:"description,omitempty"`
}

type DeleteUnitRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}
