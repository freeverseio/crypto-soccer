package gql

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/awa/go-iap/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

const GooglePackage = "com.freeverse.phoenix"
const GoogleProductID = "coinpack_45"

func (b *Resolver) SubmitPlayerPurchase(args struct {
	Input input.SubmitPlayerPurchaseInput
}) (graphql.ID, error) {
	log.Debugf("SubmitPlayerPurchase %v", args)

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

	// You need to prepare a public key for your Android app's in app billing
	// at https://console.developers.google.com.
	jsonKey, err := ioutil.ReadFile("jsonKey.json")
	if err != nil {
		return result, err
	}
	client, err := playstore.New(jsonKey)
	if err != nil {
		return result, err
	}
	ctx := context.Background()
	purchase, err := client.VerifyProduct(ctx, GooglePackage, GoogleProductID, string(args.Input.PurchaseId))
	if err != nil {
		return result, err
	}

	log.Infof("%+v", purchase)

	return result, errors.New("not implemented")
}
