package engine

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/contracts/router"
	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/matchevents"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	mp.SerializedEvents = big.NewInt(0)
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
		} else {
			event.TeamID = b.VisitorTeam.TeamID
		}
		if computedEvent.PrimaryPlayerID != "" {
			event.PrimaryPlayerID.String = computedEvent.PrimaryPlayerID
			event.PrimaryPlayerID.Valid = true
		}
		if computedEvent.SecondaryPlayerID != "" {
			event.SecondaryPlayerID.String = computedEvent.SecondaryPlayerID
			event.SecondaryPlayerID.Valid = true
		}
		event.TimezoneIdx = int(b.TimezoneIdx)
		event.CountryIdx = int(b.CountryIdx)
		event.LeagueIdx = int(b.LeagueIdx)
		event.MatchDayIdx = int(b.MatchDayIdx)
		event.MatchIdx = int(b.MatchIdx)
		event.Minute = int(computedEvent.Minute)
		var err error
		if event.Type, err = matchevents.MarchEventTypeByMatchEvent(computedEvent.Type); err != nil {
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
	version2StartEpoch := viper.GetInt64("patch.engine_version_2")

	var err error
	if b.StartTime.Cmp(big.NewInt(version2StartEpoch)) < 0 {
		err = b.play1stHalfV1(contracts)
	} else {
		err = b.play1stHalfV2(contracts)
	}
	if err != nil {
		b.State = storage.MatchCancelled
		b.StateExtra = err.Error()
	} else {
		b.State = storage.MatchHalf
	}
	return err
}

