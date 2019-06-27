package sqlite3

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

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
	file, err := os.Open("../../../postgres/sql/00_schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	script, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	_, err = storage.db.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}
	return &storage, nil
}
