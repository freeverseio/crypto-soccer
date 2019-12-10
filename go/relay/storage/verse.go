package storage

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

type Verse struct {
	ID      int
	StartAt time.Time
}

func (b *Storage) GetLastVerse() (*Verse, error) {
	log.Debug("[DBMS] get last verse")
	rows, err := b.tx.Query("SELECT id, start_at FROM verses ORDER BY id DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Unexistent")
	}
	verse := Verse{}
	err = rows.Scan(
		&verse.ID,
		&verse.StartAt,
	)
	return &verse, err
}

func (b *Storage) CloseVerse() error {
	log.Debug("[DBMS] close verse")
	currentVerse, err := b.GetLastVerse()
	if err != nil {
		return err
	}
	_, err = b.tx.Exec("INSERT INTO verses (id) VALUES ($1);", currentVerse.ID+1)
	return err
}
