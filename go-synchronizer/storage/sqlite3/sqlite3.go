package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite3 struct {
	db *sql.DB
}

func New() (*Sqlite3, error) {
	var err error
	storage := Sqlite3{}
	storage.db, err = sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		return nil, err
	}
	return &storage, nil
}
