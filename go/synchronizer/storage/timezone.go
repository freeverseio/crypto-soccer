package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type Timezone struct {
	TimezoneIdx uint8
}

func TimezoneCount(tx *sql.Tx) (uint64, error) {
	rows, err := tx.Query("SELECT COUNT(*) FROM timezones;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	rows.Scan(&count)
	return count, nil
}

func (b *Timezone) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Create timezone %v", b)
	_, err := tx.Exec("INSERT INTO timezones (timezone_idx) VALUES ($1);",
		b.TimezoneIdx,
	)
	return err
}
