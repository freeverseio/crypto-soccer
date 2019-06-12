package storage

import (
	"testing"
)

func TestInit(t *testing.T) {
	url := "ffff"
	err := Init(url)
	if err == nil {
		t.Error("I can connect with " + url)
	}

	url = "postgres://freeverse:freeverse@localhost/cryptosoccer"
	err = Init(url)
	if err != nil {
		t.Error(err)
	}
}