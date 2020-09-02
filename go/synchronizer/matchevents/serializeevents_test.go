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

	assister := uint(13)
	eventsLog, _ = router.SetAssister(eventsLog, round, assister)
	newVal, _ = router.GetAssister(eventsLog, round)
	assert.Equal(t, newVal, assister)

	isGoal := true
	eventsLog, _ = router.SetIsGoal(eventsLog, round, isGoal)
	newBool, _ := router.GetIsGoal(eventsLog, round)
	assert.Equal(t, newBool, isGoal)

	managesToShoot := true
	eventsLog, _ = router.SetManagesToShoot(eventsLog, round, isGoal)
	newBool, _ = router.GetManagesToShoot(eventsLog, round)
	assert.Equal(t, newBool, managesToShoot)

}
