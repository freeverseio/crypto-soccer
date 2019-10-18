package storage

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/rlp"
)

func uint642bytes(v uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, v)
	b := buf[:n]
	return b
}

func bytes2uint64(v []byte) uint64 {
	x, n := binary.Uvarint(v)
	if n <= 0 {
		panic("bytes2uint64")
	}
	return x
}

func (s *Storage) UserActions(leagueIdx uint64) (*UserActionsEntry, error) {
	key := []byte(prefixUserAction)
	key = append(key, uint642bytes(leagueIdx)...)

	value, err := s.db.Get(key, nil)

	var entry UserActionsEntry
	err = rlp.DecodeBytes(value, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *Storage) PutUserActions(leagueIdx uint64, entry *UserActionsEntry) error {

	var err error

	key := []byte(prefixUserAction)
	key = append(key, uint642bytes(leagueIdx)...)

	var value []byte
	if value, err = rlp.EncodeToBytes(entry); err != nil {
		return err
	}

	return s.db.Put(key, value, nil)
}
