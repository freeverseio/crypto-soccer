package gql

import (
	"errors"

	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

type CancelAuctionInput struct {
	ID graphql.ID
}

func (b *Resolver) CancelAuction(args struct{ Input CancelAuctionInput }) (graphql.ID, error) {
	if b.c != nil {
		select {
		case b.c <- args.Input:
		default:
			log.Warning("channel is full")
			return "ciao", errors.New("channel is full")
		}
	}
	return "cippo", nil
}
