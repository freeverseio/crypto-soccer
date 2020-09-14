package gql

import (
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

//

func (b *Resolver) CreateBid(args struct{ Input input.CreateBidInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] create bid %+v", args.Input)

	if b.ch == nil {
		return graphql.ID(""), errors.New("internal error: no channel")
	}

	isOwner, err := args.Input.IsSignerOwnerOfTeam(b.contracts)
	if err != nil {
		return graphql.ID(""), err
	}
	if !isOwner {
		return graphql.ID(""), fmt.Errorf("signer is not the owner of teamId %v", args.Input.TeamId)
	}

	tx, err := b.service.Begin()
	if err != nil {
		return graphql.ID(""), err
	}

	if err := createBid(tx, args.Input); err != nil {
		tx.Rollback()
		return graphql.ID(""), err
	}

	return args.Input.AuctionId, tx.Commit()
}

func createBid(tx storage.Tx, in input.CreateBidInput) error {
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