func (b *Match) play1stHalfV1(contracts contracts.Contracts) error {
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
	var BCError uint8
	isHomeBotOrZombie := b.HomeTeam.IsBot() || b.HomeTeam.IsZombie
	isVisitorBotOrZombie := b.VisitorTeam.IsBot() || b.VisitorTeam.IsZombie
	newSkills, logsAndEvents, BCError, err := contracts.PlayAndEvolve.Play1stHalfAndEvolve(
		&bind.CallOpts{},
		b.Seed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{homeTeamID, visitorTeamID},
		[2]*big.Int{homeTactic, visitorTactic},
		matchLogs,
		[5]bool{is2ndHalf, isHomeStadium, isPlayoff, isHomeBotOrZombie, isVisitorBotOrZombie},
		[2]*big.Int{homeAssignedTP, visitorAssignedTP},
	)
	if err != nil {
		return errors.Wrap(err, "failed play1stHalfAndEvolve")
	}
	if BCError != 0 {
		errMsg := fmt.Sprintf("BLOCKCHAIN ERROR!!!! Play1stHalfAndEvolve: Blockchain returned error code: %v", BCError)
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	if err = b.processMatchEvents(
		contracts,
		logsAndEvents[:],
		is2ndHalf,
	); err != nil {
		return errors.Wrap(err, "failed processing match events")
	}

	b.HomeTeam.SetSkills(contracts, newSkills[0])
	b.VisitorTeam.SetSkills(contracts, newSkills[1])
	b.HomeTeam.MatchLog = logsAndEvents[0].String()
	b.VisitorTeam.MatchLog = logsAndEvents[1].String()
	return nil
}

func (b *Match) play1stHalfV2(contracts contracts.Contracts) error {
	is2ndHalf := false
	homeTeamID, _ := new(big.Int).SetString(b.HomeTeam.TeamID, 10)
	visitorTeamID, _ := new(big.Int).SetString(b.VisitorTeam.TeamID, 10)
	homeTactic, _ := new(big.Int).SetString(b.HomeTeam.Tactic, 10)
	visitorTactic, _ := new(big.Int).SetString(b.VisitorTeam.Tactic, 10)
	matchLogs := [2]*big.Int{}
	matchLogs[0], _ = new(big.Int).SetString(b.HomeTeam.MatchLog, 10)
	matchLogs[1], _ = new(big.Int).SetString(b.VisitorTeam.MatchLog, 10)
	// If we cannot encode a the team's TP assignment, we set it to zero, and see if the match can be played anyway.
	// If it can, great, the team(s) with incorrect TP will not evolve, that's all.
	// If it cannot, it will cancel, and return a valid 0-0
	homeAssignedTP, err := b.HomeTeam.EncodeAssignedTrainingPoints(contracts)
	if err != nil {
		log.Warningf("failed calculating home assignedTP: %v", err)
		homeAssignedTP = big.NewInt(0)
	}
	visitorAssignedTP, err := b.VisitorTeam.EncodeAssignedTrainingPoints(contracts)
	if err != nil {
		log.Warningf("failed calculating visitor assignedTP: %v", err)
		visitorAssignedTP = big.NewInt(0)
	}
	var BCError uint8
	isHomeBotOrZombie := b.HomeTeam.IsBot() || b.HomeTeam.IsZombie
	isVisitorBotOrZombie := b.VisitorTeam.IsBot() || b.VisitorTeam.IsZombie
	newSkills, logsAndEvents, BCError, err := contracts.PlayAndEvolve.Play1stHalfAndEvolve(
		&bind.CallOpts{},
		b.Seed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{homeTeamID, visitorTeamID},
		[2]*big.Int{homeTactic, visitorTactic},
		matchLogs,
		[5]bool{is2ndHalf, isHomeStadium, isPlayoff, isHomeBotOrZombie, isVisitorBotOrZombie},
		[2]*big.Int{homeAssignedTP, visitorAssignedTP},
	)
	// We have two types of possible errors returned:
	// - Virtual Machine errors (never expected), err != nil => halt
	// - Controlled error that was properly dealt by BC code, BCError != nil, returns CANCELLED match with valid log and events
	if err != nil {
		return errors.Wrap(err, "failed play1stHalfAndEvolve")
	}
	if BCError != 0 {
		// no events returned, no need to process them. Just log
		log.Warningf("GAME CANCELLED!!!! Play1stHalfAndEvolve: Solidity code returned error code: %v", BCError)
	} else {
		if err = b.processMatchEvents(
			contracts,
			logsAndEvents[:],
			is2ndHalf,
		); err != nil {
			return errors.Wrap(err, "failed processing match events")
		}
	}
	b.HomeTeam.SetSkills(contracts, newSkills[0])
	b.VisitorTeam.SetSkills(contracts, newSkills[1])
	b.HomeTeam.MatchLog = logsAndEvents[0].String()
	b.VisitorTeam.MatchLog = logsAndEvents[1].String()
	return nil
}

func (b *Match) play2ndHalfV1(contracts contracts.Contracts) error {
	is2ndHalf := true
	homeTeamID, _ := new(big.Int).SetString(b.HomeTeam.TeamID, 10)
	visitorTeamID, _ := new(big.Int).SetString(b.VisitorTeam.TeamID, 10)
	homeTactic, _ := new(big.Int).SetString(b.HomeTeam.Tactic, 10)
	visitorTactic, _ := new(big.Int).SetString(b.VisitorTeam.Tactic, 10)
	matchLogs := [2]*big.Int{}
	matchLogs[0], _ = new(big.Int).SetString(b.HomeTeam.MatchLog, 10)
	matchLogs[1], _ = new(big.Int).SetString(b.VisitorTeam.MatchLog, 10)
	var BCError uint8
	isHomeBotOrZombie := b.HomeTeam.IsBot() || b.HomeTeam.IsZombie
	isVisitorBotOrZombie := b.VisitorTeam.IsBot() || b.VisitorTeam.IsZombie
	newSkills, logsAndEvents, BCError, err := contracts.PlayAndEvolve.Play2ndHalfAndEvolve(
		&bind.CallOpts{},
		b.Seed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{homeTeamID, visitorTeamID},
		[2]*big.Int{homeTactic, visitorTactic},
		matchLogs,
		[5]bool{is2ndHalf, isHomeStadium, isPlayoff, isHomeBotOrZombie, isVisitorBotOrZombie},
	)
	if err != nil {
		return errors.Wrap(err, "failed play2ndHalfAndEvolve")
	}
	if BCError != 0 {
		errMsg := fmt.Sprintf("BLOCKCHAIN ERROR!!!! play2ndHalfAndEvolve: Blockchain returned error code: %v", BCError)
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}
	if err = b.processMatchEvents(
		contracts,
		logsAndEvents[:],
		is2ndHalf,
	); err != nil {
		return errors.Wrap(err, "failed processing match events")
	}
	b.HomeTeam.SetSkills(contracts, newSkills[0])
	b.VisitorTeam.SetSkills(contracts, newSkills[1])
	b.HomeTeam.MatchLog = logsAndEvents[0].String()
	b.VisitorTeam.MatchLog = logsAndEvents[1].String()
	b.updateStats()
	return nil
}

func (b *Match) Play2ndHalf(contracts contracts.Contracts) error {
	version2StartEpoch := viper.GetInt64("patch.engine_version_2")

	var err error
	if b.StartTime.Cmp(big.NewInt(version2StartEpoch)) < 0 {
		err = b.play2ndHalfV1(contracts)
	} else {
		err = b.play2ndHalfV2(contracts)
	}
	if err != nil {
		b.State = storage.MatchCancelled
		b.StateExtra = err.Error()
	} else {
		b.State = storage.MatchEnd
	}
	return err
}

func (b *Match) play2ndHalfV2(contracts contracts.Contracts) error {
	is2ndHalf := true
	homeTeamID, _ := new(big.Int).SetString(b.HomeTeam.TeamID, 10)
	visitorTeamID, _ := new(big.Int).SetString(b.VisitorTeam.TeamID, 10)
	homeTactic, _ := new(big.Int).SetString(b.HomeTeam.Tactic, 10)
	visitorTactic, _ := new(big.Int).SetString(b.VisitorTeam.Tactic, 10)
	matchLogs := [2]*big.Int{}
	matchLogs[0], _ = new(big.Int).SetString(b.HomeTeam.MatchLog, 10)
	matchLogs[1], _ = new(big.Int).SetString(b.VisitorTeam.MatchLog, 10)
	var BCError uint8
	isHomeBotOrZombie := b.HomeTeam.IsBot() || b.HomeTeam.IsZombie
	isVisitorBotOrZombie := b.VisitorTeam.IsBot() || b.VisitorTeam.IsZombie
	newSkills, logsAndEvents, BCError, err := contracts.PlayAndEvolve.Play2ndHalfAndEvolve(
		&bind.CallOpts{},
		b.Seed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{homeTeamID, visitorTeamID},
		[2]*big.Int{homeTactic, visitorTactic},
		matchLogs,
		[5]bool{is2ndHalf, isHomeStadium, isPlayoff, isHomeBotOrZombie, isVisitorBotOrZombie},
	)
	// We have two types of possible errors returned:
	// - Virtual Machine errors (never expected), err != nil => halt
	// - Controlled error that was properly dealt by BC code, BCError != nil, returns CANCELLED match with valid log and events
	if err != nil {
		return errors.Wrap(err, "failed play2ndHalfAndEvolve")
	}
	if BCError != 0 {
		// no events returned, no need to process them. Just log
		log.Warningf("GAME CANCELLED!!!! Play2ndHalfAndEvolve: Solidity code returned error code: %v", BCError)
	} else {
		// there is no need to do this when we stop storing deserialized stuff
		if err = b.processMatchEvents(
			contracts,
			logsAndEvents[:],
			is2ndHalf,
		); err != nil {
			return errors.Wrap(err, "failed processing match events")
		}
	}
	b.HomeTeam.SetSkills(contracts, newSkills[0])
	b.VisitorTeam.SetSkills(contracts, newSkills[1])
	b.HomeTeam.MatchLog = logsAndEvents[0].String()
	b.VisitorTeam.MatchLog = logsAndEvents[1].String()

	serializedEvents, err := router.SerializeEventsFromPlayHalf(logsAndEvents[2:])
	if err != nil {
		return err
	}
	b.SerializedEvents = serializedEvents
	b.updateStats()
	return nil
}

func (b *Match) Skills() [2][25]*big.Int {
	return [2][25]*big.Int{b.HomeTeam.Skills(), b.VisitorTeam.Skills()}
}

func (b *Match) processMatchEvents(
	contracts contracts.Contracts,
	logsAndEvents []*big.Int,
	is2ndHalf bool,
) error {
	decodedHomeMatchLog, err := contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[0], is2ndHalf)
	if err != nil {
		return errors.Wrap(err, "failed decoding home match log")
	}
	decodedVisitorMatchLog, err := contracts.Utils.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[1], is2ndHalf)
	if err != nil {
		return errors.Wrap(err, "failed decoding visitor match log")
	}
	homeTactic, _ := new(big.Int).SetString(b.HomeTeam.Tactic, 10)
	visitorTactic, _ := new(big.Int).SetString(b.VisitorTeam.Tactic, 10)
	events, err := matchevents.NewMatchEvents(
		contracts,
		b.Seed,
		b.HomeTeam.TeamID,
		b.VisitorTeam.TeamID,
		b.HomeTeam.PlayerIDs(),
		b.VisitorTeam.PlayerIDs(),
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
	b.HomeGoals = uint8(decodedHomeMatchLog[2])
	b.VisitorGoals = uint8(decodedVisitorMatchLog[2])
	b.HomeTeam.TrainingPoints = uint16(decodedHomeMatchLog[3])
	b.VisitorTeam.TrainingPoints = uint16(decodedVisitorMatchLog[3])
	b.HomeTeamSumSkills = uint32(decodedHomeMatchLog[0])
	b.VisitorTeamSumSkills = uint32(decodedVisitorMatchLog[0])
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
