package engine

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

type Player struct {
	sto storage.Player
}

func NewPlayer() *Player {
	player := Player{}
	player.sto.EncodedSkills = big.NewInt(0)
	return &player
}

func NewPlayerFromStorage(stoPlayer storage.Player) *Player {
	player := Player{}
	player.sto = stoPlayer
	return &player
}

func (b *Player) SetSkills(skills *big.Int) {
	b.sto.EncodedSkills = new(big.Int).Set(skills)
}

func (b Player) Skills() *big.Int {
	return new(big.Int).Set(b.sto.EncodedSkills)
}

func (b Player) DumpState() string {
	return fmt.Sprintf("skills: %v", b.sto.EncodedSkills)
}

func (b Player) ToStorage(contracts contracts.Contracts) (storage.Player, error) {
	opts := &bind.CallOpts{}
	var err error
	SK_SHO := uint8(0)
	SK_SPE := uint8(1)
	SK_PAS := uint8(2)
	SK_DEF := uint8(3)
	SK_END := uint8(4)
	defence, err := contracts.Assets.GetSkill(opts, b.sto.EncodedSkills, SK_DEF)
	if err != nil {
		return b.sto, err
	}
	speed, err := contracts.Assets.GetSkill(opts, b.sto.EncodedSkills, SK_SPE)
	if err != nil {
		return b.sto, err
	}
	pass, err := contracts.Assets.GetSkill(opts, b.sto.EncodedSkills, SK_PAS)
	if err != nil {
		return b.sto, err
	}
	shoot, err := contracts.Assets.GetSkill(opts, b.sto.EncodedSkills, SK_SHO)
	if err != nil {
		return b.sto, err
	}
	endurance, err := contracts.Assets.GetSkill(opts, b.sto.EncodedSkills, SK_END)
	if err != nil {
		return b.sto, err
	}
	b.sto.Defence = defence.Uint64()
	b.sto.Speed = speed.Uint64()
	b.sto.Pass = pass.Uint64()
	b.sto.Shoot = shoot.Uint64()
	b.sto.Endurance = endurance.Uint64()
	return b.sto, nil
}
