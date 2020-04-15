package gql

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) CreateAuction(args struct{ Input input.CreateAuctionInput }) (graphql.ID, error) {
	log.Debugf("CreateAuction %v", args)

	id := args.Input.ID()

	if b.ch == nil {
		return graphql.ID(id), errors.New("internal error: no channel")
	}

	isValid, err := args.Input.VerifySignature()
	if err != nil {
		return graphql.ID(id), err
	}
	if !isValid {
		return graphql.ID(id), errors.New("Invalid signature")
	}

	signerAddress, err := args.Input.SignerAddress()
	if err != nil {
		return graphql.ID(id), err
	}

	playerId, _ := new(big.Int).SetString(args.Input.PlayerId, 10)
	owner, err := b.contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	if err != nil {
		return graphql.ID(id), err
	}

	if signerAddress != owner {
		return graphql.ID(id), fmt.Errorf("signer %v is not the owner (%v) of playerId %v", signerAddress.Hex(), owner.Hex(), args.Input.PlayerId)
	}

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return graphql.ID(id), errors.New("channel is full")
	}

	return graphql.ID(id), nil
}
