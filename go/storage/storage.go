package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const MAX_PARAMS_IN_PG_STMT = 65535

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

type StorageService interface {
	TeamService() TeamStorageService
	MatchService() MatchStorageService
}
