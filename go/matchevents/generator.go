package matchevents

import (
	"errors"
	"hash/fnv"
	"math"
	"math/big"
	"strconv"
)

func int_hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func GenerateRnd(seed *big.Int, salt string, max_val uint64) uint64 {
	var result uint64 = int_hash(seed.String() + salt)
	if max_val == 0 {
		return result
	}
	return result % max_val
}

// output event order: (minute, eventType, managesToShoot, isGoal, player1, player2)
func addEventsInRound(seed *big.Int, matchEvents []*big.Int, NULL int16) ([][6]int16, []uint64) {
	var events [][6]int16
	nEvents := (len(matchEvents) - 2) / 5
	deltaMinutes := float64(45.0 / (nEvents * 1.0))
	deltaMinutesInt := uint64(math.Floor(deltaMinutes))

	lastMinute := uint64(0)
	var rounds2mins []uint64
	for e := 0; e < nEvents; e++ {
		// compute minute
		salt := "a" + strconv.Itoa(int(e))
		minute := uint64(math.Floor(float64(e)*deltaMinutes)) + GenerateRnd(seed, salt, deltaMinutesInt)
		if minute <= lastMinute {
			minute = lastMinute + 1
		}
		lastMinute = minute
		rounds2mins = append(rounds2mins, minute)
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
			salt := "b" + strconv.Itoa(int(e))
			thisEvent[4] = int16(1 + GenerateRnd(seed, salt, 9))
			thisEvent[5] = NULL
		}
		events = append(events, thisEvent)
	}
	return events, rounds2mins
}

func addCardsAndInjuries(events [][6]int16, seed *big.Int, matchLog [15]uint32, rounds2mins []uint64, NULL int16, NOONE int16) [][6]int16 {
	// matchLog[4,5,6] = outOfGamePlayer, outOfGameType, outOfGameRound
	// note that outofgame is a number from 0 to 13, and that NO OUT OF GAME = 14
	outOfGamePlayer := int16(matchLog[4])
	thereWasAnOutOfGame := outOfGamePlayer < NOONE
	if thereWasAnOutOfGame {
		var typeOfEvent int16
		if matchLog[5] == 1 {
			typeOfEvent = 5
		} else if matchLog[5] == 2 {
			typeOfEvent = 4
		} else if matchLog[5] == 3 {
			typeOfEvent = 3
		}
		minute := int16(rounds2mins[matchLog[6]])
		thisEvent := [6]int16{minute, typeOfEvent, NULL, NULL, outOfGamePlayer, NULL}
		events = append(events, thisEvent)
	}

	yellowCardPlayer := int16(matchLog[7])
	if yellowCardPlayer < 14 {
		maxMinute := int16(45)
		if yellowCardPlayer == outOfGamePlayer {
			maxMinute = outOfGamePlayer
		}
		salt := "c" + strconv.Itoa(int(yellowCardPlayer))
		minute := int16(GenerateRnd(seed, salt, uint64(maxMinute)))
		typeOfEvent := int16(2)
		thisEvent := [6]int16{minute, typeOfEvent, NULL, NULL, yellowCardPlayer, NULL}
		events = append(events, thisEvent)
	}
	yellowCardPlayer = int16(matchLog[8])
	if yellowCardPlayer < 14 {
		maxMinute := int16(45)
		if yellowCardPlayer == outOfGamePlayer {
			maxMinute = outOfGamePlayer
		}
		salt := "d" + strconv.Itoa(int(yellowCardPlayer))
		minute := int16(GenerateRnd(seed, salt, uint64(maxMinute)))
		typeOfEvent := int16(2)
		thisEvent := [6]int16{minute, typeOfEvent, NULL, NULL, yellowCardPlayer, NULL}
		events = append(events, thisEvent)
	}
	return events
}

func addSubstitutions(events [][6]int16, seed *big.Int, matchLog [15]uint32, rounds2mins []uint64, lineup [14]uint8, substitutions [3]uint8, subsRounds [3]uint8, NULL int16) [][6]int16 {
	// matchLog:	9,10,11 ingameSubs, ...0: no change required, 1: change happened, 2: change could not happen
	for i := 0; i < 3; i++ {
		subHappened := matchLog[9+i] == 1
		if subHappened {
			minute := int16(rounds2mins[subsRounds[i]])
			leavingPlayer := int16(substitutions[i])
			enteringPlayer := int16(lineup[11+i])
			typeOfEvent := int16(6)
			thisEvent := [6]int16{minute, typeOfEvent, NULL, NULL, leavingPlayer, enteringPlayer}
			events = append(events, thisEvent)
		}
	}
	return adjustSubstitutions(events)
}

