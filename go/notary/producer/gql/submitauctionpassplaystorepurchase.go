package gql

import (
	"context"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/playstore"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/graph-gophers/graphql-go"
	log "github.com/sirupsen/logrus"
)

func (b *Resolver) SubmitAuctionPassPlayStorePurchase(args struct {
	Input input.SubmitAuctionPassPlayStorePurchaseInput
}) (graphql.ID, error) {
	log.Infof("[notary|producer|gql] submit AuctionPass IAP %+v", args.Input)

	result := graphql.ID("")

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

	tx, err := b.service.Begin()
	if err != nil {
		return result, err
	}
	if err := submitAuctionPassPlayStorePurchase(tx, args.Input); err != nil {
		tx.Rollback()
		return result, err
	}

	return graphql.ID(data.PurchaseToken), tx.Commit()
}

func submitAuctionPassPlayStorePurchase(
	service storage.Tx,
	in input.SubmitAuctionPassPlayStorePurchaseInput,
) error {
	log.Debugf("SubmitAuctionPassPlayStorePurchase %+v", in)

	data, err := playstore.DataFromReceipt(in.Receipt)
	if err != nil {
		return err
	}

	order := storage.NewAuctionPassPlaystoreOrder()
	order.OrderId = data.OrderId
	order.PackageName = data.PackageName
	order.ProductId = data.ProductId
	order.PurchaseToken = data.PurchaseToken
	order.TeamId = string(in.TeamId)
	owner, err := in.SignerAddress()
	if err != nil {
		return err
	}
	order.Owner = string(owner.Hex())
	order.Signature = in.Signature

	if err := service.AuctionPassPlayStoreInsert(*order); err != nil {
		return err
	}

	return nil
}
