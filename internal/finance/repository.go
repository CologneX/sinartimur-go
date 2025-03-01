package finance

import (
	"database/sql"
	"sinartimur-go/internal/inventory"
)

type InventoryRepository interface {
	GetAll(req GetFinanceTransactionRequest) ([]GetFinanceTransactionResponse, int, error)
}

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) InventoryRepository {
	return &ProductRepositoryImpl{db: db}
}

// GetAll fetches all product
func (r *ProductRepositoryImpl) GetAll(req GetFinanceTransactionRequest) ([]GetFinanceTransactionResponse, int, error) {
	var transactions []GetFinanceTransactionResponse
	var totalItems int
	// Query for fetching all products with category and unit name
	query := "Select P.Id, P.Name,P.Description, P.Price, C.Name As Category, U.Name As Unit, P.Created_At, P.Updated_At From Product P Join Category C On P.Category_Id = C.Id Join Unit U On P.Unit_Id = U.Id Where P.Deleted_At Is Null"

	countQuery := "Select Count(Id) From financial_transaction Where Deleted_At Is Null"

	if req.UserID != "" {
		query += " AND u.id ILIKE '%" + req.UserID + "%'"
		countQuery += " AND id ILIKE '%" + req.UserID + "%'"
	}

	if req.Category != "" {
		query += " AND p.category_id = '" + req.Category + "'"
		countQuery += " AND category_id = '" + req.Category + "'"
	}

	if req.Unit != "" {
		query += " AND p.unit_id = '" + req.Unit + "'"
		countQuery += " AND unit_id = '" + req.Unit + "'"
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}

	countRow := r.db.QueryRow(countQuery)
	err = countRow.Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var transaction GetFinanceTransactionResponse
		err = rows.Scan(&transaction.ID, &transaction.Name, &transaction.Description, &transaction.Price, &transaction.Category, &transaction.Unit, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		products = append(products, transaction)
	}

	return products, totalItems, nil
}
