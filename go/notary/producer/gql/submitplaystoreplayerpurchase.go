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
		return result, fmt.Errorf("orderId %v has an invalid playerId %v", data.OrderId, args.Input.PlayerId)
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
	players, err := worldplayer.CreateWorldPlayerBatch(
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
