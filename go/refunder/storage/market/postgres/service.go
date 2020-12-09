package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/refunder/storage/universe"
	_ "github.com/lib/pq"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) universe.Service {
	return &Service{
		db: db,
	}
}

func (b *Service) Begin() (universe.Tx, error) {
	var err error
	tx, err := b.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
