package gql

import (
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) ProcessAuction(args struct{ Input input.ProcessAuctionInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] process auction %+v", args.Input)

	return args.Input.Id, b.push(args.Input)
}
