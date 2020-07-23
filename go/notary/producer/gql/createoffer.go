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

	id, err := args.Input.ID()
	if err != nil {
		return graphql.ID(""), err
	}

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature()
	if err != nil {
		return id, err
	}
	if !isValid {
		return id, errors.New("Invalid signature")
	}

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return id, err
	}

	if isOwner {
		return id, fmt.Errorf("signer is the owner of playerId %v you can't make an offer for your player", args.Input.PlayerId)
	}

	seller, err := args.Input.GetOwner(b.contracts)
	args.Input.Seller = seller.Hex()
	fmt.Printf("We make it to before isvalid, %v\n", args.Input)
	isValidForBlockchain, err := args.Input.IsValidForBlockchain(b.contracts)
	if err != nil {
		fmt.Printf("isvalid error: %v\n", isValidForBlockchain)

		return id, err
	}
	if !isValidForBlockchain {
		return id, fmt.Errorf("blockchain says no")
	}

	return id, b.push(args.Input)
}
