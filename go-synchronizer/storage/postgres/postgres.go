package postgres

import (
	"database/sql"
	"math/big"
)

type PostgresStorage struct {
	db *sql.DB
}

func New(url string) (*PostgresStorage, error) {
	var err error
	storage := &PostgresStorage{}
	storage.db, err = sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (b *PostgresStorage) GetBlockNumber() (*big.Int, error) {
	return big.NewInt(0), nil
}

func (b *PostgresStorage) SetBlockNumber(value *big.Int) error {
	return nil
}
