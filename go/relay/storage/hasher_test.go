package storage_test

import (
	"encoding/hex"
	"testing"
)

func TestHash(t *testing.T) {
	err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Rollback()
	hash, err := db.Hash(0)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(hash) != "5df6e0e2761359d30a8275058e299fcc0381534545f55cf43e41983f5d4c9456" {
		t.Fatalf("Wrong result %v", hex.EncodeToString(hash))
	}
}
