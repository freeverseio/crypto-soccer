package storage

import (
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewPostgres(url string) (*Storage, error) {
	var err error
	storage := &Storage{}
	storage.db, err = sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func NewSqlite3(schemaFile string) (*Storage, error) {
	var err error
	storage := Storage{}
	storage.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	if err := storage.db.Ping(); err != nil {
		return nil, err
	}
	file, err := os.Open(schemaFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	script, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	_, err = storage.db.Exec(string(script))
	if err != nil {
		return nil, err
	}
	return &storage, nil
}
