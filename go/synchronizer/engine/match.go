package engine

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/matchevents"
	"github.com/pkg/errors"
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
			if computedEvent.PrimaryPlayer >= 0 && int(computedEvent.PrimaryPlayer) < len(b.HomeTeam.Players) {
				event.PrimaryPlayerID = b.HomeTeam.Players[computedEvent.PrimaryPlayer].PlayerId.String()
			}
			if computedEvent.SecondaryPlayer >= 0 && computedEvent.SecondaryPlayer < 25 {
				event.SecondaryPlayerID.String = b.HomeTeam.Players[computedEvent.SecondaryPlayer].PlayerId.String()
				event.SecondaryPlayerID.Valid = true
			}
		} else if computedEvent.Team == 1 {
			event.TeamID = b.VisitorTeam.TeamID
			if computedEvent.PrimaryPlayer >= 0 && int(computedEvent.PrimaryPlayer) < len(b.VisitorTeam.Players) {
				event.PrimaryPlayerID = b.VisitorTeam.Players[computedEvent.PrimaryPlayer].PlayerId.String()
			}
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
		b.StateExtra = err.Error()
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
	homeAssignedTP, err := b.HomeTeam.EncodeAssignedTrainingPoints(contracts)
	if err != nil {
		return errors.Wrap(err, "failed calculating home assignedTP")
	}
	visitorAssignedTP, err := b.VisitorTeam.EncodeAssignedTrainingPoints(contracts)
	if err != nil {
		return errors.Wrap(err, "failed calculating visitor assignedTP")
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
		return errors.Wrap(err, "failed play1stHalfAndEvolve")
	}
	decodedHomeMatchLog, err := contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[0], is2ndHalf)
	if err != nil {
		return errors.Wrap(err, "failed decoding home match log")
	}
	decodedVisitorMatchLog, err := contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[1], is2ndHalf)
	if err != nil {
		return errors.Wrap(err, "failed decoding visitor match log")
	}
	if err = b.processMatchEvents(
		contracts,
		logsAndEvents[:],
		decodedHomeMatchLog,
		decodedVisitorMatchLog,
		is2ndHalf,
	); err != nil {
		return errors.Wrap(err, "failed processing match events")
	}
	b.HomeTeam.SetSkills(contracts, newSkills[0])
	b.VisitorTeam.SetSkills(contracts, newSkills[1])
	b.HomeTeam.MatchLog = logsAndEvents[0].String()
	b.VisitorTeam.MatchLog = logsAndEvents[1].String()
	b.HomeGoals = uint8(decodedHomeMatchLog[2])
	b.VisitorGoals = uint8(decodedVisitorMatchLog[2])
	b.HomeTeam.TrainingPoints = uint16(decodedHomeMatchLog[3])
	b.VisitorTeam.TrainingPoints = uint16(decodedVisitorMatchLog[3])
	return nil
}

func (b *Match) Play2ndHalf(contracts contracts.Contracts) error {
	err := b.play2ndHalf(contracts)
	if err != nil {
		b.State = storage.MatchCancelled
		b.StateExtra = err.Error()
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
		return errors.Wrap(err, "failed play2ndHalfAndEvolve")
	}
	decodedHomeMatchLog, err := contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[0], is2ndHalf)
	if err != nil {
		return errors.Wrap(err, "failed decoding home match log")
	}
	decodedVisitorMatchLog, err := contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[1], is2ndHalf)
	if err != nil {
		return errors.Wrap(err, "failed decoding visitor match log")
	}
	b.HomeTeam.SetSkills(contracts, newSkills[0])
	b.VisitorTeam.SetSkills(contracts, newSkills[1])
	b.HomeTeam.MatchLog = logsAndEvents[0].String()
	b.VisitorTeam.MatchLog = logsAndEvents[1].String()
	b.HomeGoals = uint8(decodedHomeMatchLog[2])
	b.VisitorGoals = uint8(decodedVisitorMatchLog[2])
	b.HomeTeam.TrainingPoints = uint16(decodedHomeMatchLog[3])
	b.VisitorTeam.TrainingPoints = uint16(decodedVisitorMatchLog[3])
	if err = b.processMatchEvents(
		contracts,
		logsAndEvents[:],
		decodedHomeMatchLog,
		decodedVisitorMatchLog,
		is2ndHalf,
	); err != nil {
		return errors.Wrap(err, "failed processing match events")
	}
	b.updateStats()
	return nil
}

func (b *Match) Skills() [2][25]*big.Int {
	return [2][25]*big.Int{b.HomeTeam.Skills(), b.VisitorTeam.Skills()}
}

func (b *Match) processMatchEvents(
	contracts contracts.Contracts,
	logsAndEvents []*big.Int,
	decodedHomeMatchLog [15]uint32,
	decodedVisitorMatchLog [15]uint32,
	is2ndHalf bool,
) error {
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
		decodedHomeMatchLog,
		decodedVisitorMatchLog,
		is2ndHalf,
	)
	if err != nil {
		return err
	}
	b.Events = append(b.Events, events...)
	return nil
}

func (b Match) ToJson() []byte {
	s, _ := json.MarshalIndent(b, "", "\t")
	return s
}

func NewMatchFromJson(input []byte) (*Match, error) {
	match := Match{}
	if err := json.Unmarshal(input, &match); err != nil {
		return nil, err
	}
	return &match, nil
}

func (b Match) Hash() []byte {
	h := sha256.New()
	h.Write(b.ToJson())
	return h.Sum(nil)
}
