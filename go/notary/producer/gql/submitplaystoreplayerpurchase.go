package gql

import (
	"context"
	"errors"

	"github.com/awa/go-iap/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

const GooglePackage = "com.freeverse.phoenix"

func (b *Resolver) SubmitPlayStorePlayerPurchase(args struct {
	Input input.SubmitPlayStorePlayerPurchaseInput
}) (graphql.ID, error) {
	log.Debugf("SubmitPlayStorePlayerPurchase %v", args)

	result := graphql.ID(args.Input.PlayerId)

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

	client, err := playstore.New(b.googleCredentials)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	purchase, err := client.VerifyProduct(
		ctx,
		string(args.Input.PackageName),
		string(args.Input.ProductId),
		args.Input.PurchaseToken,
	)
	if err != nil {
		return result, err
	}

	log.Infof("%+v", purchase)

	return result, errors.New("not implemented")
}
