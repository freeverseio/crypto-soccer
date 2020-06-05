package worldplayer

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/utils"
	"github.com/graph-gophers/graphql-go"
)

const nGoalKeepers = 3
const nDefenders = 9
const nMidfielders = 9
const nAttackers = 9

type WorldPlayer struct {
	playerId          graphql.ID
	name              string
	dayOfBirth        int32
	preferredPosition string
	defence           int32
	speed             int32
	pass              int32
	shoot             int32
	endurance         int32
	potential         int32
	validUntil        string
	countryOfBirth    string
	race              string
	productId         string
}

func NewWorldPlayer(
	playerId graphql.ID,
	name string,
	dayOfBirth int32,
	preferredPosition string,
	defence int32,
	speed int32,
	pass int32,
	shoot int32,
	endurance int32,
	potential int32,
	validUntil string,
	countryOfBirth string,
	race string,
	productId string,
) *WorldPlayer {
	player := WorldPlayer{}
	player.playerId = playerId
	player.name = name
	player.dayOfBirth = dayOfBirth
	player.preferredPosition = preferredPosition
	player.defence = defence
	player.speed = speed
	player.pass = pass
	player.shoot = shoot
	player.endurance = endurance
	player.potential = potential
	player.validUntil = validUntil
	player.countryOfBirth = countryOfBirth
	player.race = race
	player.productId = productId
	return &player
}

func (b WorldPlayer) PlayerId() graphql.ID {
	return b.playerId
}

func (b WorldPlayer) Name() string {
	return b.name
}

func (b WorldPlayer) CountryOfBirth() string {
	return b.countryOfBirth
}

func (b WorldPlayer) Race() string {
	return b.race
}

func (b WorldPlayer) ValidUntil() string {
	return b.validUntil
}

func (b WorldPlayer) DayOfBirth() int32 {
	return b.dayOfBirth
}

func (b WorldPlayer) PreferredPosition() string {
	return b.preferredPosition
}

func (b WorldPlayer) Defence() int32 {
	return b.defence
}

func (b WorldPlayer) Speed() int32 {
	return b.speed
}

func (b WorldPlayer) Pass() int32 {
	return b.pass
}

func (b WorldPlayer) Shoot() int32 {
	return b.shoot
}

func (b WorldPlayer) Endurance() int32 {
	return b.endurance
}

func (b WorldPlayer) Potential() int32 {
	return b.potential
}

func (b WorldPlayer) ProductId() string {
	return b.productId
}

func CreateWorldPlayerBatch(
	contr contracts.Contracts,
	namesdb *names.Generator,
	value int64,
	maxPotential uint8,
	teamId string,
	epoch int64,
) ([]*WorldPlayer, error) {
	result := []*WorldPlayer{}

	epochDays := epoch / (3600 * 24)
	epochWeeks := epochDays / 7
	id, _ := new(big.Int).SetString(teamId, 10)
	if id == nil {
		return nil, errors.New("invalid teamId")
	}

	timezone, countryIdxInTZ, _, err := contr.Market.DecodeTZCountryAndVal(&bind.CallOpts{}, id)
	if err != nil {
		return nil, err
	}

	playerValue := big.NewInt(value)
	worldPlayers, err := contr.Privileged.CreateBuyNowPlayerIdBatch(
		&bind.CallOpts{},
		playerValue,
		maxPotential,
		id,
		[4]uint8{
			nGoalKeepers,
			nDefenders,
			nMidfielders,
			nAttackers,
		},
		big.NewInt(epochDays),
		timezone,
		countryIdxInTZ,
	)
	if err != nil {
		return result, err
	}

	for i := range worldPlayers.PlayerIdArray {
		isSellable, err := isSellable(contr, worldPlayers.PlayerIdArray[i])
		if err != nil {
			return nil, err
		}
		if !isSellable {
			continue
		}

		playerId := graphql.ID(worldPlayers.PlayerIdArray[i].String())
		leftishness := worldPlayers.BirthTraitsArray[i][contracts.BirthTraitsLeftishnessIdx]
		forwardness := worldPlayers.BirthTraitsArray[i][contracts.BirthTraitsForwardnessIdx]
		generation := uint8(0)
		name, countryOfBirth, race, err := namesdb.GeneratePlayerFullName(worldPlayers.PlayerIdArray[i], generation, timezone, countryIdxInTZ.Uint64())
		if err != nil {
			return nil, err
		}
		dayOfBirth := int32(worldPlayers.DayOfBirthArray[i])
		preferredPosition, err := utils.PreferredPosition(forwardness, leftishness)
		if err != nil {
			return nil, err
		}
		defence := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsDefenceIdx])
		pass := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsPassIdx])
		speed := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsSpeedIdx])
		shoot := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsShootIdx])
		endurance := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsEnduranceIdx])
		potential := int32(worldPlayers.BirthTraitsArray[i][contracts.BirthTraitsPotentialIdx])
		validUntil := strconv.FormatInt((epochWeeks+1)*24*3600*7, 10) // valid 1 week
		productId := "player_tier_0"
		worldPlayer := NewWorldPlayer(
			playerId,
			name,
			dayOfBirth,
			preferredPosition,
			defence,
			speed,
			pass,
			shoot,
			endurance,
			potential,
			validUntil,
			countryOfBirth,
			race,
			productId,
		)
		result = append(result, worldPlayer)
	}

	return result, nil
}

func isSellable(contr contracts.Contracts, playerId *big.Int) (bool, error) {
	teamId, err := contr.Market.GetCurrentTeamIdFromPlayerId(
		&bind.CallOpts{},
		playerId,
	)
	if err != nil {
		return false, err
	}
	return teamId.Cmp(big.NewInt(contracts.AccademyTeamId)) == 0, nil
}
