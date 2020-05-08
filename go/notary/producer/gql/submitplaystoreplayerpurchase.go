package gql

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

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

	value := int64(1000) // TODO: value is forced to be 1000

	// check if the player is valid
	players, err := CreateWorldPlayerBatch(
		b.contracts,
		b.namesdb,
		value,
		string(args.Input.TeamId),
		time.Now().Unix(),
	)
	if err != nil {
		return result, err
	}

	i := sort.Search(len(players), func(i int) bool {
		return players[i].PlayerId() == args.Input.PlayerId
	})
	if i >= len(players) {
		return result, fmt.Errorf("orderId %v has an invalid playerId %v", purchase.OrderId, args.Input.PlayerId)
	}

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return result, errors.New("channel is full")
	}

	return args.Input.PlayerId, errors.New("not implemented")
}
