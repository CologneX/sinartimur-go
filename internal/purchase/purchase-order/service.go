package purchase_order

import (
	"database/sql"
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
func (s *PurchaseOrderService) Create(req CreatePurchaseOrderRequest, userID string) *dto.APIError {
	// Validation for credit payment method
	if req.PaymentMethod == "credit" && req.PaymentDueDate == "" {
		return &dto.APIError{
			StatusCode: 400,
			Details: map[string]string{
				"payment_due_date": "Payment due date is required for credit payments",
			},
		}
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.Create(req, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to create purchase order: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
			},
		}
	}

	return nil
}

// GetPurchaseOrderDetail fetches details of a purchase order
func (s *PurchaseOrderService) GetPurchaseOrderDetail(id string) (*PurchaseOrderDetailResponse, *dto.APIError) {
	po, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Purchase order not found: " + err.Error(),
			},
		}
	}

	return po, nil
}

// ReceivePurchaseOrder handles receiving items for a purchase order
func (s *PurchaseOrderService) ReceivePurchaseOrder(id string, userID string, req []ReceivedItemRequest) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository to complete the purchase order
	if err := s.repo.CompletePurchaseOrder(id, req, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to receive purchase order: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
			},
		}
	}

	return nil
}

// CreateReturn handles creating a purchase order return
func (s *PurchaseOrderService) CreateReturn(req CreatePurchaseOrderReturnRequest, userID string) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Process return with transaction
	if err := s.repo.CreateReturn(req, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to create return: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
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
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Process return cancellation with transaction
	if err := s.repo.CancelReturn(id, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to cancel return: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
			},
		}
	}

	return nil
}

// Update handles updating a purchase order
func (s *PurchaseOrderService) Update(req UpdatePurchaseOrderRequest) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.Update(req, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to update purchase order: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
			},
		}
	}

	return nil
}

// Delete handles deleting a purchase order
func (s *PurchaseOrderService) Delete(id string) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.Delete(id, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to delete purchase order: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
			},
		}
	}

	return nil
}

// CheckPurchaseOrder handles checking a purchase order
func (s *PurchaseOrderService) CheckPurchaseOrder(id string, userID string) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.CheckPurchaseOrder(id, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to check purchase order: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
			},
		}
	}

	return nil
}

// CancelPurchaseOrder handles cancelling a purchase order
func (s *PurchaseOrderService) CancelPurchaseOrder(id string, userID string) *dto.APIError {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.CancelPurchaseOrder(id, userID, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to cancel purchase order: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
			},
		}
	}

	return nil
}

// GetAllPurchaseOrder fetches all purchase orders
func (s *PurchaseOrderService) GetAllPurchaseOrder(req GetPurchaseOrderRequest) ([]GetPurchaseOrderResponse, int, *dto.APIError) {
	orders, totalItems, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to fetch purchase orders: " + err.Error(),
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
				"general": "Failed to fetch purchase order returns: " + err.Error(),
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
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.AddPurchaseOrderItem(orderID, req, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to add purchase order item: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
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
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.UpdatePurchaseOrderItem(req, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to update purchase order item: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
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
				"general": "Failed to start transaction: " + err.Error(),
			},
		}
	}
	defer tx.Rollback()

	// Call repository with transaction
	if err := s.repo.RemovePurchaseOrderItem(id, tx); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to remove purchase order item: " + err.Error(),
			},
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return &dto.APIError{
			StatusCode: 500,
			Details: map[string]string{
				"general": "Failed to commit transaction: " + err.Error(),
			},
		}
	}

	return nil
}
