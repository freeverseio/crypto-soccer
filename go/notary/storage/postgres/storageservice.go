package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

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

func New(url string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
