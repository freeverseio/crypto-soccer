package engine

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/matchevents"

	log "github.com/sirupsen/logrus"
)

type Match struct {
	Seed            [32]byte
	StartTime       *big.Int
	HomeTeam        Team
	VisitorTeam     Team
	HomeGoals       uint8
	VisitorGoals    uint8
	HomeMatchLog    *big.Int
	VisitorMatchLog *big.Int
	Events          matchevents.MatchEvents
}

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

func (b *Match) Play1stHalf(contracts contracts.Contracts) error {
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := false
	matchSeed, err := b.generateMatchSeed(contracts)
	if err != nil {
		return err
	}
	matchLogs, err := contracts.Engine.PlayHalfMatch(
		&bind.CallOpts{},
		matchSeed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
		[2]*big.Int{b.HomeMatchLog, b.VisitorMatchLog},
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	)
	if err != nil {
		return err
	}
	b.HomeMatchLog = matchLogs[0]
	b.VisitorMatchLog = matchLogs[1]
	b.HomeGoals, b.VisitorGoals, err = b.getGoals(contracts, matchLogs)
	if err != nil {
		return err
	}
	if err = b.HomeTeam.Evolve(contracts, b.HomeMatchLog, b.StartTime, is2ndHalf); err != nil {
		return err
	}
	if err = b.VisitorTeam.Evolve(contracts, b.VisitorMatchLog, b.StartTime, is2ndHalf); err != nil {
		return err
	}
	if err = b.processMatchEvents(contracts, is2ndHalf); err != nil {
		return err
	}
	return nil
}

func (b *Match) Play2ndHalf(contracts contracts.Contracts) error {
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := true
	matchSeed, err := b.generateMatchSeed(contracts)
	if err != nil {
		return err
	}
	logs, err := contracts.Evolution.Play2ndHalfAndEvolve(
		&bind.CallOpts{},
		matchSeed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
		[2]*big.Int{b.HomeMatchLog, b.VisitorMatchLog},
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	)
	if err != nil {
		return err
	}
	if err = b.processMatchEvents(contracts, is2ndHalf); err != nil {
		return err
	}
	b.HomeGoals, b.VisitorGoals, err = b.getGoals(contracts, logs)
	if err != nil {
		return err
	}
	b.HomeMatchLog = logs[0]
	b.VisitorMatchLog = logs[1]
	if err = b.HomeTeam.Evolve(contracts, logs[0], b.StartTime, is2ndHalf); err != nil {
		return err
	}
	if err = b.VisitorTeam.Evolve(contracts, logs[1], b.StartTime, is2ndHalf); err != nil {
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

func (b *Match) generateMatchSeed(contracts contracts.Contracts) (*big.Int, error) {
	matchSeed, err := contracts.Engine.GenerateMatchSeed(&bind.CallOpts{}, b.Seed, b.HomeTeam.TeamID, b.VisitorTeam.TeamID)
	if err != nil {
		return nil, err
	}
	z := new(big.Int)
	z.SetBytes(matchSeed[:])
	return z, nil
}

func (b *Match) processMatchEvents(contracts contracts.Contracts, is2ndHalf bool) error {
	isHomeStadium := true
	isPlayoff := false
	matchSeed, err := b.generateMatchSeed(contracts)
	if err != nil {
		return err
	}
	seedAndStartTimeAndEvents, err := contracts.Matchevents.PlayHalfMatch(
		&bind.CallOpts{},
		matchSeed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
		[2]*big.Int{b.HomeMatchLog, b.VisitorMatchLog},
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	)
	if err != nil {
		return err
	}

	events := seedAndStartTimeAndEvents[:]
	log0, err := contracts.Utilsmatchlog.FullDecodeMatchLog(&bind.CallOpts{}, seedAndStartTimeAndEvents[0], is2ndHalf)
	if err != nil {
		return err
	}
	log1, err := contracts.Utilsmatchlog.FullDecodeMatchLog(&bind.CallOpts{}, seedAndStartTimeAndEvents[1], is2ndHalf)
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
		matchSeed,
		log0,
		log1,
		events,
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

func (b *Match) updateTeamLeaderBoard() error {
	log.Warning("TODO Match::updateTeamLeaderBoard uninplemented")
	return nil

	// b.HomeTeam.GoalsForward += uint32(*b.Match.HomeGoals)
	// b.HomeTeam.GoalsAgainst += uint32(*b.Match.VisitorGoals)
	// b.VisitorTeam.GoalsForward += uint32(*b.Match.VisitorGoals)
	// b.VisitorTeam.GoalsAgainst += uint32(*b.Match.HomeGoals)

	// deltaGoals := int(*b.Match.HomeGoals) - int(*b.Match.VisitorGoals)
	// if deltaGoals > 0 {
	// 	b.HomeTeam.W++
	// 	b.VisitorTeam.L++
	// 	b.HomeTeam.Points += 3
	// } else if deltaGoals < 0 {
	// 	b.HomeTeam.L++
	// 	b.VisitorTeam.W++
	// 	b.VisitorTeam.Points += 3
	// } else {
	// 	b.HomeTeam.D++
	// 	b.VisitorTeam.D++
	// 	b.HomeTeam.Points++
	// 	b.VisitorTeam.Points++
	// }

	// return nil
}
