package matchevents

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
)

type MatchEvents []MatchEvent

func NewMatchEvents(
	contracts contracts.Contracts,
	verseSeed [32]byte,
	homeTeamID string,
	visitorTeamID string,
	homeTactic *big.Int,
	visitorTactic *big.Int,
	logsAndEvents []*big.Int,
	is2ndHalf bool,
) (MatchEvents, error) {
	log0, err := contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[0], is2ndHalf)
	if err != nil {
		return nil, err
	}
	log1, err := contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[1], is2ndHalf)
	if err != nil {
		return nil, err
	}
	decodedTactic0, err := contracts.Engine.DecodeTactics(&bind.CallOpts{}, homeTactic)
	if err != nil {
		return nil, err
	}
	decodedTactic1, err := contracts.Engine.DecodeTactics(&bind.CallOpts{}, visitorTactic)
	if err != nil {
		return nil, err
	}
	events, err := Generate(
		verseSeed,
		homeTeamID,
		visitorTeamID,
		log0,
		log1,
		logsAndEvents,
		decodedTactic0.Lineup,
		decodedTactic1.Lineup,
		decodedTactic0.Substitutions,
		decodedTactic1.Substitutions,
		decodedTactic0.SubsRounds,
		decodedTactic1.SubsRounds,
		is2ndHalf,
	)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func Generate(
	verseSeed [32]byte,
	teamId0 string,
	teamId1 string,
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
) (MatchEvents, error) {
	NULL := int16(-1)
	NOONE := int16(14)
	var emptyEvents []MatchEvent

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

	seed := new(big.Int).SetUint64(int_hash(string(verseSeed[:]) + "_" + teamId0 + "_" + teamId1))

	// Compute main events: per-round, and cards & injuries
	events, rounds2mins := addEventsInRound(seed, blockchainEvents, NULL)
	events = addCardsAndInjuries(0, events, seed, matchLog0, rounds2mins, NULL, NOONE)
	events = addCardsAndInjuries(1, events, seed, matchLog0, rounds2mins, NULL, NOONE)
	events = addSubstitutions(0, events, seed, matchLog0, rounds2mins, lineup0, substitutions0, subsRounds0, NULL)
	events = addSubstitutions(1, events, seed, matchLog1, rounds2mins, lineup1, substitutions1, subsRounds1, NULL)

	if is2ndHalf {
		for e := range events {
			events[e].Minute += 45
		}
	}

	return events, nil
}

func addCardsAndInjuries(team int16, events []MatchEvent, seed *big.Int, matchLog [15]uint32, rounds2mins []uint64, NULL int16, NOONE int16) []MatchEvent {
	// matchLog[4,5,6] = outOfGamePlayer, outOfGameType, outOfGameRound
	// note that outofgame is a number from 0 to 13, and that NO OUT OF GAME = 14
	// eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
	outOfGamePlayer := int16(matchLog[4])
	thereWasAnOutOfGame := outOfGamePlayer < NOONE
	outOfGameMinute := int16(0)
	if thereWasAnOutOfGame {
		var typeOfEvent int16
		if matchLog[5] == 1 {
			typeOfEvent = EVNT_SOFT
		} else if matchLog[5] == 2 {
			typeOfEvent = EVNT_HARD
		} else if matchLog[5] == 3 {
			typeOfEvent = EVNT_RED
		}
		outOfGameMinute = int16(rounds2mins[matchLog[6]])
		thisEvent := MatchEvent{outOfGameMinute, typeOfEvent, team, false, false, outOfGamePlayer, NULL}
		events = append(events, thisEvent)
	}

	// First yellow card:
	yellowCardPlayer := int16(matchLog[7])
	firstYellowCoincidesWithRed := false
	if yellowCardPlayer < 14 {
		maxMinute := int16(45)
		if yellowCardPlayer == outOfGamePlayer {
			if outOfGameMinute > 0 {
				maxMinute = outOfGameMinute - 1
			} else {
				maxMinute = outOfGameMinute
			}
			firstYellowCoincidesWithRed = true
		}
		salt := "c" + strconv.Itoa(int(yellowCardPlayer))
		minute := int16(GenerateRnd(seed, salt, uint64(maxMinute)))
		typeOfEvent := EVNT_YELLOW
		thisEvent := MatchEvent{minute, typeOfEvent, team, false, false, yellowCardPlayer, NULL}
		events = append(events, thisEvent)
	}

	// Second second yellow card:
	yellowCardPlayer = int16(matchLog[8])
	if yellowCardPlayer < 14 {
		maxMinute := int16(45)
		typeOfEvent := EVNT_YELLOW
		if yellowCardPlayer == outOfGamePlayer {
			if firstYellowCoincidesWithRed {
				minute := outOfGameMinute
				thisEvent := MatchEvent{minute, typeOfEvent, team, false, false, yellowCardPlayer, NULL}
				events = append(events, thisEvent)
				return events
			} else {
				maxMinute = outOfGamePlayer
			}
		}
		salt := "d" + strconv.Itoa(int(yellowCardPlayer))
		minute := int16(GenerateRnd(seed, salt, uint64(maxMinute)))
		thisEvent := MatchEvent{minute, typeOfEvent, team, false, false, yellowCardPlayer, NULL}
		events = append(events, thisEvent)
	}
	return events
}

