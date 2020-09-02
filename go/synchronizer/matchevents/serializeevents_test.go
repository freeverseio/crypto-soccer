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
	eventsLog, _ = router.SetTeamThatAttacks(eventsLog, round, teamThatAttacks)
	newVal, _ := router.GetTeamThatAttacks(eventsLog, round)
	assert.Equal(t, newVal, teamThatAttacks)
}
