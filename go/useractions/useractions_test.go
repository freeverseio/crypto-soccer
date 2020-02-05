package useractions_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/useractions"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestMarshal(t *testing.T) {
	var ua useractions.UserActions
	data, err := ua.Marshal()
	assert.NilError(t, err)
	assert.Equal(t, string(data), `{"tactics":[],"trainings":[]}`)
	ua.Tactics = append(ua.Tactics, storage.Tactic{Verse: 3, TeamID: "ciao"})
	ua.Trainings = append(ua.Trainings, storage.Training{Verse: 5, TeamID: "pippo"})
	data, err = ua.Marshal()
	assert.NilError(t, err)
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
	assert.NilError(t, err)
	assert.Equal(t, cif, "QmRo9oYwcfJ8BbYJCZKX3JPv7j6izWi2pqePfNpCVfvmYw")
	training := storage.Training{}
	training.TeamID = "pippo"
	ua.Trainings = append(ua.Trainings, training)
	cif, err = ua.PushToIpfs("localhost:5001")
	assert.NilError(t, err)
	assert.Equal(t, cif, "QmWeiipZSst2SKyaM35W7Gc4oTqcYWVBMSu3BtfpPE6eKy")
	var ua2 useractions.UserActions
	err = ua2.PullFromIpfs("localhost:5001", cif)
	assert.NilError(t, err)
	assert.Assert(t, ua2.Equal(&ua))
}
