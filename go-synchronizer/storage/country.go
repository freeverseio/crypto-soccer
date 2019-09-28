package storage

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type Country struct {
	TimezoneIdx uint8
	CountryIdx  uint32
}

func (b *Storage) CountryCount() (uint64, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM countries;")
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

func (b *Storage) CountryCreate(country Country) error {
	log.Debugf("[DBMS] Create country %v", country)
	_, err := b.db.Exec("INSERT INTO countries (timezone_idx, country_idx) VALUES ($1, $2);",
		country.TimezoneIdx,
		country.CountryIdx,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Storage) GetCountry(timezone_id uint8, idx uint16) (Country, error) {
	country := Country{}
	rows, err := b.db.Query("SELECT timezone_idx, country_idx FROM countries WHERE (timezone_idx = $1 AND country_idx = $2);", timezone_id, idx)
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
