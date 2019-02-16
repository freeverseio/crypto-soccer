package storage

import (
	"github.com/ethereum/go-ethereum/rlp"
)

func (s *Storage) SavePoint() (*SavePointEntry, error) {

	key := []byte(prefixSavepoint)
	value, err := s.db.Get(key, nil)

	var entry SavePointEntry
	err = rlp.DecodeBytes(value, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *Storage) SetSavePoint(savePoint SavePointEntry) error {

	var err error

	gkey := []byte(prefixGlobals)
	var gvalue []byte

	if gvalue, err = rlp.EncodeToBytes(savePoint); err != nil {
		return err
	}

	return s.db.Put(gkey, gvalue, nil)
}
