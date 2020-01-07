package match_test

import (
	"fmt"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"gotest.tools/assert"
)

func TestCreateDummyPlayer(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Age       uint16
		Defence   uint16
		Speed     uint16
		Endurance uint16
		Pass      uint16
		Shoot     uint16
	}{
		{0, 0, 0, 0, 0, 0},
		{21, 40, 32, 22, 44, 66},
		{18, 4000, 3254, 2242, 4432, 6655},
		{50, 4000, 325, 242, 432, 655},
	}

	for _, ts := range cases {
		t.Run(fmt.Sprintf("%v", ts), func(t *testing.T) {
			player := match.CreateDummyPlayer(t, bc.Contracts, ts.Age, ts.Defence, ts.Speed, ts.Endurance, ts.Pass, ts.Shoot)
			value, err := player.Defence(bc.Contracts.Assets)
			assert.NilError(t, err)
			assert.Equal(t, value, ts.Defence)
			value, err = player.Speed(bc.Contracts.Assets)
			assert.NilError(t, err)
			assert.Equal(t, value, ts.Speed)
			value, err = player.Endurance(bc.Contracts.Assets)
			assert.NilError(t, err)
			assert.Equal(t, value, ts.Endurance)
			value, err = player.Pass(bc.Contracts.Assets)
			assert.NilError(t, err)
			assert.Equal(t, value, ts.Pass)
			value, err = player.Shoot(bc.Contracts.Assets)
			assert.NilError(t, err)
			assert.Equal(t, value, ts.Shoot)
			value, err = player.BirthDayUnix(bc.Contracts.Assets)
			assert.NilError(t, err)
			assert.Equal(t, match.PlayerAge(value), ts.Age)
		})
	}
}
