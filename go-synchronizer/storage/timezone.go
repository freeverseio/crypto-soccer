package storage

import log "github.com/sirupsen/logrus"

type Timezone struct {
	TimezoneIdx uint8
}

func (b *Storage) TimezoneCount() (uint64, error) {
	rows, err := b.db.Query("SELECT COUNT(*) FROM timezones;")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var count uint64
	rows.Scan(&count)
	return count, nil
}

func (b *Storage) TimezoneCreate(timezone Timezone) error {
	log.Infof("[DBMS] Adding timezone %v", timezone)
	_, err := b.db.Exec("INSERT INTO timezones (timezone_idx) VALUES ($1);",
		timezone.TimezoneIdx,
	)
	return err
}
