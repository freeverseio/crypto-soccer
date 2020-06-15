package gql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/worldplayer"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SubmitPlayStorePlayerPurchase(args struct {
	Input input.SubmitPlayStorePlayerPurchaseInput
}) (graphql.ID, error) {
	log.Infof("SubmitPlayStorePlayerPurchase %v", args)

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

	purchase, err := client.GetPurchase(
		ctx,
		data.PackageName,
		data.ProductId,
		data.PurchaseToken,
	)
	if err != nil {
		return result, err
	}

	validator := playstore.NewPurchaseValidator(*purchase)
	if validator.IsCanceled() {
		return result, errors.New("cancelled order")
	}
	if validator.IsAcknowledged() {
		return result, errors.New("already acknowledged order")
	}

	worldPlayerService := worldplayer.NewWorldPlayerService(b.contracts, b.namesdb)
	worldPlayer, err := worldPlayerService.GetWorldPlayer(
		string(args.Input.PlayerId),
		string(args.Input.TeamId),
		time.Now().Unix(),
	)
	if err != nil {
		return result, err
	}
	if worldPlayer == nil {
		return result, fmt.Errorf("orderId %v has an invalid playerId %v", data.OrderId, args.Input.PlayerId)
	}
	if worldPlayer.ProductId() != data.ProductId {
		return result, fmt.Errorf("orderId %v has an productId mismatch %v != %v", data.OrderId, worldPlayer.ProductId(), data.ProductId)
	}

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return result, errors.New("channel is full")
	}

	return graphql.ID(data.PurchaseToken), nil
}
