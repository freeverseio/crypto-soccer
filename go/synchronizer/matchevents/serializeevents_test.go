package matchevents_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts/router"
	"gotest.tools/assert"
)

func Test1(t *testing.T) {
	eventsLog := big.NewInt(0)
	round := uint(3)
	teamThatAttacks := uint(1)
	eventsLog2, err := router.SetTeamThatAttacks(eventsLog, round, teamThatAttacks)
	assert.NilError(t, err)
	assert.Equal(t, eventsLog2, eventsLog)
}
