package postgres

import (
	"database/sql"
)

type StorageHistoryService struct {
	StorageService
}

func NewStorageHistoryService(db *sql.DB) *StorageHistoryService {
	return &StorageHistoryService{*NewStorageService(db)}
}
