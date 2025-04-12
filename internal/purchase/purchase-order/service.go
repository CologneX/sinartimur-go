package purchase_order

import (
	"sinartimur-go/pkg/dto"
)

type PurchaseOrderService struct {
	repo Repository
}

// NewPurchaseOrderService creates a new instance of PurchaseOrderService
func NewPurchaseOrderService(repo Repository) *PurchaseOrderService {
	return &PurchaseOrderService{repo: repo}
}

// GetPurchaseOrderDetail fetches the details of a purchase purchase-order by ID
func (s *PurchaseOrderService) GetPurchaseOrderDetail(id string) (*PurchaseOrderDetailResponse, *dto.APIError) {
	// get the purchase purchase-order by ID
	purchaseOrder, err := s.repo.GetByID(id)
	if err != nil {
		return nil, &dto.APIError{
			StatusCode: 404,
			Details: map[string]string{
				"general": "Purchase purchase-order tidak ditemukan",
			},
		}
	}
	return purchaseOrder, nil
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
