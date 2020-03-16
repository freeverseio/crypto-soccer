package engine

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/matchevents"
)

type Match struct {
	storage.Match
	StartTime   *big.Int
	HomeTeam    Team
	VisitorTeam Team
	Events      matchevents.MatchEvents
}

const isHomeStadium = true
const isPlayoff = false

func NewMatch() *Match {
	var mp Match
	mp.StartTime = big.NewInt(0)
	mp.HomeTeam = *NewTeam()
	mp.VisitorTeam = *NewTeam()
	mp.State = storage.MatchBegin
	return &mp
}

func NewMatchFromStorage(
	sMatch storage.Match,
	sHomeTeam storage.Team,
	sVisitorTeam storage.Team,
	sHomePlayers []*storage.Player,
	sVisitorPlayers []*storage.Player,
) *Match {
	match := NewMatch()
	match.Match = sMatch
	match.HomeTeam.Team = sHomeTeam
	match.VisitorTeam.Team = sVisitorTeam
	for _, player := range sHomePlayers {
		match.HomeTeam.Players[player.ShirtNumber].Player = *player
	}
	for _, player := range sVisitorPlayers {
		match.VisitorTeam.Players[player.ShirtNumber].Player = *player
	}
	return match
}

func (b *Match) updateStats() {
	b.HomeTeam.GoalsForward += uint32(b.HomeGoals)
	b.HomeTeam.GoalsAgainst += uint32(b.VisitorGoals)
	b.VisitorTeam.GoalsForward += uint32(b.VisitorGoals)
	b.VisitorTeam.GoalsAgainst += uint32(b.HomeGoals)
	deltaGoals := int(b.HomeGoals) - int(b.VisitorGoals)
	if deltaGoals > 0 {
		b.HomeTeam.W++
		b.VisitorTeam.L++
		b.HomeTeam.Points += 3
	} else if deltaGoals < 0 {
		b.HomeTeam.L++
		b.VisitorTeam.W++
		b.VisitorTeam.Points += 3
	} else {
		b.HomeTeam.D++
		b.VisitorTeam.D++
		b.HomeTeam.Points++
		b.VisitorTeam.Points++
	}
}

func (b Match) ToStorage(contracts contracts.Contracts, tx *sql.Tx, blockNumber uint64) error {
	if err := b.HomeTeam.ToStorage(contracts, tx, blockNumber); err != nil {
		return err
	}
	if err := b.VisitorTeam.ToStorage(contracts, tx, blockNumber); err != nil {
		return err
	}
	for _, computedEvent := range b.Events {
		event := storage.MatchEvent{}
		if computedEvent.Team == 0 {
			event.TeamID = b.HomeTeam.TeamID
			event.PrimaryPlayerID = b.HomeTeam.Players[computedEvent.PrimaryPlayer].PlayerId.String()
			if computedEvent.SecondaryPlayer >= 0 && computedEvent.SecondaryPlayer < 25 {
				event.SecondaryPlayerID.String = b.HomeTeam.Players[computedEvent.SecondaryPlayer].PlayerId.String()
				event.SecondaryPlayerID.Valid = true
			}
		} else if computedEvent.Team == 1 {
			event.TeamID = b.VisitorTeam.TeamID
			event.PrimaryPlayerID = b.VisitorTeam.Players[computedEvent.PrimaryPlayer].PlayerId.String()
			if computedEvent.SecondaryPlayer >= 0 && computedEvent.SecondaryPlayer < 25 {
				event.SecondaryPlayerID.String = b.VisitorTeam.Players[computedEvent.SecondaryPlayer].PlayerId.String()
				event.SecondaryPlayerID.Valid = true
			}
		} else {
			return fmt.Errorf("Wrong match event team %v", computedEvent.Team)
		}
		event.TimezoneIdx = int(b.TimezoneIdx)
		event.CountryIdx = int(b.CountryIdx)
		event.LeagueIdx = int(b.LeagueIdx)
		event.MatchDayIdx = int(b.MatchDayIdx)
		event.MatchIdx = int(b.MatchIdx)
		event.Minute = int(computedEvent.Minute)
		var err error
		if event.Type, err = storage.MarchEventTypeByMatchEvent(computedEvent.Type); err != nil {
			return err
		}
		event.ManageToShoot = computedEvent.ManagesToShoot
		event.IsGoal = computedEvent.IsGoal
		if err = event.Insert(tx); err != nil {
			return err
		}
	}
	return b.Update(tx, blockNumber)
}

