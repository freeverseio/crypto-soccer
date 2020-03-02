package engine

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

type Player struct {
	storage.Player
}

func NewPlayer() *Player {
	player := Player{}
	player.EncodedSkills = big.NewInt(0)
	return &player
}

func (b Player) IsNil() bool {
	return b.EncodedSkills.Cmp(big.NewInt(0)) == 0
}

func NewPlayerFromStorage(stoPlayer storage.Player) *Player {
	player := Player{}
	player.Player = stoPlayer
	return &player
}

func (b *Player) SetSkills(contracts contracts.Contracts, skills *big.Int) {
	b.EncodedSkills = new(big.Int).Set(skills)
	opts := &bind.CallOpts{}

	// var err error
	SK_SHO := uint8(0)
	SK_SPE := uint8(1)
	SK_PAS := uint8(2)
	SK_DEF := uint8(3)
	SK_END := uint8(4)
	decodedSkills, err := contracts.Utils.FullDecodeSkills(opts, b.EncodedSkills)
	if err != nil {
		panic("TODO: return err" + err.Error())
	}
	b.Potential = uint64(decodedSkills.BirthTraits[0])
	b.Defence = uint64(decodedSkills.Skills[SK_DEF])
	b.Speed = uint64(decodedSkills.Skills[SK_SPE])
	b.Pass = uint64(decodedSkills.Skills[SK_PAS])
	b.Shoot = uint64(decodedSkills.Skills[SK_SHO])
	b.Endurance = uint64(decodedSkills.Skills[SK_END])
	b.RedCard = decodedSkills.Aligned1stSubst1stRedCardLastGame[2]
	b.InjuryMatchesLeft = decodedSkills.GenerationGamesNonStopInjuryWeeks[2]
}

func (b Player) Skills() *big.Int {
	return new(big.Int).Set(b.EncodedSkills)
}
