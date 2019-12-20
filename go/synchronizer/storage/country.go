package storage

import (
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"
)

type Country struct {
	TimezoneIdx uint8
	CountryIdx  uint32
}

func CountryCount(tx *sql.Tx) (uint32, error) {
	rows, err := tx.Query("SELECT COUNT(*) FROM countries;")
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

func CountryInTimezoneCount(tx *sql.Tx, timezoneIdx uint8) (uint32, error) {
	rows, err := tx.Query("SELECT COUNT(*) FROM countries WHERE timezone_idx = $1;", timezoneIdx)
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

func (b *Country) CountryCreate(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create country %v", b)
	_, err := tx.Exec("INSERT INTO countries (timezone_idx, country_idx) VALUES ($1, $2);",
		b.TimezoneIdx,
		b.CountryIdx,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetCountry(tx *sql.Tx, timezone_id uint8, idx uint32) (Country, error) {
	country := Country{}
	rows, err := tx.Query("SELECT timezone_idx, country_idx FROM countries WHERE (timezone_idx = $1 AND country_idx = $2);", timezone_id, idx)
	if err != nil {
		return country, err
	}
	defer rows.Close()
	if !rows.Next() {
		return country, errors.New("Unexistent country")
	}
	err = rows.Scan(
		&country.TimezoneIdx,
		&country.CountryIdx,
	)
	if err != nil {
		return country, err
	}
	return country, nil
}
