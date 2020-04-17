package input

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
)

type CreateBidInput struct {
	Signature  string
	Auction    graphql.ID
	ExtraPrice int32
	Rnd        int32
	TeamId     string
}

func (b CreateBidInput) Hash(
	contracts contracts.Contracts,
) (common.Hash, error) {
	teamId, _ := new(big.Int).SetString(b.TeamId, 10)
	if teamId == nil {
		return common.Hash{}, errors.New("invalid teamId")
	}
	isOffer2StartAuction := false
	var auctionHashMsg [32]byte
	copy(auctionHashMsg[:], b.Auction)
	hash, err := signer.HashBidMessage2(
		contracts.Market,
		auctionHashMsg,
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

func (b CreateBidInput) VerifySignature(
	contracts contracts.Contracts,
	auction storage.Auction,
) (bool, error) {
	hash, err := b.Hash(contracts)
	if err != nil {
		return false, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return false, err
	}
	return signer.VerifySignature(hash[:], sign)
}
