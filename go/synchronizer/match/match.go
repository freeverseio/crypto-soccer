package match

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

type Match struct {
	MatchSeed         *big.Int
	StartTime         *big.Int
	HomeTeam          *Team
	VisitorTeam       *Team
	HomeGoals         uint8
	VisitorGoals      uint8
	homeMatchLog      *big.Int
	visitorMatchLog   *big.Int
	NOOUTOFGAMEPLAYER uint8
	REDCARD           uint8
	SOFTINJURY        uint8
	HARDINJURY        uint8
}

func NewMatch(contracts *contracts.Contracts) (*Match, error) {
	match := Match{}
	var err error
	match.MatchSeed = big.NewInt(0)
	match.StartTime = big.NewInt(0)
	if match.HomeTeam, err = NewTeam(contracts); err != nil {
		return nil, err
	}
	match.homeMatchLog = big.NewInt(0)
	if match.VisitorTeam, err = NewTeam(contracts); err != nil {
		return nil, err
	}
	match.visitorMatchLog = big.NewInt(0)
	return &match, nil
}

func (b *Match) Play1stHalf(contracts *contracts.Contracts) error {
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := false
	matchBools := [3]bool{is2ndHalf, isHomeStadium, isPlayoff}
	logs, err := contracts.Engine.PlayHalfMatch(
		&bind.CallOpts{},
		b.MatchSeed,
		b.StartTime,
		[2][25]*big.Int{b.HomeTeam.State(), b.VisitorTeam.State()},
		[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
		[2]*big.Int{b.homeMatchLog, b.visitorMatchLog},
		matchBools,
	)
	if err != nil {
		return err
	}

	homeLog := logs[0]
	visitorLog := logs[1]
	if b.HomeGoals, err = contracts.Evolution.GetNGoals(&bind.CallOpts{}, homeLog); err != nil {
		return err
	}
	if b.VisitorGoals, err = contracts.Evolution.GetNGoals(&bind.CallOpts{}, visitorLog); err != nil {
		return err
	}
	// if err = b.UpdatePlayedByHalf(contracts, b.HomeTeamPlayers, is2ndHalf, b.HomeTeam.tactic, homeLog); err != nil {
	// 	return err
	// }
	// if err = b.UpdatePlayedByHalf(contracts, b.VisitorTeamPlayers, is2ndHalf, b.VisitorTeam.tactic, visitorLog); err != nil {
	// 	return err
	// }
	return nil
}

// func (b *Match) Play2ndHalf(contracts *contracts.Contracts) error {
// 	isHomeStadium := true
// 	isPlayoff := false
// 	is2ndHalf := true
// 	matchBools := [3]bool{is2ndHalf, isHomeStadium, isPlayoff}
// 	logs, err := contracts.Evolution.Play2ndHalfAndEvolve(
// 		&bind.CallOpts{},
// 		b.MatchSeed,
// 		b.StartTime,
// 		[2][25]*big.Int{b.HomeTeam.State(), b.VisitorTeam.State()},
// 		[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
// 		[2]*big.Int{b.HomeMatchLog, b.VisitorMatchLog},
// 		matchBools,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	homeLog := logs[0]
// 	visitorLog := logs[1]
// 	if b.HomeGoals, err = contracts.Evolution.GetNGoals(&bind.CallOpts{}, homeLog); err != nil {
// 		return err
// 	}
// 	if b.VisitorGoals, err = contracts.Evolution.GetNGoals(&bind.CallOpts{}, visitorLog); err != nil {
// 		return err
// 	}
// 	if err = b.UpdatePlayedByHalf(contracts, b.HomeTeamPlayers, is2ndHalf, b.HomeTactic, homeLog); err != nil {
// 		return err
// 	}
// 	if err = b.UpdatePlayedByHalf(contracts, b.VisitorTeamPlayers, is2ndHalf, b.VisitorTactic, visitorLog); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (b *Match) UpdatePlayedByHalf(
	contracts *contracts.Contracts,
	players []*storage.Player,
	is2ndHalf bool,
	tactic *big.Int,
	matchLog *big.Int,
) error {
	decodedTactic, err := contracts.Leagues.DecodeTactics(&bind.CallOpts{}, tactic)
	if err != nil {
		return err
	}
	outOfGamePlayer, err := contracts.Engineprecomp.GetOutOfGamePlayer(&bind.CallOpts{}, matchLog, is2ndHalf)
	if err != nil {
		return err
	}
	outOfGameType, err := contracts.Engineprecomp.GetOutOfGameType(&bind.CallOpts{}, matchLog, is2ndHalf)
	if err != nil {
		return err
	}
	for _, player := range players {
		wasAligned, err := contracts.Engine.WasPlayerAlignedEndOfLastHalf(
			&bind.CallOpts{},
			player.ShirtNumber,
			tactic,
			matchLog,
		)
		if err != nil {
			return err
		}
		player.EncodedSkills, err = contracts.Evolution.SetAlignedEndOfLastHalf(
			&bind.CallOpts{},
			player.EncodedSkills,
			wasAligned,
		)
		if err != nil {
			return err
		}
		if outOfGamePlayer.Int64() != int64(b.NOOUTOFGAMEPLAYER) {
			if outOfGamePlayer.Int64() < 0 || int(outOfGamePlayer.Int64()) >= len(decodedTactic.Lineup) {
				return fmt.Errorf("out of game player unknown %v, tactics %v, matchlog %v", outOfGamePlayer.Int64(), tactic, matchLog)
			}
			if player.ShirtNumber == decodedTactic.Lineup[outOfGamePlayer.Int64()] {
				switch outOfGameType.Int64() {
				case int64(b.REDCARD):
					player.RedCardMatchesLeft = 2
				case int64(b.SOFTINJURY):
					player.InjuryMatchesLeft = 3
				case int64(b.HARDINJURY):
					player.InjuryMatchesLeft = 7
				default:
					return fmt.Errorf("out of game type unknown %v", outOfGameType)
				}
			}
		}
		if is2ndHalf {
			if player.RedCardMatchesLeft > 0 {
				player.RedCardMatchesLeft--
			}
			if player.InjuryMatchesLeft > 0 {
				player.InjuryMatchesLeft--
			}
		}
		// log.Infof("encoded skills %v, redCard %v, injuries %v", player.EncodedSkills, player.RedCardMatchesLeft, player.InjuryMatchesLeft)
		if player.EncodedSkills, err = contracts.Evolution.SetRedCardLastGame(&bind.CallOpts{}, player.EncodedSkills, player.RedCardMatchesLeft != 0); err != nil {
			return err
		}
		if player.EncodedSkills, err = contracts.Evolution.SetInjuryWeeksLeft(&bind.CallOpts{}, player.EncodedSkills, player.InjuryMatchesLeft); err != nil {
			return err
		}
	}
	return nil
}

func GetTeamState(players []*storage.Player) ([25]*big.Int, error) {
	var state [25]*big.Int
	for i := 0; i < 25; i++ {
		state[i] = big.NewInt(0)
	}
	for i := 0; i < len(players); i++ {
		player := players[i]
		playerSkills := player.EncodedSkills
		shirtNumber := player.ShirtNumber
		state[shirtNumber] = playerSkills
	}
	return state, nil
}
