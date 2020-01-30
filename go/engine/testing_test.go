package engine_test

import (
	"fmt"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/engine"
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
			player := engine.CreateDummyPlayer(t, *bc.Contracts, ts.Age, ts.Defence, ts.Speed, ts.Endurance, ts.Pass, ts.Shoot)
			assert.Equal(t, player.Defence(), ts.Defence)
			assert.Equal(t, player.Speed(), ts.Speed)
			assert.Equal(t, player.Endurance(), ts.Endurance)
			assert.Equal(t, player.Pass(), ts.Pass)
			assert.Equal(t, player.Shoot(), ts.Shoot)
			assert.Equal(t, engine.PlayerAge(player.BirthDayUnix()), ts.Age)
		})
	}
}
