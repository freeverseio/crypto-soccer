package storage

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
)

type StakerEntry struct {
	HashIndex uint64
}

func (s *Storage) Staker(staker common.Address) (*StakerEntry, error) {
	key := []byte(prefixStakers)
	key = append(key, staker.Bytes()...)

	value, err := s.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	var entry StakerEntry
	err = rlp.DecodeBytes(value, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *Storage) HasStaker(staker common.Address) (bool, error) {

	key := []byte(prefixStakers)
	key = append(key, staker.Bytes()...)
	_, err := s.db.Get(key, nil)

	if err == leveldb.ErrNotFound {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func (s *Storage) SetStaker(staker common.Address, entry *StakerEntry) error {

	var err error

	key := []byte(prefixStakers)
	key = append(key, staker.Bytes()...)

	var value []byte
	if value, err = rlp.EncodeToBytes(entry); err != nil {
		return err
	}

	return s.db.Put(key, value, nil)
}
