package useractions_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/useractions"
)

func TestMarshal(t *testing.T) {
	var ua useractions.UserActions
	data, err := ua.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"tactics":[],"trainings":[]}` {
		t.Fatalf("Wrong %v", string(data))
	}
}

func TestUnmarshal(t *testing.T) {
	var ua useractions.UserActions
	err := ua.Unmarshal([]byte(`{"tactics":[],"trainings":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(ua.Tactics) != 0 {
		t.Fatal("Tactics not empty")
	}
	if len(ua.Trainings) != 0 {
		t.Fatal("Training not empty")
	}
}

func TestIpfsPushAndPull(t *testing.T) {
	var ua useractions.UserActions
	tactic := storage.Tactic{TeamID: "ciao"}
	training := storage.Training{TeamID: "pippo"}
	ua.Tactics = append(ua.Tactics, &tactic)
	ua.Trainings = append(ua.Trainings, &training)
	cif, err := ua.IpfsPush("localhost:5001")
	if err != nil {
		t.Fatal(err)
	}
	if cif != "QmPuMad3rkusEKVbUrPAaY2wbKKfQ8qT1Gtd7apcc8gcsx" {
		t.Fatalf("Wrong cif %v", cif)
	}
	var ua2 useractions.UserActions
	err = ua2.IpfsPull("localhost:5001", cif)
	if err != nil {
		t.Fatal(err)
	}
	if ua2.Trainings[0].TeamID != training.TeamID {
		t.Fatal("Differents")
	}
	if ua2.Tactics[0].TeamID != tactic.TeamID {
		t.Fatal("Differents")
	}
}
