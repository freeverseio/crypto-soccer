package useractions_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/useractions"
	"gotest.tools/golden"
)

func TestMarshal(t *testing.T) {
	var ua useractions.UserActions
	data, err := ua.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"verse":0,"tactics":[],"trainings":[]}` {
		t.Fatalf("Wrong %v", string(data))
	}
	ua.Tactics = append(ua.Tactics, storage.Tactic{Verse: 3, TeamID: "ciao"})
	ua.Trainings = append(ua.Trainings, storage.Training{Verse: 5, TeamID: "pippo"})
	data, err = ua.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	json.Indent(&out, data, "", "\t")
	golden.Assert(t, out.String(), t.Name()+".golden")
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
	tactic := storage.Tactic{}
	tactic.TeamID = "ciao"
	ua.Tactics = append(ua.Tactics, tactic)
	cif, err := ua.PushToIpfs("localhost:5001")
	if err != nil {
		t.Fatal(err)
	}
	if cif != "QmYnN4c9nKJijrK2RdRAVfVRjL5FW8ygAiUVWCV6bwDFBv" {
		t.Fatalf("Wrong cif %v", cif)
	}
	training := storage.Training{}
	training.TeamID = "pippo"
	ua.Trainings = append(ua.Trainings, training)
	cif, err = ua.PushToIpfs("localhost:5001")
	if err != nil {
		t.Fatal(err)
	}
	if cif != "QmVLduJ2FboB1yFqMhB6UkkyMee6QWQuu6mNjPPBhXW2iW" {
		t.Fatalf("Wrong cif %v", cif)
	}
	var ua2 useractions.UserActions
	err = ua2.PullFromIpfs("localhost:5001", cif)
	if err != nil {
		t.Fatal(err)
	}
	if !ua2.Equal(&ua) {
		t.Fatalf("Expected %v but %v", ua, ua2)
	}
}
