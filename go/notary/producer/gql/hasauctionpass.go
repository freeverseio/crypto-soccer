package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) HasAuctionPass(args struct{ Input input.HasAuctionPassInput }) (*bool, error) {
	log.Debugf("HasAuctionPass %v", args)

	tx, err := b.service.Begin()
	owner := string(args.Input.Owner)
	storageAuctionPass, err := tx.AuctionPass(owner)
	auctionPassExists := false
	if err != nil {
		tx.Rollback()
		return &auctionPassExists, err
	}

	if storageAuctionPass == nil {
		tx.Rollback()
		return &auctionPassExists, errors.New("auciton pass not exists for this owner")
	}

	auctionPassExists = owner == storageAuctionPass.Owner

	return &auctionPassExists, tx.Commit()
}
