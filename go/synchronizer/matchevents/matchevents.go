package matchevents

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
)

type MatchEvents []MatchEvent

func IntHash(s string) uint64 {
	h := sha256.Sum256([]byte(s))
	// retain only the first 8 bytes and convert to uint64
	biggy := new(big.Int).SetBytes(h[:])
	return new(big.Int).Rsh(biggy, 192).Uint64()
}

func NewMatchEvents(
	contracts contracts.Contracts,
	verseSeed [32]byte,
	homeTeamID string,
	visitorTeamID string,
	homeTeamPlayerIDs [25]*big.Int,
	visitorTeamPlayerIDs [25]*big.Int,
	homeTactic *big.Int,
	visitorTactic *big.Int,
	logsAndEvents []*big.Int,
	homeDecodedMatchLog [15]uint32,
	visitorDecodedMatchLog [15]uint32,
	is2ndHalf bool,
) (MatchEvents, error) {
	if len(logsAndEvents) < 2 {
		return nil, errors.New("logAndEvents len < 2")
	}
	decodedTactic0, err := contracts.Utils.DecodeTactics(&bind.CallOpts{}, homeTactic)
	if err != nil {
		return nil, err
	}
	decodedTactic1, err := contracts.Utils.DecodeTactics(&bind.CallOpts{}, visitorTactic)
	if err != nil {
		return nil, err
	}

	events, err := Generate(
		verseSeed,
		homeTeamID,
		visitorTeamID,
		homeDecodedMatchLog,
		visitorDecodedMatchLog,
		logsAndEvents,
		RemoveFreeShirtsFromLineUp(decodedTactic0.Lineup, homeTeamPlayerIDs),
		RemoveFreeShirtsFromLineUp(decodedTactic1.Lineup, visitorTeamPlayerIDs),
		decodedTactic0.Substitutions,
		decodedTactic1.Substitutions,
		decodedTactic0.SubsRounds,
		decodedTactic1.SubsRounds,
		is2ndHalf,
	)
	if err != nil {
		return nil, err
	}

	if err := events.populateWithPlayerID(homeTeamPlayerIDs, visitorTeamPlayerIDs); err != nil {
		return nil, err
	}

	return events, nil
}

// This function makes sure that all players in lineUp exist in the Universe.
// To avoid, for example, selling a player after setting the lineUp.
// It sets to NOONE all lineUp entries pointing to playerIds that are larger than 2.
// (recall that playerID = 0 if not set, or = 1 if sold)
func RemoveFreeShirtsFromLineUp(lineUp [14]uint8, playerIDs [25]*big.Int) [14]uint8 {
	NO_LINEUP := uint8(25)
	MIN_PLAYERID := new(big.Int).SetUint64(2)
	for l := 0; l < len(lineUp); l++ {
		if lineUp[l] < NO_LINEUP {
			playerId := playerIDs[lineUp[l]]
			if playerId == nil || playerId.Cmp(MIN_PLAYERID) != 1 {
				lineUp[l] = NO_LINEUP
			}
		}
	}
	return lineUp
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
	PENALTY := int16(100)
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

	// before: verseSeedStr := string(verseSeed[:])
	// now: hex.EncodeToString(verseSeed[:]), because this is precisely the string exposed to frontend from the backend DB
	verseSeedStr := hex.EncodeToString(verseSeed[:])

	seed0 := new(big.Int).SetUint64(IntHash(verseSeedStr + "_0_" + teamId0 + "_" + teamId1))
	seed1 := new(big.Int).SetUint64(IntHash(verseSeedStr + "_1_" + teamId0 + "_" + teamId1))
	seed2 := new(big.Int).SetUint64(IntHash(verseSeedStr + "_2_" + teamId0 + "_" + teamId1))

	// There are mainly 3 types of events to reports, which are in different parts of the inputs:
	// - per-round (always 12 per half)
	// - cards & injuries
	// - substitutions
	events, rounds2mins := addEventsInRound(seed0, blockchainEvents, lineup0, lineup1, NULL, NOONE, PENALTY)
	events, err = addCardsAndInjuries(0, events, seed1, matchLog0, rounds2mins, lineup0, NULL, NOONE)
	if err != nil {
		return events, err
	}
	events, err = addCardsAndInjuries(1, events, seed2, matchLog1, rounds2mins, lineup1, NULL, NOONE)
	if err != nil {
		return events, err
	}
	events = addSubstitutions(0, events, matchLog0, rounds2mins, lineup0, substitutions0, subsRounds0, NULL, NOONE)
	events = addSubstitutions(1, events, matchLog1, rounds2mins, lineup1, substitutions1, subsRounds1, NULL, NOONE)

	if is2ndHalf {
		for e := range events {
			events[e].Minute += 145
		}
	}

	return events, nil
}

