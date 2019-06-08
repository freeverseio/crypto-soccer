package storage

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Log("Creating a storage")

	storage, err := New("./test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll("./test")

	count := 0
	iter := storage.db.NewIterator(nil, nil)
	for iter.Next() {
		count++
	}

	if count != 0 {
		t.Error("New storage is not empty")
	}
	iter.Release()
	err = iter.Error()	
	if (err != nil) {
		t.Error("iter Error: ", err)
	}
}