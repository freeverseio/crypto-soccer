package engine

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

type Player struct {
	sto storage.Player
}

func NewNullPlayer() *Player {
	player := Player{}
	player.sto.EncodedSkills = big.NewInt(0)
	return &player
}

func NewPlayerFromStorage(stoPlayer storage.Player) *Player {
	player := Player{}
	player.sto = stoPlayer
	return &player
}

func (b Player) ToStorage(contracts contracts.Contracts) (storage.Player, error) {
	return b.sto, nil
}

func (b Player) IsNull() bool {
	return b.sto.EncodedSkills.Cmp(big.NewInt(0)) == 0
}

func (b *Player) SetSkills(contracts contracts.Contracts, skills *big.Int) error {
	b.sto.EncodedSkills = new(big.Int).Set(skills)
	return b.decodeSkills(contracts)
}

func NewPlayer(
	contracts contracts.Contracts,
	playerID *big.Int,
	defence uint16,
	speed uint16,
	endurance uint16,
	pass uint16,
	shoot uint16,
	dayOfBirthUnix uint16,
	generation uint8,
	potential uint8,
	forwardness uint8,
	leftishness uint8,
	aggressiveness uint8,
	alignedEndOfLastHalf bool,
	redCardLastGame bool,
	gamesNonStopping uint8,
	injuryWeeksLeft uint8,
	substitutedLastHalf bool,
) (*Player, error) {
	var err error
	player := Player{}
	sumSkills := uint32(defence) + uint32(speed) + uint32(endurance) + uint32(pass) + uint32(shoot)
	player.sto.EncodedSkills, err = contracts.Engine.EncodePlayerSkills(
		&bind.CallOpts{},
		[5]uint16{shoot, speed, pass, defence, endurance},
		big.NewInt(int64(dayOfBirthUnix)),
		generation,
		playerID,
		[4]uint8{potential, forwardness, leftishness, aggressiveness},
		alignedEndOfLastHalf,
		redCardLastGame,
		gamesNonStopping,
		injuryWeeksLeft,
		substitutedLastHalf,
		sumSkills,
	)
	if err != nil {
		return nil, err
	}
	err = player.decodeSkills(contracts)
	return &player, err
}

func NewPlayerFromSkills(contracts contracts.Contracts, skills string) (*Player, error) {
	var player Player
	player.sto.EncodedSkills, _ = new(big.Int).SetString(skills, 10)
	if err := player.decodeSkills(contracts); err != nil {
		return nil, err
	}
	return &player, nil
}

func (b Player) DumpState() string {
	return fmt.Sprintf("skills: %v", b.sto.EncodedSkills)
}

func (b Player) Skills() *big.Int {
	return new(big.Int).Set(b.sto.EncodedSkills)
}

func (b Player) Defence() uint16 {
	return uint16(b.sto.Defence)
}

func (b Player) Speed() uint16 {
	return uint16(b.sto.Speed)
}

func (b Player) Pass() uint16 {
	return uint16(b.sto.Pass)
}

func (b Player) Shoot() uint16 {
	return uint16(b.sto.Shoot)
}

func (b Player) Endurance() uint16 {
	return uint16(b.sto.Endurance)
}

func (b Player) Potential() uint16 {
	return uint16(b.sto.Potential)
}

// TODO change name to DayOfBirth()
func (b Player) BirthDayUnix() uint16 {
	return uint16(b.sto.DayOfBirth)
}

func PlayerAge(birthDayUnix uint16) uint16 {
	nowInDays := time.Now().Unix() / 3600 / 24
	age := uint16((nowInDays - int64(birthDayUnix)) * 7 / 365)
	return age
}

func (b *Player) decodeSkills(contracts contracts.Contracts) error {
	opts := &bind.CallOpts{}
	var err error
	defence, err := contracts.Assets.GetDefence(opts, b.sto.EncodedSkills)
	if err != nil {
		return err
	}
	speed, err := contracts.Assets.GetSpeed(opts, b.sto.EncodedSkills)
	if err != nil {
		return err
	}
	pass, err := contracts.Assets.GetPass(opts, b.sto.EncodedSkills)
	if err != nil {
		return err
	}
	shoot, err := contracts.Assets.GetShoot(opts, b.sto.EncodedSkills)
	if err != nil {
		return err
	}
	endurance, err := contracts.Assets.GetEndurance(opts, b.sto.EncodedSkills)
	if err != nil {
		return err
	}
	potential, err := contracts.Assets.GetPotential(opts, b.sto.EncodedSkills)
	if err != nil {
		return err
	}
	dayOfBirth, err := contracts.Assets.GetBirthDay(opts, b.sto.EncodedSkills)
	if err != nil {
		return err
	}
	b.sto.Defence = defence.Uint64()
	b.sto.Speed = speed.Uint64()
	b.sto.Pass = pass.Uint64()
	b.sto.Shoot = shoot.Uint64()
	b.sto.Endurance = endurance.Uint64()
	b.sto.Potential = potential.Uint64()
	b.sto.DayOfBirth = dayOfBirth.Uint64()
	return nil
}
