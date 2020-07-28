package gql

import (
	"context"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SubmitPlayStorePlayerPurchase(args struct {
	Input input.SubmitPlayStorePlayerPurchaseInput
}) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] submit IAP %+v", args.Input)

	result := graphql.ID("")

	if b.ch == nil {
		return result, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.IsValidSignature()
	if err != nil {
		return result, err
	}
	if !isValid {
		return result, errors.New("Invalid signature")
	}

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return result, err
	}
	if !isOwner {
		return result, errors.New("Not team owner")
	}

	data, err := playstore.DataFromReceipt(args.Input.Receipt)
	if err != nil {
		return result, err
	}

	ctx := context.Background()
	client, err := playstore.NewGoogleClientService(b.googleCredentials)
	if err != nil {
		return result, err
	}

	_, err = client.GetPurchase(
		ctx,
		data.PackageName,
		data.ProductId,
		data.PurchaseToken,
	)
	if err != nil {
		return result, err
	}

	return graphql.ID(data.PurchaseToken), b.push(args.Input)
}
