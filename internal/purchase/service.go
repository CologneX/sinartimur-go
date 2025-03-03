package purchase

import "sinartimur-go/pkg/dto"

// SupplierService is the service for the Supplier domain
type SupplierService struct {
	repo SupplierRepository
}

// NewSupplierService creates a new instance of SupplierService
func NewSupplierService(repo SupplierRepository) *SupplierService {
	return &SupplierService{repo: repo}
}

// GetAllSuppliers fetches all suppliers with pagination
func (s *SupplierService) GetAllSuppliers(req GetSupplierRequest) ([]GetSupplierResponse, int, *dto.APIError) {
	suppliers, totalItems, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return suppliers, totalItems, nil
}

// GetSupplierByID fetches a supplier by ID
func (s *SupplierService) GetSupplierByID(id string) (*GetSupplierResponse, *dto.APIError) {
	supplier, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Supplier tidak ditemukan",
			},
		}
	}

	return supplier, nil
}

// CreateSupplier creates a new supplier
func (s *SupplierService) CreateSupplier(req CreateSupplierRequest) (*GetSupplierResponse, *dto.APIError) {
	// Check if supplier with same name already exists
	existing, err := s.repo.GetByName(req.Name)
	if err == nil && existing != nil {
		return nil, &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"name": "Supplier dengan nama ini sudah terdaftar",
			},
		}
	}

	supplier, err := s.repo.Create(req)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return supplier, nil
}

// UpdateSupplier updates an existing supplier
func (s *SupplierService) UpdateSupplier(req UpdateSupplierRequest) (*GetSupplierResponse, *dto.APIError) {
	// Check if supplier exists
	_, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Supplier tidak ditemukan",
			},
		}
	}

	// If name is changing, check if new name is already taken
	if req.Name != "" {
		existing, err := s.repo.GetByName(req.Name)
		if err == nil && existing != nil && existing.ID != req.ID {
			return nil, &dto.APIError{
				StatusCode: 400,
				Details: map[string]string{
					"name": "Supplier dengan nama ini sudah terdaftar",
				},
			}
		}
	}

	supplier, err := s.repo.Update(req)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return supplier, nil
}

// DeleteSupplier deletes a supplier
func (s *SupplierService) DeleteSupplier(id string) *dto.APIError {
	err := s.repo.Delete(id)
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}

// PurchaseOrderService is the service for the Purchase Order domain
type PurchaseOrderService struct {
	repo PurchaseOrderRepository
}

// NewPurchaseOrderService creates a new instance of PurchaseOrderService
func NewPurchaseOrderService(repo PurchaseOrderRepository) *PurchaseOrderService {
	return &PurchaseOrderService{repo: repo}
}

// GetAllPurchaseOrders fetches all purchase orders with filtering and pagination
func (s *PurchaseOrderService) GetAllPurchaseOrders(req GetPurchaseOrderRequest) ([]GetPurchaseOrderResponse, int, *dto.APIError) {
	orders, totalItems, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return orders, totalItems, nil
}

// GetPurchaseOrderByID fetches a purchase order by ID
func (s *PurchaseOrderService) GetPurchaseOrderByID(id string) (*GetPurchaseOrderResponse, *dto.APIError) {
	order, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"id": err.Error(),
			},
		}
	}

	return order, nil
}

// GetPurchaseOrderDetailByID fetches a purchase order with its items by ID
func (s *PurchaseOrderService) GetPurchaseOrderDetailByID(id string) (*PurchaseOrderDetailResponse, *dto.APIError) {
	orderDetail, err := s.repo.GetDetailByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"id": err.Error(),
			},
		}
	}

	return orderDetail, nil
}

// CreatePurchaseOrder creates a new purchase order with its items
// Update the service's CreatePurchaseOrder method
func (s *PurchaseOrderService) CreatePurchaseOrder(req CreatePurchaseOrderRequest, userID string) (*GetPurchaseOrderResponse, *dto.APIError) {
	// Create the purchase order
	order, err := s.repo.Create(req, userID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return order, nil
}

// UpdatePurchaseOrder updates an existing purchase order (only date as specified)
func (s *PurchaseOrderService) UpdatePurchaseOrder(req UpdatePurchaseOrderRequest) (*GetPurchaseOrderResponse, *dto.APIError) {
	order, err := s.repo.Update(req)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return order, nil
}

// ReceivePurchaseOrder receives a purchase order
func (s *PurchaseOrderService) ReceivePurchaseOrder(id string) *dto.APIError {
	order, err := s.repo.GetByID(id)
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"id": err.Error(),
			},
		}
	}

	if order.Status == "received" {
		return &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"general": "Purchase Order sudah diterima",
			},
		}
	}

	order, err = s.repo.Update(
		UpdatePurchaseOrderRequest{
			ID:     id,
			Status: "received",
		})

	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Kesalahan Server",
			},
		}
	}

	return nil
}

// CancelPurchaseOrder cancels a purchase order
func (s *PurchaseOrderService) CancelPurchaseOrder(id string) *dto.APIError {
	err := s.repo.Cancel(id)
	if err != nil {
		return &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"id": err.Error(),
			},
		}
	}

	return nil
}

// PurchaseOrderDetailService is the service for the Purchase Order Detail domain
type PurchaseOrderDetailService struct {
	repo PurchaseOrderDetailRepository
}

// NewPurchaseOrderDetailService creates a new instance of PurchaseOrderDetailService
func NewPurchaseOrderDetailService(repo PurchaseOrderDetailRepository) *PurchaseOrderDetailService {
	return &PurchaseOrderDetailService{repo: repo}
}

// GetDetailsByOrderID fetches all details for a purchase order with pagination
func (s *PurchaseOrderDetailService) GetDetailsByOrderID(orderID string, page, pageSize int) ([]GetPurchaseOrderItemResponse, int, *dto.APIError) {
	items, totalItems, err := s.repo.GetAllByOrderID(orderID, page, pageSize)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return items, totalItems, nil
}

// GetDetailByID fetches a single purchase order detail by ID
func (s *PurchaseOrderDetailService) GetDetailByID(id string) (*GetPurchaseOrderItemResponse, *dto.APIError) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"id": err.Error(),
			},
		}
	}

	return item, nil
}

// CreateDetail adds a new item to an existing purchase order
func (s *PurchaseOrderDetailService) CreateDetail(orderID string, req CreatePurchaseOrderItemRequest) (*GetPurchaseOrderItemResponse, *dto.APIError) {
	if req.ProductID != "" {
		productExists, err := s.repo.CheckProductByID(req.ProductID)
		if err != nil {
			return nil, &dto.APIError{
				StatusCode: 500,
				Details: map[string]string{
					"general": err.Error(),
				},
			}
		}

		if !productExists {
			return nil, &dto.APIError{
				StatusCode: 400,
				Details: map[string]string{
					"product_id": "Product tidak ditemukan atau sudah dihapus",
				},
			}
		}
	}

	item, err := s.repo.Create(orderID, req)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return item, nil
}

// UpdateDetail modifies an existing purchase order detail
func (s *PurchaseOrderDetailService) UpdateDetail(req UpdatePurchaseOrderItemRequest) (*GetPurchaseOrderItemResponse, *dto.APIError) {

	item, err := s.repo.Update(req)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return item, nil
}

// DeleteDetail removes a detail from a purchase order
func (s *PurchaseOrderDetailService) DeleteDetail(id string) *dto.APIError {
	err := s.repo.Delete(id)
	if err != nil {
		return &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}
