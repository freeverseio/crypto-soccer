package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/refunder"
	_ "github.com/lib/pq"
)

type UniverseService struct {
	db *sql.DB
}

func NewUniverseService(db *sql.DB) refunder.UniverseService {
	return &UniverseService{
		db: db,
	}
}

func (b *UniverseService) Begin() (refunder.UniverseTx, error) {
	var err error
	tx, err := b.db.Begin()
	if err != nil {
		return nil, err
	}
	return &UniverseTx{tx}, nil
}

type UniverseTx struct {
	Tx *sql.Tx
}

func (b UniverseTx) Rollback() error {
	return b.Tx.Rollback()
}

func (b UniverseTx) Commit() error {
	return b.Tx.Commit()
}
