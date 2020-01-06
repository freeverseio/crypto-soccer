package match

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
)

type Team struct {
	Players [25]*Player
	tactic  *big.Int
}

func NewTeam(
	contracts *contracts.Contracts,
) (*Team, error) {
	var team Team
	for i := range team.Players {
		team.Players[i] = NewNullPlayer()
	}
	var err error
	if team.tactic, err = DefaultTactic(contracts); err != nil {
		return nil, err
	}
	return &team, nil
}

func (b Team) Skills() [25]*big.Int {
	var states [25]*big.Int
	for i := range states {
		states[i] = b.Players[i].skills
	}
	return states
}

func DefaultTactic(contracts *contracts.Contracts) (*big.Int, error) {
	substitutions := [3]uint8{11, 11, 11}
	substitutionsMinute := [3]uint8{2, 3, 4}
	formation := [14]uint8{0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 25, 25, 25}
	extraAttack := [10]bool{false, false, false, false, false, false, false, false, false, false}
	tacticID := uint8(1)
	tactic, err := contracts.Engine.EncodeTactics(
		&bind.CallOpts{},
		substitutions,
		substitutionsMinute,
		formation,
		extraAttack,
		tacticID,
	)
	if err != nil {
		return nil, err
	}
	return tactic, nil
}
