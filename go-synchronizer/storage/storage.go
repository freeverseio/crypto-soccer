package storage

import (
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

func Init(url string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}

    if err = db.Ping(); err != nil {
        return err;
	}

	return nil
}