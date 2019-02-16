package storage

import (
	"errors"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	prefixGlobals   = "G"
	prefixSavepoint = "S"
)

var (
	ErrKeyNotExists     = errors.New("key does not exists in db")
	ErrInconsistentSize = errors.New("inconsistent db and IPFS datasize")
)

// Storage manages the application state
type Storage struct {
	mutex *sync.Mutex
	db    *leveldb.DB
}

// New creates a new storage path.
func New(path string) (*Storage, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db:    db,
		mutex: &sync.Mutex{},
	}, nil
}
