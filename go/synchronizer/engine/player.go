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

func (b *Player) SetPlayerId(playerId *big.Int) {
	b.PlayerId = playerId
}

func (b *Player) SetSkills(c contracts.Contracts, skills *big.Int) {
	b.EncodedSkills = new(big.Int).Set(skills)
	opts := &bind.CallOpts{}

	decodedSkills, err := c.Utils.FullDecodeSkills(opts, b.EncodedSkills)
	if err != nil {
		panic("TODO: return err" + err.Error())
	}
	b.Potential = uint64(decodedSkills.BirthTraits[contracts.BirthTraitsPotentialIdx])
	b.Defence = uint64(decodedSkills.Skills[contracts.SkillsDefenceIdx])
	b.Speed = uint64(decodedSkills.Skills[contracts.SkillsSpeedIdx])
	b.Pass = uint64(decodedSkills.Skills[contracts.SkillsPassIdx])
	b.Shoot = uint64(decodedSkills.Skills[contracts.SkillsShootIdx])
	b.Endurance = uint64(decodedSkills.Skills[contracts.SkillsEnduranceIdx])
	b.RedCard = decodedSkills.Aligned1stSubst1stRedCardLastGame[2]
	b.InjuryMatchesLeft = decodedSkills.GenerationGamesNonStopInjuryWeeks[2]
	b.Tiredness = int(decodedSkills.GenerationGamesNonStopInjuryWeeks[1])
}

func (b Player) Skills() *big.Int {
	return new(big.Int).Set(b.EncodedSkills)
}
