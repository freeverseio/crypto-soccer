package relay

// TODO: use go-soccer/storage

import (
	"errors"
	"time"
)

var db = make(map[string]*UserEntry)

// Action - ...
type Action struct {
	Type  interface{}
	Value interface{}
}

// UserEntry - ...
type UserEntry struct {
	ID     int64
	Nonce  uint64
	Action Action
}

// AddUserEntry - adds user to db
func AddUserEntry(account string) error {
	_, ok := db[account]
	if ok {
		return errors.New("User already exist")
	}
	db[account] = &UserEntry{ID: time.Now().Unix(), Nonce: 0}
	return nil
}

// GetUserEntry - adds user to db
func GetUserEntry(account string) *UserEntry {
	if entry, ok := db[account]; ok == true {
		return entry
	}
	return nil
}
