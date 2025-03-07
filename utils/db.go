package utils

import (
	"database/sql"
	"fmt"
	"strings"
)

// QueryBuilder helps construct SQL queries dynamically
type QueryBuilder struct {
	Query  strings.Builder
	params []interface{}
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
		qb.params = append(qb.params, value)
		qb.count++
	}
	return qb
}

// AddPagination adds pagination parameters
func (qb *QueryBuilder) AddPagination(pageSize, page int) *QueryBuilder {
	qb.Query.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", qb.count, qb.count+1))
	qb.params = append(qb.params, pageSize, (page-1)*pageSize)
	qb.count += 2
	return qb
}

// Build returns the final query and parameters
func (qb *QueryBuilder) Build() (string, []interface{}) {
	return qb.Query.String(), qb.params
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
