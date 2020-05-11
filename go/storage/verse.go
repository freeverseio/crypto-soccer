package storage

import "database/sql"

type Verse struct {
	VerseNumber int64
	Root        string
}

func VerseByNumber(tx *sql.Tx, verseNumber int64) (*Verse, error) {
	rows, err := tx.Query(`SELECT root FROM verses WHERE verse_number = $1`, verseNumber)
	if err != nil {
		return nil, err
	}

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
