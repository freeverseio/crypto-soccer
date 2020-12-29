package gql

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) DismissPlayer(args struct{ Input input.DismissPlayerInput }) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] dismiss player %+v", args.Input)

	id := graphql.ID(args.Input.PlayerId)

	if b.ch == nil {
		return id, errors.New("internal error: no channel")
	}

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return id, err
	}
	if !isOwner {
		return id, errors.New("not player owner")
	}

	tx, err := b.service.Begin()
	if err != nil {
		return graphql.ID(""), err
	}

	isPlayerOnSale, auctionID, err := args.Input.IsPlayerOnSale(tx)
	if err != nil {
		return id, err
	}

	if isPlayerOnSale {
		if err := tx.AuctionCancel(string(auctionID)); err != nil {
			tx.Rollback()
			return id, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return id, err
	}

	return id, b.push(args.Input)
}
