package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	_ "github.com/lib/pq"
)

type StorageService struct {
	db *sql.DB
}

type Tx struct {
	tx *sql.Tx
}

func NewStorageService(db *sql.DB) *StorageService {
	return &StorageService{
		db: db,
	}
}

func (b *StorageService) Begin() (storage.Tx, error) {
	var err error
	tx, err := b.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func (b *Tx) Rollback() error {
	return b.tx.Rollback()
}

func (b *Tx) Commit() error {
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
