package names

import (
	"errors"
	"hash/fnv"
	"math"
	"math/big"

	_ "github.com/mattn/go-sqlite3"
)

func int_hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func GenerateRnd(seed *big.Int, max_val uint64, nLayers int) uint64 {
	var iterated_seed uint64
	iterated_seed = int_hash(seed.String())
	for i := 1; i < nLayers; i++ {
		iterated_seed = int_hash(big.NewInt(int64(iterated_seed)).String())
	}
	if max_val == 0 {
		return iterated_seed
	} else {
		return iterated_seed % max_val
	}
}

// inputs:
//	 	seed(the same one used to compute the match)
// 		matchEvents: uint256[2+5*ROUNDS_PER_MATCH], where ROUNDS = 12 => 62 numbers
// outputs:
//		an array of variable size, where each entry is an array of 6 uint16
//			0: minute
// 			1: eventType (0 = team 0 attacks, 1 = team 1 attacks, 2 = redCard, 3 = yellowCard, 4 = injury, 5 = substitutions)
// 				see: getInGameSubsHappened
// 			2: managesToShoot
// 			3: isGoal
// 			4: primary player (0...11):
// 				(type == 0, 1) && managesToShoot 	: shooter
// 				(type == 0, 1) && !managesToShoot 	: tackler
// 				(type == 2) 						: redCarded
// 				(type == 3) 						: yellowCarded
// 				(type == 4) 						: injured
// 				(type == 5) 						: getting out of field
// 			5: secondary player (0...11):
// 				(type == 0, 1) && managesToShoot 	: assister
// 				(type == 0, 1) && !managesToShoot 	: null
// 				(type == 2) 						: null
// 				(type == 3) 						: null
// 				(type == 4) 						: null
// 				(type == 5) 						: getting inside field
func GenerateMatchEvents(seed *big.Int, matchEvents []*big.Int) ([][6]int16, error) {
	NULL := int16(-1)
	var events [][6]int16
	// matchLog0 := matchEvents[0]
	// matchLog1 := matchEvents[1]
	if (len(matchEvents)-2)%5 != 0 {
		return events, errors.New("the length of matchEvents should be 2 + a multiple of")
	}
	nEvents := (len(matchEvents) - 2) / 5
	deltaMinutes := float64(45.0 / (nEvents * 1.0))
	deltaMinutesInt := uint64(math.Floor(deltaMinutes))

	lastMinute := uint64(0)
	for e := 2; e < nEvents; e++ {
		// compute minute
		minute := uint64(math.Floor(float64(e)*deltaMinutes)) + GenerateRnd(seed, deltaMinutesInt, 1)
		if minute <= lastMinute {
			minute = lastMinute + 1
		}
		lastMinute = minute
		// parse type of event and data
		teamThatAttacks := matchEvents[2+5*e]
		managesToShoot := matchEvents[2+5*e+1]
		shooter := matchEvents[2+5*e+2]
		isGoal := matchEvents[2+5*e+3]
		assister := matchEvents[2+5*e+4]
		var thisEvent [6]int16
		thisEvent[0] = int16(minute)
		thisEvent[1] = int16(teamThatAttacks.Int64())
		thisEvent[2] = int16(managesToShoot.Int64())
		thisEvent[3] = int16(isGoal.Int64())
		if managesToShoot.Int64() == 1 {
			thisEvent[4] = int16(shooter.Int64())
			thisEvent[5] = int16(assister.Int64())
		} else {
			thisEvent[4] = int16(1 + GenerateRnd(seed, 9, 2))
			thisEvent[5] = NULL
		}
		events = append(events, thisEvent)
	}
	return events, nil
}
