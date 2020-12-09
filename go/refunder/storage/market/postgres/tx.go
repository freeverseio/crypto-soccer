package postgres

import "database/sql"

type Tx struct {
	Tx *sql.Tx
}

func (b Tx) Rollback() error {
	return b.Tx.Rollback()
}

func (b Tx) Commit() error {
	return b.Tx.Commit()
}
