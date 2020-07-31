package postgres

import "database/sql"

type StorageService struct {
	db *sql.DB
}

func NewStorageService(db *sql.DB) *StorageService {
	return &StorageService{
		db: db,
	}
}

func (b *StorageService) DB() *sql.DB {
	return b.db
}
