package matchevents

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type MatchEvent struct {
	Minute            int16 `json:"minute"`
	Type              int16 `json:"type"`
	Team              int16 `json:"team"`
	ManagesToShoot    bool  `json:"managestoshoot"`
	IsGoal            bool  `json:"isgoal"`
	PrimaryPlayer     int16 `json:"primaryplayer"`
	SecondaryPlayer   int16 `json:"secondaryplayer"`
	PrimaryPlayerID   string
	SecondaryPlayerID string
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

func MarchEventTypeByMatchEvent(event int16) (storage.MatchEventType, error) {
	switch event {
	case EVNT_ATTACK:
		return storage.Attack, nil
	case EVNT_YELLOW:
		return storage.YellowCard, nil
	case EVNT_RED:
		return storage.RedCard, nil
	case EVNT_SOFT:
		return storage.InjurySoft, nil
	case EVNT_HARD:
		return storage.InjuryHard, nil
	case EVNT_SUBST:
		return storage.Substitution, nil
	default:
		return "", fmt.Errorf("Unknown match event %v", event)
	}
}

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
