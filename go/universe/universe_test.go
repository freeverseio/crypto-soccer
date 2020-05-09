package universe_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/universe"
	"gotest.tools/assert"
)

func TestUniverseSize(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)
	blockNumber := uint64(0)
	player := storage.Player{}
	player.TeamId = teamID
	assert.NilError(t, player.Insert(tx, blockNumber))

	u, err := universe.NewFromStorage(tx, int(timezoneIdx)+1)
	assert.NilError(t, err)
	assert.Equal(t, u.Size(), 0)
	hash, err := u.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(hash[:]), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")

	u, err = universe.NewFromStorage(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, u.Size(), 1)
	hash, err = u.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(hash[:]), "652c24123076f570589e89c4325f113f2c95730c14106b49865a4c4fd286d1cf")
}

func TestUniverseTimezoneEncodedSkills(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)
	blockNumber := uint64(0)
	player := storage.Player{}
	player.PlayerId = big.NewInt(1)
	player.EncodedSkills = big.NewInt(5)
	player.TeamId = teamID
	assert.NilError(t, player.Insert(tx, blockNumber))

	universe, err := universe.NewFromStorage(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, universe.Size(), 1)
	hash, err := universe.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(hash[:]), "1bad926e42913bc44be162d2d426ad416a77342347fc036c758f297abbf58740")
}

func TestUniverseTimezoneEncodedState(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)
	blockNumber := uint64(0)
	player := storage.Player{}
	player.PlayerId = big.NewInt(1)
	player.EncodedState = big.NewInt(5)
	player.TeamId = teamID
	assert.NilError(t, player.Insert(tx, blockNumber))

	universe, err := universe.NewFromStorage(tx, int(timezoneIdx))
	assert.NilError(t, err)
	assert.Equal(t, universe.Size(), 1)
	hash, err := universe.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(hash[:]), "d87c2b8160619ad83073b161249233c293d5afb601476b97255eb7b9f28e465b")
}
