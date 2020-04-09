package gql

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type CancelAuctionInput struct {
	Signature Sign
}

func (b *Resolver) CancelAuction(args struct{ Input CancelAuctionInput }) (Sign, error) {
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
