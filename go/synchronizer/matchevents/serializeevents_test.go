package matchevents_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts/router"
	"gotest.tools/assert"
)

func TestSerializationEvents(t *testing.T) {
	eventsLog := big.NewInt(0)
	round := uint(3)

	teamThatAttacks := uint(1)
	eventsLog, _ = router.SetTeamThatAttacks(eventsLog, round, teamThatAttacks)
	newVal, _ := router.GetTeamThatAttacks(eventsLog, round)
	assert.Equal(t, newVal, teamThatAttacks)

	shooter := uint(14)
	eventsLog, _ = router.SetShooter(eventsLog, round, shooter)
	newVal, _ = router.GetShooter(eventsLog, round)
	assert.Equal(t, newVal, shooter)

}
