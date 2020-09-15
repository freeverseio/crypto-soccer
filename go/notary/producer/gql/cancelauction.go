package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CancelAuction(args struct{ Input input.CancelAuctionInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] cancel auction %+v", args.Input)

	id := args.Input.AuctionId

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}

	if string(args.Input.AuctionId) == "" {
		return id, errors.New("empty AuctionId when trying to cancel an auction")
	}

	auction, err := tx.Auction(string(args.Input.AuctionId))
	if err != nil {
		tx.Rollback()
		return id, err
	}
	if auction == nil {
		tx.Rollback()
		return id, errors.New("could not find an auction to cancel")
	}

	signer, err := args.Input.SignerAddress()
	if err != nil {
		tx.Rollback()
		return id, err
	}

	if signer.Hex() != auction.Seller {
		return id, errors.New("Signer of CancelAuction is not the Seller")
	}

	if err := tx.AuctionCancel(string(args.Input.AuctionId)); err != nil {
		tx.Rollback()
		return id, err
	}

	return id, tx.Commit()
}
