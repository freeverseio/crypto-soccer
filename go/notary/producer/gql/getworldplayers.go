package gql

import (
	"encoding/hex"
	"errors"
	"math/big"

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

	result := []*WorldPlayer{}

	if b.ch == nil {
		return result, errors.New("internal error: no channel")
	}

	hash, err := args.Input.Hash()
	if err != nil {
		return result, err
	}
	sign, err := hex.DecodeString(args.Input.Signature)
	if err != nil {
		return result, err
	}

	isValid, err := input.VerifySignature(hash, sign)
	if err != nil {
		return result, err
	}
	if !isValid {
		return result, errors.New("Invalid signature")
	}

	sender, err := input.AddressFromSignature(hash, sign)
	if err != nil {
		return result, err
	}
	log.Infof("TODO sender is %v", sender.Hex())

	playerValue := big.NewInt(3000)
	seed, _ := new(big.Int).SetString(string(args.Input.TeamId), 10)
	if seed == nil {
		return result, errors.New("Invalid TeamId")
	}
	worldPlayers, err := b.contracts.Privileged.CreateBuyNowPlayerIdBatch(
		&bind.CallOpts{},
		playerValue,
		seed,
		[4]uint8{
			nGoalKeepers,
			nDefenders,
			nMidfielders,
			nAttackers,
		},
	)
	if err != nil {
		return result, err
	}

	for i := range worldPlayers.PlayerIdArray {
		playerId := graphql.ID(worldPlayers.PlayerIdArray[i].String())
		name := "TODO"
		dayOfBirth := int32(worldPlayers.DayOfBirthArray[i])
		preferredPosition := "TODO"
		defence := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsDefenceIdx])
		pass := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsPassIdx])
		speed := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsSpeedIdx])
		shoot := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsShootIdx])
		endurance := int32(worldPlayers.SkillsVecArray[i][contracts.SkillsEnduranceIdx])
		potential := int32(worldPlayers.BirthTraitsArray[i][contracts.BirthTraitsPotentialIdx])
		validUntil := "TODO"
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
