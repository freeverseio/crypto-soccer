package gql

import (
	"context"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/googleplaystoreutils"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SubmitPlayStorePlayerPurchase(args struct {
	Input input.SubmitPlayStorePlayerPurchaseInput
}) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] submit IAP %+v", args.Input)

	result := graphql.ID("")

	isOwner, err := args.Input.IsSignerOwner(b.contracts)
	if err != nil {
		return result, err
	}
	if !isOwner {
		return result, errors.New("Not team owner")
	}

	data, err := googleplaystoreutils.DataFromReceipt(args.Input.Receipt)
	if err != nil {
		return result, err
	}

	ctx := context.Background()
	client, err := googleplaystoreutils.NewGoogleClientService(b.googleCredentials)
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

	tx, err := b.service.Begin()
	if err != nil {
		return result, err
	}
	if err := submitPlayStorePlayerPurchase(tx, args.Input); err != nil {
		tx.Rollback()
		return result, err
	}

	return graphql.ID(data.PurchaseToken), tx.Commit()
}

func submitPlayStorePlayerPurchase(
	service storage.Tx,
	in input.SubmitPlayStorePlayerPurchaseInput,
) error {
	log.Debugf("SubmitPlayStorePlayerPurchase %+v", in)

	data, err := googleplaystoreutils.DataFromReceipt(in.Receipt)
	if err != nil {
		return err
	}

	order := storage.NewPlaystoreOrder()
	order.OrderId = data.OrderId
	order.PackageName = data.PackageName
	order.ProductId = data.ProductId
	order.PurchaseToken = data.PurchaseToken
	order.PlayerId = string(in.PlayerId)
	order.TeamId = string(in.TeamId)
	order.Signature = in.Signature

	if err := service.PlayStoreInsert(*order); err != nil {
		return err
	}

	return nil
}
