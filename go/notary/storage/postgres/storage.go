package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Storage struct {
	db *sql.DB
}

func New(url string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Info("[DBMS] ... connected")
	return db, nil
}
