package gql

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type CreateBidInput struct {
	Signature  Sign
	Auction    Sign
	ExtraPrice int
	Rnd        int
	TeamId     string
}

func (b *Resolver) CreateBid(args struct{ Input CreateBidInput }) (Sign, error) {
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
