package input

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/graph-gophers/graphql-go"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type CreateOfferInput struct {
	Signature   string
	PlayerId    string
	CurrencyId  int32
	Price       int32
	Rnd         int32
	ValidUntil  string
	BuyerTeamId string
	Seller      string
}

func (b CreateOfferInput) ID(contracts contracts.Contracts) (graphql.ID, error) {
	teamId, _ := new(big.Int).SetString(b.BuyerTeamId, 10)
	if teamId == nil {
		return graphql.ID(""), errors.New("invalid teamId")
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return graphql.ID(""), errors.New("invalid validUntil")
	}
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return graphql.ID(""), errors.New("invalid playerId")
	}
	dummyValidUntil := int64(0)
	// the validUntil of the offer becomes the offerValidUntil of the auction
	// the seller adds the new validUntil
	auctionId, err := signer.ComputeAuctionId(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		dummyValidUntil,
		validUntil,
		playerId,
	)
	if err != nil {
		return graphql.ID(""), errors.New("invalid validUntil")
	}
	return graphql.ID(auctionId.String()[2:]), nil
}

func (b CreateOfferInput) Hash(contracts contracts.Contracts) (common.Hash, error) {
	teamId, _ := new(big.Int).SetString(b.BuyerTeamId, 10)
	if teamId == nil {
		return common.Hash{}, errors.New("invalid teamId")
	}
	validUntil, err := strconv.ParseInt(b.ValidUntil, 10, 64)
	if err != nil {
		return common.Hash{}, errors.New("invalid validUntil")
	}
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}
	dummyValidUntil := int64(0)
	// the validUntil of the offer becomes the offerValidUntil of the auction
	// the seller adds the new validUntil
	auctionId, err := signer.ComputeAuctionId(
		uint8(b.CurrencyId),
		big.NewInt(int64(b.Price)),
		big.NewInt(int64(b.Rnd)),
		dummyValidUntil,
		validUntil,
		playerId,
	)
	if err != nil {
		return common.Hash{}, errors.New("invalid validUntil")
	}
	dummyExtraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)

	hash, err := signer.HashBidMessageFromAuctionId(
		contracts.Market,
		auctionId,
		dummyExtraPrice,
		dummyRnd,
		teamId,
	)
	if err != nil {
		return common.Hash{}, err
	}
	return common.Hash(hash), nil
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
	return helper.AddressFromHashAndSignature(hash, sign)
}

func (b CreateOfferInput) SignerAlreadyOwnsPlayer(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress(contracts)
	if err != nil {
		return false, err
	}
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return false, errors.New("invalid BuyerTeamId")
	}
	owner, err := contracts.Market.GetOwnerTeam(&bind.CallOpts{}, playerId)
	if err != nil {
		return false, err
	}
	return signerAddress == owner, nil
}

func (b CreateOfferInput) GetOwnerOfRequestedPlayer(contracts contracts.Contracts) (common.Address, error) {
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

func (b CreateOfferInput) IsAddrTeamOwner(contracts contracts.Contracts, addr common.Address) (bool, error) {
	teamId, _ := new(big.Int).SetString(b.BuyerTeamId, 10)
	if teamId == nil {
		return false, errors.New("invalid teamId")
	}
	owner, err := contracts.Market.GetOwnerTeam(&bind.CallOpts{}, teamId)
	if err != nil {
		return false, err
	}
	return addr == owner, nil
}

func (b CreateOfferInput) IsSignerTeamOwner(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress(contracts)
	if err != nil {
		return false, err
	}
	return b.IsAddrTeamOwner(contracts, signerAddress)
}

func (b CreateOfferInput) IsPlayerFrozen(contracts contracts.Contracts) (bool, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return false, errors.New("invalid playerId")
	}
	isFrozen, err := contracts.Market.IsPlayerFrozenInAnyMarket(&bind.CallOpts{}, playerId)
	if err != nil {
		return false, err
	}
	return isFrozen, nil
}

func (b CreateOfferInput) IsPlayerOnSale(tx storage.Tx) (bool, error) {
	playerId, _ := new(big.Int).SetString(b.PlayerId, 10)
	if playerId == nil {
		return false, errors.New("invalid playerId")
	}

	auctions, err := tx.AuctionsByPlayerId(b.PlayerId)
	if err != nil {
		return false, err
	}

	isOnSale := false

	for _, auction := range auctions {
		if (auction.State != storage.AuctionCancelled) && (auction.State != storage.AuctionEnded) && (auction.State != storage.AuctionFailed) {
			isOnSale = true
		}
	}
	return isOnSale, nil
}
