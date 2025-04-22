package purchase_order

import (
	"database/sql"
	"sinartimur-go/internal/product"
	"sinartimur-go/pkg/dto"
)

// PurchaseOrderService handles business logic for purchase orders
type PurchaseOrderService struct {
	repo Repository
	db   *sql.DB
}

// NewPurchaseOrderService creates a new purchase order service
func NewPurchaseOrderService(repo Repository, db *sql.DB) *PurchaseOrderService {
	return &PurchaseOrderService{
		repo: repo,
		db:   db,
	}
}

// Create handles creating a purchase order
func (s *PurchaseOrderService) Create(req CreatePurchaseOrderRequest, userID string) (*CreatePurchaseOrderResponse, *dto.APIError) {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	purchaseOrderID, err := s.repo.Create(req, userID, tx)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Retrieve the created purchase order
	purchaseOrder, err := s.repo.GetByID(purchaseOrderID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Convert to CreatePurchaseOrderResponse
	response := &CreatePurchaseOrderResponse{
		GetPurchaseOrderResponse: GetPurchaseOrderResponse{
			ID:           purchaseOrder.ID,
			SerialID:     purchaseOrder.SerialID,
			SupplierID:   purchaseOrder.SupplierID,
			SupplierName: purchaseOrder.SupplierName,
			OrderDate:    purchaseOrder.OrderDate,
			Status:       purchaseOrder.Status,
			TotalAmount:  purchaseOrder.TotalAmount,
			CreatedBy:    purchaseOrder.CreatedBy,
			CreatedAt:    purchaseOrder.CreatedAt,
			UpdatedAt:    purchaseOrder.UpdatedAt,
			ItemCount:    len(purchaseOrder.Items),
		},
	}

	return response, nil
}

// GetPurchaseOrderDetail fetches details of a purchase order
func (s *PurchaseOrderService) GetPurchaseOrderDetail(id string) (*GetPurchaseOrderDetailResponse, *dto.APIError) {
	po, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return po, nil
}

// ReceivePurchaseOrder handles receiving items for a purchase order
// func (s *PurchaseOrderService) ReceivePurchaseOrder(id string, userID string, req []ReceivedItemRequest) *dto.APIError {
// 	// Start transaction
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return &dto.APIError{
// 			StatusCode: 500,
// 			Details: map[string]string{
// 				"general": err.Error(),
// 			},
// 		}
// 	}
// 	defer tx.Rollback()

// 	// Call repository to complete the purchase order
// 	if err := s.repo.CompletePurchaseOrder(id, req, userID, tx); err != nil {
// 		return &dto.APIError{
// 			StatusCode: 500,
// 			Details: map[string]string{
// 				"general": err.Error(),
// 			},
// 		}
// 	}

// 	// Commit transaction
// 	if err := tx.Commit(); err != nil {
// 		return &dto.APIError{
// 			StatusCode: 500,
// 			Details: map[string]string{
// 				"general": err.Error(),
// 			},
// 		}
// 	}

// 	return nil
// }

// CreateReturn handles creating a purchase order return
func (s *PurchaseOrderService) CreateReturn(req CreatePurchaseOrderReturnRequest, userID string) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Process return with transaction
	if err := s.repo.CreateReturn(req, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}

// CancelReturn handles cancelling a purchase order return
func (s *PurchaseOrderService) CancelReturn(id string, userID string) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Process return cancellation with transaction
	if err := s.repo.CancelReturn(id, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}

// Update handles updating a purchase order
func (s *PurchaseOrderService) Update(req UpdatePurchaseOrderRequest) (*UpdatePurchaseOrderResponse, *dto.APIError) {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	orderID, err := s.repo.Update(req, tx)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Retrieve the updated purchase order
	purchaseOrder, err := s.repo.GetByID(orderID)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Convert to UpdatePurchaseOrderResponse
	response := &UpdatePurchaseOrderResponse{
		GetPurchaseOrderResponse: GetPurchaseOrderResponse{
			ID:           purchaseOrder.ID,
			SerialID:     purchaseOrder.SerialID,
			SupplierID:   purchaseOrder.SupplierID,
			SupplierName: purchaseOrder.SupplierName,
			OrderDate:    purchaseOrder.OrderDate,
			Status:       purchaseOrder.Status,
			TotalAmount:  purchaseOrder.TotalAmount,
			CreatedBy:    purchaseOrder.CreatedBy,
			CreatedAt:    purchaseOrder.CreatedAt,
			UpdatedAt:    purchaseOrder.UpdatedAt,
			ItemCount:    len(purchaseOrder.Items),
		},
	}

	return response, nil
}

// CheckPurchaseOrder handles checking a purchase order
func (s *PurchaseOrderService) CheckPurchaseOrder(id string, userID string) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.CheckPurchaseOrder(id, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}

// CancelPurchaseOrder handles cancelling a purchase order
func (s *PurchaseOrderService) CancelPurchaseOrder(id string, userID string) (*CancelPurchaseOrderResponse, *dto.APIError) {
	// Get purchase order before cancellation to return its details later
	purchaseOrder, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.CancelPurchaseOrder(id, userID, tx); err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	response := &CancelPurchaseOrderResponse{
		GetPurchaseOrderResponse: GetPurchaseOrderResponse{
			ID:           purchaseOrder.ID,
			SerialID:     purchaseOrder.SerialID,
			SupplierID:   purchaseOrder.SupplierID,
			SupplierName: purchaseOrder.SupplierName,
			OrderDate:    purchaseOrder.OrderDate,
			Status:       "cancelled",
			TotalAmount:  purchaseOrder.TotalAmount,
			CreatedBy:    purchaseOrder.CreatedBy,
			CreatedAt:    purchaseOrder.CreatedAt,
			UpdatedAt:    purchaseOrder.UpdatedAt,
			ItemCount:    len(purchaseOrder.Items),
		},
	}

	return response, nil
}

// GetAllPurchaseOrder fetches all purchase orders
func (s *PurchaseOrderService) GetAllPurchaseOrder(req GetPurchaseOrderRequest) ([]GetPurchaseOrderResponse, int, *dto.APIError) {
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

// GetAllReturns fetches all purchase order returns
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

// AddPurchaseOrderItem adds an item to a purchase order
func (s *PurchaseOrderService) AddPurchaseOrderItem(orderID string, req CreatePurchaseOrderItemRequest) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.AddPurchaseOrderItem(orderID, req, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}

// UpdatePurchaseOrderItem updates a purchase order item
func (s *PurchaseOrderService) UpdatePurchaseOrderItem(req UpdatePurchaseOrderItemRequest) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.UpdatePurchaseOrderItem(req, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}

// RemovePurchaseOrderItem removes a purchase order item
func (s *PurchaseOrderService) RemovePurchaseOrderItem(id string) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.RemovePurchaseOrderItem(id, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}

// GetAllProducts fetches all products
func (s *PurchaseOrderService) GetAllProducts(req product.GetProductRequest) ([]product.GetProductResponse, int, *dto.APIError) {
	products, totalItems, err := s.repo.GetProducts(req)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	return products, totalItems, nil
}

// CompleteFullPurchaseOrder handles completing an entire purchase order at once
func (s *PurchaseOrderService) CompleteFullPurchaseOrder(req CompletePurchaseOrderRequest, userID string) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository to complete the purchase order
	if err := s.repo.CompleteFullPurchaseOrder(req.PurchaseOrderID, req.StorageID, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": err.Error(),
			},
		}
	}

	return nil
}
