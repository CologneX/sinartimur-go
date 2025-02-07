package products

import (
	"database/sql"
)

type WageRepository interface {
}

type WageRepositoryImpl struct {
	db *sql.DB
}

func NewWageRepository(db *sql.DB) WageRepository {
	return &WageRepositoryImpl{db: db}
}
