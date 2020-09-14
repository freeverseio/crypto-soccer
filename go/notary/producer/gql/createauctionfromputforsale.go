package gql

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CreateAuctionFromPutForSale(args struct {
	Input input.CreatePutPlayerForSaleInput
}) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] create auction %+v", args.Input)

	id, err := args.Input.ID()
	if err != nil {
		return graphql.ID(""), err
	}

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	isOwner, err := args.Input.IsSignerOwnerOfPlayer(b.contracts)
	if err != nil {
		return id, err
	}
	if !isOwner {
		return id, fmt.Errorf("signer is not the owner of playerId %v", args.Input.PlayerId)
	}

	isValidForBlockchain, err := args.Input.IsValidForBlockchainFreeze(b.contracts)
	if err != nil {
		return id, err
	}
	if !isValidForBlockchain {
		return id, fmt.Errorf("blockchain says no")
	}

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}
	if err := CreateAuctionFromPutForSale(tx, args.Input); err != nil {
		tx.Rollback()
		return id, err
	}

	return id, tx.Commit()
}

func CreateAuctionFromPutForSale(tx storage.Tx, in input.CreatePutPlayerForSaleInput) error {
	auction := storage.NewAuction()
	id, err := in.ID()
	if err != nil {
		return err
	}
	auction.ID = string(id)
	auction.Rnd = int64(in.Rnd)
	auction.PlayerID = in.PlayerId
	auction.CurrencyID = int(in.CurrencyId)
	auction.Price = int64(in.Price)
	if auction.ValidUntil, err = strconv.ParseInt(in.ValidUntil, 10, 64); err != nil {
		return fmt.Errorf("invalid validUntil %v", in.ValidUntil)
	}
	auction.OfferValidUntil = int64(0)
	auction.Signature = in.Signature
	auction.State = storage.AuctionStarted
	auction.StateExtra = ""
	auction.PaymentURL = ""
	signerAddress, err := in.SignerAddress()
	if err != nil {
		return err
	}
	auction.Seller = signerAddress.Hex()
	if err = tx.AuctionInsert(*auction); err != nil {
		return err
	}

	return nil
}
