package matchevents_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts/router"
	"gotest.tools/assert"
)

func bool2int64(b bool) int64 {
	if b {
		return int64(1)
	}
	return int64(0)
}

func TestSerializationEvents(t *testing.T) {
	eventsLog := big.NewInt(0)
	round := uint(3)

	teamThatAttacks := uint(1)
	eventsLog, _ = router.SetTeamThatAttacks(eventsLog, round, teamThatAttacks)
	newVal, _ := router.GetTeamThatAttacks(eventsLog, round)
	assert.Equal(t, newVal, teamThatAttacks)

	managesToShoot := true
	eventsLog, _ = router.SetManagesToShoot(eventsLog, round, managesToShoot)
	newBool, _ := router.GetManagesToShoot(eventsLog, round)
	assert.Equal(t, newBool, managesToShoot)

	shooter := uint(14)
	eventsLog, _ = router.SetShooter(eventsLog, round, shooter)
	newVal, _ = router.GetShooter(eventsLog, round)
	assert.Equal(t, newVal, shooter)

	isGoal := true
	eventsLog, _ = router.SetIsGoal(eventsLog, round, isGoal)
	newBool, _ = router.GetIsGoal(eventsLog, round)
	assert.Equal(t, newBool, isGoal)

	assister := uint(13)
	eventsLog, _ = router.SetAssister(eventsLog, round, assister)
	newVal, _ = router.GetAssister(eventsLog, round)
	assert.Equal(t, newVal, assister)
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
		managesToShoot = append(managesToShoot, r%2 == 1)
		shooter = append(shooter, 14-r%5)
		isGoal = append(isGoal, r%2 == 0)
		assister = append(assister, 14-r%4)
	}
	eventsLog := big.NewInt(0)
	for r := uint(0); r < N_ROUNDS; r++ {
		eventsLog, _ = router.SetTeamThatAttacks(eventsLog, r, teamThatAttacks[r])
		eventsLog, _ = router.SetManagesToShoot(eventsLog, r, managesToShoot[r])
		eventsLog, _ = router.SetShooter(eventsLog, r, shooter[r])
		eventsLog, _ = router.SetIsGoal(eventsLog, r, isGoal[r])
		eventsLog, _ = router.SetAssister(eventsLog, r, assister[r])
	}
	assert.Equal(t, eventsLog.String(), "4666666260639135463621236164240476920764")

	for r := uint(0); r < N_ROUNDS; r++ {
		val, _ := router.GetTeamThatAttacks(eventsLog, r)
		assert.Equal(t, val, teamThatAttacks[r])
		val2, _ := router.GetManagesToShoot(eventsLog, r)
		assert.Equal(t, val2, managesToShoot[r])
		val, _ = router.GetShooter(eventsLog, r)
		assert.Equal(t, val, shooter[r])
		val2, _ = router.GetIsGoal(eventsLog, r)
		assert.Equal(t, val2, isGoal[r])
		val, _ = router.GetAssister(eventsLog, r)
		assert.Equal(t, val, assister[r])
	}

	eventsLog2, _ := router.EncodeMatchEvents(teamThatAttacks, shooter, assister, isGoal, managesToShoot)
	assert.Equal(t, eventsLog2.Cmp(eventsLog), 0)

	teamThatAttacks2, managesToShoot2, shooter2, isGoal2, assister2, err2 := router.DecodeMatchEvents(eventsLog, N_ROUNDS)
	assert.Equal(t, err2, nil)
	assert.Equal(t, len(teamThatAttacks2), N_ROUNDS)
	for r := uint(0); r < N_ROUNDS; r++ {
		assert.Equal(t, teamThatAttacks2[r], teamThatAttacks[r])
		assert.Equal(t, managesToShoot2[r], managesToShoot[r])
		assert.Equal(t, shooter2[r], shooter[r])
		assert.Equal(t, isGoal2[r], isGoal[r])
		assert.Equal(t, assister2[r], assister[r])
	}

	var eventsFromPlayHalf []*big.Int
	for r := uint(0); r < N_ROUNDS; r++ {
		eventsFromPlayHalf = append(eventsFromPlayHalf, big.NewInt(int64(teamThatAttacks[r])))
		eventsFromPlayHalf = append(eventsFromPlayHalf, big.NewInt(bool2int64(managesToShoot2[r])))
		eventsFromPlayHalf = append(eventsFromPlayHalf, big.NewInt(int64(shooter[r])))
		eventsFromPlayHalf = append(eventsFromPlayHalf, big.NewInt(bool2int64(isGoal[r])))
		eventsFromPlayHalf = append(eventsFromPlayHalf, big.NewInt(int64(assister[r])))
	}
	serializedEvents, err := router.SerializeEventsFromPlayHalf(eventsFromPlayHalf)
	assert.Equal(t, err, nil)
	assert.Equal(t, eventsLog.Cmp(serializedEvents), 0)
}