func (b *Match) Play1stHalf(contracts contracts.Contracts) error {
	err := b.play1stHalf(contracts)
	if err != nil {
		b.State = storage.MatchCancelled
	} else {
		b.State = storage.MatchHalf
	}
	return err
}

func (b *Match) play1stHalf(contracts contracts.Contracts) error {
	is2ndHalf := false
	homeTeamID, _ := new(big.Int).SetString(b.HomeTeam.TeamID, 10)
	visitorTeamID, _ := new(big.Int).SetString(b.VisitorTeam.TeamID, 10)
	homeTactic, _ := new(big.Int).SetString(b.HomeTeam.Tactic, 10)
	visitorTactic, _ := new(big.Int).SetString(b.VisitorTeam.Tactic, 10)
	matchLogs := [2]*big.Int{}
	matchLogs[0], _ = new(big.Int).SetString(b.HomeTeam.MatchLog, 10)
	matchLogs[1], _ = new(big.Int).SetString(b.VisitorTeam.MatchLog, 10)
	homeAssignedTP, err := b.HomeTeam.CalculateAssignedTrainingPoints(contracts)
	if err != nil {
		return err
	}
	visitorAssignedTP, err := b.VisitorTeam.CalculateAssignedTrainingPoints(contracts)
	if err != nil {
		return err
	}
	newSkills, logsAndEvents, err := contracts.PlayAndEvolve.Play1stHalfAndEvolve(
		&bind.CallOpts{},
		b.Seed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{homeTeamID, visitorTeamID},
		[2]*big.Int{homeTactic, visitorTactic},
		matchLogs,
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
		[2]*big.Int{homeAssignedTP, visitorAssignedTP},
	)
	if err != nil {
		return err
	}
	b.HomeTeam.SetSkills(contracts, newSkills[0])
	b.VisitorTeam.SetSkills(contracts, newSkills[1])
	b.HomeTeam.MatchLog = logsAndEvents[0].String()
	b.VisitorTeam.MatchLog = logsAndEvents[1].String()
	b.HomeGoals, b.VisitorGoals, err = b.getGoals(contracts, [2]*big.Int{logsAndEvents[0], logsAndEvents[1]})
	if err != nil {
		return err
	}
	if err = b.processMatchEvents(contracts, logsAndEvents[:], is2ndHalf); err != nil {
		return err
	}
	b.State = storage.MatchHalf
	return nil
}

func (b *Match) Play2ndHalf(contracts contracts.Contracts) error {
	err := b.play2ndHalf(contracts)
	if err != nil {
		b.State = storage.MatchCancelled
	} else {
		b.State = storage.MatchEnd
	}
	return err
}

func (b *Match) play2ndHalf(contracts contracts.Contracts) error {
	is2ndHalf := true
	homeTeamID, _ := new(big.Int).SetString(b.HomeTeam.TeamID, 10)
	visitorTeamID, _ := new(big.Int).SetString(b.VisitorTeam.TeamID, 10)
	homeTactic, _ := new(big.Int).SetString(b.HomeTeam.Tactic, 10)
	visitorTactic, _ := new(big.Int).SetString(b.VisitorTeam.Tactic, 10)
	matchLogs := [2]*big.Int{}
	matchLogs[0], _ = new(big.Int).SetString(b.HomeTeam.MatchLog, 10)
	matchLogs[1], _ = new(big.Int).SetString(b.VisitorTeam.MatchLog, 10)
	newSkills, logsAndEvents, err := contracts.PlayAndEvolve.Play2ndHalfAndEvolve(
		&bind.CallOpts{},
		b.Seed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{homeTeamID, visitorTeamID},
		[2]*big.Int{homeTactic, visitorTactic},
		matchLogs,
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	)
	if err != nil {
		return err
	}
	b.HomeGoals, b.VisitorGoals, err = b.getGoals(contracts, [2]*big.Int{logsAndEvents[0], logsAndEvents[1]})
	if err != nil {
		return err
	}
	b.HomeTeam.SetSkills(contracts, newSkills[0])
	b.VisitorTeam.SetSkills(contracts, newSkills[1])
	b.HomeTeam.MatchLog = logsAndEvents[0].String()
	b.VisitorTeam.MatchLog = logsAndEvents[1].String()
	if err = b.processMatchEvents(contracts, logsAndEvents[:], is2ndHalf); err != nil {
		return err
	}
	b.updateStats()
	if err = b.updateTrainingPoints(contracts); err != nil {
		return err
	}
	return nil
}

