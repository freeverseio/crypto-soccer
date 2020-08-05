package gql

import (
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CreateAuction(args struct{ Input input.CreateAuctionInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] create auction %+v", args.Input)

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
	if !isOwner {
		return id, fmt.Errorf("signer is not the owner of playerId %v", args.Input.PlayerId)
	}

	isValidForBlockchain, err := args.Input.IsValidForBlockchain(b.contracts)
	if err != nil {
		return id, err
	}
	if !isValidForBlockchain {
		return id, fmt.Errorf("blockchain says no")
	}

	return id, b.push(args.Input)
}