// make sure that if a player that enters via a substitution appears in any other action (goal, pass, cards & injuries),
// then the substitution time must take place at least before that minute.
func adjustSubstitutions(events [][6]int16) [][6]int16 {
	adjustedEvents := events
	for e := 0; e < len(events); e++ {
		if events[e][1] == 6 {
			enteringPlayer := events[e][5]
			enteringMin := events[e][0]
			for e2 := 0; e2 < len(events); e2++ {
				if (e != e2) && (enteringPlayer == events[e2][4]) && (enteringMin >= events[e2][0]-1)  {
					adjustedEvents[e][0] = events[e2][0]-1
				}
			}
		}
	}
	return adjustedEvents
}

// INPUTS:
//	 	seed(the same one used to compute the match)
// 		matchLog decoded: uint32[15] with entries as specified below
// 		matchEvents: uint256[2+5*ROUNDS_PER_MATCH], where ROUNDS = 12 => 62 numbers
// 		bool is2ndHalf
// INPUTS.MATCHLOG:
//  	0 teamSumSkills
//  	1 winner: 0 = home, 1 = away, 2 = draw
//  	2 nGoals
//  	3 trainingPoints
//  	4 uint8 memory outOfGamePlayer
//  	5 uint8 memory typesOutOfGames,
//     		 injuryHard:  1
//     		 injuryLow:   2
//     		 redCard:     3
//  	6 uint8 memory outOfGameRounds
//  	7,8 uint8[2] memory yellowCards
//  	9,10,11 uint8[3] memory ingameSubs, ...0: no change required, 1: change happened, 2: change could not happen
//  	12,13,14 uint8[3] memory halfTimeSubstitutions: 0...10 the player in the starting 11 that was changed during half time
// OUTPUTS:
//		an array of variable size, where each entry is an array of 6 uint16
//			0: minute
// 			1: eventType (0 = team 0 attacks, 1 = team 1 attacks, 2 = yellowCard, 3 = redCard, 4 = injurySoft, 5 = injuryHard, 6 = substitutions)
// 				see: getInGameSubsHappened
// 			2: managesToShoot
// 			3: isGoal
// 			4: primary player (0...11):
// 				(type == 0, 1) && managesToShoot 	: shooter
// 				(type == 0, 1) && !managesToShoot 	: tackler
// 				(type == 2) 						: yellowCarded
// 				(type == 3) 						: redCarded
// 				(type == 4,5) 						: injured
// 				(type == 6) 						: getting out of field
// 			5: secondary player (0...11):
// 				(type == 0, 1) && managesToShoot 	: assister
// 				(type == 0, 1) && !managesToShoot 	: null
// 				(type == 2) 						: null
// 				(type == 3) 						: null
// 				(type == 4,5) 						: null
// 				(type == 6) 						: getting inside field

func GenerateMatchEvents(seed *big.Int, matchLog [15]uint32, matchEvents []*big.Int, lineup [14]uint8, substitutions [3]uint8, subsRounds [3]uint8, is2ndHalf bool) ([][6]int16, error) {
	NULL := int16(-1)
	NOONE := int16(14)
	var emptyEvents [][6]int16

	// check 1:
	if (len(matchEvents)-2)%5 != 0 {
		return emptyEvents, errors.New("the length of matchEvents should be 2 + a multiple of")
	}

	// check 2:
	outOfGamePlayer := int16(matchLog[4])
	thereWasAnOutOfGame := outOfGamePlayer < NOONE
	if thereWasAnOutOfGame && (matchLog[5] > 3 || matchLog[5] == 0) {
		return emptyEvents, errors.New("received an incorrect matchlog entry for matchLog")
	}

	// Compute main events: per-round, and cards & injuries
	events, rounds2mins := addEventsInRound(seed, matchEvents, NULL)
	events = addCardsAndInjuries(events, seed, matchLog, rounds2mins, NULL, NOONE)
	events = addSubstitutions(events, seed, matchLog, rounds2mins, lineup, substitutions, subsRounds, NULL)

	return events, nil
}
