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
const value = int64(1000)
const maxPotential = uint8(9)

type WorldPlayerService struct {
	contracts    contracts.Contracts
	namesdb      *names.Generator
	distribution []WorldPlayersTier
}

func NewWorldPlayerService(contracts contracts.Contracts, namesdb *names.Generator) *WorldPlayerService {
	return &WorldPlayerService{
		contracts:    contracts,
		namesdb:      namesdb,
		distribution: GenerateBatchDistribution(),
	}
}

func (b WorldPlayerService) CreateBatch(teamId string, epoch int64) ([]*WorldPlayer, error) {
	result := []*WorldPlayer{}

	epochDays := epoch / (3600 * 24)
	epochWeeks := epochDays / 7
	id, _ := new(big.Int).SetString(teamId, 10)
	if id == nil {
		return nil, errors.New("invalid teamId")
	}

	timezone, countryIdxInTZ, _, err := b.contracts.Market.DecodeTZCountryAndVal(&bind.CallOpts{}, id)
	if err != nil {
		return nil, err
	}

	playerValue := big.NewInt(value)
	worldPlayers, err := b.contracts.Privileged.CreateBuyNowPlayerIdBatch(
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
		isSellable, err := b.isSellable(worldPlayers.PlayerIdArray[i])
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
		name, countryOfBirth, race, err := b.namesdb.GeneratePlayerFullName(worldPlayers.PlayerIdArray[i], generation, timezone, countryIdxInTZ.Uint64())
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

func (b WorldPlayerService) isSellable(playerId *big.Int) (bool, error) {
	teamId, err := b.contracts.Market.GetCurrentTeamIdFromPlayerId(
		&bind.CallOpts{},
		playerId,
	)
	if err != nil {
		return false, err
	}
	return teamId.Cmp(big.NewInt(contracts.AccademyTeamId)) == 0, nil
}
