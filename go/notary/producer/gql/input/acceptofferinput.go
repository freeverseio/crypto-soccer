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
	Signature       string
	PlayerId        string
	CurrencyId      int32
	Price           int32
	Rnd             int32
	ValidUntil      string
	OfferValidUntil string
	OfferId         graphql.ID
}

func (b AcceptOfferInput) AuctionID() (graphql.ID, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return graphql.ID(""), errors.New("invalid playerId")
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return graphql.ID(""), err
	}
	offerValidUntil, err := strconv.ParseInt(b.OfferValidUntil, 10, 64)
	if err != nil {
		return graphql.ID(""), err
	}
	auctionId, err := signer.ComputeAuctionId(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		validUntil,
		offerValidUntil,
		playerId,
	)
	return graphql.ID(auctionId.String()[2:]), nil
}

func (b AcceptOfferInput) SellerDigest() (common.Hash, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return common.Hash{}, err
	}
	offerValidUntil, err := strconv.ParseInt(b.OfferValidUntil, 10, 64)
	if err != nil {
		return common.Hash{}, err
	}

	sellerDigest, err := signer.ComputePutAssetForSaleDigest(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		validUntil,
		offerValidUntil,
		playerId,
	)
	return sellerDigest, err
}

func (b AcceptOfferInput) VerifySignature() (bool, error) {
	hash, err := b.SellerDigest()
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
	hash, err := b.SellerDigest()
	if err != nil {
		return common.Address{}, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return helper.AddressFromHashAndSignature(hash, sign)
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
	if err != nil {
		return false, err
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return false, errors.New("invalid valid until")
	}
	offerValidUntil, err := strconv.ParseInt(b.OfferValidUntil, 10, 64)
	if err != nil {
		return false, errors.New("invalid valid until")
	}
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return false, errors.New("invalid playerId")
	}
	auctionId, err := signer.ComputeAuctionId(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		validUntil,
		offerValidUntil,
		playerId,
	)
	isValid, err := contracts.Market.AreFreezePlayerRequirementsOK(
		&bind.CallOpts{},
		auctionId,
		playerId,
		sig,
		sigV,
		uint32(validUntil),
		uint32(offerValidUntil),
	)
	if err != nil {
		return false, err
	}

	return isValid, nil
}
