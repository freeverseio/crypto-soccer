package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CancelAuction(args struct{ Input input.CancelAuctionInput }) (graphql.ID, error) {
	log.Debugf("CancelAuction %v", args)

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

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return id, errors.New("channel is full")
	}
	return id, nil
}
