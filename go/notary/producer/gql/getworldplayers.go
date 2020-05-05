package gql

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

const nGoalKeepers = 3
const nDefenders = 9
const nMidfielders = 9
const nAttackers = 9

func (b *Resolver) GetWorldPlayers(args struct{ Input input.GetWorldPlayersInput }) ([]*WorldPlayer, error) {
	log.Debugf("GetWorldPlayers %v", args)

	if b.ch == nil {
		return nil, errors.New("internal error: no channel")
	}

	hash, err := args.Input.Hash()
	if err != nil {
		return nil, err
	}
	sign, err := hex.DecodeString(args.Input.Signature)
	if err != nil {
		return nil, err
	}

	isValid, err := input.VerifySignature(hash, sign)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, errors.New("Invalid signature")
	}

	sender, err := input.AddressFromSignature(hash, sign)
	if err != nil {
		return nil, err
	}
	log.Infof("TODO check sender is %v", sender.Hex())

	value := int64(3000) // TODO

	return CreateWorldPlayerBatch(
		b.contracts,
		value,
		string(args.Input.TeamId),
		time.Now().Unix(),
	)
}

func CreateWorldPlayerBatch(
	contr contracts.Contracts,
	value int64,
	teamId string,
	epoch int64,
) ([]*WorldPlayer, error) {
	result := []*WorldPlayer{}
	epochDays := epoch / (3600 * 24)
	epochWeeks := epochDays / 7
	seed, _ := new(big.Int).SetString(teamId, 10)
	if seed == nil {
		return nil, errors.New("invalid teamId")
	}

	playerValue := big.NewInt(value)
	worldPlayers, err := contr.Privileged.CreateBuyNowPlayerIdBatch(
		&bind.CallOpts{},
		playerValue,
		seed,
		[4]uint8{
			nGoalKeepers,
			nDefenders,
			nMidfielders,
			nAttackers,
		},
		big.NewInt(epochDays),
	)
	if err != nil {
		return result, err
	}

	for i := range worldPlayers.PlayerIdArray {
		playerId := graphql.ID(worldPlayers.PlayerIdArray[i].String())
		name := "" // TODO
		dayOfBirth := int32(worldPlayers.DayOfBirthArray[i])
		preferredPosition := "" // TODO
		defence := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsDefenceIdx])
		pass := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsPassIdx])
		speed := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsSpeedIdx])
		shoot := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsShootIdx])
		endurance := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsEnduranceIdx])
		potential := int32(worldPlayers.BirthTraitsArray[i][contracts.BirthTraitsPotentialIdx])
		validUntil := strconv.FormatInt((epochWeeks+1)*24*3600*7, 10) // valid 1 week
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
		)
		result = append(result, worldPlayer)
	}

	return result, nil
}