// output event order: (minute, eventType, managesToShoot, isGoal, player1, player2)
// eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
func addEventsInRound(seed *big.Int, blockchainEvents []*big.Int, NULL int16) ([]MatchEvent, []uint64) {
	var events []MatchEvent
	nEvents := (len(blockchainEvents) - 2) / 5
	deltaMinutes := float64(45.0 / ((nEvents - 1) * 1.0))
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
		var thisEvent MatchEvent
		thisEvent.Minute = int16(minute)
		thisEvent.Type = int16(EVNT_ATTACK)
		thisEvent.Team = int16(teamThatAttacks.Int64())
		thisEvent.ManagesToShoot = managesToShoot.Int64() != 0
		thisEvent.IsGoal = isGoal.Int64() != 0
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

func addSubstitutions(team int16, events []MatchEvent, seed *big.Int, matchLog [15]uint32, rounds2mins []uint64, lineup [14]uint8, substitutions [3]uint8, subsRounds [3]uint8, NULL int16) []MatchEvent {
	// matchLog:	9,10,11 ingameSubs, ...0: no change required, 1: change happened, 2: change could not happen
	// halftimesubs: 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
	for i := 0; i < 3; i++ {
		subHappened := matchLog[9+i] == 1
		if subHappened {
			minute := int16(rounds2mins[subsRounds[i]])
			leavingPlayer := int16(substitutions[i])
			enteringPlayer := int16(lineup[11+i])
			typeOfEvent := EVNT_SUBST
			thisEvent := MatchEvent{minute, typeOfEvent, team, false, false, leavingPlayer, enteringPlayer}
			events = append(events, thisEvent)
		}
	}
	return adjustSubstitutions(team, events)
}

// make sure that if a player that enters via a substitution appears in any other action (goal, pass, cards & injuries),
// then the substitution time must take place at least before that minute.
func adjustSubstitutions(team int16, events []MatchEvent) []MatchEvent {
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

func (b MatchEvents) HomeGoals() uint8 {
	var counter uint8
	for _, event := range b {
		if event.Team == 0 && event.IsGoal {
			counter++
		}
	}
	return counter
}

func (b MatchEvents) VisitorGoals() uint8 {
	var counter uint8
	for _, event := range b {
		if event.Team == 1 && event.IsGoal {
			counter++
		}
	}
	return counter
}

func (b MatchEvents) DumpState() string {
	var state string
	for i, event := range b {
		state += fmt.Sprintf("Events[%d]: %+v\n", i, event)
	}
	return state
}
