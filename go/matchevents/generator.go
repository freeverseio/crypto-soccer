package matchevents

import (
	"errors"
	"hash/fnv"
	"math"
	"math/big"
	"strconv"
)

type Matchevent struct {
	Minute          int16 `json:"minute"`
	Type            int16 `json:"type"`
	Team            int16 `json:"team"`
	ManagesToShoot  int16 `json:"managestoshoot"`
	IsGoal          int16 `json:"isgoal"`
	PrimaryPlayer   int16 `json:"primaryplayer"`
	SecondaryPlayer int16 `json:"secondaryplayer"`
}

// eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
const (
	EVNT_ATTACK int16 = 0
	EVNT_YELLOW int16 = 1
	EVNT_RED    int16 = 2
	EVNT_SOFT   int16 = 3
	EVNT_HARD   int16 = 4
	EVNT_SUBST  int16 = 5
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
// eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
func addEventsInRound(seed *big.Int, blockchainEvents []*big.Int, NULL int16) ([]Matchevent, []uint64) {
	var events []Matchevent
	nEvents := (len(blockchainEvents) - 2) / 5
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
		teamThatAttacks := blockchainEvents[2+5*e]
		managesToShoot := blockchainEvents[2+5*e+1]
		shooter := blockchainEvents[2+5*e+2]
		isGoal := blockchainEvents[2+5*e+3]
		assister := blockchainEvents[2+5*e+4]
		var thisEvent Matchevent
		thisEvent.Minute = int16(minute)
		thisEvent.Type = int16(EVNT_ATTACK)
		thisEvent.Team = int16(teamThatAttacks.Int64())
		thisEvent.ManagesToShoot = int16(managesToShoot.Int64())
		thisEvent.IsGoal = int16(isGoal.Int64())
		if managesToShoot.Int64() == 1 {
			thisEvent.PrimaryPlayer = int16(shooter.Int64())
			thisEvent.SecondaryPlayer = int16(assister.Int64())
		} else {
			salt := "b" + strconv.Itoa(int(e))
			thisEvent.PrimaryPlayer = int16(1 + GenerateRnd(seed, salt, 9))
			thisEvent.SecondaryPlayer = NULL
		}
		events = append(events, thisEvent)
	}
	return events, rounds2mins
}

func addCardsAndInjuries(team int16, events []Matchevent, seed *big.Int, matchLog [15]uint32, rounds2mins []uint64, NULL int16, NOONE int16) []Matchevent {
	// matchLog[4,5,6] = outOfGamePlayer, outOfGameType, outOfGameRound
	// note that outofgame is a number from 0 to 13, and that NO OUT OF GAME = 14
	// eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
	outOfGamePlayer := int16(matchLog[4])
	thereWasAnOutOfGame := outOfGamePlayer < NOONE
	if thereWasAnOutOfGame {
		var typeOfEvent int16
		if matchLog[5] == 1 {
			typeOfEvent = EVNT_HARD
		} else if matchLog[5] == 2 {
			typeOfEvent = EVNT_SOFT
		} else if matchLog[5] == 3 {
			typeOfEvent = EVNT_RED
		}
		minute := int16(rounds2mins[matchLog[6]])
		thisEvent := Matchevent{minute, typeOfEvent, team, NULL, NULL, outOfGamePlayer, NULL}
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
		typeOfEvent := EVNT_YELLOW
		thisEvent := Matchevent{minute, typeOfEvent, team, NULL, NULL, yellowCardPlayer, NULL}
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
		typeOfEvent := EVNT_YELLOW
		thisEvent := Matchevent{minute, typeOfEvent, team, NULL, NULL, yellowCardPlayer, NULL}
		events = append(events, thisEvent)
	}
	return events
}

func addSubstitutions(team int16, events []Matchevent, seed *big.Int, matchLog [15]uint32, rounds2mins []uint64, lineup [14]uint8, substitutions [3]uint8, subsRounds [3]uint8, NULL int16) []Matchevent {
	// matchLog:	9,10,11 ingameSubs, ...0: no change required, 1: change happened, 2: change could not happen
	for i := 0; i < 3; i++ {
		subHappened := matchLog[9+i] == 1
		if subHappened {
			minute := int16(rounds2mins[subsRounds[i]])
			leavingPlayer := int16(substitutions[i])
			enteringPlayer := int16(lineup[11+i])
			typeOfEvent := EVNT_SUBST
			thisEvent := Matchevent{minute, typeOfEvent, team, NULL, NULL, leavingPlayer, enteringPlayer}
			events = append(events, thisEvent)
		}
	}
	return adjustSubstitutions(team, events)
}

