package storage

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

const CurrentVerse = 0

// Verse represents a row from 'public.verses'.
type Verse struct {
	ID      int       `json:"id"`       // id
	StartAt time.Time `json:"start_at"` // start_at
}

func VerseById(tx *sql.Tx, id int) (*Verse, error) {
	log.Debugf("[DBMS] get verse %v", id)
	rows, err := tx.Query("SELECT start_at FROM verses WHERE id=$1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Unexistent verse")
	}
	verse := Verse{}
	verse.ID = id
	err = rows.Scan(
		&verse.StartAt,
	)
	return &verse, err
}

func LastVerse(tx *sql.Tx) (*Verse, error) {
	log.Debug("[DBMS] get last verse")
	rows, err := tx.Query("SELECT id, start_at FROM verses ORDER BY id DESC LIMIT 1")
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

func CloseVerse(tx *sql.Tx) error {
	log.Debug("[DBMS] close verse")
	currentVerse, err := LastVerse(tx)
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO verses (id) VALUES ($1);", currentVerse.ID+1)
	return err
}
