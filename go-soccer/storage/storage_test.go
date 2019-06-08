package storage

import "testing"

func TestNew(t *testing.T) {
	t.Log("Creating a storage")

	path := "./"
	storage, err := New(path)
	if err != nil {
		t.Error(err)
	}

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