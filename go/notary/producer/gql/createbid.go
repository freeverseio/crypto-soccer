package gql

import (
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CreateBid(args struct{ Input input.CreateBidInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] create bid %+v", args.Input)

	id, err := args.Input.ID(b.contracts)
	if err != nil {
		return graphql.ID(""), err
	}

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature(b.contracts)
	if err != nil {
		return graphql.ID(id), err
	}
	if !isValid {
		return graphql.ID(id), errors.New("Invalid signature")
	}

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}
	if err := createBid(tx, args.Input); err != nil {
		tx.Rollback()
		return id, err
	}

	return id, tx.Commit()
}

func createBid(tx storage.Tx, in input.CreateBidInput) error {
	auction, err := tx.Auction(string(in.AuctionId))
	if err != nil {
		return err
	}

	if auction == nil {
		return fmt.Errorf("No auction for bid %v", in)
	}

	bid := storage.NewBid()
	bid.AuctionID = string(in.AuctionId)
	bid.ExtraPrice = int64(in.ExtraPrice)
	bid.Rnd = int64(in.Rnd)
	bid.TeamID = in.TeamId
	bid.Signature = in.Signature
	bid.State = storage.BidAccepted
	bid.StateExtra = ""
	bid.PaymentID = ""
	bid.PaymentURL = ""
	bid.PaymentDeadline = 0

	return tx.BidInsert(*bid)
}
