package storage

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	prefixGlobals    = "G"
	prefixStakers    = "S"
	prefixUserAction = "U"
)

// Storage manages the application state
type Storage struct {
	mutex *sync.Mutex
	db    *leveldb.tx
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
