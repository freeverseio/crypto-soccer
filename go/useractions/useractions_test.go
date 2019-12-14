package useractions_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/useractions"
)

func TestMarshal(t *testing.T) {
	var ua useractions.UserActions
	data, err := ua.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"tactics":[],"training":[]}` {
		t.Fatalf("Wrong %v", string(data))
	}
}

func TestUnmarshal(t *testing.T) {
	var ua useractions.UserActions
	err := ua.Unmarshal([]byte(`{"tactics":[],"training":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(ua.Tactics) != 0 {
		t.Fatal("Tactics not empty")
	}
	if len(ua.Training) != 0 {
		t.Fatal("Training not empty")
	}
}

// func TestNewUserActionsByVerse(t *testing.T) {
// 	ua, err := storage.
// }
