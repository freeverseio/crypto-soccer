package useractions_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/relay/storage"
	sync "github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
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
	ua.Tactics = append(ua.Tactics, storage.Tactic{Verse: 3, TeamID: "ciao"})
	ua.Trainings = append(ua.Trainings, storage.Training{Verse: 5, TeamID: "pippo"})
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

func TestIpfsPushAndPull(t *testing.T) {
	t.Parallel()
	var ua useractions.UserActions
	tactic := storage.Tactic{}
	tactic.TeamID = "ciao"
	ua.Tactics = append(ua.Tactics, tactic)
	cif, err := ua.ToIpfs("localhost:5001")
	assert.NilError(t, err)
	assert.Equal(t, cif, "QmRo9oYwcfJ8BbYJCZKX3JPv7j6izWi2pqePfNpCVfvmYw")
	training := storage.Training{}
	training.TeamID = "pippo"
	ua.Trainings = append(ua.Trainings, training)
	cif, err = ua.ToIpfs("localhost:5001")
	assert.NilError(t, err)
	assert.Equal(t, cif, "QmWeiipZSst2SKyaM35W7Gc4oTqcYWVBMSu3BtfpPE6eKy")
	ua2, err := useractions.NewFromIpfs("localhost:5001", cif)
	assert.NilError(t, err)
	assert.Assert(t, ua2.Equal(&ua))
}

func TestUserActionsPullFromStorageNoUserActions(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	verse := uint64(0)
	timezone := 4
	ua, err := useractions.NewFromStorage(tx, verse, timezone)
	assert.NilError(t, err)
	assert.Equal(t, len(ua.Tactics), 0)
	assert.Equal(t, len(ua.Trainings), 0)
}

func TestUserActionsPullFromStorage(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	verse := uint64(6)
	tz := sync.Timezone{}
	assert.NilError(t, tz.Insert(tx))
	country := sync.Country{}
	assert.NilError(t, country.Insert(tx))
	league := sync.League{}
	assert.NilError(t, league.Insert(tx))
	team := sync.Team{}
	team.TeamID = big.NewInt(0)
	assert.NilError(t, team.Insert(tx))
	timezone := 4
	training := storage.Training{}
	training.TeamID = "0"
	training.Verse = verse
	training.Timezone = timezone
	assert.NilError(t, training.Insert(tx))
	tactic := storage.Tactic{}
	tactic.Verse = verse
	tactic.Timezone = timezone
	tactic.TeamID = "0"
	assert.NilError(t, tactic.Insert(tx))
	ua, err := useractions.NewFromStorage(tx, verse, timezone)
	assert.NilError(t, err)
	assert.Equal(t, len(ua.Tactics), 1)
	assert.Equal(t, len(ua.Trainings), 1)

	training.Verse = verse - 1
	assert.NilError(t, training.Insert(tx))
	tactic.Verse = verse + 1
	assert.NilError(t, tactic.Insert(tx))
	ua, err = useractions.NewFromStorage(tx, verse, timezone)
	assert.NilError(t, err)
	assert.Equal(t, len(ua.Tactics), 1)
	assert.Equal(t, len(ua.Trainings), 1)

	team.TeamID = big.NewInt(43)
	assert.NilError(t, team.Insert(tx))
	training.Verse = verse
	training.Timezone = timezone + 1
	training.TeamID = "43"
	assert.NilError(t, training.Insert(tx))
	tactic.Verse = verse
	tactic.Timezone = timezone + 1
	tactic.TeamID = "43"
	assert.NilError(t, tactic.Insert(tx))
	ua, err = useractions.NewFromStorage(tx, verse, timezone)
	assert.NilError(t, err)
	assert.Equal(t, len(ua.Tactics), 1)
	assert.Equal(t, len(ua.Trainings), 1)

}
