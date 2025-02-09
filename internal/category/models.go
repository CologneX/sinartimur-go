package category

import "github.com/google/uuid"

// Category is the model for the Category table.
type Category struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description string    `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at" validate:"rfc3339"`
	UpdatedAt   string    `json:"updated_at" validate:"rfc3339"`
	DeletedAt   *string   `json:"deleted_at,omitempty" validate:"omitempty,rfc3339"`
}

// CreateCategoryRequest is the payload for creating a new category.
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description,omitempty"`
}

// UpdateCategoryRequest is the payload for updating an existing category.
type UpdateCategoryRequest struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description string    `json:"description,omitempty"`
}

// GetCategoryResponse is the response payload for the GetCategory endpoint.
type GetCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at" validate:"rfc3339"`
	UpdatedAt   string    `json:"updated_at" validate:"rfc3339"`
}

// GetCategoryRequest is the payload for fetching categories.
type GetCategoryRequest struct {
	Name string `json:"name,omitempty"`
}

// DeleteCategoryRequest is the payload for deleting a category.
type DeleteCategoryRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}