func addCardsAndInjuries(team int16, events []MatchEvent, seed *big.Int, matchLog [15]uint32, rounds2mins []uint64, lineUp [14]uint8, NULL int16, NOONE int16) ([]MatchEvent, error) {

	// matchLog[4,5,6] = outOfGamePlayer, outOfGameType, outOfGameRound
	// note that outofgame is a number from 0 to 13, and that NO OUT OF GAME = 14
	// eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
	if matchLog[5] > 3 {
		return events, errors.New("typeOfEvent larger than 3")
	}
	if matchLog[6] >= uint32(len(rounds2mins)) {
		return events, errors.New("outOfGameRound larger than allowed")
	}

	outOfGamePlayer := int16(matchLog[4])
	// convert player in the lineUp to shirtNum before storing it as match event:
	primaryPlayer := toShirtNum(uint8(outOfGamePlayer), lineUp, NULL, NOONE)
	thereWasAnOutOfGame := primaryPlayer != NULL
	outOfGameMinute := int16(0)
	if thereWasAnOutOfGame {
		if matchLog[5] == 0 {
			return events, errors.New("typeOfEvent = 0 is not allowed if thereWasAnOutOfGame")
		}
		var typeOfEvent int16
		if matchLog[5] == 1 {
			typeOfEvent = EVNT_SOFT
		} else if matchLog[5] == 2 {
			typeOfEvent = EVNT_HARD
		} else if matchLog[5] == 3 {
			typeOfEvent = EVNT_RED
		}
		outOfGameMinute = int16(rounds2mins[matchLog[6]])
		thisEvent := MatchEvent{outOfGameMinute, typeOfEvent, team, false, false, primaryPlayer, NULL, "", ""}
		events = append(events, thisEvent)
	}

	// First yellow card:
	yellowCardPlayer := int16(matchLog[7])
	// convert player in the lineUp to shirtNum before storing it as match event:
	primaryPlayer = toShirtNum(uint8(yellowCardPlayer), lineUp, NULL, NOONE)
	thereWasYellowCard := primaryPlayer != NULL
	firstYellowCoincidesWithRed := false
	if thereWasYellowCard {
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
		thisEvent := MatchEvent{minute, typeOfEvent, team, false, false, primaryPlayer, NULL, "", ""}
		events = append(events, thisEvent)
	}

	// Second second yellow card:
	yellowCardPlayer = int16(matchLog[8])
	// convert player in the lineUp to shirtNum before storing it as match event:
	primaryPlayer = toShirtNum(uint8(yellowCardPlayer), lineUp, NULL, NOONE)
	thereWasYellowCard = primaryPlayer != NULL
	if thereWasYellowCard {
		maxMinute := int16(45)
		typeOfEvent := EVNT_YELLOW
		if yellowCardPlayer == outOfGamePlayer {
			if firstYellowCoincidesWithRed {
				minute := outOfGameMinute
				thisEvent := MatchEvent{minute, typeOfEvent, team, false, false, primaryPlayer, NULL, "", ""}
				events = append(events, thisEvent)
				return events, nil
			} else {
				maxMinute = outOfGamePlayer
			}
		}
		salt := "d" + strconv.Itoa(int(yellowCardPlayer))
		minute := int16(GenerateRnd(seed, salt, uint64(maxMinute)))
		// convert player in the lineUp to shirtNum before storing it as match event:
		thisEvent := MatchEvent{minute, typeOfEvent, team, false, false, primaryPlayer, NULL, "", ""}
		events = append(events, thisEvent)
	}
	return events, nil
}

// output event order: (minute, eventType, managesToShoot, isGoal, player1, player2)
// eventType (0 = normal event, 1 = yellowCard, 2 = redCard, 3 = injurySoft, 4 = injuryHard, 5 = substitutions)
func addEventsInRound(seed *big.Int, blockchainEvents []*big.Int, lineup0 [14]uint8, lineup1 [14]uint8, NULL int16, NOONE int16, PENALTY int16) ([]MatchEvent, []uint64) {
	var events []MatchEvent
	nEvents := (len(blockchainEvents) - 2) / 5
	deltaMinutes := uint64(45 / (nEvents - 1))

	lineUps := [2][len(lineup1)]uint8{lineup0, lineup1}
	lastMinute := uint64(0)
	var rounds2mins []uint64
	for e := 0; e < nEvents; e++ {
		// compute minute
		salt := "a" + strconv.Itoa(int(e))
		minute := uint64(e)*deltaMinutes + GenerateRnd(seed, salt, deltaMinutes)
		if minute <= lastMinute {
			minute = lastMinute + 1
		}
		lastMinute = minute
		rounds2mins = append(rounds2mins, minute)
		// parse type of event and data
		// note that both "shooter" and "assister" referred to the lineup players (0...10)
		// so they need to be converted to shirtNums by using lineUp
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
			// select the players from the team that attacks:
			thisEvent.PrimaryPlayer = toShirtNum(uint8(shooter.Int64()), lineUps[thisEvent.Team], NULL, NOONE)
			if int16(assister.Int64()) == PENALTY {
				thisEvent.SecondaryPlayer = PENALTY
			} else {
				thisEvent.SecondaryPlayer = toShirtNum(uint8(assister.Int64()), lineUps[thisEvent.Team], NULL, NOONE)
			}
		} else {
			salt := "b" + strconv.Itoa(int(e))
			// select the player from the team that defends:
			thisEvent.PrimaryPlayer = toShirtNum(uint8(1+GenerateRnd(seed, salt, 9)), lineUps[1-thisEvent.Team], NULL, NOONE)
			thisEvent.SecondaryPlayer = NULL
		}
		events = append(events, thisEvent)
	}
	endOfGameMin := rounds2mins[len(rounds2mins)-1] + 1
	rounds2mins = append(rounds2mins, endOfGameMin)
	return events, rounds2mins
}

