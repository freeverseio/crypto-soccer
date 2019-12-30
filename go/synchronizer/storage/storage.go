package storage

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func New(url string) (*sql.DB, error) {
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
	log.Info("[DBMS] ... connected")
	return db, nil
}
