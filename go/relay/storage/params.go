package storage

import (
	"database/sql"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/log"
)

func (b *Storage) GetVerse() (*big.Int, error) {
	log.Debug("[DBMS] GetVerse")
	rows, err := b.db.Query("SELECT value FROM params WHERE name='verse';")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Unexistent verse entry")
	}
	var value sql.NullString
	err = rows.Scan(&value)
	if err != nil {
		return nil, err
	}
	verse, _ := new(big.Int).SetString(value.String, 10)

	return verse, nil
}

func (b *Storage) SetVerse(verse *big.Int) error {
	log.Debug("[DBMS] SetVerse %v", verse.String())
	_, err := b.db.Exec("UPDATE params SET value=$1 WHERE name='verse';", verse.String())
	return err
}
