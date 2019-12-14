package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
)

func TestMarshal(t *testing.T) {
	var ua storage.UserActions
	data, err := ua.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"tactics":null,"training":null}` {
		t.Fatalf("Wrong %v", string(data))
	}
}

func TestUnmarshal(t *testing.T) {
	var ua storage.UserActions
	err := ua.Unmarshal([]byte(`{"tactics":null,"training":null}`))
	if err != nil {
		t.Fatal(err)
	}
	if ua.Tactics != nil {
		t.Fatal("Tactics not empty")
	}
	if ua.Training != nil {
		t.Fatal("Training not empty")
	}
}
