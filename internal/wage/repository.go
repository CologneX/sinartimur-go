package wage

import (
	"database/sql"
	"github.com/google/uuid"
	"sinartimur-go/internal/employee"
	"sinartimur-go/utils"
	"strconv"
)

type WageRepository interface {
	Create(request CreateWageRequest) error
	CreateDetail(request WageDetailRequest) error
	Delete(request DeleteWageRequest) error
	//DeleteDetail(request DeleteWageDetailRequest) error
	UpdateDetail(request UpdateWageDetailRequest) error
	GetAll(request GetWageRequest) ([]GetWageResponse, int, error)
	GetWageDetailByWageID(wageID string) ([]*GetWageDetail, error)
	GetByID(id string) (*GetWageResponse, error)
	GetDetailByID(id string) (*WageDetail, error)
	GetEmployeeByID(id string) (*employee.GetEmployeeResponse, error)
	GetByEmployeeIdAndMonthYear(employeeId string, month int, year int) (*GetWageResponse, error)
}

type WageRepositoryImpl struct {
	db *sql.DB
}

func NewWageRepository(db *sql.DB) WageRepository {
	return &WageRepositoryImpl{db: db}
}

// GetEmployeeByID fetches employee by id
func (r *WageRepositoryImpl) GetEmployeeByID(id string) (*employee.GetEmployeeResponse, error) {
	var emp employee.GetEmployeeResponse

	err := r.db.QueryRow("Select Id, Name, Position, Phone, Nik, Hired_Date, Created_At, Updated_At From Employee Where Id = $1", id).Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Phone, &emp.Nik, &emp.HiredDate, &emp.CreatedAt, &emp.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &emp, nil
}

// Create creates a new wage
func (r *WageRepositoryImpl) Create(request CreateWageRequest) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		var wageID uuid.UUID

		// Get total amount
		var totalAmount float64
		for _, detail := range request.WageDetail {
			totalAmount += detail.Amount
		}

		// Insert into Wages
		if err := tx.QueryRow("Insert Into Wage (Employee_Id, Month, Year, Total_Amount) Values ($1, $2, $3, $4) Returning Id", request.EmployeeId, request.Month, request.Year, totalAmount).Scan(&wageID); err != nil {
			return err
		}

		// Insert into Wage_Details
		for _, detail := range request.WageDetail {
			if _, err := tx.Exec("Insert Into Wage_Detail (Wage_Id, Component_Name, Description, Amount) Values ($1, $2, $3, $4)", wageID, detail.ComponentName, detail.Description, detail.Amount); err != nil {
				return err
			}
		}

		return nil
	})
}

// CreateDetail creates a new wage Detail
func (r *WageRepositoryImpl) CreateDetail(request WageDetailRequest) error {
	_, err := r.db.Exec("Insert Into Wage_Detail (Component_Name, Description, Amount) Values ($1, $2, $3)", request.ComponentName, request.Description, request.Amount)
	if err != nil {
		return err
	}

	return nil
}

// Delete soft deletes a wage
func (r *WageRepositoryImpl) Delete(request DeleteWageRequest) error {
	_, err := r.db.Exec("Update Wage Set Deleted_At = Now() Where Id = $1", request.ID)
	if err != nil {
		return err
	}

	return nil
}

//// DeleteDetail deletes a wage Detail
//func (r *WageRepositoryImpl) DeleteDetail(request DeleteWageDetailRequest) error {
//	_, err := r.db.Exec("Delete From Wage_Details Where Id = $1", request.ID)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// UpdateDetail updates a wage
func (r *WageRepositoryImpl) UpdateDetail(request UpdateWageDetailRequest) error {
	return utils.WithTransaction(r.db, func(tx *sql.Tx) error {
		// Get total amount
		var totalAmount float64
		for _, detail := range request.WageDetail {
			totalAmount += detail.Amount
		}

		// UpdateDetail Wages
		if _, err := tx.Exec("Update Wage Set Total_Amount = $1, Updated_At = Now() Where Id = $2", totalAmount, request.ID); err != nil {
			return err
		}

		// Delete existing wage details
		if _, err := tx.Exec("Delete From Wage_Detail Where Wage_Id = $1", request.ID); err != nil {
			return err
		}

		// Insert new wage details
		for _, detail := range request.WageDetail {
			if _, err := tx.Exec("Insert Into Wage_Detail (Wage_Id, Component_Name, Description, Amount) Values ($1, $2, $3, $4)", request.ID, detail.ComponentName, detail.Description, detail.Amount); err != nil {
				return err
			}
		}

		return nil
	})
}

