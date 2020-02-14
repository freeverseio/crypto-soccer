package engine

import (
	"database/sql"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

const NoOutOfGamePlayer = uint8(14)
const RedCard = uint8(3)
const SoftInjury = uint8(1)
const HardInjury = uint8(2)

type Team struct {
	storage.Team
	Players    [25]*Player
	AssignedTP *big.Int
}

func NewTeam() *Team {
	var team Team
	team.Team = *storage.NewTeam()
	for i := range team.Players {
		team.Players[i] = NewPlayer()
	}
	team.AssignedTP = big.NewInt(0)
	return &team
}

func (b Team) ToStorage(contracts contracts.Contracts, tx *sql.Tx) error {
	for _, player := range b.Players {
		if player.IsNil() {
			continue
		}
		stoPlayer, err := player.ToStorage(contracts)
		if err != nil {
			return err
		}
		if err := stoPlayer.Update(tx); err != nil {
			return err
		}
	}
	return b.Update(tx)
}

func (b Team) Skills() [25]*big.Int {
	var skills [25]*big.Int
	for i := range skills {
		skills[i] = b.Players[i].Skills()
	}
	return skills
}

func DefaultTactic() *big.Int {
	tactic, _ := new(big.Int).SetString("340596594427581673436941882753025", 10)
	return tactic
}

func (b *Team) SetSkills(contracts contracts.Contracts, skills [25]*big.Int) {
	for i := range skills {
		b.Players[i].SetSkills(skills[i])
	}
}
