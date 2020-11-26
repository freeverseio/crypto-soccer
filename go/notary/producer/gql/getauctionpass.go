package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/auctionpass"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) GetAuctionPass(args struct{ Input input.GetAuctionPassInput }) (*auctionpass.AuctionPass, error) {
	log.Debugf("GetAuctionPass %v", args)

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("not owner of the team")
	}

	ownerAddress, err := args.Input.SignerAddress()
	owner := string(ownerAddress.Hex())

	tx, err := b.service.Begin()

	storageAuctionPass, err := tx.AuctionPass(owner)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if storageAuctionPass == nil {
		tx.Rollback()
		return nil, errors.New("auciton pass not exists for this owner")
	}

	auctionPass :service= auctionpass.NewAuctionPass(storageAuctionPass.Owner, storageAuctionPass.PurchasedForTeamId, storageAuctionPass.ProductId, storageAuctionPass.Ack)

	return auctionPass, tx.Commit()
}
