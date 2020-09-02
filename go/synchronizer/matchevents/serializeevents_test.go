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

func TestSerializationEventsArray(t *testing.T) {
	const N_ROUNDS = 12
	var teamThatAttacks []uint
	var shooter []uint
	var assister []uint
	var isGoal []bool
	var managesToShoot []bool

	for r := uint(0); r < N_ROUNDS; r++ {
		teamThatAttacks = append(teamThatAttacks, r%2)
		shooter = append(shooter, 14-r%3)
		assister = append(assister, 14-r%4)
		isGoal = append(isGoal, r%2 == 1)
		managesToShoot = append(managesToShoot, r%2 == 0)
	}
	eventsLog := big.NewInt(0)
	for r := uint(0); r < N_ROUNDS; r++ {
		eventsLog, _ = router.SetTeamThatAttacks(eventsLog, r, teamThatAttacks[r])
		eventsLog, _ = router.SetShooter(eventsLog, r, shooter[r])
		eventsLog, _ = router.SetAssister(eventsLog, r, assister[r])
		eventsLog, _ = router.SetIsGoal(eventsLog, r, isGoal[r])
		eventsLog, _ = router.SetManagesToShoot(eventsLog, r, managesToShoot[r])
	}
	for r := uint(0); r < N_ROUNDS; r++ {
		val, _ := router.GetTeamThatAttacks(eventsLog, r)
		assert.Equal(t, val, teamThatAttacks[r])
		val, _ = router.GetShooter(eventsLog, r)
		assert.Equal(t, val, shooter[r])
		val, _ = router.GetAssister(eventsLog, r)
		assert.Equal(t, val, assister[r])
		val2, _ := router.GetIsGoal(eventsLog, r)
		assert.Equal(t, val2, isGoal[r])
		val2, _ = router.GetManagesToShoot(eventsLog, r)
		assert.Equal(t, val2, managesToShoot[r])
	}

	eventsLog2, _ := router.EncodeMatchEvents(teamThatAttacks, shooter, assister, isGoal, managesToShoot)
	assert.Equal(t, eventsLog2.Cmp(eventsLog), 0)

}
