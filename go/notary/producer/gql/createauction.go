package gql

import (
	"errors"

	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

type CreateAuctionInput struct {
	Signature  string
	PlayerId   string
	CurrencyId int
	Price      int
	Rnd        int
	ValidUntil string
}

func (b *Resolver) CreateAuction(args struct{ Input CreateAuctionInput }) (graphql.ID, error) {
	if b.c != nil {
		select {
		case b.c <- args.Input:
		default:
			log.Warning("channel is full")
			return graphql.ID("ciao"), errors.New("channel is full")
		}
	}
	return graphql.ID("cippo"), nil
}
