package input

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
)

type CreateOfferInput struct {
	Signature  string
	PlayerId   string
	CurrencyId int32
	Price      int32
	Rnd        int32
	ValidUntil string
	TeamId     string
	Seller     string
}

func (b CreateOfferInput) Hash(contracts contracts.Contracts) (common.Hash, error) {
	teamId, _ := new(big.Int).SetString(b.TeamId, 10)
	if teamId == nil {
		return common.Hash{}, errors.New("invalid teamId")
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return common.Hash{}, errors.New("invalid validUntil")
	}
	playerId, err := strconv.ParseInt(b.PlayerId, 10, 64)
	if err != nil {
		return common.Hash{}, errors.New("invalid playerId")
	}
	dummyRnd := int64(0)

	hash, err := signer.HashBidMessage(
		contracts.Market,
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		validUntil,
		big.NewInt(int64(playerId)),
		big.NewInt(0),
		big.NewInt(dummyRnd),
		teamId,
		true,
	)
	if err != nil {
		return common.Hash{}, err
	}
	return common.Hash(hash), nil
}

func (b CreateOfferInput) VerifySignature(contracts contracts.Contracts) (bool, error) {
	hash, err := b.Hash(contracts)
	if err != nil {
		return false, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return false, err
	}
	return helper.VerifySignature(hash, sign)
}

func (b CreateOfferInput) SignerAddress(contracts contracts.Contracts) (common.Address, error) {
	hash, err := b.Hash(contracts)
	if err != nil {
		return common.Address{}, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return helper.AddressFromSignature(hash, sign)
}

func (b CreateOfferInput) IsSignerOwner(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress(contracts)
	if err != nil {
		return false, err
	}
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return false, errors.New("invalid playerId")
	}
	owner, err := contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	if err != nil {
		return false, err
	}
	return signerAddress == owner, nil
}

func (b CreateOfferInput) GetOwner(contracts contracts.Contracts) (common.Address, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return common.Address{}, errors.New("invalid playerId")
	}
	owner, err := contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	if err != nil {
		return common.Address{}, err
	}
	return owner, nil
}

func (b CreateOfferInput) IsSignerTeamOwner(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress(contracts)
	if err != nil {
		return false, err
	}
	teamId, _ := new(big.Int).SetString(b.TeamId, 10)
	if teamId == nil {
		return false, errors.New("invalid teamId")
	}
	owner, err := contracts.Market.GetOwnerTeam(&bind.CallOpts{}, teamId)
	if err != nil {
		return false, err
	}
	return signerAddress == owner, nil
}
