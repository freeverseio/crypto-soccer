package storage

import (
	"crypto/sha256"
	"errors"
)

func (b *Storage) HashVerse(id int) ([]byte, error) {
	h := sha256.New()
	tacticsHash, err := b.hashVerseTactics(id)
	if err != nil {
		return nil, nil
	}
	h.Write(tacticsHash)
	return h.Sum(nil), nil
}

func (b *Storage) hashVerseTactics(id int) ([]byte, error) {
	verse, err := b.GetVerse(id)
	if err != nil {
		return nil, err
	}
	if verse == nil {
		return nil, errors.New("Unexistent verse")
	}
	prevVerse, err := b.GetVerse(id - 1)
	if err != nil {
		return nil, err
	}
	if prevVerse == nil {
		return nil, errors.New("Unexistent prevVerse")
	}
	rows, err := b.tx.Query("SELECT * FROM TACTICS WHERE (created_at > $1) AND (created_at <= $2)", prevVerse.StartAt, verse.StartAt)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	return h.Sum(nil), nil
}
