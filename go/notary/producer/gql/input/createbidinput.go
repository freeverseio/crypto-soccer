package input

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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

func (b CreateBidInput) ID() string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s%s%d%d%s", b.Signature, string(b.Auction), b.ExtraPrice, b.Rnd, b.TeamId)))
	return hex.EncodeToString(h.Sum(nil))
}

func (b CreateBidInput) Hash(
	contracts contracts.Contracts,
	auction storage.Auction,
) (common.Hash, error) {
	teamId, _ := new(big.Int).SetString(b.TeamId, 10)
	if teamId == nil {
		return common.Hash{}, errors.New("invalid teamId")
	}
	playerId, _ := new(big.Int).SetString(auction.PlayerID, 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}
	isOffer2StartAuction := false
	hash, err := signer.HashBidMessage(
		contracts.Market,
		uint8(auction.CurrencyID),
		big.NewInt(int64(auction.Price)),
		big.NewInt(int64(auction.Rnd)),
		auction.ValidUntil,
		playerId,
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
	hash, err := b.Hash(contracts, auction)
	if err != nil {
		return false, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return false, err
	}
	return signer.VerifySignature(hash[:], sign)
}
