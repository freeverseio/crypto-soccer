package sqlite3

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite3 struct {
	db *sql.DB
}

func New() (*Sqlite3, error) {
	var err error
	storage := Sqlite3{}
	storage.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	if err := storage.db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}
	return &storage, nil
}
