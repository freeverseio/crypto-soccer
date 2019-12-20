package storage

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Storage struct {
	db *sql.DB
	tx *sql.Tx
}

func (b *Storage) Begin() error {
	var err error
	b.tx, err = b.db.Begin()
	return err
}

func (b *Storage) Commit() error {
	return b.tx.Commit()
}

func (b *Storage) Rollback() error {
	return b.tx.Rollback()
}

func NewPostgres(url string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	for db.Ping() != nil {
		const pause = 5
		log.Errorf("[DBMS] Failed to connect to DBMS: %v", url)
		log.Infof("[DBMS] wainting %v sec ...", pause)
		time.Sleep(pause * time.Second)
	}
	return db, nil
}
