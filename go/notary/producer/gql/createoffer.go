package gql

import (
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
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
		return id, fmt.Errorf("signer is not the owner of teamId %v", args.Input.TeamId)
	}

	tx, err := b.db.Begin()
	if err != nil {
		return id, err
	}

	isPlayerOnSale, err := args.Input.IsPlayerOnSale(b.contracts, tx)
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

	seller, err := args.Input.GetOwner(b.contracts)
	args.Input.Seller = seller.Hex()

	return id, b.push(args.Input)
}
