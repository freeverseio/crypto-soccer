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
	if err := tx.AuctionCancel(string(args.Input.AuctionId)); err != nil {
		tx.Rollback()
		return id, err
	}

	return id, tx.Commit()
}
