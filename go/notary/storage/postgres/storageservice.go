package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type StorageService struct {
	db *sql.DB
	tx *sql.Tx
}

func NewStorageService(db *sql.DB) *StorageService {
	return &StorageService{
		db: db,
	}
}

func (b *StorageService) Begin() error {
	var err error
	b.tx, err = b.db.Begin()
	return err
}

func (b *StorageService) Rollback() error {
	return b.tx.Rollback()
}

func (b *StorageService) Commit() error {
	return b.tx.Commit()
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
