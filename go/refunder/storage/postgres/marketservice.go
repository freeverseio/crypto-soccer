package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/refunder/storage"
	_ "github.com/lib/pq"
)

type MarketService struct {
	db *sql.DB
}

func NewMarketService(db *sql.DB) storage.MarketService {
	return &MarketService{
		db: db,
	}
}

func (b *MarketService) Begin() (storage.MarketTx, error) {
	var err error
	tx, err := b.db.Begin()
	if err != nil {
		return nil, err
	}
	return &MarketTx{tx}, nil
}

type MarketTx struct {
	Tx *sql.Tx
}

func (b MarketTx) Rollback() error {
	return b.Tx.Rollback()
}

func (b MarketTx) Commit() error {
	return b.Tx.Commit()
}
