package storage

import (
	"crypto/sha256"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (b *Storage) HashVerse(id int) ([]byte, error) {
	h := sha256.New()
	if id <= 0 {
		return h.Sum(nil), nil
	}
	end, err := b.GetVerse(id)
	if err != nil {
		return nil, err
	}
	if end == nil {
		return nil, fmt.Errorf("Unexistent verse %v", id)
	}
	start, err := b.GetVerse(id - 1)
	if err != nil {
		return nil, err
	}
	if start == nil {
		return nil, fmt.Errorf("Unexistent previous verse %v", id-1)
	}
	tacticsHash, err := b.hashVerseTactics(start, end)
	if err != nil {
		return nil, err
	}
	h.Write(tacticsHash)
	return h.Sum(nil), nil
}

func (b *Storage) hashVerseTactics(start *Verse, end *Verse) ([]byte, error) {
	h := sha256.New()
	rows, err := b.GetRowsTactic(start, end)
	if err != nil {
		return nil, err
	}
	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	readCols := make([]interface{}, len(colNames))
	writeCols := make([]string, len(colNames))
	for i := range writeCols {
		readCols[i] = &writeCols[i]
	}
	for rows.Next() {
		err := rows.Scan(readCols...)
		if err != nil {
			return nil, err
		}
		for _, s := range writeCols {
			log.Infof("here %v", s)
			h.Write([]byte(s))
		}
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	return h.Sum(nil), nil
}
