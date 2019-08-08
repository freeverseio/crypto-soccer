package storage

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type League struct {
	Id          uint64
	Name        string
	TimezoneUTC int
}

func (b *Storage) LeagueCount() (uint64, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM leagues;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *Storage) LeagueAdd(league League) error {
	log.Infof("[DBMS] Adding league %v", league)
	_, err := b.db.Exec("INSERT INTO leagues (id, name, timezoneUTC) VALUES ($1, $2, $3);",
		league.Id,
		league.Name,
		league.TimezoneUTC,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Storage) GetLeague(id uint64) (League, error) {
	league := League{}
	rows, err := b.db.Query("SELECT id, name, timezoneUTC FROM leagues WHERE (id = $1);", id)
	if err != nil {
		return league, err
	}
	defer rows.Close()
	if !rows.Next() {
		return league, errors.New("Unexistent league")
	}
	err = rows.Scan(
		&league.Id,
		&league.Name,
		&league.TimezoneUTC,
	)
	if err != nil {
		return league, err
	}
	return league, nil
}
