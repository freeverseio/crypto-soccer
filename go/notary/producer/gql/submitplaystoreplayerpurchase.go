package gql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/awa/go-iap/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func GetOrderId(
	credentials []byte,
	packageName string,
	productID string,
	token string,
) (string, error) {
	client, err := playstore.New(credentials)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	purchase, err := client.VerifyProduct(
		ctx,
		packageName,
		productID,
		token,
	)
	if err != nil {
		return "", err
	}

	if purchase.PurchaseType != nil {
		if *purchase.PurchaseType == 0 { // Test
			log.Infof("[TEST] OrderId %v", purchase.OrderId)
			return purchase.OrderId, nil
		}
		return purchase.OrderId, fmt.Errorf("orderId %v with unknown purchase type %v", purchase.OrderId, *purchase.PurchaseType)
	}

	if purchase.AcknowledgementState == 1 { // Acknowledged
		return purchase.OrderId, fmt.Errorf("OrderId %v is already acknowledged", purchase.OrderId)
	}
	if purchase.AcknowledgementState != 0 { // unknown state
		return purchase.OrderId, fmt.Errorf("OrderId %v state %v unknown", purchase.OrderId, purchase.AcknowledgementState)
	}

	if purchase.ConsumptionState == 1 { // consumed
		return purchase.OrderId, fmt.Errorf("OrderId %v is already consumed", purchase.OrderId)
	}
	if purchase.ConsumptionState != 0 { // unknown state
		return purchase.OrderId, fmt.Errorf("orderId %v consuption state %v unknown", purchase.OrderId, purchase.ConsumptionState)
	}

	if purchase.PurchaseState == 1 {
		return purchase.OrderId, fmt.Errorf("orderId %v is cancelled", purchase.OrderId)
	}
	if purchase.PurchaseState == 2 {
		return purchase.OrderId, fmt.Errorf("orderId %v is pending", purchase.OrderId)
	}
	if purchase.PurchaseState != 0 {
		return purchase.OrderId, fmt.Errorf("orderId %v purchase state %v unknown", purchase.OrderId, purchase.PurchaseState)
	}
	return purchase.OrderId, nil
}

func (b *Resolver) SubmitPlayStorePlayerPurchase(args struct {
	Input input.SubmitPlayStorePlayerPurchaseInput
}) (graphql.ID, error) {
	log.Debugf("SubmitPlayStorePlayerPurchase %v", args)

	result := graphql.ID(args.Input.PlayerId)

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

	orderId, err := GetOrderId(
		b.googleCredentials,
		string(args.Input.PackageName),
		string(args.Input.ProductId),
		args.Input.PurchaseToken,
	)
	if err != nil {
		return result, err
	}

	value := int64(1000)     // TODO: value is forced to be 1000
	maxPotential := uint8(9) // TODO: value is forced to be 9

	isValidPlayer, err := b.IsValidPlayer(
		string(args.Input.PlayerId),
		value,
		maxPotential,
		string(args.Input.TeamId),
		time.Now().Unix(),
	)
	if err != nil {
		return result, err
	}
	if !isValidPlayer {
		return result, fmt.Errorf("orderId %v has an invalid playerId %v", orderId, args.Input.PlayerId)
	}

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return result, errors.New("channel is full")
	}

	return args.Input.PlayerId, nil
}

func (b Resolver) IsValidPlayer(
	playerId string,
	value int64,
	maxPotential uint8,
	teamId string,
	epoch int64,
) (bool, error) {
	players, err := CreateWorldPlayerBatch(
		b.contracts,
		b.namesdb,
		value,
		maxPotential,
		teamId,
		epoch,
	)
	if err != nil {
		return false, err
	}

	for _, player := range players {
		if string(player.PlayerId()) == playerId {
			return true, nil
		}
	}

	return false, nil
}
