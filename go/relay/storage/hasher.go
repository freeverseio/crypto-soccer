package storage

import "crypto/sha256"

func (b *Storage) Hash(verse int) ([]byte, error) {
	h := sha256.New()
	tacticsHash, err := b.hashTactics(verse)
	if err != nil {
		return nil, nil
	}
	h.Write(tacticsHash)
	return h.Sum(nil), nil
}

func (b *Storage) hashTactics(verse int) ([]byte, error) {
	h := sha256.New()
	return h.Sum(nil), nil
}
