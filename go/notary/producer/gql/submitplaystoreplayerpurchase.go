package gql

import (
	"context"
	"errors"
	"fmt"

	"github.com/awa/go-iap/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

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

	if purchase.PurchaseType != nil {
		if *purchase.PurchaseType == 0 { // Test
			return result, nil
		} else {
			return result, fmt.Errorf("orderId %v with unknown purchase type %v", purchase.OrderId, *purchase.PurchaseType)
		}
	}

	if purchase.AcknowledgementState == 1 { // Acknowledged
		return result, fmt.Errorf("OrderId %v is already acknowledged", purchase.OrderId)
	}
	if purchase.AcknowledgementState != 0 { // unknown state
		return result, fmt.Errorf("OrderId %v state %v unknown", purchase.OrderId, purchase.AcknowledgementState)
	}

	if purchase.ConsumptionState == 1 { // consumed
		return result, fmt.Errorf("OrderId %v is already consumed", purchase.OrderId)
	}
	if purchase.ConsumptionState != 0 { // unknown state
		return result, fmt.Errorf("orderId %v consuption state %v unknown", purchase.OrderId, purchase.ConsumptionState)
	}

	if purchase.PurchaseState == 1 {
		return result, fmt.Errorf("orderId %v is cancelled", purchase.OrderId)
	}
	if purchase.PurchaseState == 2 {
		return result, fmt.Errorf("orderId %v is pending", purchase.OrderId)
	}
	if purchase.PurchaseState != 0 {
		return result, fmt.Errorf("orderId %v purchase state %v unknown", purchase.OrderId, purchase.PurchaseState)
	}

	log.Infof("%+v", purchase)

	// check if the player is valid

	return result, errors.New("not implemented")
}
