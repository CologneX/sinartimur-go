package products

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

// Unit is the model for the Unit table.
type Unit struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=50"`
	Description string    `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at" validate:"rfc3339"`
	UpdatedAt   string    `json:"updated_at" validate:"rfc3339"`
	DeletedAt   *string   `json:"deleted_at,omitempty" validate:"omitempty,rfc3339"`
}

// GetProductResponse is the response payload for the GetProduct endpoint.
type GetProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Sku         string    `json:"sku"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    Category  `json:"category"`
	Unit        Unit      `json:"unit"`
	CreatedAt   string    `json:"created_at" validate:"rfc3339"`
	UpdatedAt   string    `json:"updated_at" validate:"rfc3339"`
}

// Product is the model for the Products table.
type Product struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Sku         string    `json:"sku" validate:"omitempty,max=50"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price" validate:"required,numeric,gt=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	UnitID      uuid.UUID `json:"unit_id" validate:"required"`
	CreatedAt   string    `json:"created_at" validate:"rfc3339"`
	UpdatedAt   string    `json:"updated_at" validate:"rfc3339"`
	DeletedAt   *string   `json:"deleted_at,omitempty" validate:"omitempty,rfc3339"`
}

// CreateProductRequest is the payload used for creating a new product.
type CreateProductRequest struct {
	Name        string    `json:"name" validate:"required,max=255"`
	Sku         string    `json:"sku" validate:"omitempty,max=50"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price" validate:"required,numeric,gt=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	UnitID      uuid.UUID `json:"unit_id" validate:"required"`
}

// UpdateProductRequest is the payload used for updating an existing product.
type UpdateProductRequest struct {
	Name        string    `json:"name" validate:"required,max=255"`
	Sku         string    `json:"sku" validate:"omitempty,max=50"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price" validate:"required,numeric,gt=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	UnitID      uuid.UUID `json:"unit_id" validate:"required"`
}

// CreateCategoryRequest is the payload for creating a new category.
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description,omitempty"`
}

// UpdateCategoryRequest is the payload for updating an existing category.
type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description,omitempty"`
}

// CreateUnitRequest is the payload for creating a new unit.
type CreateUnitRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description,omitempty"`
}

// UpdateUnitRequest is the payload for updating an existing unit.
type UpdateUnitRequest struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description,omitempty"`
}
