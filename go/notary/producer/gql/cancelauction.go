package gql

import (
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CancelAuction(args struct{ Input input.CancelAuctionInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] cancel auction %+v", args.Input)

	id := args.Input.AuctionId

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature()
	if err != nil {
		return id, err
	}
	if !isValid {
		return id, errors.New("Invalid signature")
	}

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}
	if err := cancelAuction(tx, args.Input); err != nil {
		tx.Rollback()
		return id, err
	}

	return id, tx.Commit()
}

func cancelAuction(service storage.Tx, in input.CancelAuctionInput) error {
	auction, err := service.Auction(string(in.AuctionId))
	if err != nil {
		return err
	}
	if auction == nil {
		return errors.New("trying to cancel a nil auction")
	}
	if auction.State != storage.AuctionStarted {
		return fmt.Errorf("not possible to cancel an auction in state %v", auction.State)
	}

	auction.State = storage.AuctionCancelled
	return service.AuctionUpdate(*auction)
}
