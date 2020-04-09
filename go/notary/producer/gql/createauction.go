package gql

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type AuctionInput struct {
	PlayerId   string
	CurrencyId int
	Price      int
	Rnd        int
	ValidUntil string
	Signature  string
}

func (b *Resolver) CreateAuction(args struct{ Input AuctionInput }) (UUID, error) {
	if b.c != nil {
		select {
		case b.c <- args.Input:
		default:
			log.Warning("channel is full")
			return UUID("ciao"), errors.New("channel is full")
		}
	}
	return UUID("reciao"), nil
}
