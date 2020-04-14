package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CancelAuction(args struct{ Input input.CancelAuctionInput }) (graphql.ID, error) {
	if b.ch != nil {
		select {
		case b.ch <- args.Input:
		default:
			log.Warning("channel is full")
			return "ciao", errors.New("channel is full")
		}
	}
	return "cippo", nil
}
