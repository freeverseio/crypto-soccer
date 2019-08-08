package storage

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type Country struct {
	Id          uint64
	Name        string
	TimezoneUTC int
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

func (b *Storage) CountryAdd(country Country) error {
	log.Infof("[DBMS] Adding country %v", country)
	_, err := b.db.Exec("INSERT INTO countries (id, name, timezoneUTC) VALUES ($1, $2, $3);",
		country.Id,
		country.Name,
		country.TimezoneUTC,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Storage) GetCountry(id uint64) (Country, error) {
	country := Country{}
	rows, err := b.db.Query("SELECT id, name, timezoneUTC FROM countries WHERE (id = $1);", id)
	if err != nil {
		return country, err
	}
	defer rows.Close()
	if !rows.Next() {
		return country, errors.New("Unexistent country")
	}
	err = rows.Scan(
		&country.Id,
		&country.Name,
		&country.TimezoneUTC,
	)
	if err != nil {
		return country, err
	}
	return country, nil
}
