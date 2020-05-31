package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type Verse struct {
	VerseNumber int64
	Root        string
}

func VerseByNumber(tx *sql.Tx, verseNumber int64) (*Verse, error) {
	log.Debugf("[DBMS] VerseByNumber %v", verseNumber)
	rows, err := tx.Query(`SELECT root FROM verses WHERE verse_number = $1`, verseNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	verse := Verse{}
	verse.VerseNumber = verseNumber
	if err := rows.Scan(&verse.Root); err != nil {
		return nil, err
	}

	return &verse, nil
}

func (b Verse) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] Verse Insert %v", b)

	if _, err := tx.Exec(`
		INSERT INTO verses (
			verse_number,
			root
		) VALUES ($1, $2);`,
		b.VerseNumber,
		b.Root,
	); err != nil {
		return err
	}
	return nil
}