// make sure that if a player that enters via a substitution appears in any other action (goal, pass, cards & injuries),
// then the substitution time must take place at least before that minute.
func adjustSubstitutions(team int16, events []Matchevent) []Matchevent {
	adjustedEvents := events
	for e := 0; e < len(events); e++ {
		if (events[e].Type == EVNT_SUBST) && (events[e].Team == team) {
			enteringPlayer := events[e].SecondaryPlayer
			enteringMin := events[e].Minute
			for e2 := 0; e2 < len(events); e2++ {
				if (e != e2) && (events[e2].Team == team) && (enteringPlayer == events[e2].PrimaryPlayer) && (enteringMin >= events[e2].Minute-1) {
					adjustedEvents[e].Minute = events[e2].Minute - 1
				}
			}
		}
	}
	return adjustedEvents
}

// INPUTS:
//	 	seed(the same one used to compute the match)
// 		matchLog decoded: uint32[15] with entries as specified below
// 		blockchainEvents: uint256[2+5*ROUNDS_PER_MATCH], where ROUNDS = 12 => 62 numbers
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
//		an array of variable size, where each entry is a Matchevent struct
//			0: minute
// 			1: eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
// 				see: getInGameSubsHappened
// 			2: team: 0, 1
// 			3: managesToShoot
// 			4: isGoal
// 			5: primary player (0...11):
// 				(type == 0, 1) && managesToShoot 	: shooter
// 				(type == 0, 1) && !managesToShoot 	: tackler
// 				(type == 2) 						: yellowCarded
// 				(type == 3) 						: redCarded
// 				(type == 4,5) 						: injured
// 				(type == 6) 						: getting out of field
// 			6: secondary player (0...11):
// 				(type == 0, 1) && managesToShoot 	: assister
// 				(type == 0, 1) && !managesToShoot 	: null
// 				(type == 2) 						: null
// 				(type == 3) 						: null
// 				(type == 4,5) 						: null
// 				(type == 6) 						: getting inside field

func checkOutOfGameData(matchLog [15]uint32, NOONE int16) error {
	outOfGamePlayer := int16(matchLog[4])
	thereWasAnOutOfGame := outOfGamePlayer < NOONE
	if thereWasAnOutOfGame && (matchLog[5] > 3 || matchLog[5] == 0) {
		return errors.New("received an incorrect matchlog entry for matchLog")
	}
	return nil
}

func GenerateMatchEvents(
	seed *big.Int,
	matchLog0 [15]uint32,
	matchLog1 [15]uint32,
	blockchainEvents []*big.Int,
	lineup0 [14]uint8,
	lineup1 [14]uint8,
	substitutions0 [3]uint8,
	substitutions1 [3]uint8,
	subsRounds0 [3]uint8,
	subsRounds1 [3]uint8,
	is2ndHalf bool,
) ([]Matchevent, error) {
	NULL := int16(-1)
	NOONE := int16(14)
	var emptyEvents []Matchevent

	// check 1:
	if (len(blockchainEvents)-2)%5 != 0 {
		return emptyEvents, errors.New("the length of blockchainEvents should be 2 + a multiple of")
	}

	// check 2:
	err := checkOutOfGameData(matchLog0, NOONE)
	if err != nil {
		return emptyEvents, err
	}
	err = checkOutOfGameData(matchLog1, NOONE)
	if err != nil {
		return emptyEvents, err
	}

	// Compute main events: per-round, and cards & injuries
	events, rounds2mins := addEventsInRound(seed, blockchainEvents, NULL)
	events = addCardsAndInjuries(0, events, seed, matchLog0, rounds2mins, NULL, NOONE)
	events = addCardsAndInjuries(1, events, seed, matchLog0, rounds2mins, NULL, NOONE)
	events = addSubstitutions(0, events, seed, matchLog0, rounds2mins, lineup0, substitutions0, subsRounds0, NULL)
	events = addSubstitutions(1, events, seed, matchLog1, rounds2mins, lineup1, substitutions1, subsRounds1, NULL)

	return events, nil
}
