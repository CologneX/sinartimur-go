package unit

import (
	"database/sql"
)

type UnitRepository interface {
	GetAll(req GetUnitRequest) ([]GetUnitResponse, error)
	GetByID(id string) (*GetUnitResponse, error)
	GetByName(name string) (*GetUnitResponse, error)
	Create(req CreateUnitRequest) (*GetUnitResponse, error)
	Update(req UpdateUnitRequest) (*GetUnitResponse, error)
	Delete(req DeleteUnitRequest) error
}

type UnitRepositoryImpl struct {
	db *sql.DB
}

func NewUnitRepository(db *sql.DB) UnitRepository {
	return &UnitRepositoryImpl{db: db}
}

// GetAll fetches all units
func (r *UnitRepositoryImpl) GetAll(req GetUnitRequest) ([]GetUnitResponse, error) {
	var units []GetUnitResponse
	rows, err := r.db.Query("SELECT id, name, description, created_at, updated_at FROM unit WHERE deleted_at is null AND name ILIKE $1", "%"+req.Name+"%")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var unit GetUnitResponse
		err = rows.Scan(&unit.ID, &unit.Name, &unit.Description, &unit.CreatedAt, &unit.UpdatedAt)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}

	return units, nil
}

// GetByID fetches a unit by ID
func (r *UnitRepositoryImpl) GetByID(id string) (*GetUnitResponse, error) {
	var unit GetUnitResponse
	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM unit WHERE id = $1 AND deleted_at is null", id).Scan(&unit.ID, &unit.Name, &unit.Description, &unit.CreatedAt, &unit.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// GetByName fetches a unit by name
func (r *UnitRepositoryImpl) GetByName(name string) (*GetUnitResponse, error) {
	var unit GetUnitResponse
	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM unit WHERE name = $1 AND deleted_at is null", name).Scan(&unit.ID, &unit.Name, &unit.Description, &unit.CreatedAt, &unit.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// Create creates a new unit
func (r *UnitRepositoryImpl) Create(req CreateUnitRequest) (*GetUnitResponse, error) {
	var unit GetUnitResponse
	err := r.db.QueryRow("INSERT INTO unit (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at", req.Name, req.Description).Scan(&unit.ID, &unit.Name, &unit.Description, &unit.CreatedAt, &unit.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// Update updates an existing unit
func (r *UnitRepositoryImpl) Update(req UpdateUnitRequest) (*GetUnitResponse, error) {
	var unit GetUnitResponse
	err := r.db.QueryRow("UPDATE unit SET name = $1, description = $2, updated_at = now() WHERE id = $3 RETURNING id, name, description, created_at, updated_at", req.Name, req.Description, req.ID).Scan(&unit.ID, &unit.Name, &unit.Description, &unit.CreatedAt, &unit.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// Delete marks a unit as deleted
func (r *UnitRepositoryImpl) Delete(req DeleteUnitRequest) error {
	_, err := r.db.Exec("UPDATE unit SET deleted_at = now() WHERE id = $1", req.ID)
	if err != nil {
		return err
	}
	return nil
}
