package gql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/worldplayer"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SubmitPlayStorePlayerPurchase(args struct {
	Input input.SubmitPlayStorePlayerPurchaseInput
}) (*worldplayer.WorldPlayer, error) {
	log.Infof("SubmitPlayStorePlayerPurchase %v", args)

	if b.ch == nil {
		return nil, errors.New("internal error: no channel")
	}

	isValid, err := args.Input.IsValidSignature()
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, errors.New("Invalid signature")
	}

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("Not team owner")
	}

	data, err := playstore.DataFromReceipt(args.Input.Receipt)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client, err := playstore.NewGoogleClientService(b.googleCredentials)
	if err != nil {
		return nil, err
	}

	purchase, err := client.GetPurchase(
		ctx,
		data.PackageName,
		data.ProductId,
		data.PurchaseToken,
	)
	if err != nil {
		return nil, err
	}

	validator := playstore.NewPurchaseValidator(*purchase)
	if validator.IsCanceled() {
		return nil, errors.New("cancelled order")
	}
	if validator.IsAcknowledged() {
		return nil, errors.New("already acknowledged order")
	}

	worldPlayerService := worldplayer.NewWorldPlayerService(b.contracts, b.namesdb)
	worldPlayer, err := worldPlayerService.GetWorldPlayer(
		string(args.Input.PlayerId),
		string(args.Input.TeamId),
		time.Now().Unix(),
	)
	if err != nil {
		return nil, err
	}
	if worldPlayer == nil {
		return nil, fmt.Errorf("orderId %v has an invalid playerId %v", data.OrderId, args.Input.PlayerId)
	}
	if worldPlayer.ProductId() != data.ProductId {
		return nil, fmt.Errorf("orderId %v has an productId mismatch %v != %v", data.OrderId, worldPlayer.ProductId(), data.ProductId)
	}

	select {
	case b.ch <- args.Input:
	default:
		log.Warning("channel is full")
		return nil, errors.New("channel is full")
	}

	return worldPlayer, nil
}
