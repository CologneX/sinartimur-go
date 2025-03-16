package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// QueryBuilder helps construct SQL queries dynamically
type QueryBuilder struct {
	Query  strings.Builder
	Params []interface{}
	count  int
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(baseQuery string) *QueryBuilder {
	qb := &QueryBuilder{
		count: 1,
	}
	qb.Query.WriteString(baseQuery)
	return qb
}

// AddFilter adds a filter condition with parameter
func (qb *QueryBuilder) AddFilter(condition string, value interface{}) *QueryBuilder {
	if value != nil && value != "" {
		qb.Query.WriteString(fmt.Sprintf(" AND %s $%d", condition, qb.count))
		qb.Params = append(qb.Params, value)
		qb.count++
	}
	return qb
}

// AddPagination adds pagination parameters
func (qb *QueryBuilder) AddPagination(pageSize, page int) *QueryBuilder {
	qb.Query.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", qb.count, qb.count+1))
	qb.Params = append(qb.Params, pageSize, (page-1)*pageSize)
	qb.count += 2
	return qb
}

// Build returns the final query and parameters
func (qb *QueryBuilder) Build() (string, []interface{}) {
	return qb.Query.String(), qb.Params
}

// WithTransaction executes a function within a transaction
func WithTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if errComm := tx.Commit(); errComm != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// GenerateNextSerialID generates the next serial ID for a document type
// Note: This should be called within a transaction to ensure atomicity
func GenerateNextSerialID(tx *sql.Tx, documentType string) (string, error) {
	now := time.Now()
	year, month, day := now.Year(), int(now.Month()), now.Day()

	// Try to update existing counter for today
	var counter int
	err := tx.QueryRow(`
        UPDATE document_counter 
        SET counter = counter + 1, last_updated = NOW()
        WHERE document_type = $1 AND year = $2 AND month = $3 AND day = $4
        RETURNING counter
    `, documentType, year, month, day).Scan(&counter)

	// If no rows exist for today, insert a new counter
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow(`
            INSERT INTO document_counter (document_type, year, month, day, counter)
            VALUES ($1, $2, $3, $4, 1)
            ON CONFLICT (document_type, year, month, day) 
            DO UPDATE SET counter = document_counter.counter + 1
            RETURNING counter
        `, documentType, year, month, day).Scan(&counter)
	}

	if err != nil {
		return "", fmt.Errorf("gagal mengupdate penghitung dokumen: %w", err)
	}

	// Format the serial number: XX-YYYYMMDD-NNNN
	// XX = document type, YYYYMMDD = date, NNNN = counter (padded with zeros)
	return fmt.Sprintf(
		"%s-%04d%02d%02d-%04d",
		documentType,
		year, month, day,
		counter,
	), nil
}
