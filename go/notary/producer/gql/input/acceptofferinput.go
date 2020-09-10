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
	"github.com/graph-gophers/graphql-go"
)

type AcceptOfferInput struct {
	Signature                           string
	PlayerId                            string
	CurrencyId                          int32
	Price                               int32
	Rnd                                 int32
	ValidUntil                          string
	AuctionDurationAfterOfferIsAccepted string
	OfferId                             graphql.ID
}

func (b AcceptOfferInput) AuctionID() (graphql.ID, error) {
	hash, err := b.Hash()
	if err != nil {
		return graphql.ID(""), err
	}
	return graphql.ID(hash.String()[2:]), nil
}
func (b AcceptOfferInput) Hash() (common.Hash, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}
	validUntil, err := strconv.ParseUint(b.ValidUntil, 10, 32)
	if err != nil {
		return common.Hash{}, err
	}
	auctionDurationAfterOfferIsAccepted, err := strconv.ParseUint(b.AuctionDurationAfterOfferIsAccepted, 10, 32)
	if err != nil {
		return common.Hash{}, err
	}

	sellPlayerDigest, err := signer.ComputeSellPlayerDigest(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		uint32(validUntil),
		uint32(auctionDurationAfterOfferIsAccepted),
		playerId,
	)
	return sellPlayerDigest, err
}

func (b AcceptOfferInput) VerifySignature() (bool, error) {
	hash, err := b.Hash()
	if err != nil {
		return false, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return false, err
	}
	return helper.VerifySignature(hash, sign)
}

func (b AcceptOfferInput) SignerAddress() (common.Address, error) {
	hash, err := b.Hash()
	if err != nil {
		return common.Address{}, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return helper.AddressFromSignature(hash, sign)
}

func (b AcceptOfferInput) IsSignerOwner(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress()
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

func (b AcceptOfferInput) IsValidForBlockchain(contracts contracts.Contracts) (bool, error) {
	var err error
	var sig [2][32]byte
	var sigV uint8
	sig[0], sig[1], sigV, err = signer.RSV(b.Signature)
	if err != nil {
		return false, err
	}

	sellerHiddenPrice, err := signer.HashPrivateMsg(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
	)
	if err != nil {
		return false, err
	}

	validUntil, err := strconv.ParseUint(b.ValidUntil, 10, 32)
	if err != nil {
		return false, errors.New("invalid valid until")
	}
	auctionDurationAfterOfferIsAccepted, err := strconv.ParseUint(b.AuctionDurationAfterOfferIsAccepted, 10, 32)
	if err != nil {
		return false, errors.New("invalid auctionDurationAfterOfferIsAccepted")
	}

	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return false, errors.New("invalid playerId")
	}
	isValid, _, err := contracts.Market.AreFreezePlayerRequirementsOK(
		&bind.CallOpts{},
		sellerHiddenPrice,
		playerId,
		sig,
		sigV,
		uint32(validUntil),
		uint32(auctionDurationAfterOfferIsAccepted),
	)
	if err != nil {
		return false, err
	}

	return isValid, nil
}
