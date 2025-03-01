package inventory

type Storage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type Inventory struct {
	ID              string `json:"id"`
	ProductID       string `json:"product_id"`
	StorageID       string `json:"storage_id"`
	Quantity        int    `json:"quantity"`
	MinimumQuantity int    `json:"minimum_quantity"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type InventoryLog struct {
	ID          string `json:"id"`
	InventoryID string `json:"inventory_id"`
	UserID      string `json:"user_id"`
	Action      string `json:"action"`
	Quantity    int    `json:"quantity"`
	LogDate     string `json:"log_date"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GetStorageResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetStorageRequest struct {
	Name string `json:"name"`
}

type CreateStorageRequest struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type UpdateStorageRequest struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type DeleteStorageRequest struct {
	ID string `json:"id" validate:"required"`
}

type GetInventoryResponse struct {
	ID          string `json:"id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	StorageID   string `json:"storage_id"`
	StorageName string `json:"storage_name"`
	Quantity    int    `json:"quantity"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GetInventoryRequest struct {
	ProductName string `json:"product_name"`
	StorageName string `json:"storage_name"`
}

type CreateInventoryRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	StorageID string `json:"storage_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,gte=0"`
}

type UpdateInventoryRequest struct {
	ID       string `json:"id" validate:"required,uuid"`
	Quantity int    `json:"quantity" validate:"required,gte=0"`
}

//type GetInventoryLogResponse struct {
//	ID            string `json:"id"`
//	InventoryID   string `json:"inventory_id"`
//	InventoryName string `json:"inventory_name"`
//	UserID        string `json:"user_id"`
//}
