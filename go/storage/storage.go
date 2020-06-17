package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

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

type StorageDumpService interface {
	Dump(fileName string) error
}

type StorageService interface {
	Team(teamId string) (*Team, error)
	Insert(team Team) error
	UpdateName(teamId string, name string) error
	UpdateManagerName(teamId string, name string) error

	Dump(fileName string) error
}
