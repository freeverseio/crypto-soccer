package gql

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type ConsumePromoInput struct {
	Signature string
	PlayerId  string
	TeamId    string
}

func (b *Resolver) ConsumePromo(args struct{ Input ConsumePromoInput }) (bool, error) {
	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return false, err
	}

	if !isOwner {
		return false, errors.New("signer is not the owner of the team")
	}

	if b.ch != nil {
		select {
		case b.ch <- args.Input:
		default:
			log.Warning("ConsumePromo: channel is full, discarding value")
			return false, errors.New("channel is full, discarding value")
		}
	}
	return true, nil
}
