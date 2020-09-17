package input

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/graph-gophers/graphql-go"
)

type CreateBidInput struct {
	Signature  string
	AuctionId  graphql.ID
	ExtraPrice int32
	Rnd        int32
	TeamId     string
}

func (b CreateBidInput) Hash(contracts contracts.Contracts) (common.Hash, error) {
	teamId, _ := new(big.Int).SetString(b.TeamId, 10)
	if teamId == nil {
		return common.Hash{}, errors.New("invalid teamId")
	}
	auctionId := common.HexToHash(string(b.AuctionId))
	hash, err := signer.HashBidMessageFromAuctionId(
		contracts.Market,
		auctionId,
		big.NewInt(int64(b.ExtraPrice)),
		big.NewInt(int64(b.Rnd)),
		teamId,
	)
	if err != nil {
		return common.Hash{}, err
	}
	return common.Hash(hash), nil
}

func (b CreateBidInput) SignerAddress(contracts contracts.Contracts) (common.Address, error) {
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

func (b CreateBidInput) IsSignerOwnerOfTeam(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress(contracts)
	if err != nil {
		return false, err
	}
	teamId, _ := new(big.Int).SetString(b.TeamId, 10)
	if teamId == nil {
		return false, errors.New("invalid teamd")
	}
	owner, err := contracts.Market.GetOwnerTeam(&bind.CallOpts{}, teamId)
	if err != nil {
		return false, err
	}
	return signerAddress == owner, nil
}
