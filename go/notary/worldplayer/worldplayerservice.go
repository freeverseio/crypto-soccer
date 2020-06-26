package worldplayer

import (
	"errors"
	"hash/fnv"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/names"
	"github.com/freeverseio/crypto-soccer/go/utils"
	"github.com/graph-gophers/graphql-go"
)

type WorldPlayerService struct {
	contracts contracts.Contracts
	namesdb   *names.Generator
}

func NewWorldPlayerService(contracts contracts.Contracts, namesdb *names.Generator) *WorldPlayerService {
	return &WorldPlayerService{
		contracts: contracts,
		namesdb:   namesdb,
	}
}

func GetSeedForWorldPlayers(teamId string, epoch int64) string {
	epochHalfDays := epoch / (3600 * 24 * 2)
	return teamId + strconv.FormatUint(uint64(epochHalfDays), 10)
}

func (b WorldPlayerService) CreateBatch(teamId string, epoch int64) ([]*WorldPlayer, error) {
	seed := GetSeedForWorldPlayers(teamId, epoch)
	distribution := GenerateBatchDistribution(seed)

	batch := []*WorldPlayer{}
	for _, tier := range distribution {
		batchByTier, err := b.createBatchByTier(
			teamId,
			epoch,
			tier,
		)
		if err != nil {
			return nil, err
		}
		batch = append(batch, batchByTier...)
	}
	return batch, nil
}

func (b WorldPlayerService) GetWorldPlayer(
	playerId string,
	teamId string,
	epoch int64,
) (*WorldPlayer, error) {
	players, err := b.CreateBatch(teamId, epoch)
	if err != nil {
		return nil, err
	}
	for _, player := range players {
		if string(player.PlayerId()) == playerId {
			return player, nil
		}
	}
	return nil, nil
}

func int_hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func GenerateRnd(seed *big.Int, salt string, max_val uint64) uint64 {
	var result uint64 = int_hash(seed.String() + salt)
	if max_val == 0 {
		return result
	}
	return result % max_val
}

func (b WorldPlayerService) createBatchByTier(
	teamId string,
	epoch int64,
	tier WorldPlayersTier,
) ([]*WorldPlayer, error) {
	result := []*WorldPlayer{}

	epochDays := epoch / (3600 * 24)
	id, _ := new(big.Int).SetString(teamId, 10)
	if id == nil {
		return nil, errors.New("invalid teamId")
	}

	timezone, countryIdxInTZ, _, err := b.contracts.Market.DecodeTZCountryAndVal(&bind.CallOpts{}, id)
	if err != nil {
		return nil, err
	}

	worldPlayers, err := b.contracts.Privileged.CreateBuyNowPlayerIdBatch(
		&bind.CallOpts{},
		big.NewInt(tier.Value),
		tier.MaxPotential,
		id,
		[4]uint8{
			tier.GoalKeepersCount,
			tier.DefendersCount,
			tier.MidfieldersCount,
			tier.AttackersCount,
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
		validUntil := strconv.FormatInt(epoch+tier.Duration, 10)
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
			tier.ProductId,
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
