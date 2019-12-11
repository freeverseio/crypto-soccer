package storage

import (
	"crypto/sha256"
	"fmt"
)

func (b *Storage) HashVerse(id int) ([]byte, error) {
	h := sha256.New()
	if id <= 0 {
		return h.Sum(nil), nil
	}
	tacticsHash, err := b.hashVerseTactics(id)
	if err != nil {
		return nil, err
	}
	h.Write(tacticsHash)
	return h.Sum(nil), nil
}

func (b *Storage) hashVerseTactics(id int) ([]byte, error) {
	if id <= 0 {
		return nil, fmt.Errorf("Can't hash between verse %v and verse %v", id-1, id)
	}
	verse, err := b.GetVerse(id)
	if err != nil {
		return nil, err
	}
	if verse == nil {
		return nil, fmt.Errorf("Unexistent verse %v", id)
	}
	prevVerse, err := b.GetVerse(id - 1)
	if err != nil {
		return nil, err
	}
	if prevVerse == nil {
		return nil, fmt.Errorf("Unexistent previous verse %v", id-1)
	}
	_, err = b.tx.Query("SELECT * FROM TACTICS WHERE (created_at > $1) AND (created_at <= $2)", prevVerse.StartAt, verse.StartAt)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	return h.Sum(nil), nil
}
