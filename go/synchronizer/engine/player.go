package engine

import (
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

func (b Player) ToStorage(contracts contracts.Contracts) (storage.Player, error) {
	opts := &bind.CallOpts{}
	var err error
	SK_SHO := uint8(0)
	SK_SPE := uint8(1)
	SK_PAS := uint8(2)
	SK_DEF := uint8(3)
	SK_END := uint8(4)
	decodedSkills, err := contracts.Utils.FullDecodeSkills(opts, b.sto.EncodedSkills)
	if err != nil {
		return b.sto, err
	}
	b.sto.Defence = uint64(decodedSkills.Skills[SK_DEF])
	b.sto.Speed = uint64(decodedSkills.Skills[SK_SPE])
	b.sto.Pass = uint64(decodedSkills.Skills[SK_PAS])
	b.sto.Shoot = uint64(decodedSkills.Skills[SK_SHO])
	b.sto.Endurance = uint64(decodedSkills.Skills[SK_END])
	b.sto.RedCard = decodedSkills.Aligned1stSubst1stRedCardLastGame[2]
	b.sto.InjuryMatchesLeft = decodedSkills.GenerationGamesNonStopInjuryWeeks[2]
	return b.sto, nil
}
