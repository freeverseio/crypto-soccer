package gql

import (
	"errors"

	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

type CreateBidInput struct {
	Signature  string
	Auction    graphql.ID
	ExtraPrice int
	Rnd        int
	TeamId     string
}

func (b *Resolver) CreateBid(args struct{ Input CreateBidInput }) (graphql.ID, error) {
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
