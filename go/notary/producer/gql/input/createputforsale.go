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

type CreatePutPlayerForSaleInput struct {
	Signature  string
	PlayerId   string
	CurrencyId int32
	Price      int32
	Rnd        int32
	ValidUntil string
}

func (b CreatePutPlayerForSaleInput) ID() (common.Hash, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return common.Hash{}, err
	}
	offerValidUntil := int64(0)
	auctionId, err := signer.ComputeAuctionId(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		validUntil,
		offerValidUntil,
		playerId,
	)
	return auctionId, nil
}

func (b CreatePutPlayerForSaleInput) SellerDigest() (common.Hash, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return common.Hash{}, err
	}
	offerValidUntil := int64(0)

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

func (b CreatePutPlayerForSaleInput) SignerAddress() (common.Address, error) {
	sellerDigest, err := b.SellerDigest()
	if err != nil {
		return common.Address{}, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return helper.AddressFromHashAndSignature(sellerDigest, sign)
}

func (b CreatePutPlayerForSaleInput) IsSignerOwnerOfPlayer(contracts contracts.Contracts) (bool, error) {
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

func (b CreatePutPlayerForSaleInput) ValidForBlockchainFreeze(contracts contracts.Contracts) (bool, error) {
	var err error
	var sig [2][32]byte
	var sigV uint8
	sig[0], sig[1], sigV, err = signer.RSV(b.Signature)
	if err != nil {
		return false, err
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return false, errors.New("invalid valid until")
	}
	offerValidUntil := int64(0)

	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return false, errors.New("invalid playerId")
	}
	sellerHiddenPrice, err := signer.HidePrice(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
	)
	if err != nil {
		return false, errors.New("invalid valid auctionId")
	}

	isValid, err := contracts.Market.AreFreezePlayerRequirementsOK(
		&bind.CallOpts{},
		sellerHiddenPrice,
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
