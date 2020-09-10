package gql

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CreateOffer(args struct{ Input input.CreateOfferInput }) (graphql.ID, error) {
	log.Debugf("CreateOffer %v", args)

	id, err := args.Input.ID(b.contracts)
	if err != nil {
		return graphql.ID(""), err
	}

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature(b.contracts)
	if err != nil {
		return id, err
	}
	if !isValid {
		return id, errors.New("Invalid signature")
	}

	isPlayerOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return id, err
	}

	if isPlayerOwner {
		return id, fmt.Errorf("signer is the owner of playerId %v you can't make an offer for your player", args.Input.PlayerId)
	}

	isTeamOwner, err := args.Input.IsSignerTeamOwner(b.contracts)
	if err != nil {
		return id, err
	}

	if !isTeamOwner {
		return id, fmt.Errorf("signer is not the owner of teamId %v", args.Input.BuyerTeamId)
	}

	tx, err := b.service.Begin()
	if err != nil {
		return id, err
	}

	isPlayerOnSale, err := args.Input.IsPlayerOnSale(tx)
	if err != nil {
		return id, err
	}

	if isPlayerOnSale {
		return id, fmt.Errorf("Player is already on sale %v", args.Input.PlayerId)
	}

	isPlayerFrozen, err := args.Input.IsPlayerFrozen(b.contracts)
	if err != nil {
		return id, err
	}

	if isPlayerFrozen {
		return id, fmt.Errorf("Player is frozen %v", args.Input.PlayerId)
	}

	if err := createOffer(tx, args.Input, b.contracts); err != nil {
		tx.Rollback()
		return id, err
	}

	return id, tx.Commit()
}

func createOffer(service storage.Tx, in input.CreateOfferInput, contracts contracts.Contracts) error {
	offer := storage.NewOffer()
	id, err := in.ID(contracts)
	if err != nil {
		return err
	}
	offer.ID = string(id)
	offer.Rnd = int64(in.Rnd)
	offer.PlayerID = in.PlayerId
	offer.CurrencyID = int(in.CurrencyId)
	offer.Price = int64(in.Price)
	if offer.ValidUntil, err = strconv.ParseInt(in.ValidUntil, 10, 32); err != nil {
		fmt.Printf("%d of type %T", offer.ValidUntil, offer.ValidUntil)
	}
	offer.Signature = in.Signature
	offer.State = storage.OfferStarted
	offer.StateExtra = ""
	signerAddress, err := in.SignerAddress(contracts)
	if err != nil {
		return err
	}
	offer.Buyer = signerAddress.Hex()
	offer.Seller = in.Seller
	offer.BuyerTeamID = in.BuyerTeamId
	if err = service.OfferInsert(*offer); err != nil {
		return err
	}

	return nil
}
