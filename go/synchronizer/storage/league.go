package storage

import (
	log "github.com/sirupsen/logrus"
)

type League struct {
	TimezoneIdx uint8
	CountryIdx  uint32
	LeagueIdx   uint32
}

func (b *Storage) LeagueCount() (uint32, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM leagues;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint32
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *Storage) LeagueInCountryCount(timezoneIdx uint8, countryIdx uint32) (uint32, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM leagues WHERE (timezone_idx = $1 AND country_idx = $2);", timezoneIdx, countryIdx)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint32
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *Storage) LeagueCreate(league League) error {
	log.Debugf("[DBMS] Create league %v", league)
	_, err := b.db.Exec("INSERT INTO leagues (timezone_idx, country_idx, league_idx) VALUES ($1, $2, $3);",
		league.TimezoneIdx,
		league.CountryIdx,
		league.LeagueIdx,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Storage) GetLeague(id uint32) (*League, error) {
	rows, err := b.db.Query("SELECT timezone_idx, country_idx, league_idx FROM leagues WHERE (league_idx = $1);", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var league League
	err = rows.Scan(
		&league.TimezoneIdx,
		&league.CountryIdx,
		&league.LeagueIdx,
	)
	if err != nil {
		return nil, err
	}
	return &league, nil
}
