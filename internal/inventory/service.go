package inventory

import (
	"database/sql"
	"errors"
	"sinartimur-go/pkg/dto"
	"strings"
	"time"
)

// StorageService is the service for the Storage domain
type StorageService struct {
	repo StorageRepository
}

// NewStorageService creates a new instance of StorageService
func NewStorageService(repo StorageRepository) *StorageService {
	return &StorageService{repo: repo}
}

// GetAllStorages fetches all storage locations with filtering and pagination
func (s *StorageService) GetAllStorages(req GetStorageRequest) ([]GetStorageResponse, int, *dto.APIError) {
	storages, totalItems, err := s.repo.GetAllStorages(req)
	if err != nil {
		return nil, 0, dto.NewAPIError(500, map[string]string{
			"general": "Gagal mengambil data gudang",
		})
	}
	return storages, totalItems, nil
}

// GetStorageByID fetches a storage location by ID
func (s *StorageService) GetStorageByID(id string) (*GetStorageResponse, *dto.APIError) {
	storage, err := s.repo.GetStorageByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dto.NewAPIError(404, map[string]string{
				"general": "Gudang tidak ditemukan",
			})
		}
		return nil, dto.NewAPIError(500, map[string]string{
			"general": "Gagal mengambil data gudang",
		})
	}
	return storage, nil
}

// CreateStorage creates a new storage location
func (s *StorageService) CreateStorage(req CreateStorageRequest) (*GetStorageResponse, *dto.APIError) {
	// Check if storage with same name already exists
	existing, err := s.repo.GetStorageByName(req.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, dto.NewAPIError(500, map[string]string{
			"general": "Gagal memeriksa nama gudang",
		})
	}

	if existing != nil {
		return nil, dto.NewAPIError(400, map[string]string{
			"name": "Nama gudang sudah digunakan",
		})
	}

	storage, err := s.repo.CreateStorage(req)
	if err != nil {
		return nil, dto.NewAPIError(500, map[string]string{
			"general": "Gagal membuat gudang baru",
		})
	}

	return storage, nil
}

// UpdateStorage updates an existing storage location
func (s *StorageService) UpdateStorage(req UpdateStorageRequest) (*GetStorageResponse, *dto.APIError) {
	// Check if storage exists
	_, err := s.repo.GetStorageByID(req.ID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dto.NewAPIError(404, map[string]string{
				"general": "Gudang tidak ditemukan",
			})
		}
		return nil, dto.NewAPIError(500, map[string]string{
			"general": "Gagal memeriksa keberadaan gudang",
		})
	}

	// Check if name is already taken by another storage
	existing, err := s.repo.GetStorageByName(req.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, dto.NewAPIError(500, map[string]string{
			"general": "Gagal memeriksa nama gudang",
		})
	}

	if existing != nil && existing.ID != req.ID.String() {
		return nil, dto.NewAPIError(400, map[string]string{
			"name": "Nama gudang sudah digunakan",
		})
	}

	storage, err := s.repo.UpdateStorage(req)
	if err != nil {
		return nil, dto.NewAPIError(500, map[string]string{
			"general": "Gagal mengupdate gudang",
		})
	}

	return storage, nil
}

// DeleteStorage deletes a storage location
func (s *StorageService) DeleteStorage(id string) *dto.APIError {
	// Check if storage exists
	_, err := s.repo.GetStorageByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NewAPIError(404, map[string]string{
				"general": "Gudang tidak ditemukan",
			})
		}
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal memeriksa keberadaan gudang",
		})
	}

	if err := s.repo.DeleteStorage(id); err != nil {
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal menghapus gudang",
		})
	}

	return nil
}

// MoveBatch moves products from one storage to another
func (s *StorageService) MoveBatch(req MoveBatchRequest, userID string) *dto.APIError {
	// Validate source storage exists
	_, err := s.repo.GetStorageByID(req.SourceStorageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NewAPIError(404, map[string]string{
				"source_storage_id": "Gudang sumber tidak ditemukan",
			})
		}
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal memeriksa gudang sumber",
		})
	}

	// Validate target storage exists
	_, err = s.repo.GetStorageByID(req.TargetStorageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NewAPIError(404, map[string]string{
				"target_storage_id": "Gudang tujuan tidak ditemukan",
			})
		}
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal memeriksa gudang tujuan",
		})
	}

	// Check if batch exists in source
	_, err = s.repo.GetBatchInStorage(req.BatchID, req.SourceStorageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NewAPIError(404, map[string]string{
				"batch_id": "Batch produk tidak ditemukan di gudang sumber",
			})
		}
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal memeriksa ketersediaan batch",
		})
	}

	// Perform the move operation
	err = s.repo.MoveBatch(req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "kuantitas tidak mencukupi") {
			return dto.NewAPIError(400, map[string]string{
				"quantity": "Kuantitas tidak mencukupi di gudang sumber",
			})
		}
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal memindahkan batch: " + err.Error(),
		})
	}

	return nil
}

// GetInventoryLogs fetches inventory logs with filtering and pagination
func (s *StorageService) GetInventoryLogs(req GetInventoryLogsRequest) ([]GetInventoryLogResponse, int, *dto.APIError) {
	logs, totalItems, err := s.repo.GetInventoryLogs(req)
	if err != nil {
		return nil, 0, dto.NewAPIError(500, map[string]string{
			"general": "Gagal mengambil log inventaris: " + err.Error(),
		})
	}
	return logs, totalItems, nil
}

// RefreshInventoryLogView refreshes the materialized view
func (s *StorageService) RefreshInventoryLogView() *dto.APIError {
	err := s.repo.RefreshInventoryLogView()
	if err != nil {
		return dto.NewAPIError(500, map[string]string{
			"general": "Gagal memperbarui data log inventaris: " + err.Error(),
		})
	}
	return nil
}

// GetInventoryLogLastRefreshed fetches the last time the inventory log was refreshed
func (s *StorageService) GetInventoryLogLastRefreshed() (*time.Time, *dto.APIError) {
	lastRefreshed, err := s.repo.GetInventoryLogLastRefreshed()
	if err != nil {
		return nil, dto.NewAPIError(500, map[string]string{
			"general": "Gagal mendapatkan informasi waktu refresh terakhir: " + err.Error(),
		})
	}
	return lastRefreshed, nil
}
