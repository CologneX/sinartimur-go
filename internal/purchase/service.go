package purchase

import (
	"sinartimur-go/pkg/dto"
)

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
func (s *SupplierService) CreateSupplier(req CreateSupplierRequest) *dto.APIError {
	// Check if supplier with same name already exists
	existing, err := s.repo.GetByName(req.Name)
	if err == nil && existing != nil {
		return &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"name": "Supplier dengan nama ini sudah terdaftar",
			},
		}
	}

	err = s.repo.Create(req)
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

// UpdateSupplier updates an existing supplier
func (s *SupplierService) UpdateSupplier(req UpdateSupplierRequest) *dto.APIError {
	// Check if supplier exists
	_, err := s.repo.GetByID(req.ID)
	if err != nil {
		return &dto.APIError{
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
			return &dto.APIError{
				StatusCode: 400,
				Details: map[string]string{
					"name": "Supplier dengan nama ini sudah terdaftar",
				},
			}
		}
	}

	err = s.repo.Update(req)
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

type PurchaseOrderService struct {
	repo PurchaseOrderRepository
}

// NewPurchaseOrderService creates a new instance of PurchaseOrderService
func NewPurchaseOrderService(repo PurchaseOrderRepository) *PurchaseOrderService {
	return &PurchaseOrderService{repo: repo}
}

func (s *PurchaseOrderService) GetAllPurchaseOrder(req GetPurchaseOrderRequest) ([]GetPurchaseOrderResponse, int, *dto.APIError) {
	// get the request parameters from url search params
	purchaseOrders, totalItems, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	return purchaseOrders, totalItems, nil
}

func (s *PurchaseOrderService) Create(req CreatePurchaseOrderRequest, userID string) *dto.APIError {
	err := s.repo.Create(req, userID)
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

func (s *PurchaseOrderService) Update(req UpdatePurchaseOrderRequest) *dto.APIError {
	err := s.repo.Update(req)
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

func (s *PurchaseOrderService) Delete(id string) *dto.APIError {
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

func (s *PurchaseOrderService) ReceivePurchaseOrder(id string, userID string, req []ReceivedItemRequest) *dto.APIError {
	err := s.repo.ReceiveOrder(id, req, userID)
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

func (s *PurchaseOrderService) CheckPurchaseOrder(id string, userID string) *dto.APIError {
	err := s.repo.CheckOrder(id, userID)
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

func (s *PurchaseOrderService) CancelPurchaseOrder(id string, userID string) *dto.APIError {
	err := s.repo.CancelOrder(id, userID)
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

func (s *PurchaseOrderService) CreateReturn(req CreatePurchaseOrderReturnRequest, userID string) *dto.APIError {
	err := s.repo.CreateReturn(req, userID)
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

func (s *PurchaseOrderService) GetAllReturns(req GetPurchaseOrderReturnRequest) ([]GetPurchaseOrderReturnResponse, int, *dto.APIError) {
	returns, totalItems, err := s.repo.GetAllReturns(req)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	return returns, totalItems, nil
}

func (s *PurchaseOrderService) CancelReturn(id string, userID string) *dto.APIError {
	err := s.repo.CancelReturn(id, userID)
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

func (s *PurchaseOrderService) RemovePurchaseOrderItem(id string) *dto.APIError {
	err := s.repo.RemoveItem(id)
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

func (s *PurchaseOrderService) AddPurchaseOrderItem(orderID string, req CreatePurchaseOrderItemRequest) *dto.APIError {
	err := s.repo.AddItem(orderID, req)
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

func (s *PurchaseOrderService) UpdatePurchaseOrderItem(req UpdatePurchaseOrderItemRequest) *dto.APIError {
	err := s.repo.UpdateItem(req)
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
