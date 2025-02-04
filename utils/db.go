package utils

import (
	"database/sql"
	"fmt"
)

// WithTransaction wraps a function execution within a database transaction.
func WithTransaction(db *sql.DB, fn func(*sql.Tx) error) (retErr error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if retErr != nil {
			if err := tx.Rollback(); err != nil {
				retErr = fmt.Errorf("transaction rollback failed: %w", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				retErr = fmt.Errorf("transaction commit failed: %w", err)
			}
		}
	}()

	retErr = fn(tx)
	return retErr
}
