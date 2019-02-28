package storage

import (
	"github.com/ethereum/go-ethereum/rlp"
)

// Globals gets the globals from the db.
func (s *Storage) Globals() (*GlobalsEntry, error) {

	key := []byte(prefixGlobals)
	value, err := s.db.Get(key, nil)

	var entry GlobalsEntry
	err = rlp.DecodeBytes(value, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// SetGlobals in the storage.
func (s *Storage) SetGlobals(globals GlobalsEntry) error {

	var err error

	gkey := []byte(prefixGlobals)
	var gvalue []byte

	if gvalue, err = rlp.EncodeToBytes(globals); err != nil {
		return err
	}

	return s.db.Put(gkey, gvalue, nil)
}