func (b *Match) updateTrainingPoints(contracts contracts.Contracts) error {
	var err error
	matchLog, _ := new(big.Int).SetString(b.HomeTeam.MatchLog, 10)
	if b.HomeTeam.TrainingPoints, err = contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, matchLog); err != nil {
		return err
	}
	matchLog, _ = new(big.Int).SetString(b.VisitorTeam.MatchLog, 10)
	if b.VisitorTeam.TrainingPoints, err = contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, matchLog); err != nil {
		return err
	}
	return nil
}

func (b *Match) getGoals(contracts contracts.Contracts, logs [2]*big.Int) (homeGoals uint8, VisitorGoals uint8, err error) {
	homeGoals, err = contracts.Evolution.GetNGoals(
		&bind.CallOpts{},
		logs[0],
	)
	if err != nil {
		return homeGoals, VisitorGoals, err
	}
	VisitorGoals, err = contracts.Evolution.GetNGoals(
		&bind.CallOpts{},
		logs[1],
	)
	return homeGoals, VisitorGoals, err
}

func (b *Match) Skills() [2][25]*big.Int {
	return [2][25]*big.Int{b.HomeTeam.Skills(), b.VisitorTeam.Skills()}
}

func (b *Match) processMatchEvents(contracts contracts.Contracts, logsAndEvents []*big.Int, is2ndHalf bool) error {
	homeTactic, _ := new(big.Int).SetString(b.HomeTeam.Tactic, 10)
	visitorTactic, _ := new(big.Int).SetString(b.VisitorTeam.Tactic, 10)
	events, err := matchevents.NewMatchEvents(
		contracts,
		b.Seed,
		b.HomeTeam.TeamID,
		b.VisitorTeam.TeamID,
		homeTactic,
		visitorTactic,
		logsAndEvents,
		is2ndHalf,
	)
	if err != nil {
		return err
	}
	b.Events = append(b.Events, events...)
	return nil
}

func (b Match) ToString() string {
	var result string
	result += fmt.Sprintf("seed = '0x%v';", hex.EncodeToString(b.Seed[:]))
	result += fmt.Sprintf("startTime = '%v';", b.StartTime)
	result += fmt.Sprintf("matchLog0 = '%v';", b.HomeTeam.MatchLog)
	result += fmt.Sprintf("teamId0 = '%v';", b.HomeTeam.TeamID)
	result += fmt.Sprintf("tactic0 = '%v';", b.HomeTeam.Tactic)
	// result += fmt.Sprintf("assignedTP0 = '%v';", b.HomeTeam.AssignedTP)
	result += "players0 = ["
	for _, player := range b.HomeTeam.Players {
		result += fmt.Sprintf("'%v',", player.EncodedSkills)
	}
	result += "];"
	result += fmt.Sprintf("matchLog1 = '%v';", b.VisitorTeam.MatchLog)
	result += fmt.Sprintf("teamId1 = '%v';", b.VisitorTeam.TeamID)
	result += fmt.Sprintf("tactic1 = '%v';", b.VisitorTeam.Tactic)
	// result += fmt.Sprintf("assignedTP1 = '%v';", b.VisitorTeam.AssignedTP)
	result += "players1 = ["
	for _, player := range b.VisitorTeam.Players {
		result += fmt.Sprintf("'%v',", player.EncodedSkills)
	}
	result += "];"
	return result
}
