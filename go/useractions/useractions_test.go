package useractions_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/useractions"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestMarshal(t *testing.T) {
	t.Parallel()
	var ua useractions.UserActions
	data, err := ua.Marshal()
	assert.NilError(t, err)
	assert.Equal(t, string(data), `{"tactics":[],"trainings":[]}`)
	ua.Tactics = append(ua.Tactics, storage.Tactic{TeamID: "ciao"})
	ua.Trainings = append(ua.Trainings, storage.Training{TeamID: "pippo"})
	data, err = ua.Marshal()
	assert.NilError(t, err)
	var out bytes.Buffer
	json.Indent(&out, data, "", "\t")
	golden.Assert(t, out.String(), t.Name()+".golden")
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
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

func TestUserActionsPullFromStorageNoUserActions(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	timezone := 4
	ua, err := useractions.NewFromStorage(tx, timezone)
	assert.NilError(t, err)
	assert.Equal(t, len(ua.Tactics), 0)
	assert.Equal(t, len(ua.Trainings), 0)
}

func TestUserActionsPullFromStorage(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	training := storage.NewTraining()
	training.TeamID = teamID
	assert.NilError(t, training.Insert(tx))
	tactic := storage.Tactic{}
	tactic.TeamID = teamID
	assert.NilError(t, tactic.Insert(tx))
	ua, err := useractions.NewFromStorage(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, len(ua.Tactics), 1)
	assert.Equal(t, len(ua.Trainings), 1)
}
