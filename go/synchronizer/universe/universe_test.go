package universe_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/universe"
	"gotest.tools/assert"
)

func TestUniverseSize(t *testing.T) {
	player := storage.Player{}

	u := universe.Universe{}

	assert.Equal(t, u.Size(), 0)

	assert.Error(t, u.Append(player), "encodedSkills is nil")
	player.EncodedSkills = big.NewInt(656)
	assert.Error(t, u.Append(player), "encodedState is nil")
	player.EncodedState = big.NewInt(56)
	assert.NilError(t, u.Append(player))

	assert.Equal(t, u.Size(), 1)
}

func TestUniverseEmptyHash(t *testing.T) {
	u := universe.Universe{}
	hash, err := u.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(hash[:]), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
}

func TestUniverseAppendOrderDoNotChangeHash(t *testing.T) {
	a := storage.Player{}
	a.PlayerId = big.NewInt(4)
	a.EncodedSkills = big.NewInt(1)
	a.EncodedState = big.NewInt(2)
	b := storage.Player{}
	b.PlayerId = big.NewInt(5)
	b.EncodedSkills = big.NewInt(3)
	b.EncodedState = big.NewInt(4)

	expected := "9f64a747e1b97f131fabb6b447296c9b6f0201e79fb3c5356e6c77e89b6a806a"

	t.Run("AB", func(t *testing.T) {
		u := universe.Universe{}
		assert.NilError(t, u.Append(a))
		assert.NilError(t, u.Append(b))
		hash, err := u.Hash()
		assert.NilError(t, err)
		assert.Equal(t, hex.EncodeToString(hash[:]), expected)
	})

	t.Run("BA", func(t *testing.T) {
		u := universe.Universe{}
		assert.NilError(t, u.Append(b))
		assert.NilError(t, u.Append(a))
		hash, err := u.Hash()
		assert.NilError(t, err)
		assert.Equal(t, hex.EncodeToString(hash[:]), expected)
	})
}
