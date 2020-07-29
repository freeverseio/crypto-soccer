package gql

import (
	"errors"
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CreateOffer(args struct{ Input input.CreateOfferInput }) (bool, error) {
	log.Debugf("CreateOffer %v", args)

	if b.ch == nil {
		return false, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature(b.contracts)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, errors.New("Invalid signature")
	}

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return false, err
	}

	if isOwner {
		return false, fmt.Errorf("signer is the owner of playerId %v you can't make an offer for your player", args.Input.PlayerId)
	}

	seller, err := args.Input.GetOwner(b.contracts)
	args.Input.Seller = seller.Hex()

	return true, b.push(args.Input)
}
