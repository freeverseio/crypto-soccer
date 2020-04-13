package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/signer"
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
	if b.c == nil {
		return graphql.ID(""), errors.New("internal error: no channel")
	}

	hash, err := signer.HashSellMessage(
		b.market,
		args.Input.CurrencyId,
		args.Input.Price,
		args.Input.Rnd,
		args.Input.ValidUntil,
		args.Input.PlayerId,
	)

	select {
	case b.c <- args.Input:
	default:
		log.Warning("channel is full")
		return graphql.ID("ciao"), errors.New("channel is full")
	}

	return graphql.ID(args.Input.Signature), nil
}
