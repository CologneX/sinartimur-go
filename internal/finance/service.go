package finance

import (
	"sinartimur-go/pkg/dto"
	"time"
)

// FinanceService is the service for the Product domain.
type FinanceService struct {
	repo FinanceTransactionRepository
}

// NewFinanceTransactionService creates a new instance of FinanceService
func NewFinanceTransactionService(repo FinanceTransactionRepository) *FinanceService {
	return &FinanceService{repo: repo}
}

func (s *FinanceService) RefreshFinanceTransactionView() *dto.APIError {
	err := s.repo.RefreshFinanceTransactionView()
	if err != nil {
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal memperbarui data transaksi keuangan: " + err.Error(),
		})
	}
	return nil
}

func (s *FinanceService) GetFinanceTransactionViewLastRefreshed() (*time.Time, *dto.APIError) {
	lastRefreshed, err := s.repo.GetFinanceTransactionViewLastRefreshed()
	if err != nil {
		return nil, dto.NewAPIError(500, map[string]string{
			"general": "Gagal mendapatkan informasi waktu refresh terakhir: " + err.Error(),
		})
	}
	return lastRefreshed, nil
}

// CreateFinanceTransaction handles creating a new finance transaction
func (s *FinanceService) CreateFinanceTransaction(req CreateFinanceTransactionRequest, userID string) *dto.APIError {
	// Validate data
	if req.Amount <= 0 {
		return dto.NewAPIError(400, map[string]string{
			"amount": "Jumlah transaksi harus lebih besar dari 0",
		})
	}

	if req.Type == "" {
		return dto.NewAPIError(400, map[string]string{
			"type": "Tipe transaksi harus diisi",
		})
	}

	if req.Description == "" {
		return dto.NewAPIError(400, map[string]string{
			"description": "Deskripsi transaksi harus diisi",
		})
	}

	// Create the transaction
	err := s.repo.Create(req, userID)
	if err != nil {
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal membuat transaksi keuangan: " + err.Error(),
		})
	}

	return nil
}

// GetAllFinanceTransactions retrieves all finance transactions with pagination and filtering
func (s *FinanceService) GetAllFinanceTransactions(req GetFinanceTransactionRequest) ([]GetFinanceTransactionResponse, int, *dto.APIError) {
	// Fetch transactions from repository
	transactions, totalItems, err := s.repo.GetAll(req)
	if err != nil {
		return nil, 0, dto.NewAPIError(500, map[string]string{
			"general": "Gagal mengambil data transaksi keuangan: " + err.Error(),
		})
	}

	return transactions, totalItems, nil
}

// GetFinanceTransactionByID retrieves a single finance transaction by ID
func (s *FinanceService) GetFinanceTransactionByID(id string) (*GetFinanceTransactionResponse, *dto.APIError) {
	transaction, err := s.repo.GetByID(id)
	if err != nil {
		return nil, dto.NewAPIError(404, map[string]string{
			"general": err.Error(),
		})
	}

	return transaction, nil
}

// CancelFinanceTransaction cancels/soft deletes a finance transaction
func (s *FinanceService) CancelFinanceTransaction(req CancelFinanceTransactionRequest, userID string) *dto.APIError {
	// Check if transaction exists
	transaction, err := s.repo.GetByID(req.ID)
	if err != nil {
		return dto.NewAPIError(404, map[string]string{
			"general": "Transaksi tidak ditemukan",
		})
	}

	// Check if it's a system transaction
	if transaction.IsSystem {
		return dto.NewAPIError(403, map[string]string{
			"general": "Transaksi yang dibuat sistem tidak dapat dibatalkan",
		})
	}

	// Cancel the transaction
	err = s.repo.Cancel(req, userID)
	if err != nil {
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal membatalkan transaksi: " + err.Error(),
		})
	}

	return nil
}

// GetFinanceTransactionSummary retrieves financial summary for a date range
func (s *FinanceService) GetFinanceTransactionSummary(startDate, endDate time.Time) (*FinanceTransactionSummary, *dto.APIError) {
	// Validate date range
	if startDate.IsZero() || endDate.IsZero() {
		return nil, dto.NewAPIError(400, map[string]string{
			"date_range": "Rentang tanggal harus diisi",
		})
	}

	if endDate.Before(startDate) {
		return nil, dto.NewAPIError(400, map[string]string{
			"date_range": "Tanggal akhir tidak boleh sebelum tanggal awal",
		})
	}

	// Get summary from repository
	summary, err := s.repo.GetSummary(startDate, endDate)
	if err != nil {
		return nil, dto.NewAPIError(500, map[string]string{
			"general": "Gagal mengambil log transaksi: " + err.Error(),
		})
	}

	return summary, nil
}
