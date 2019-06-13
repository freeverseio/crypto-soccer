package storage

import (
	"testing"
)

func TestInit(t *testing.T) {
	if db != nil {
		t.Error("db variable is not nil")
	}

	url := "postgres://freeverse:freeverse@localhost/cryptosoccer"
	err := Init(url)
	if err != nil {
		t.Error(err)
	}

	if db == nil {
		t.Error("db variable is not initialized")
	}
}