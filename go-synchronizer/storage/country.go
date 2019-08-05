package storage

import log "github.com/sirupsen/logrus"

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
	log.Infof("[DBMS] Adding team %v", country)
	return nil
}