// GetAll fetches all wages
// GetAll fetches all wages
func (r *WageRepositoryImpl) GetAll(request GetWageRequest) ([]GetWageResponse, int, error) {
	var wages []GetWageResponse
	var totalItems int

	query := `Select W.Id, W.Employee_Id, E.Name, W.Total_Amount, W.Month, W.Year, W.Created_At, W.Updated_At
			  From Wage W
			  Join Employee E On W.Employee_Id = E.Id
			  Where W.Deleted_At Is Null`
	countQuery := `Select Count(*)
				   From Wage W
				   Join Employee E On W.Employee_Id = E.Id
				   Where W.Deleted_At Is Null`

	var args []interface{}
	argIndex := 1

	if request.EmployeeId != "" {
		query += ` AND W.Employee_Id = $` + strconv.Itoa(argIndex)
		countQuery += ` AND W.Employee_Id = $` + strconv.Itoa(argIndex)
		args = append(args, request.EmployeeId)
		argIndex++
	}
	if request.Month != 0 {
		query += ` AND W.Month = $` + strconv.Itoa(argIndex)
		countQuery += ` AND W.Month = $` + strconv.Itoa(argIndex)
		args = append(args, request.Month)
		argIndex++
	}
	if request.Year != 0 {
		query += ` AND W.Year = $` + strconv.Itoa(argIndex)
		countQuery += ` AND W.Year = $` + strconv.Itoa(argIndex)
		args = append(args, request.Year)
		argIndex++
	}

	if request.SortBy != "" {
		query += ` ORDER BY ` + request.SortBy
		if request.SortOrder != "" {
			query += ` ` + request.SortOrder
		}
	} else {
		query += ` ORDER BY W.Created_At DESC`
	}

	query += ` LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, request.PageSize, (request.Page-1)*request.PageSize)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var wage GetWageResponse
		err := rows.Scan(&wage.ID, &wage.EmployeeId, &wage.EmployeeName, &wage.TotalAmount, &wage.Month, &wage.Year, &wage.CreatedAt, &wage.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		wages = append(wages, wage)
	}

	err = r.db.QueryRow(countQuery, args[:argIndex-1]...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	return wages, totalItems, nil
}

// GetWageDetailByWageID fetches wage details by wage id
func (r *WageRepositoryImpl) GetWageDetailByWageID(wageID string) ([]*GetWageDetail, error) {
	rows, err := r.db.Query("Select Id, Component_Name, Description, Amount, Created_At, Updated_At From Wage_Detail Where Wage_Id = $1", wageID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		errCl := rows.Close()
		if errCl != nil {
			err = errCl
		}
	}(rows)

	var details []*GetWageDetail
	for rows.Next() {
		var detail GetWageDetail
		err = rows.Scan(&detail.ID, &detail.ComponentName, &detail.Description, &detail.Amount, &detail.CreatedAt, &detail.UpdatedAt)
		if err != nil {
			return nil, err
		}
		details = append(details, &detail)
	}

	return details, nil
}

// GetByID fetches wage by id
func (r *WageRepositoryImpl) GetByID(id string) (*GetWageResponse, error) {
	var wage GetWageResponse

	err := r.db.QueryRow("Select W.Id, W.Employee_Id, E.Name, W.Total_Amount, W.Month, W.Year, W.Created_At, W.Updated_At From Wage W Join Employee E On W.Employee_Id = E.Id Where W.Id = $1 And W.Deleted_At Is Null", id).Scan(&wage.ID, &wage.EmployeeId, &wage.EmployeeName, &wage.TotalAmount, &wage.Month, &wage.Year, &wage.CreatedAt, &wage.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &wage, nil
}

// GetDetailByID fetches wage Detail by id
func (r *WageRepositoryImpl) GetDetailByID(id string) (*WageDetail, error) {
	var wageDetail WageDetail

	err := r.db.QueryRow("Select Id, Wage_Id, Component_Name, Description, Amount, Created_At, Updated_At From Wage_Detail Where Id = $1 And Deleted_At Is Null", id).Scan(&wageDetail.ID, &wageDetail.WageId, &wageDetail.ComponentName, &wageDetail.Description, &wageDetail.Amount, &wageDetail.CreatedAt, &wageDetail.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &wageDetail, nil
}

// GetByEmployeeIdAndMonthYear fetches wage by employee id, month, and year
func (r *WageRepositoryImpl) GetByEmployeeIdAndMonthYear(employeeId string, month int, year int) (*GetWageResponse, error) {
	var wage GetWageResponse

	err := r.db.QueryRow("Select Id, Employee_Id, Total_Amount, Month, Year, Created_At, Updated_At From Wage Where Employee_Id = $1 And Month = $2 And Year = $3 And Deleted_At Is Null", employeeId, month, year).Scan(&wage.ID, &wage.EmployeeId, &wage.TotalAmount, &wage.Month, &wage.Year, &wage.CreatedAt, &wage.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &wage, nil
}