func toShirtNum(posInLineUp uint8, lineUp [14]uint8, NULL int16, NOONE int16) int16 {
	if int16(posInLineUp) < NOONE {
		return preventNoPlayer(int16(lineUp[posInLineUp]), NULL)
	} else {
		return NULL
	}
}

func preventNoPlayer(inPlayer int16, NULL int16) int16 {
	if inPlayer < 25 {
		return inPlayer
	} else {
		return NULL
	}
}

func addSubstitutions(team int16, events []MatchEvent, matchLog [15]uint32, rounds2mins []uint64, lineup [14]uint8, substitutions [3]uint8, subsRounds [3]uint8, NULL int16, NOONE int16) []MatchEvent {
	// matchLog:	9,10,11 ingameSubs, ...0: no change required, 1: change happened, 2: change could not happen
	// halftimesubs: 0 means no subs, and we store here p+1 (where p = player in the starting 11 that was substituted)
	for i := 0; i < 3; i++ {
		subHappened := matchLog[9+i] == 1
		if subHappened {
			minute := int16(rounds2mins[subsRounds[i]])
			leavingPlayer := toShirtNum(substitutions[i], lineup, NULL, NOONE)
			enteringPlayer := toShirtNum(uint8(11+i), lineup, NULL, NOONE)
			typeOfEvent := EVNT_SUBST
			thisEvent := MatchEvent{minute, typeOfEvent, team, false, false, leavingPlayer, enteringPlayer, "", ""}
			events = append(events, thisEvent)
		}
	}
	return adjustSubstitutions(team, events, NULL)
}

// make sure that if a player that enters via a substitution appears in any other action (goal, pass, cards & injuries),
// then the substitution time must take place at least before that minute.
func adjustSubstitutions(team int16, events []MatchEvent, NULL int16) []MatchEvent {
	adjustedEvents := events
	for e := 0; e < len(events); e++ {
		if (events[e].Type == EVNT_SUBST) && (events[e].Team == team) {
			enteringPlayer := events[e].SecondaryPlayer
			if enteringPlayer != NULL {
				enteringMin := events[e].Minute
				for e2 := 0; e2 < len(events); e2++ {
					if (e != e2) && (events[e2].Team == team) && (enteringPlayer == events[e2].PrimaryPlayer) && (enteringMin >= events[e2].Minute-1) {
						adjustedEvents[e].Minute = events[e2].Minute - 1
					}
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

func (b *MatchEvents) populateWithPlayerID(
	homeTeamPlayerIDs [25]*big.Int,
	visitorTeamPlayerIDs [25]*big.Int,
) error {
	for i := range *b {
		var primaryPlayerTeam [25]*big.Int
		var secondaryPlayerTeam [25]*big.Int
		var tacklerPlayerTeam [25]*big.Int
		if (*b)[i].Team == 0 {
			primaryPlayerTeam = homeTeamPlayerIDs
			secondaryPlayerTeam = homeTeamPlayerIDs
			tacklerPlayerTeam = visitorTeamPlayerIDs
		} else {
			primaryPlayerTeam = visitorTeamPlayerIDs
			secondaryPlayerTeam = visitorTeamPlayerIDs
			tacklerPlayerTeam = homeTeamPlayerIDs
		}

		if (*b)[i].PrimaryPlayer != -1 {
			if (*b)[i].Type == EVNT_ATTACK && !(*b)[i].ManagesToShoot {
				if tacklerPlayerTeam[(*b)[i].PrimaryPlayer] == nil {
					return fmt.Errorf("inconsistent event %+v", (*b)[i])
				}
				(*b)[i].PrimaryPlayerID = tacklerPlayerTeam[(*b)[i].PrimaryPlayer].String()
			} else {
				if primaryPlayerTeam[(*b)[i].PrimaryPlayer] == nil {
					return fmt.Errorf("inconsistent event %+v", (*b)[i])
				}
				(*b)[i].PrimaryPlayerID = primaryPlayerTeam[(*b)[i].PrimaryPlayer].String()
			}
		}

		switch (*b)[i].SecondaryPlayer {
		case -1: // no seconday player
		case contracts.PenaltyPlayerId:
			(*b)[i].SecondaryPlayerID = fmt.Sprintf("%d", contracts.PenaltyPlayerId)
		default:
			if secondaryPlayerTeam[(*b)[i].SecondaryPlayer] == nil {
				return fmt.Errorf("inconsistent event %+v", (*b)[i])
			}
			(*b)[i].SecondaryPlayerID = secondaryPlayerTeam[(*b)[i].SecondaryPlayer].String()
		}
	}

	return nil
}
