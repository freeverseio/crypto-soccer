package gql

import (
	"errors"
	"time"

	"github.com/freeverseio/crypto-soccer/go/storage/postgres"

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

	tx, err := b.db.Begin()
	if err != nil {
		return false, err
	}

	service := postgres.NewTeamStorageService(tx)
	promoTimeout, err := service.TeamPromoTimeout(args.Input.TeamId)
	if err != nil {
		return false, err
	}

	epoch := time.Now().Unix()
	if promoTimeout < uint32(epoch) {
		return false, errors.New("team has no promo")
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
