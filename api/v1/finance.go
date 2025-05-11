package v1

import (
	"net/http"
	"sinartimur-go/internal/finance"
	"sinartimur-go/pkg/dto"
	"sinartimur-go/utils"
	"time"
)

// CreateFinanceTransactionHandler handles requests to create a new financial transaction
func CreateFinanceTransactionHandler(financialService *finance.FinanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get user ID from context
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusUnauthorized, map[string]string{
				"general": "Tidak terautentikasi",
			}))
			return
		}

		// Decode and validate the request
		var req finance.CreateFinanceTransactionRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, errors))
			return
		}

		// Create the transaction
		apiErr := financialService.CreateFinanceTransaction(req, userID)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		// Refresh the materialized view after creating a transaction
		_ = financialService.RefreshFinanceTransactionView()

		utils.WriteJSON(w, http.StatusCreated, map[string]string{
			"message": "Transaksi keuangan berhasil dibuat",
		})
	}
}

// GetAllFinanceTransactionsHandler handles requests to get all financial transactions with pagination and filtering
func GetAllFinanceTransactionsHandler(financialService *finance.FinanceService) http.HandlerFunc {
	return utils.NewPaginatedHandler(func(w http.ResponseWriter, r *http.Request, page, pageSize int, sortBy, sortOrder string) {
		// Create request with pagination and filtering parameters
		var req finance.GetFinanceTransactionRequest

		// Parse additional filter parameters
		if userID := r.URL.Query().Get("user_id"); userID != "" {
			req.UserID = userID
		}

		if transactionType := r.URL.Query().Get("type"); transactionType != "" {
			req.Type = transactionType
		}

		if purchaseOrderID := r.URL.Query().Get("purchase_order_id"); purchaseOrderID != "" {
			req.PurchaseOrderID = purchaseOrderID
		}

		if salesOrderID := r.URL.Query().Get("sales_order_id"); salesOrderID != "" {
			req.SalesOrderID = salesOrderID
		}

		// Parse boolean parameters
		if isSystem := r.URL.Query().Get("is_system"); isSystem != "" {
			isSystemVal := isSystem == "true"
			req.IsSystem = &isSystemVal
		}

		// Parse date ranges
		if startDate := r.URL.Query().Get("start_date"); startDate != "" {
			req.StartDate = startDate
		}

		if endDate := r.URL.Query().Get("end_date"); endDate != "" {
			req.EndDate = endDate
		}

		// Set pagination parameters
		req.Page = page
		req.PageSize = pageSize
		req.SortBy = sortBy
		req.SortOrder = sortOrder

		// Get transactions from service
		transactions, totalItems, apiErr := financialService.GetAllFinanceTransactions(req)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		// Get last refresh time
		lastRefreshed, _ := financialService.GetFinanceTransactionViewLastRefreshed()

		// Add last refresh time to response
		response := map[string]interface{}{
			"data":           transactions,
			"last_refreshed": lastRefreshed,
		}

		utils.WritePaginationJSON(w, http.StatusOK, page, totalItems, pageSize, response)
	})
}

// CancelFinanceTransactionHandler handles requests to cancel a financial transaction
func CancelFinanceTransactionHandler(financialService *finance.FinanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID := r.Context().Value("userId").(string)

		// Decode and validate the request
		var req finance.CancelFinanceTransactionRequest
		errors := utils.DecodeAndValidate(r, &req)
		if errors != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, errors))
			return
		}

		// Cancel the transaction
		apiErr := financialService.CancelFinanceTransaction(req, userID)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		// Refresh the materialized view after canceling a transaction
		_ = financialService.RefreshFinanceTransactionView()

		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "Transaksi keuangan berhasil dibatalkan",
		})
	}
}

// GetFinanceTransactionSummaryHandler handles requests to get a summary of financial transactions
func GetFinanceTransactionSummaryHandler(financialService *finance.FinanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse date range parameters
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		if startDateStr == "" || endDateStr == "" {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"date_range": "Rentang tanggal harus diisi",
			}))
			return
		}

		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"start_date": "Format tanggal awal tidak valid",
			}))
			return
		}

		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			utils.ErrorJSON(w, dto.NewAPIError(http.StatusBadRequest, map[string]string{
				"end_date": "Format tanggal akhir tidak valid",
			}))
			return
		}

		// Get financial summary
		summary, apiErr := financialService.GetFinanceTransactionSummary(startDate, endDate)
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, summary)
	}
}

// RefreshFinanceTransactionViewHandler handles manual refresh requests for the finance transaction materialized view
func RefreshFinanceTransactionViewHandler(financialService *finance.FinanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiErr := financialService.RefreshFinanceTransactionView()
		if apiErr != nil {
			utils.ErrorJSON(w, apiErr)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.WriteMessage("Log finansial berhasil diperbaharui"))
	}
}
