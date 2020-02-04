package engine

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/matchevents"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"

	log "github.com/sirupsen/logrus"
)

type Match struct {
	storage.Match
	Seed        [32]byte
	StartTime   *big.Int
	HomeTeam    Team
	VisitorTeam Team
	Events      matchevents.MatchEvents
}

const isHomeStadium = true
const isPlayoff = false

func (b Match) DumpState() string {
	var state string
	state += fmt.Sprintf("Seed: %v\n", hex.EncodeToString(b.Seed[:]))
	state += fmt.Sprintf("StartTime: %v\n", b.StartTime)
	state += fmt.Sprintf("HomeTeam: %v\n", b.HomeTeam.DumpState())
	state += fmt.Sprintf("VisitorTeam: %v\n", b.VisitorTeam.DumpState())
	state += fmt.Sprintf("HomeGoals: %v\n", b.HomeGoals)
	state += fmt.Sprintf("VisitorGoals: %v\n", b.VisitorGoals)
	state += fmt.Sprintf("HomeMatchLog: %v\n", b.HomeMatchLog)
	state += fmt.Sprintf("VisitorMatchLog: %v\n", b.VisitorMatchLog)
	state += b.Events.DumpState()
	return state
}

func NewMatch() *Match {
	var mp Match
	mp.StartTime = big.NewInt(0)
	mp.HomeTeam = *NewTeam()
	mp.VisitorTeam = *NewTeam()
	mp.HomeMatchLog = big.NewInt(0)
	mp.VisitorMatchLog = big.NewInt(0)
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
		match.HomeTeam.Players[player.ShirtNumber].sto = *player
	}
	for _, player := range sVisitorPlayers {
		match.VisitorTeam.Players[player.ShirtNumber].sto = *player
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

func (b Match) ToStorage(contracts contracts.Contracts, tx *sql.Tx) error {
	if err := b.HomeTeam.ToStorage(contracts, tx); err != nil {
		return err
	}
	if err := b.VisitorTeam.ToStorage(contracts, tx); err != nil {
		return err
	}
	for _, computedEvent := range b.Events {
		event := storage.MatchEvent{}
		if computedEvent.Team == 0 {
			event.TeamID = b.HomeTeam.TeamID.String()
			event.PrimaryPlayerID = b.HomeTeam.Players[computedEvent.PrimaryPlayer].sto.PlayerId.String()
			if computedEvent.SecondaryPlayer >= 0 && computedEvent.SecondaryPlayer < 25 {
				event.SecondaryPlayerID.String = b.HomeTeam.Players[computedEvent.SecondaryPlayer].sto.PlayerId.String()
				event.SecondaryPlayerID.Valid = true
			}
		} else if computedEvent.Team == 1 {
			event.TeamID = b.VisitorTeam.TeamID.String()
			event.PrimaryPlayerID = b.VisitorTeam.Players[computedEvent.PrimaryPlayer].sto.PlayerId.String()
			if computedEvent.SecondaryPlayer >= 0 && computedEvent.SecondaryPlayer < 25 {
				event.SecondaryPlayerID.String = b.VisitorTeam.Players[computedEvent.SecondaryPlayer].sto.PlayerId.String()
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
	return b.Update(tx)
}

func (b *Match) Play1stHalf(contracts contracts.Contracts) error {
	is2ndHalf := false
	assignedTPs := big.NewInt(int64(0))
	newSkills, logsAndEvents, err := contracts.PlayAndEvolve.Play1stHalfAndEvolve(
		&bind.CallOpts{},
		b.Seed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{b.HomeTeam.TeamID, b.VisitorTeam.TeamID},
		[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
		[2]*big.Int{b.HomeMatchLog, b.VisitorMatchLog},
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
		[2]*big.Int{assignedTPs, assignedTPs},
	)
	if err != nil {
		return err
	}
	b.HomeTeam.SetSkills(contracts, newSkills[0])
	b.VisitorTeam.SetSkills(contracts, newSkills[1])
	b.HomeMatchLog = logsAndEvents[0]
	b.VisitorMatchLog = logsAndEvents[1]
	b.HomeGoals, b.VisitorGoals, err = b.getGoals(contracts, [2]*big.Int{logsAndEvents[0], logsAndEvents[1]})
	if err != nil {
		return err
	}
	if err = b.processMatchEvents(contracts, logsAndEvents[:], is2ndHalf); err != nil {
		return err
	}
	return nil
}

func (b *Match) Play2ndHalf(contracts contracts.Contracts) error {
	is2ndHalf := true
	newSkills, logsAndEvents, err := contracts.PlayAndEvolve.Play2ndHalfAndEvolve(
		&bind.CallOpts{},
		b.Seed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{b.HomeTeam.TeamID, b.VisitorTeam.TeamID},
		[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
		[2]*big.Int{b.HomeMatchLog, b.VisitorMatchLog},
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
	b.HomeMatchLog = logsAndEvents[0]
	b.VisitorMatchLog = logsAndEvents[1]
	if err = b.processMatchEvents(contracts, logsAndEvents[:], is2ndHalf); err != nil {
		return err
	}
	b.updateStats()
	if b.HomeTeam.TrainingPoints, err = contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, b.HomeMatchLog); err != nil {
		return err
	}
	if b.VisitorTeam.TrainingPoints, err = contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, b.VisitorMatchLog); err != nil {
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
	log0, err := contracts.Utilsmatchlog.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[0], is2ndHalf)
	if err != nil {
		return err
	}
	log1, err := contracts.Utilsmatchlog.FullDecodeMatchLog(&bind.CallOpts{}, logsAndEvents[1], is2ndHalf)
	if err != nil {
		return err
	}
	log.Debugf("Full decoded match log 0: %v", log0)
	log.Debugf("Full decoded match log 1: %v", log1)
	decodedTactics0, err := contracts.Assets.DecodeTactics(&bind.CallOpts{}, b.HomeTeam.tactic)
	if err != nil {
		return err
	}
	decodedTactics1, err := contracts.Assets.DecodeTactics(&bind.CallOpts{}, b.VisitorTeam.tactic)
	if err != nil {
		return err
	}
	log.Debugf("Decoded tactics 0: %v", decodedTactics0)
	log.Debugf("Decoded tactics 1: %v", decodedTactics1)

	generatedEvents, err := matchevents.Generate(
		b.Seed,
		b.HomeTeam.TeamID,
		b.VisitorTeam.TeamID,
		log0,
		log1,
		logsAndEvents,
		decodedTactics0.Lineup,
		decodedTactics1.Lineup,
		decodedTactics0.Substitutions,
		decodedTactics1.Substitutions,
		decodedTactics0.SubsRounds,
		decodedTactics1.SubsRounds,
		is2ndHalf,
	)
	if err != nil {
		return err
	}
	b.Events = append(b.Events, generatedEvents...)
	return nil
}
