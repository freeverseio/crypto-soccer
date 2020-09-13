package storage

import (
	"errors"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/notary/signer"
)

type AuctionState string

const (
	AuctionAssetFrozen        AuctionState = "asset_frozen"
	AuctionPaying             AuctionState = "paying"
	AuctionWithdrableBySeller AuctionState = "withadrable_by_seller"
	AuctionWithdrableByBuyer  AuctionState = "withadrable_by_buyer"
	AuctionStarted            AuctionState = "started"
	AuctionEnded              AuctionState = "ended"
	AuctionCancelled          AuctionState = "cancelled"
	AuctionFailed             AuctionState = "failed"
	AuctionValidation         AuctionState = "validation"
)

type Auction struct {
	ID              string
	PlayerID        string
	CurrencyID      int
	Price           int64
	Rnd             int64
	ValidUntil      int64
	OfferValidUntil int64
	Signature       string
	State           AuctionState
	StateExtra      string
	PaymentURL      string
	Seller          string
}

func NewAuction() *Auction {
	auction := Auction{}
	auction.State = AuctionStarted
	return &auction
}

func (b Auction) ComputeID() (string, error) {
	playerId, ok := new(big.Int).SetString(b.PlayerID, 10)
	if !ok {
		return "", errors.New("invalid playerId")
	}
	id, err := signer.ComputeAuctionId(
		uint8(b.CurrencyID),
		big.NewInt(b.Price),
		big.NewInt(b.Rnd),
		b.ValidUntil,
		b.OfferValidUntil,
		playerId,
	)
	if err != nil {
		return "", errors.New("invalid auctio id")
	}
	return id.String()[2:], nil
}
