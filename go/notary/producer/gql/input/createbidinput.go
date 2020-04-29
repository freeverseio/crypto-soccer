package input

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/contracts"
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

func (b CreateBidInput) ID(contracts contracts.Contracts) (graphql.ID, error) {
	hash, err := b.Hash(contracts)
	if err != nil {
		return graphql.ID(""), err
	}
	return graphql.ID(hash.String()[2:]), nil
}

func (b CreateBidInput) Hash(contracts contracts.Contracts) (common.Hash, error) {
	teamId, _ := new(big.Int).SetString(b.TeamId, 10)
	if teamId == nil {
		return common.Hash{}, errors.New("invalid teamId")
	}
	isOffer2StartAuction := false
	auctionHash := common.HexToHash(string(b.AuctionId))
	hash, err := signer.HashBidMessage2(
		contracts.Market,
		auctionHash,
		big.NewInt(int64(b.ExtraPrice)),
		big.NewInt(int64(b.Rnd)),
		teamId,
		isOffer2StartAuction,
	)
	if err != nil {
		return common.Hash{}, err
	}
	return common.Hash(hash), nil
}

func (b CreateBidInput) VerifySignature(contracts contracts.Contracts) (bool, error) {
	hash, err := b.Hash(contracts)
	if err != nil {
		return false, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return false, err
	}
	return VerifySignature(hash, sign)
}
